package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ProtonMail/go-appdir"
	"gopkg.in/yaml.v2"
)

var (
	defaultBlockerAddress   = ":5678"
	defaultBlockerEnabled   = true
	defaultAutostartEnabled = true
	defaultGUIEnabled       = true
	defaultLogEnabled       = true
	defaultProxyAddress     = ""
	defaultBlocklistPath    = filepath.Join(configDir(), "blocklist")
	defaultBlacklistPath    = filepath.Join(configDir(), "blacklist")
	defaultWhitelistPath    = filepath.Join(configDir(), "whitelist")
	defaultLogPath          = logFile()
	defaultUpdateInterval   = 24 * time.Hour
)

// Config holds the settings for the application
type Config struct {
	BlockerAddress   string
	BlocklistPath    string
	BlacklistPath    string
	WhitelistPath    string
	AutostartEnabled bool
	GUIEnabled       bool
	LogEnabled       bool
	LogPath          string
	ProxyAddress     string
	UpdateInterval   time.Duration
}

type fileConfig struct {
	BlockerAddress   *string        `yaml:"address,omitempty"`
	BlocklistPath    *string        `yaml:"blocklist,omitempty"`
	BlacklistPath    *string        `yaml:"blacklist,omitempty"`
	WhitelistPath    *string        `yaml:"whitelist,omitempty"`
	AutostartEnabled *bool          `yaml:"autostart,omitempty"`
	GUIEnabled       *bool          `yaml:"gui,omitempty"`
	LogEnabled       *bool          `yaml:"log,omitempty"`
	LogPath          *string        `yaml:"logfile,omitempty"`
	ProxyAddress     *string        `yaml:"proxy,omitempty"`
	UpdateInterval   *time.Duration `ỳaml:"updateInterval,omitempty"`
}

func (fc *fileConfig) toConfig() *Config {
	c := &Config{}
	if fc.BlockerAddress != nil {
		c.BlockerAddress = *fc.BlockerAddress
	}
	if fc.BlocklistPath != nil {
		c.BlocklistPath = *fc.BlocklistPath
	}
	if fc.BlacklistPath != nil {
		c.BlacklistPath = *fc.BlacklistPath
	}
	if fc.WhitelistPath != nil {
		c.WhitelistPath = *fc.WhitelistPath
	}
	if fc.AutostartEnabled != nil {
		c.AutostartEnabled = *fc.AutostartEnabled
	}
	if fc.GUIEnabled != nil {
		c.GUIEnabled = *fc.GUIEnabled
	}
	if fc.LogEnabled != nil {
		c.LogEnabled = *fc.LogEnabled
	}
	if fc.LogPath != nil {
		c.LogPath = *fc.LogPath
	}
	if fc.ProxyAddress != nil {
		c.ProxyAddress = *fc.ProxyAddress
	}
	if fc.UpdateInterval != nil {
		c.UpdateInterval = *fc.UpdateInterval
	}
	return c
}

func (c *Config) String() string {
	return fmt.Sprintf(`Config {
  BlockerAddress:   %v,
  BlocklistPath:    %v,
  BlacklistPath:    %v,
  WhitelistPath:    %v,
  AutostartEnabled: %v,
  GUIEnabled:       %v,
  LogEnabled:       %v,
  LogPath:          %v,
  ProxyAddress:     %v,
  UpdateInterval:   %v,
}`, c.BlockerAddress, c.BlocklistPath, c.BlacklistPath, c.WhitelistPath,
		c.AutostartEnabled, c.GUIEnabled, c.LogEnabled, c.LogPath, c.ProxyAddress, c.UpdateInterval)
}

func defaultConfig(config *fileConfig) {
	if config.BlockerAddress == nil {
		config.BlockerAddress = &defaultBlockerAddress
	}
	if config.BlocklistPath == nil {
		config.BlocklistPath = &defaultBlocklistPath
	}
	if config.BlacklistPath == nil {
		config.BlacklistPath = &defaultBlacklistPath
	}
	if config.WhitelistPath == nil {
		config.WhitelistPath = &defaultWhitelistPath
	}
	if config.AutostartEnabled == nil {
		config.AutostartEnabled = &defaultAutostartEnabled
	}
	if config.GUIEnabled == nil {
		config.GUIEnabled = &defaultGUIEnabled
	}
	if config.LogEnabled == nil {
		config.LogEnabled = &defaultLogEnabled
	}
	if config.LogPath == nil {
		config.LogPath = &defaultLogPath
	}
	if config.ProxyAddress == nil {
		config.ProxyAddress = &defaultProxyAddress
	}
	if config.UpdateInterval == nil {
		config.UpdateInterval = &defaultUpdateInterval
	}
}

// parseFile parses a yaml config.
func parseFile(c *fileConfig, filePath string) {
	content, _ := ioutil.ReadFile(filePath)
	parseFileContent(c, content)
}

func parseFileContent(c *fileConfig, content []byte) {
	yaml.Unmarshal(content, c)
}

func isFlagPassed(flags *flag.FlagSet, name string) bool {
	found := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// parseFlags parses the flags and sets the values for keys in the config
// that are passed in the flags.
func parseFlags(config *Config, args []string) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	blockerAddress := flags.String("address", "", "address to run blocker")
	blocklistPath := flags.String("blocklist", "", "path to blocklist file")
	blacklistPath := flags.String("blacklist", "", "path to blacklist file")
	whitelistPath := flags.String("whitelist", "", "path to whitelist file")
	autostartEnabled := flags.Bool("autostart", true, "autostart enabled")
	guiEnabled := flags.Bool("gui", true, "start app with GUI")
	logEnabled := flags.Bool("log", true, "logging enabled")
	logFile := flags.String("logfile", "", "path to log file")
	proxyAddress := flags.String("proxy", "", "upstream proxy address")
	updateInterval := flags.Duration("update", 0, "update interval")

	flags.Parse(args[1:])

	if isFlagPassed(flags, "address") {
		config.BlockerAddress = *blockerAddress
	}
	if isFlagPassed(flags, "blocklist") {
		config.BlocklistPath = *blocklistPath
	}
	if isFlagPassed(flags, "blocklist") {
		config.BlacklistPath = *blacklistPath
	}
	if isFlagPassed(flags, "whitelist") {
		config.WhitelistPath = *whitelistPath
	}
	if isFlagPassed(flags, "autostart") {
		config.AutostartEnabled = *autostartEnabled
	}
	if isFlagPassed(flags, "gui") {
		config.GUIEnabled = *guiEnabled
	}
	if isFlagPassed(flags, "log") {
		config.LogEnabled = *logEnabled
	}
	if isFlagPassed(flags, "logfile") {
		config.LogPath = *logFile
	}
	if isFlagPassed(flags, "proxy") {
		config.ProxyAddress = *proxyAddress
	}
	if isFlagPassed(flags, "update") {
		config.UpdateInterval = *updateInterval
	}
}

func parseConfig(args []string) *Config {
	return parseConfigFile(args, configFile())
}

func parseConfigFile(args []string, configPath string) *Config {
	fc := &fileConfig{}
	parseFile(fc, configPath)
	defaultConfig(fc)

	c := fc.toConfig()
	parseFlags(c, args)
	return c
}

func isDirExists(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func createDir(dir string) error {
	if isDirExists(dir) {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

func configDir() string {
	dirs := appdir.New(appName)
	return dirs.UserConfig()
}

func createConfigDir() error {
	return createDir(configDir())
}

func configFile() string {
	return filepath.Join(configDir(), "lycurgus.yml")
}

func createConfigFile(logFile string) (*os.File, error) {
	if err := createConfigDir(); err != nil {
		return nil, err
	}
	return os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

func cacheDir() string {
	dir, _ := os.UserCacheDir()
	return filepath.Join(dir, appName)
}

func blocklistCacheDir() string {
	return filepath.Join(cacheDir(), "blocklist")
}
