package main

import (
	"log"
	"net/http"
	"net/url"
)

const (
	appName    = "lycurgus"
	appTitle   = "Lycurgus"
	appTooltip = "Lycurgus Ad Blocker"
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
	autostartEnabled bool
	proxyAddress     string

	storage   *Storage
	blocker   *Blocker
	gui       *GUI
	autostart *Autostart

	QuitCh chan struct{}
}

// NewApp creates and initializes an App.
func NewApp(config Config) (*App, error) {
	app := &App{
		blockerAddress:   config.BlockerAddress,
		blockerEnabled:   defaultBlockerEnabled,
		autostartEnabled: config.AutostartEnabled,
		storage: &Storage{
			blocklistPath:  config.BlocklistPath,
			blacklistPath:  config.BlacklistPath,
			whitelistPath:  config.WhitelistPath,
			updateInterval: config.UpdateInterval,
		},
		QuitCh: make(chan struct{}, 1),
	}

	// set upstream proxy for default http client
	if app.proxyAddress != "" {
		proxyURL, err := url.Parse("http://" + app.proxyAddress)
		if err == nil {
			http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		}
	}

	app.blocker = NewBlocker(
		WithBlockerEnabled(app.blockerEnabled),
		WithBlockerProxyAddress(app.proxyAddress),
	)
	if err := app.LoadBlocklist(true); err != nil {
		return nil, err
	}
	if err := app.LoadBlacklist(); err != nil {
		return nil, err
	}
	if err := app.LoadWhitelist(); err != nil {
		return nil, err
	}

	gui, err := NewGUI(
		WithGUIEnabled(app.blockerEnabled),
		WithGUIAutostart(app.autostartEnabled),
	)
	if err != nil {
		return nil, err
	}
	app.gui = gui

	autostart, err := NewAutostart()
	if err != nil {
		log.Println("Error initializing autostart: ", err)
	}
	app.autostart = autostart
	if err := app.autostart.setEnabled(app.autostartEnabled); err != nil {
		log.Println("Error setting autostart: ", err)
	}

	return app, nil
}

// LoadBlocklist reads the blocklist file
// and initializes the blocker's blocklist matcher.
func (app *App) LoadBlocklist(allowCache bool) error {
	blocklist, err := app.storage.GetBlocklist(allowCache)
	if err != nil {
		return err
	}
	app.blocker.blocklist = blocklist
	return nil
}

// LoadBlacklist reads the blacklist file
// and initializes the blocker's blacklist matcher.
func (app *App) LoadBlacklist() error {
	blacklist, err := app.storage.GetBlacklist()
	if err != nil {
		return err
	}
	log.Println("Blacklist loaded")
	app.blocker.blacklist = blacklist
	return nil
}

// LoadWhitelist reads the whitelist file
// and initializes the blocker's whitelist matcher.
func (app *App) LoadWhitelist() error {
	whitelist, err := app.storage.GetWhitelist()
	if err != nil {
		return err
	}
	log.Println("Whitelist loaded")
	app.blocker.whitelist = whitelist
	return nil
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
				log.Println("Autostart set to: ", enabled)
			case <-app.gui.UpdateCh:
				if err := app.LoadBlocklist(false); err != nil {
					log.Println("Error reloading blocklist: ", err)
				}
				if err := app.LoadBlacklist(); err != nil {
					log.Println("Error reloading blacklist: ", err)
				}
				if err := app.LoadWhitelist(); err != nil {
					log.Println("Error reloading whitelist: ", err)
				}
			case <-app.gui.QuitCh:
				app.QuitCh <- struct{}{}
				return
			}
		}
	}()
	app.gui.Run()
}
