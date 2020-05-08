package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
}

type fileConfig struct {
	BlockerAddress   *string `yaml:"address,omitempty"`
	BlocklistPath    *string `yaml:"blocklist,omitempty"`
	BlacklistPath    *string `yaml:"blacklist,omitempty"`
	WhitelistPath    *string `yaml:"whitelist,omitempty"`
	AutostartEnabled *bool   `yaml:"autostart,omitempty"`
	GUIEnabled       *bool   `yaml:"gui,omitempty"`
	LogEnabled       *bool   `yaml:"log,omitempty"`
	LogPath          *string `yaml:"logfile,omitempty"`
	ProxyAddress     *string `yaml:"proxy,omitempty"`
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
}`, c.BlockerAddress, c.BlocklistPath, c.BlacklistPath, c.WhitelistPath,
		c.AutostartEnabled, c.GUIEnabled, c.LogEnabled, c.LogPath, c.ProxyAddress)
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
}

// parseFile parses a yaml config.
func parseFile(c *fileConfig, filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading config: ", err)
		return
	}
	parseFileContent(c, content)
}

func parseFileContent(c *fileConfig, content []byte) {
	err := yaml.Unmarshal(content, c)
	if err != nil {
		log.Println("Error unmarshaling config: ", err)
	}
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

func logDir() string {
	dirs := appdir.New(appName)
	return dirs.UserLogs()
}

func createLogDir() error {
	return createDir(logDir())
}

func logFile() string {
	return filepath.Join(logDir(), "lycurgus.log")
}

func createLogFile(logFile string) (*os.File, error) {
	if err := createLogDir(); err != nil {
		return nil, err
	}
	return os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}
