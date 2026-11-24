package main

import (
	"net/netip"
	"path/filepath"
	"testing"
	"time"

	"github.com/split-route/routing-rules/internal/parse"
)

func TestSaveLKGTracksChanges(t *testing.T) {
	path := filepath.Join(t.TempDir(), "src.json")
	a := parse.Result{Domains: []parse.Domain{{Kind: parse.Suffix, Name: "a.com"}}, CIDRs: []netip.Prefix{netip.MustParsePrefix("1.2.3.0/24")}}
	day := func(d int) time.Time { return time.Date(2026, 1, d, 0, 0, 0, 0, time.UTC) }

	if f, c, err := saveLKG(path, a, day(1)); err != nil || !f.Equal(day(1)) || !c.Equal(day(1)) {
		t.Fatalf("first save: %v %v %v", f, c, err)
	}
	// Same rules (different order) later: fetched moves, changed stays.
	same := parse.Result{CIDRs: a.CIDRs, Domains: a.Domains}
	if f, c, _ := saveLKG(path, same, day(2)); !f.Equal(day(2)) || !c.Equal(day(1)) {
		t.Fatalf("unchanged rules: fetched %v changed %v", f, c)
	}
	// Different rules: changed moves too.
	b := parse.Result{Domains: []parse.Domain{{Kind: parse.Suffix, Name: "b.com"}}}
	if f, c, _ := saveLKG(path, b, day(3)); !f.Equal(day(3)) || !c.Equal(day(3)) {
		t.Fatalf("changed rules: fetched %v changed %v", f, c)
	}
	if res, meta, err := loadLKG(path); err != nil || len(res.Domains) != 1 || !meta.ChangedAt.Equal(day(3)) {
		t.Fatalf("load: %v %v %v", res, meta, err)
	}
}
