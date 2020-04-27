package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
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
	if runtime.GOOS == "windows" {
		expectedBlocklists = "https://asdf.aa\r\nhttps://qwer.qq"
	}

	if blocklists != expectedBlocklists {
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
	matcher, err := loadMatcherFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if matcher == nil {
		t.Fatal("Matcher should not be nil")
	}
}

func TestLoadMathcerFromNotExistingFile(t *testing.T) {
	path := filepath.Join("testdata", "nothing")
	matcher, err := loadMatcherFromFile(path)
	if !os.IsNotExist(err) {
		t.Fatal("Error should be os.ErrNotExists")
	}
	if matcher != nil {
		t.Fatal("Matcher should be nil")
	}
}

func TestLoadMatcherFromEmptyFile(t *testing.T) {
	path := filepath.Join("testdata", "empty")
	matcher, err := loadMatcherFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if matcher != nil {
		t.Fatal("Matcher should be nil")
	}
}

func TestLoadBlacklist(t *testing.T) {

	app := &App{
		blocker: &Blocker{},
	}
	app.blacklistPath = filepath.Join("testdata", "blacklist")
	if err := app.LoadBlacklist(); err != nil {
		t.Errorf("Error should be nil; got: %v", err)
	}
	if app.blocker.blacklist == nil {
		t.Errorf("Blacklist should not be nil")
	}

	app = &App{
		blocker: &Blocker{},
	}
	app.blacklistPath = filepath.Join("testdata", "nothing")
	if err := app.LoadBlacklist(); err != nil {
		t.Errorf("Error should be nil; got: %v", err)
	}
	if app.blocker.blacklist != nil {
		t.Errorf("Blacklist should be nil")
	}
}

func TestLoadWhitelist(t *testing.T) {

	app := &App{
		blocker: &Blocker{},
	}
	app.whitelistPath = filepath.Join("testdata", "blacklist")
	if err := app.LoadWhitelist(); err != nil {
		t.Errorf("Error should be nil; got: %v", err)
	}
	if app.blocker.whitelist == nil {
		t.Errorf("Whitelist should not be nil")
	}

	app = &App{
		blocker: &Blocker{},
	}
	app.whitelistPath = filepath.Join("testdata", "nothing")
	if err := app.LoadWhitelist(); err != nil {
		t.Errorf("Error should be nil; got: %v", err)
	}
	if app.blocker.blacklist != nil {
		t.Errorf("Whitelist should be nil")
	}
}
