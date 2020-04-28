package main

import (
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	proxyAddress := flag.String("proxy", defaultBlockerAddress, "proxy address")
	blocklistPath := flag.String("blocklist", defaultBlocklistPath, "path to blocklist file")
	blacklistPath := flag.String("blacklist", defaultBlacklistPath, "path to blacklist file")
	whitelistPath := flag.String("whitelist", defaultWhitelistPath, "path to whitelist file")
	autostartEnabled := flag.Bool("autostart", defaultAutostartEnabled, "autostart enabled")
	logFile := flag.String("log", "", "path to log file")
	flag.Parse()

	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer file.Close()
		wrt := io.MultiWriter(os.Stdout, file)
		log.SetOutput(wrt)
	}

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
