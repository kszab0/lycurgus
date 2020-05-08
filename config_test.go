package main

import (
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	path := filepath.Join("testdata", "config.yml")
	c := &fileConfig{}
	parseFile(c, path)

	if *c.BlockerAddress != ":5555" {
		t.Errorf("BlockerAddress should be %v; got: %v", ":5555", *c.BlockerAddress)
	}
	if *c.AutostartEnabled != false {
		t.Errorf("AutostartEnabled should be %v; got: %v", false, *c.AutostartEnabled)
	}
	if *c.GUIEnabled != true {
		t.Errorf("GUIEnabled should be %v; got: %v", true, *c.GUIEnabled)
	}
	if *c.LogEnabled != false {
		t.Errorf("LogEnabled should be %v; got: %v", false, *c.LogEnabled)
	}

	if c.BlocklistPath != nil {
		t.Errorf("BlocklistPath should be nil; got: %v", *c.BlocklistPath)
	}
	if c.BlacklistPath != nil {
		t.Errorf("BlacklistPath should be nil; got: %v", *c.BlacklistPath)
	}
	if c.WhitelistPath != nil {
		t.Errorf("WhitelistPath should be nil; got: %v", *c.WhitelistPath)
	}
	if c.LogPath != nil {
		t.Errorf("LogPath should be nil; got: %v", *c.LogPath)
	}
	if c.ProxyAddress != nil {
		t.Errorf("ProxyAddress should be nil; got: %v", *c.ProxyAddress)
	}
}

func TestParseFileContent(t *testing.T) {

	tc := struct {
		content  []byte
		expected *Config
	}{
		content: []byte(`address: ":5555"
autostart: false
gui: true
log: fasle`),
		expected: &Config{
			BlockerAddress:   ":5555",
			AutostartEnabled: false,
			GUIEnabled:       true,
			LogEnabled:       false,
		},
	}

	c := &fileConfig{}
	parseFileContent(c, tc.content)

	if *c.BlockerAddress != tc.expected.BlockerAddress {
		t.Errorf("BlockerAddress should be %v; got: %v", tc.expected.BlockerAddress, *c.BlockerAddress)
	}
	if *c.AutostartEnabled != tc.expected.AutostartEnabled {
		t.Errorf("AutostartEnabled should be %v; got: %v", tc.expected.AutostartEnabled, *c.AutostartEnabled)
	}
	if *c.GUIEnabled != tc.expected.GUIEnabled {
		t.Errorf("GUIEnabled should be %v; got: %v", tc.expected.GUIEnabled, *c.GUIEnabled)
	}
	if *c.LogEnabled != tc.expected.LogEnabled {
		t.Errorf("LogEnabled should be %v; got: %v", tc.expected.LogEnabled, *c.LogEnabled)
	}

	if c.BlocklistPath != nil {
		t.Errorf("BlocklistPath should be nil; got: %v", *c.BlocklistPath)
	}
	if c.BlacklistPath != nil {
		t.Errorf("BlacklistPath should be nil; got: %v", *c.BlacklistPath)
	}
	if c.WhitelistPath != nil {
		t.Errorf("WhitelistPath should be nil; got: %v", *c.WhitelistPath)
	}
	if c.LogPath != nil {
		t.Errorf("LogPath should be nil; got: %v", *c.LogPath)
	}
	if c.ProxyAddress != nil {
		t.Errorf("ProxyAddress should be nil; got: %v", *c.ProxyAddress)
	}
}

func TestDefaultConfig(t *testing.T) {
	t5555 := ":5555"

	tc := struct {
		initial, expected *fileConfig
	}{
		initial: &fileConfig{
			BlockerAddress: &t5555,
		},
		expected: &fileConfig{
			BlockerAddress:   &t5555,
			BlocklistPath:    &defaultBlocklistPath,
			BlacklistPath:    &defaultBlacklistPath,
			WhitelistPath:    &defaultWhitelistPath,
			AutostartEnabled: &defaultAutostartEnabled,
			GUIEnabled:       &defaultGUIEnabled,
			LogEnabled:       &defaultLogEnabled,
			LogPath:          &defaultLogPath,
			ProxyAddress:     &defaultProxyAddress,
		},
	}

	defaultConfig(tc.initial)

	if *tc.initial.BlockerAddress != *tc.expected.BlockerAddress {
		t.Errorf("BlockerAddress should be %v; got: %v", *tc.expected.BlockerAddress, *tc.initial.BlockerAddress)
	}
	if *tc.initial.BlocklistPath != *tc.expected.BlocklistPath {
		t.Errorf("BlocklistPath should be %v; got: %v", *tc.expected.BlocklistPath, *tc.initial.BlocklistPath)
	}
	if *tc.initial.BlacklistPath != *tc.expected.BlacklistPath {
		t.Errorf("BlacklistPath should be %v; got: %v", *tc.expected.BlacklistPath, *tc.initial.BlacklistPath)
	}
	if *tc.initial.WhitelistPath != *tc.expected.WhitelistPath {
		t.Errorf("WhitelistPath should be %v; got: %v", *tc.expected.WhitelistPath, *tc.initial.WhitelistPath)
	}
	if *tc.initial.AutostartEnabled != *tc.expected.AutostartEnabled {
		t.Errorf("AutostartEnabled should be %v; got: %v", *tc.expected.AutostartEnabled, *tc.initial.AutostartEnabled)
	}
	if *tc.initial.GUIEnabled != *tc.expected.GUIEnabled {
		t.Errorf("GUIEnabled should be %v; got: %v", *tc.expected.GUIEnabled, *tc.initial.GUIEnabled)
	}
	if *tc.initial.LogEnabled != *tc.expected.LogEnabled {
		t.Errorf("LogEnabled should be %v; got: %v", *tc.expected.LogEnabled, *tc.initial.LogEnabled)
	}
	if *tc.initial.LogPath != *tc.expected.LogPath {
		t.Errorf("LogPath should be %v; got: %v", *tc.expected.LogPath, *tc.initial.LogPath)
	}
	if *tc.initial.ProxyAddress != *tc.expected.ProxyAddress {
		t.Errorf("ProxyAddress should be %v; got: %v", *tc.expected.ProxyAddress, *tc.initial.ProxyAddress)
	}
}

func TestParseFlags(t *testing.T) {

	tc := struct {
		args     []string
		expected *Config
	}{
		args: []string{"lycurgus", "--address=:8888", "--blocklist=blocklist", "--blacklist=blacklist", "--whitelist=whitelist", "--autostart=true", "--gui=true", "--log=true", "--logfile=logfile", "--proxy=:9999"},
		expected: &Config{
			BlockerAddress:   ":8888",
			BlocklistPath:    "blocklist",
			BlacklistPath:    "blacklist",
			WhitelistPath:    "whitelist",
			AutostartEnabled: true,
			GUIEnabled:       true,
			LogEnabled:       true,
			LogPath:          "logfile",
			ProxyAddress:     ":9999",
		},
	}

	c := &Config{}
	parseFlags(c, tc.args)

	if c.BlockerAddress != tc.expected.BlockerAddress {
		t.Errorf("BlockerAddress should be %v; got: %v", tc.expected.BlockerAddress, c.BlockerAddress)
	}
	if c.BlocklistPath != tc.expected.BlocklistPath {
		t.Errorf("BlocklistPath should be %v; got: %v", tc.expected.BlocklistPath, c.BlocklistPath)
	}
	if c.BlacklistPath != tc.expected.BlacklistPath {
		t.Errorf("BlacklistPath should be %v; got: %v", tc.expected.BlacklistPath, c.BlacklistPath)
	}
	if c.WhitelistPath != tc.expected.WhitelistPath {
		t.Errorf("WhitelistPath should be %v; got: %v", tc.expected.WhitelistPath, c.WhitelistPath)
	}
	if c.AutostartEnabled != tc.expected.AutostartEnabled {
		t.Errorf("AutostartEnabled should be %v; got: %v", tc.expected.AutostartEnabled, c.AutostartEnabled)
	}
	if c.GUIEnabled != tc.expected.GUIEnabled {
		t.Errorf("GUIEnabled should be %v; got: %v", tc.expected.GUIEnabled, c.GUIEnabled)
	}
	if c.LogEnabled != tc.expected.LogEnabled {
		t.Errorf("LogEnabled should be %v; got: %v", tc.expected.LogEnabled, c.LogEnabled)
	}
	if c.LogPath != tc.expected.LogPath {
		t.Errorf("LogPath should be %v; got: %v", tc.expected.LogPath, c.LogPath)
	}
	if c.ProxyAddress != tc.expected.ProxyAddress {
		t.Errorf("ProxyAddress should be %v; got: %v", tc.expected.ProxyAddress, c.ProxyAddress)
	}
}
