// Fetching sources, with the last-known-good fallback.

package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/split-route/routing-rules/internal/fetch"
	"github.com/split-route/routing-rules/internal/parse"
)

type outcome struct {
	src     Source
	res     parse.Result
	status  string // ok | stale | skipped | FAILED | disabled
	err     error
	fetched time.Time // last successful fetch (now, or from the last-known-good copy)
	changed time.Time // last time the source's rules changed
}

func fetchAll(f *fetch.Fetcher, cfg Config, cats map[string]Category, cacheDir string) []outcome {
	out := make([]outcome, len(cfg.Sources))
	sem := make(chan struct{}, max(1, cfg.Settings.Concurrency))
	var wg sync.WaitGroup
	for i, s := range cfg.Sources {
		out[i].src = s
		if !s.Enabled || !cats[s.Category].Enabled {
			out[i].status = "disabled"
			continue
		}
		wg.Add(1)
		go func(o *outcome) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			res, err := fetchOne(f, o.src)
			lkg := filepath.Join(cacheDir, "sources", o.src.Name+".json")
			if err == nil && len(res.Domains)+len(res.CIDRs) == 0 {
				err = errors.New("parsed zero rules")
			}
			if err == nil {
				o.res, o.status = res, "ok"
				fetched, changed, werr := saveLKG(lkg, res, time.Now().UTC())
				if werr != nil {
					log.Printf("::warning::%s: saving last-known-good: %v", o.src.Name, werr)
				}
				o.fetched, o.changed = fetched, changed
				return
			}
			o.err = err
			if res, meta, lerr := loadLKG(lkg); lerr == nil {
				o.res, o.status, o.fetched, o.changed = res, "stale", meta.FetchedAt, meta.ChangedAt
				log.Printf("::warning::%s: %v; using last-known-good", o.src.Name, err)
			} else {
				o.status = "FAILED"
				log.Printf("::warning::%s: %v; no last-known-good available", o.src.Name, err)
			}
		}(&out[i])
	}
	wg.Wait()
	return out
}

// transport adds file:// support so sources can point at files in this repo
// (file:///data/ir-ads.txt resolves relative to the working directory).
func transport() http.RoundTripper {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	return t
}

// monthPlaceholder in a source URL is replaced by the current UTC month
// ("2026-11"), falling back to the previous month while the new file is not
// published yet (DB-IP names its monthly dumps this way).
const monthPlaceholder = "{YYYY-MM}"

func getMonthly(ctx context.Context, f *fetch.Fetcher, url string, maxBytes int64) ([]byte, error) {
	if !strings.Contains(url, monthPlaceholder) {
		return f.Get(ctx, url, maxBytes)
	}
	now := time.Now().UTC()
	b, err := f.Get(ctx, strings.ReplaceAll(url, monthPlaceholder, now.Format("2006-01")), maxBytes)
	if err == nil {
		return b, nil
	}
	prev := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	if b, perr := f.Get(ctx, strings.ReplaceAll(url, monthPlaceholder, prev.Format("2006-01")), maxBytes); perr == nil {
		return b, nil
	}
	return nil, err
}

func fetchOne(f *fetch.Fetcher, s Source) (parse.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	url := s.URL
	body, err := getMonthly(ctx, f, url, s.MaxBytes)
	if err != nil {
		return parse.Result{}, err
	}
	opt := parse.Options{Match: parse.Exact, Filter: s.Filter, Country: s.Country, AllowTLD: s.AllowTLD}
	if s.Match == "suffix" {
		opt.Match = parse.Suffix
	}
	if s.Format == "v2fly" {
		base := url[:strings.LastIndexByte(url, '/')+1]
		opt.Include = func(name string) ([]byte, error) {
			if slices.Contains(s.SkipInclude, name) {
				return nil, nil
			}
			return f.Get(ctx, base+name, s.MaxBytes)
		}
	}
	res, err := parse.Formats[s.Format](body, opt)
	if err != nil {
		return parse.Result{}, err
	}
	if s.MaxRules > 0 {
		res = truncate(res, s.MaxRules)
	}
	return res, nil
}

// truncate keeps the first n unique rules in upstream order (upstream lists
// are often ranked, so the head is the most valuable part).
func truncate(r parse.Result, n int) parse.Result {
	var out parse.Result
	seenD := map[parse.Domain]bool{}
	for _, d := range r.Domains {
		if len(out.Domains) >= n {
			break
		}
		if !seenD[d] {
			seenD[d] = true
			out.Domains = append(out.Domains, d)
		}
	}
	seenC := map[netip.Prefix]bool{}
	for _, p := range r.CIDRs {
		if len(out.Domains)+len(out.CIDRs) >= n {
			break
		}
		if !seenC[p] {
			seenC[p] = true
			out.CIDRs = append(out.CIDRs, p)
		}
	}
	return out
}

type lkgFile struct {
	Domains   []string  `json:"domains"`
	CIDRs     []string  `json:"cidrs"`
	Hash      string    `json:"hash,omitempty"`       // sha256 of the rules, to detect upstream changes
	FetchedAt time.Time `json:"fetched_at,omitempty"` // last successful fetch
	ChangedAt time.Time `json:"changed_at,omitempty"` // last fetch whose rules differed
}

// saveLKG stores a source's rules as its last-known-good copy and returns when
// they were fetched and when they last changed (carried over from the previous
// copy if the rules are identical).
func saveLKG(path string, r parse.Result, now time.Time) (fetched, changed time.Time, err error) {
	var l lkgFile
	for _, d := range r.Domains {
		l.Domains = append(l.Domains, ruleString(d))
	}
	for _, p := range r.CIDRs {
		l.CIDRs = append(l.CIDRs, p.String())
	}
	h := sha256.New()
	for _, s := range slices.Sorted(slices.Values(slices.Concat(l.Domains, l.CIDRs))) {
		h.Write([]byte(s + "\n"))
	}
	l.Hash, l.FetchedAt, l.ChangedAt = hex.EncodeToString(h.Sum(nil)), now, now
	if prev, err := readLKGFile(path); err == nil && prev.Hash == l.Hash && !prev.ChangedAt.IsZero() {
		l.ChangedAt = prev.ChangedAt
	}
	b, err := json.Marshal(l)
	if err != nil {
		return
	}
	if err = os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	tmp := path + ".tmp"
	if err = os.WriteFile(tmp, b, 0o644); err != nil {
		return
	}
	return l.FetchedAt, l.ChangedAt, os.Rename(tmp, path)
}

func readLKGFile(path string) (lkgFile, error) {
	var l lkgFile
	b, err := os.ReadFile(path)
	if err != nil {
		return l, err
	}
	return l, json.Unmarshal(b, &l)
}

func loadLKG(path string) (parse.Result, lkgFile, error) {
	l, err := readLKGFile(path)
	if err != nil {
		return parse.Result{}, l, err
	}
	var r parse.Result
	kinds := map[string]parse.Kind{"DOMAIN-SUFFIX": parse.Suffix, "DOMAIN": parse.Exact, "DOMAIN-KEYWORD": parse.Keyword}
	for _, s := range l.Domains {
		k, n, _ := strings.Cut(s, ",")
		if kind, ok := kinds[k]; ok {
			r.Domains = append(r.Domains, parse.Domain{Kind: kind, Name: n})
		}
	}
	for _, s := range l.CIDRs {
		if p, err := netip.ParsePrefix(s); err == nil {
			r.CIDRs = append(r.CIDRs, p)
		}
	}
	if len(r.Domains)+len(r.CIDRs) == 0 {
		return r, l, errors.New("empty last-known-good")
	}
	return r, l, nil
}
