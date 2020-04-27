package main

import "testing"

func TestRegexpMatcher(t *testing.T) {
	rules := []string{
		"reddit.com",
		"^stackoverflow.*$",
		"news.ycombinator.*",
		"twitter.com",
	}

	tt := []struct {
		value    string
		expected bool
	}{
		{"reddit.com", true},
		{"www.reddit.com", true},
		{"reddit.com:443", true},
		{"www.reddit.com:443", true},
		{"ads.reddit.com", true},
		{"reddit.eu", false},
		{"stackoverflow.com", true},
		{"stackoverflow.com:443", true},
		{"ads.stackoverflow.com", false},
		{"news.ycombinator.com", true},
		{"news.ycombinator.eu", true},
		{"ycombinator.com", false},
		{"google.com", false},
		{"twitter.com", true},
		{"www.twitter.com", true},
		{"twitter.com:443", true},
		{"www.twitter.com:443", true},
	}

	matcher := regexpMatcher{}
	matcher.Load(rules)

	for _, tc := range tt {
		if tc.expected != matcher.Match(tc.value) {
			t.Errorf("Match '%s' should be: %t", tc.value, tc.expected)
		}
	}
}

func TestHashMatcher(t *testing.T) {
	rules := []string{
		"reddit.com",
		"reddit.com:443",
		"stackoverflow.com",
		"news.ycombinator.com",
	}

	tt := []struct {
		value    string
		expected bool
	}{
		{"reddit.com", true},
		{"reddit.com:443", true},
		{"ads.reddit.com", false},
		{"reddit.eu", false},
		{"stackoverflow.com", true},
		{"stackoverflow.com:443", true},
		{"ads.stackoverflow.com", false},
		{"news.ycombinator.com", true},
		{"news.ycombinator.eu", false},
		{"ycombinator.com", false},
		{"google.com", false},
	}

	matcher := hashMatcher{}
	matcher.Load(rules)

	for _, tc := range tt {
		if tc.expected != matcher.Match(tc.value) {
			t.Errorf("Match '%s' should be: %t", tc.value, tc.expected)
		}
	}
}
