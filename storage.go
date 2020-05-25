package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Storage struct {
	blocklistPath  string
	blacklistPath  string
	whitelistPath  string
	updateInterval time.Duration
}

func (s *Storage) GetBlocklist(allowCache bool) (Matcher, error) {
	if allowCache {
		if matcher, ok := getBlocklistFromCache(s.updateInterval); ok {
			return matcher, nil
		}
	}

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

func blocklistCachePath() (path string, lastUpdate time.Time, exist bool) {
	files, err := ioutil.ReadDir(blocklistCacheDir())
	if err != nil || len(files) == 0 {
		return "", time.Time{}, false
	}
	path = filepath.Join(blocklistCacheDir(), files[len(files)-1].Name())

	timestamp, err := strconv.ParseInt(filepath.Base(path), 10, 64)
	if err != nil {
		return "", time.Time{}, false
	}
	lastUpdate = time.Unix(timestamp, 0)

	return path, lastUpdate, true
}

func getBlocklistFromCache(updateInterval time.Duration) (Matcher, bool) {
	path, lastUpdate, ok := blocklistCachePath()
	if !ok {
		return nil, false
	}
	if time.Now().After(lastUpdate.Add(updateInterval)) {
		return nil, false
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	rules, err := readLines(file)
	if err != nil {
		return nil, false
	}
	matcher := &hashMatcher{}
	matcher.Load(rules)
	//log.Printf("Blocklists loaded from cache (%v)\n", len(rules))

	return matcher, true
}

func cacheBlocklist(rules []string) error {
	path := filepath.Join(blocklistCacheDir(), fmt.Sprintf("%v", time.Now().Unix()))

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		file.WriteString(fmt.Sprintf("%s\n", rule))
	}
	defer file.Close()

	return nil
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

	cacheBlocklist(rules)

	//log.Printf("Blocklists loaded (%v - %v)\n", len(urls), len(rules))
	matcher := &hashMatcher{}
	matcher.Load(rules)
	return matcher
}

func (s *Storage) GetBlacklist() (Matcher, error) {
	blacklist, err := getMatcherFromFile(s.blacklistPath)
	if err != nil {
		return nil, err
	}
	return blacklist, nil
}

func (s *Storage) GetWhitelist() (Matcher, error) {
	whitelist, err := getMatcherFromFile(s.whitelistPath)
	if err != nil {
		return nil, err
	}
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
