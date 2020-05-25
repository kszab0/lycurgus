package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Storage struct {
	blocklistPath string
	blacklistPath string
	whitelistPath string
}

func (s *Storage) GetBlocklist() (Matcher, error) {
	file, err := getBlocklists(s.blocklistPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	matcher := s.getBlocklist(file)
	return matcher, nil
}

// getBlocklists reads a blocklists file.
// If the file is not exists it returns the default blocklists.
func getBlocklists(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			rc := ioutil.NopCloser(strings.NewReader(defaultBlocklists))
			return rc, nil
		}
		return nil, err
	}
	return file, nil
}

// getBloclist parses a blocklists file
// and initializes a hashMatcher as a blocklist.
func (s *Storage) getBlocklist(r io.Reader) Matcher {
	rules := []string{}
	urls, err := readLines(r)
	if err != nil {
		return nil
	}
	for _, url := range urls {
		hosts, err := parseHostsURL(http.DefaultClient, url)
		if err != nil {
			log.Println("Error reading blocklist: ", url)
			continue
		}
		rules = append(rules, hosts...)
	}

	log.Printf("Blocklists loaded (%v - %v)\n", len(urls), len(rules))
	matcher := &hashMatcher{}
	matcher.Load(rules)
	return matcher
}

func (s *Storage) GetBlacklist() (Matcher, error) {
	blacklist, err := getMatcherFromFile(s.blacklistPath)
	if err != nil {
		return nil, err
	}
	log.Println("Blacklist loaded")
	return blacklist, nil
}

func (s *Storage) GetWhitelist() (Matcher, error) {
	whitelist, err := getMatcherFromFile(s.whitelistPath)
	if err != nil {
		return nil, err
	}
	log.Println("Whitelist loaded")
	return whitelist, nil

}

// getMatcherFromFile loads and initializes a regexpMatcher from a file.
func getMatcherFromFile(path string) (Matcher, error) {
	hosts, err := parseHostsFile(path)
	if err != nil {
		// ignore if file not exists
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	// deal with existing but empty file
	if len(hosts) <= 0 {
		return nil, nil
	}
	matcher := &regexpMatcher{}
	matcher.Load(hosts)
	return matcher, nil
}
