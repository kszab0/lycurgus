package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	appName    = "lycurgus"
	appTitle   = "Lycurgus"
	appTooltip = "Lycurgus Ad Blocker"
)

const (
	defaultBlockerAddress   = ":8080"
	defaultBlockerEnabled   = true
	defaultBlocklistPath    = "./blocklist"
	defaultBlacklistPath    = "./blacklist"
	defaultWhitelistPath    = "./whitelist"
	defaultAutostartEnabled = true
)

const defaultBlocklists = `https://adaway.org/hosts.txt
https://v.firebog.net/hosts/AdguardDNS.txt
https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt
https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt
https://hosts-file.net/ad_servers.txt
https://v.firebog.net/hosts/Easylist.txt
https://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&showintro=0&mimetype=plaintext
https://raw.githubusercontent.com/FadeMind/hosts.extras/master/UncheckyAds/hosts
https://raw.githubusercontent.com/bigdargon/hostsVN/master/hosts
https://raw.githubusercontent.com/jdlingyu/ad-wars/master/hosts`

// App holds all the states for the application
// and manages all the components.
type App struct {
	blockerAddress   string
	blockerEnabled   bool
	blocklistPath    string
	blacklistPath    string
	whitelistPath    string
	autostartEnabled bool

	blocker   *Blocker
	getter    Getter
	gui       *GUI
	autostart *Autostart
}

// AppOption is a functional option for configuring App.
type AppOption func(*App)

// WithBlockerAddress sets the blocker's http address.
func WithBlockerAddress(addr string) AppOption {
	return func(app *App) {
		app.blockerAddress = addr
	}
}

// WithBlockerEnabled sets the blocker's enabled state.
func WithBlockerEnabled(enabled bool) AppOption {
	return func(app *App) {
		app.blockerEnabled = enabled
	}
}

// WithBlocklistPath sets the file path for the blocklists.
func WithBlocklistPath(path string) AppOption {
	return func(app *App) {
		app.blocklistPath = path
	}
}

// WithBlacklistPath sets the file path for the blacklist.
func WithBlacklistPath(path string) AppOption {
	return func(app *App) {
		app.blacklistPath = path
	}
}

// WithWhitelistPath sets the file path for the whitelist.
func WithWhitelistPath(path string) AppOption {
	return func(app *App) {
		app.whitelistPath = path
	}
}

// WithAutostartEnabled sets the file path for the blacklist.
func WithAutostartEnabled(enabled bool) AppOption {
	return func(app *App) {
		app.autostartEnabled = enabled
	}
}

// NewApp creates and initializes an App.
func NewApp(opts ...AppOption) (*App, error) {
	app := &App{
		blockerAddress:   defaultBlockerAddress,
		blockerEnabled:   defaultBlockerEnabled,
		blocklistPath:    defaultBlocklistPath,
		blacklistPath:    defaultBlacklistPath,
		whitelistPath:    defaultWhitelistPath,
		autostartEnabled: defaultAutostartEnabled,
		getter:           http.DefaultClient,
	}
	for _, opt := range opts {
		opt(app)
	}

	blocker := NewBlocker(app.blockerEnabled)
	app.blocker = blocker
	if err := app.LoadBlocklist(); err != nil {
		return nil, err
	}
	if err := app.LoadBlacklist(); err != nil {
		return nil, err
	}
	if err := app.LoadWhitelist(); err != nil {
		return nil, err
	}

	gui, err := NewGUI(WithGUIEnabled(app.blockerEnabled))
	if err != nil {
		return nil, err
	}
	app.gui = gui

	autostart, err := NewAutostart()
	if err != nil {
		log.Println("Error setting autostart: ", err)
	}
	app.autostart = autostart
	if err := app.autostart.setEnabled(app.autostartEnabled); err != nil {
		log.Println("Error setting autostart: ", err)
	}

	return app, nil
}

// LoadBlocklist reads the blocklist file
// and initializes the blocker's blocklist matcher.
func (app *App) LoadBlocklist() error {
	file, err := getBlocklists(app.blocklistPath)
	if err != nil {
		return err
	}
	defer file.Close()

	app.blocker.blocklist = app.loadBlocklist(file)
	return nil
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

// loadBloclists parses a blocklists file
// and initializes a hashMatcher as a blocklist.
func (app *App) loadBlocklist(r io.Reader) Matcher {
	rules := []string{}
	lines, err := readLines(r)
	if err != nil {
		return nil
	}
	for _, line := range lines {
		hosts, err := parseHostsURL(app.getter, line)
		if err != nil {
			log.Println("Error reading url: ", line)
			continue
		}
		rules = append(rules, hosts...)
	}
	matcher := &hashMatcher{}
	matcher.Load(rules)
	return matcher
}

// LoadBlacklist reads the blacklist file
// and initializes the blocker's blacklist matcher.
func (app *App) LoadBlacklist() error {
	blacklist, err := loadMatcherFromFile(app.blacklistPath)
	if err != nil {
		// ignore if file not exists
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	app.blocker.blacklist = blacklist
	return nil
}

// LoadWhitelist reads the whitelist file
// and initializes the blocker's whitelist matcher.
func (app *App) LoadWhitelist() error {
	whitelist, err := loadMatcherFromFile(app.whitelistPath)
	if err != nil {
		// ignore if file not exists
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	app.blocker.whitelist = whitelist
	return nil
}

// loadMatcherFromFile loads and initializes a regexpMatcher from a file.
func loadMatcherFromFile(path string) (Matcher, error) {
	hosts, err := parseHostsFile(path)
	if err != nil {
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

// RunBlocker serves the Blocker.
func (app *App) RunBlocker() error {
	return http.ListenAndServe(app.blockerAddress, app.blocker)
}

// RunGUI starts the GUI.
func (app *App) RunGUI() {
	go func() {
		for {
			select {
			case enabled := <-app.gui.EnabledCh:
				app.blocker.enabled = enabled
			case enabled := <-app.gui.AutostartCh:
				if err := app.autostart.setEnabled(enabled); err != nil {
					log.Println("Error setting autostart: ", err)
				}
			case <-app.gui.QuitCh:
				return
			}
		}
	}()
	app.gui.Run()
}
