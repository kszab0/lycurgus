package main

import (
	"flag"
	"log"
)

func main() {
	proxyAddress := flag.String("proxy", defaultBlockerAddress, "proxy address")
	blocklistPath := flag.String("blocklist", defaultBlocklistPath, "path to blocklist file")
	blacklistPath := flag.String("blacklist", defaultBlacklistPath, "path to blacklist file")
	whitelistPath := flag.String("whitelist", defaultWhitelistPath, "path to whitelist file")
	autostartEnabled := flag.Bool("autostart", defaultAutostartEnabled, "autostart enabled")
	flag.Parse()

	app, err := NewApp(
		WithBlockerAddress(*proxyAddress),
		WithBlocklistPath(*blocklistPath),
		WithBlacklistPath(*blacklistPath),
		WithWhitelistPath(*whitelistPath),
		WithAutostartEnabled(*autostartEnabled),
	)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Fatal(app.RunBlocker())
	}()
	app.RunGUI()
}
