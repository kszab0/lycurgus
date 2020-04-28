package main

import (
	"testing"

	"gopkg.in/elazarl/goproxy.v1"
)

func TestNewBlocker(t *testing.T) {
	blocker := NewBlocker(WithBlockerEnabled(true))
	if blocker.enabled != defaultBlockerEnabled {
		t.Errorf("enabled should be %v", defaultBlockerEnabled)
	}
	if blocker.proxy == nil {
		t.Errorf("proxy should not be nil")
	}
}

type blocklistMatcher struct{}

func (m *blocklistMatcher) Load(rules []string) {}
func (m *blocklistMatcher) Match(host string) bool {
	if host == "blocklist.com" {
		return true
	}
	return false
}

type blacklistMatcher struct{}

func (m *blacklistMatcher) Load(rules []string) {}
func (m *blacklistMatcher) Match(host string) bool {
	if host == "blacklist.com" {
		return true
	}
	return false
}

type whitelistMatcher struct{}

func (m *whitelistMatcher) Load(rules []string) {}
func (m *whitelistMatcher) Match(host string) bool {
	if host == "whitelist.com" {
		return true
	}
	if host == "blocklist-whitelist.com" {
		return true
	}
	if host == "blacklist-whitelist.com" {
		return true
	}
	return false
}

func TestBlockerEnabled(t *testing.T) {
	host := "blacklist.com"

	blocker := NewBlocker(WithBlockerEnabled(true))
	blocker.blacklist = &blacklistMatcher{}

	if blocker.enabled != true {
		t.Errorf("enabled should be true")
	}
	resp, _ := blocker.handleConnect(host, nil)
	if resp != goproxy.RejectConnect {
		t.Errorf("response should be %v; got: %v", goproxy.RejectConnect, resp)
	}

	blocker.Toggle()

	if blocker.enabled != false {
		t.Errorf("enabled should be false after Toggle()")
	}
	resp, _ = blocker.handleConnect(host, nil)
	if resp != goproxy.OkConnect {
		t.Errorf("response should be %v; got: %v", goproxy.OkConnect, resp)
	}
}

func TestBlockerHandler(t *testing.T) {
	blocker := NewBlocker(WithBlockerEnabled(true))
	blocker.blocklist = &blocklistMatcher{}
	blocker.blacklist = &blacklistMatcher{}
	blocker.whitelist = &whitelistMatcher{}

	tt := []struct {
		host    string
		expResp *goproxy.ConnectAction
	}{
		{
			host:    "blocklist.com",
			expResp: goproxy.RejectConnect,
		},
		{
			host:    "blacklist.com",
			expResp: goproxy.RejectConnect,
		},
		{
			host:    "whitelist.com",
			expResp: goproxy.OkConnect,
		},
		{
			host:    "blocklist-whitelist.com",
			expResp: goproxy.OkConnect,
		},
		{
			host:    "blacklist-whitelist.com",
			expResp: goproxy.OkConnect,
		},
		{
			host:    "host.com",
			expResp: goproxy.OkConnect,
		},
	}

	for _, tc := range tt {
		resp, respHost := blocker.handleConnect(tc.host, nil)
		if resp != tc.expResp {
			t.Errorf("response should be %v; got: %v (%v)", tc.expResp, resp, tc.host)
		}
		if respHost != tc.host {
			t.Errorf("host should be %v; got: %v", tc.host, respHost)
		}
	}
}
