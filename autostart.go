package main

import (
	"os"

	"github.com/kszab0/go-autostart"
)

// Autostart is an application that will be started when the user logs in
type Autostart struct {
	*autostart.App
}

// NewAutostart creates and initializes an Autostart application
func NewAutostart() (*Autostart, error) {
	var appExec []string
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	appExec = []string{exePath}
	autostart := &Autostart{
		App: &autostart.App{
			Name:        appName,
			DisplayName: appTitle,
			Exec:        appExec,
		},
	}
	return autostart, nil
}

func (a *Autostart) setEnabled(enabled bool) error {
	if enabled {
		return a.Enable()
	}
	return a.Disable()
}
