// Build configuration: the sources.yaml schema and the allowlist.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/split-route/routing-rules/internal/parse"
)

type Config struct {
	Settings struct {
		UserAgent   string        `yaml:"user_agent"`
		Timeout     time.Duration `yaml:"timeout"`
		Concurrency int           `yaml:"concurrency"`
		Retries     int           `yaml:"retries"`
		MaxBytes    int64         `yaml:"max_bytes"`
		ListCapWarn int           `yaml:"list_cap_warn"`
		RepoURL     string        `yaml:"repo_url"`
		BaseURL     string        `yaml:"base_url"`
	} `yaml:"settings"`
	Categories []Category `yaml:"categories"`
	Sources    []Source   `yaml:"sources"`
}

type Category struct {
	Name    string `yaml:"name"`
	Action  string `yaml:"action"` // reject | proxy | direct | none
	Enabled bool   `yaml:"enabled"`
}

type Source struct {
	Name        string            `yaml:"name"`
	URL         string            `yaml:"url"`
	Format      string            `yaml:"format"`
	Category    string            `yaml:"category"`
	Enabled     bool              `yaml:"enabled"`
	Match       string            `yaml:"match"`
	Filter      map[string]string `yaml:"filter"`
	Country     string            `yaml:"country"`
	MaxRules    int               `yaml:"max_rules"`
	MaxBytes    int64             `yaml:"max_bytes"`
	SkipInclude []string          `yaml:"skip_includes"`
	AllowTLD    bool              `yaml:"allow_tld"`
	License     string            `yaml:"license"`
}

// minBits rejects absurdly broad prefixes from any feed (a scraped 0.0.0.0/0
// in ir-cidr would send all traffic DIRECT).
var minBits = map[bool]int{true: 8, false: 16} // is4 -> min prefix length

func catEnabled(cfg Config, name string) bool {
	for _, c := range cfg.Categories {
		if c.Name == name {
			return c.Enabled
		}
	}
	return false
}

func readAllowlist(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []string
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		l, _, _ := strings.Cut(sc.Text(), "#")
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		d, ok := parse.NormalizeDomain(l)
		if !ok {
			return nil, fmt.Errorf("%s: invalid domain %q", path, l)
		}
		out = append(out, d)
	}
	return out, sc.Err()
}
