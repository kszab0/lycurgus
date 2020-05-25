package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestGetBlocklists(t *testing.T) {
	path := filepath.Join("testdata", "blocklists")
	rc, err := getBlocklists(path)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	blocklists := string(b)

	expectedBlocklists := "https://asdf.aa\nhttps://qwer.qq"
	expectedBlocklistsWin := "https://asdf.aa\r\nhttps://qwer.qq"

	if blocklists != expectedBlocklists && blocklists != expectedBlocklistsWin {
		t.Errorf("blocklists should be: %v; got: %v", expectedBlocklists, blocklists)
	}
}

func TestGetBlocklistsDefault(t *testing.T) {
	path := filepath.Join("testdata", "nothing")
	rc, err := getBlocklists(path)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	blocklists := string(b)

	if blocklists != defaultBlocklists {
		t.Errorf("blocklists should be: %v; got: %v", defaultBlocklists, blocklists)
	}
}

func TestLoadMatcherFromFile(t *testing.T) {
	path := filepath.Join("testdata", "blacklist")
	matcher, err := getMatcherFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if matcher == nil {
		t.Fatal("Matcher should not be nil")
	}
}

func TestLoadMathcerFromNotExistingFile(t *testing.T) {
	path := filepath.Join("testdata", "nothing")
	matcher, err := getMatcherFromFile(path)
	if err != nil {
		t.Fatal("Error should be nil")
	}
	if matcher != nil {
		t.Fatal("Matcher should be nil")
	}
}

func TestLoadMatcherFromEmptyFile(t *testing.T) {
	path := filepath.Join("testdata", "empty")
	matcher, err := getMatcherFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if matcher != nil {
		t.Fatal("Matcher should be nil")
	}
}

func TestLoadBlacklist(t *testing.T) {
	storage := &Storage{}
	storage.blacklistPath = filepath.Join("testdata", "blacklist")
	matcher, err := storage.GetBlacklist()
	if err != nil {
		t.Errorf("Error should be nil; got: %v", err)
	}
	if matcher == nil {
		t.Errorf("Blacklist should not be nil")
	}

	storage = &Storage{}
	storage.blacklistPath = filepath.Join("testdata", "nothing")
	matcher, err = storage.GetBlacklist()
	if err != nil {
		t.Errorf("Error should be nil; got: %v", err)
	}
	if matcher != nil {
		t.Errorf("Blacklist should be nil")
	}
}

func TestLoadWhitelist(t *testing.T) {
	storage := &Storage{}
	storage.whitelistPath = filepath.Join("testdata", "blacklist")
	matcher, err := storage.GetWhitelist()
	if err != nil {
		t.Errorf("Error should be nil; got: %v", err)
	}
	if matcher == nil {
		t.Errorf("Whitelist should not be nil")
	}

	storage = &Storage{}
	storage.whitelistPath = filepath.Join("testdata", "nothing")
	matcher, err = storage.GetWhitelist()
	if err != nil {
		t.Errorf("Error should be nil; got: %v", err)
	}
	if matcher != nil {
		t.Errorf("Whitelist should be nil")
	}
}
