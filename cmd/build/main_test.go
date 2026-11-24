package main

import (
	"maps"
	"net/netip"
	"testing"

	"github.com/split-route/routing-rules/internal/parse"
)

func TestSourceLink(t *testing.T) {
	repo := "https://github.com/o/r/"
	for in, want := range map[string]string{
		"file:///data/ir-ads.txt":   "https://github.com/o/r/blob/HEAD/data/ir-ads.txt",
		"https://example.com/a.txt": "https://example.com/a.txt",
	} {
		if got := sourceLink(repo, in); got != want {
			t.Errorf("sourceLink(%q) = %q, want %q", in, got, want)
		}
	}
	if got := sourceLink("", "file:///data/x.txt"); got != "file:///data/x.txt" {
		t.Errorf("empty repo: got %q", got)
	}
}

func TestProvenanceStats(t *testing.T) {
	a := parse.Domain{Kind: parse.Suffix, Name: "a.com"}
	b := parse.Domain{Kind: parse.Suffix, Name: "b.com"}
	gone := parse.Domain{Kind: parse.Suffix, Name: "gone.com"} // filtered out of the final list
	srcs := []outcome{
		{src: Source{Name: "s1"}, res: parse.Result{Domains: []parse.Domain{a, b}, CIDRs: []netip.Prefix{netip.MustParsePrefix("10.0.0.0/25")}}},
		{src: Source{Name: "s2"}, res: parse.Result{Domains: []parse.Domain{b}, CIDRs: []netip.Prefix{netip.MustParsePrefix("10.0.0.0/24")}}},
		{src: Source{Name: "s3"}, res: parse.Result{Domains: []parse.Domain{gone}}},
	}
	final := catRules{domains: []parse.Domain{a, b}, cidrs: []netip.Prefix{netip.MustParsePrefix("10.0.0.0/16")}}
	got, stats := provenance(srcs, final, map[string]srcStat{})
	if want := "[\"s1\",\"s2\",\"s3\"]\n1\n3\n3\n"; string(got) != want {
		t.Errorf("prov = %q, want %q", got, want)
	}
	// Both sources overlap the merged 10.0.0.0/16, but only s2 has
	// 10.0.0.128/25: one unique prefix for s2, none for s1.
	want := map[string]srcStat{"s1": {3, 1}, "s2": {2, 1}, "s3": {0, 0}}
	if !maps.Equal(stats, want) {
		t.Errorf("stats = %v, want %v", stats, want)
	}
}
