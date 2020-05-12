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

	UpdateCh chan struct{}

	QuitCh chan struct{}
}

type menu struct {
	enabled         *systray.MenuItem
	enabledAction   *systray.MenuItem
	autostart       *systray.MenuItem
	autostartAction *systray.MenuItem
	update          *systray.MenuItem
	quit            *systray.MenuItem
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
		UpdateCh:    make(chan struct{}),
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
	//systray.SetTitle(gui.title)
	systray.SetTooltip(gui.tooltip)
	systray.SetIcon(gui.icon)

	gui.menu.enabled = systray.AddMenuItem("Enabled", "")
	gui.menu.enabled.Disable()
	gui.menu.enabledAction = systray.AddMenuItem("Disable", "")
	gui.setEnabled()
	systray.AddSeparator()

	gui.menu.autostart = systray.AddMenuItem("Autostart enabled", "")
	gui.menu.autostart.Disable()
	gui.menu.autostartAction = systray.AddMenuItem("Disable autostart", "")
	gui.setAutostart()
	systray.AddSeparator()

	gui.menu.update = systray.AddMenuItem("Update lists", "")
	systray.AddSeparator()

	gui.menu.quit = systray.AddMenuItem("Quit", "")
}

func (gui *GUI) setAutostart() {
	if gui.autostart {
		gui.menu.autostart.SetTitle("Autostart enabled")
		gui.menu.autostartAction.SetTitle("Disable autostart")
	} else {
		gui.menu.autostart.SetTitle("Autostart disabled")
		gui.menu.autostartAction.SetTitle("Enable autostart")
	}
}

func (gui *GUI) setEnabled() {
	if gui.enabled {
		gui.menu.enabled.SetTitle("Lycurgus is Enabled")
		gui.menu.enabledAction.SetTitle("Disable Lycurgus")
	} else {
		gui.menu.enabled.SetTitle("Lycurgus is Disabled")
		gui.menu.enabledAction.SetTitle("Enable Lycurgus")
	}
}

func (gui *GUI) listen() {
	for {
		select {
		case <-gui.menu.enabledAction.ClickedCh:
			gui.enabled = !gui.enabled
			gui.EnabledCh <- gui.enabled
			gui.setEnabled()
		case <-gui.menu.autostartAction.ClickedCh:
			gui.autostart = !gui.autostart
			gui.AutostartCh <- gui.autostart
			gui.setAutostart()
		case <-gui.menu.update.ClickedCh:
			gui.UpdateCh <- struct{}{}
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
