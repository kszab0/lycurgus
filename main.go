package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := parseConfig(os.Args)

	initLog(config)

	log.Printf("Starting Lycurgus version: %s built @ %s\n", Version, BuildDate)

	app, err := NewApp(*config)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Fatal(app.RunBlocker())
	}()

	if config.GUIEnabled {
		app.RunGUI()
	}

	var stopCh = make(chan os.Signal, 2)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-app.QuitCh:
	case <-stopCh:
		return
	}
}
