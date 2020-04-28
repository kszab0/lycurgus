package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	blockerAddress := flag.String("address", defaultBlockerAddress, "address to run blocker")
	blocklistPath := flag.String("blocklist", defaultBlocklistPath, "path to blocklist file")
	blacklistPath := flag.String("blacklist", defaultBlacklistPath, "path to blacklist file")
	whitelistPath := flag.String("whitelist", defaultWhitelistPath, "path to whitelist file")
	autostartEnabled := flag.Bool("autostart", defaultAutostartEnabled, "autostart enabled")
	guiEnabled := flag.Bool("gui", true, "start app with GUI")
	logFile := flag.String("log", "", "path to log file")
	proxyAddress := flag.String("proxy", "", "upstream proxy address")
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
		WithBlockerAddress(*blockerAddress),
		WithBlocklistPath(*blocklistPath),
		WithBlacklistPath(*blacklistPath),
		WithWhitelistPath(*whitelistPath),
		WithAutostartEnabled(*autostartEnabled),
		WithProxyAddress(*proxyAddress),
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Fatal(app.RunBlocker())
	}()

	if *guiEnabled {
		go app.RunGUI()
	}

	var stopCh = make(chan os.Signal, 2)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-app.QuitCh:
		return
	case <-stopCh:
		return
	}
}
