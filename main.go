package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := parseConfig(os.Args)

	if config.LogEnabled && config.LogPath != "" {
		file, err := createLogFile(config.LogPath)
		if err != nil {
			log.Printf("Error opening file: %v", err)
		}
		defer file.Close()
                log.SetOutput(file)
	}

	app, err := NewApp(*config)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Fatal(app.RunBlocker())
	}()

	if config.GUIEnabled {
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
