package main

import (
	"github.com/getlantern/systray"
)

// GUI is a tray ui for the application
type GUI struct {
	enabled   bool
	autostart bool

	title   string
	tooltip string
	icon    []byte
	menu    *menu

	EnabledCh   chan bool
	AutostartCh chan bool
	QuitCh      chan struct{}
}

type menu struct {
	enabled   *systray.MenuItem
	autostart *systray.MenuItem
	quit      *systray.MenuItem
}

// GUIOption is a functional option for configuring the GUI
type GUIOption func(*GUI)

// WithGUIEnabled sets the GUI's enabled state
func WithGUIEnabled(enabled bool) GUIOption {
	return func(gui *GUI) {
		gui.enabled = enabled
	}
}

// WithGUIAutostart sets the GUI's autostart state
func WithGUIAutostart(enabled bool) GUIOption {
	return func(gui *GUI) {
		gui.autostart = enabled
	}
}

// NewGUI creates and initializes the GUI
func NewGUI(opts ...GUIOption) (*GUI, error) {
	gui := &GUI{
		enabled:   defaultBlockerEnabled,
		autostart: defaultAutostartEnabled,

		title:   appTitle,
		tooltip: appTooltip,
		icon:    icon,
		menu:    &menu{},

		EnabledCh:   make(chan bool),
		AutostartCh: make(chan bool),
		QuitCh:      make(chan struct{}),
	}

	for _, opt := range opts {
		opt(gui)
	}

	return gui, nil
}

// Run starts the GUI
func (gui *GUI) Run() {
	start := func() {
		gui.init()
		gui.listen()
	}
	exit := func() {}
	systray.Run(start, exit)
}

func (gui *GUI) init() {
	systray.SetTitle(gui.title)
	systray.SetTooltip(gui.tooltip)
	systray.SetIcon(gui.icon)

	gui.menu.enabled = systray.AddMenuItem("Enable", "")
	if gui.enabled {
		gui.menu.enabled.Check()
	}
	gui.menu.autostart = systray.AddMenuItem("Autostart", "")
	if gui.autostart {
		gui.menu.autostart.Check()
	}
	systray.AddSeparator()
	gui.menu.quit = systray.AddMenuItem("Quit", "")
}

func (gui *GUI) listen() {
	for {
		select {
		case <-gui.menu.enabled.ClickedCh:
			gui.enabled = !gui.enabled
			gui.EnabledCh <- gui.enabled
			if gui.enabled {
				gui.menu.enabled.Check()
			} else {
				gui.menu.enabled.Uncheck()
			}
		case <-gui.menu.autostart.ClickedCh:
			gui.autostart = !gui.autostart
			gui.AutostartCh <- gui.autostart
			if gui.autostart {
				gui.menu.autostart.Check()
			} else {
				gui.menu.autostart.Uncheck()
			}
		case <-gui.menu.quit.ClickedCh:
			gui.QuitCh <- struct{}{}
			gui.Quit()
			return
		}
	}
}

// Quit terminates the GUI
func (gui *GUI) Quit() {
	systray.Quit()
}
