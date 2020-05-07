package main

import (
	"flag"
	"fmt"
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
	logEnabled := flag.Bool("log", true, "logging enabled")
	logFile := flag.String("logfile", logFile(), "path to log file")
	proxyAddress := flag.String("proxy", "", "upstream proxy address")
	flag.Parse()

	if *logEnabled && *logFile != "" {
		fmt.Println("logenabled: ", *logFile)

		file, err := createLogFile(*logFile)
		if err != nil {
			log.Printf("Error opening file: %v", err)
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

	log.Println("Lycurgus started")

	var stopCh = make(chan os.Signal, 2)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-app.QuitCh:
		return
	case <-stopCh:
		return
	}
}
