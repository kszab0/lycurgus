package main

import (
	"regexp"
	"strings"
)

// Matcher decides if an input text is matched by its loaded rules
type Matcher interface {
	// Load loads rules to be used to match input texts
	Load(rules []string)
	Match(text string) bool
}

// regexpMatcher uses regular expression rules to match input text
type regexpMatcher struct {
	regexp *regexp.Regexp
}

// Load loads regular expression rules
func (m *regexpMatcher) Load(rules []string) {
	regexes := strings.Join(rules, "|")
	m.regexp = regexp.MustCompile(regexes)
}

func (m *regexpMatcher) Match(text string) bool {
	return m.regexp.MatchString(text)
}

// hashMatcher uses string rules to match input texts exactly
// or texts with ":443" as suffix
type hashMatcher struct {
	hm map[string]int
}

func (m *hashMatcher) Load(rules []string) {
	m.hm = make(map[string]int)
	for _, rule := range rules {
		m.hm[rule] = 0
	}
}

// Match matches an input text exactly or a text with ":443" as suffix
func (m *hashMatcher) Match(rule string) bool {
	rule = strings.TrimSuffix(rule, ":443")
	_, ok := m.hm[rule]
	return ok
}
