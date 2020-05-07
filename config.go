package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/ProtonMail/go-appdir"
)

func configDir() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", appName)
	case "linux":
		baseDir := os.Getenv("XDG_CONFIG_HOME")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("HOME"), ".config")
		}
		return filepath.Join(baseDir, appName)
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), appName)
	}
	return ""
}

func isDirExists(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func createDir(dir string) error {
	if isDirExists(dir) {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

func logDir() string {
	dirs := appdir.New(appName)
	return dirs.UserLogs()
}

func createLogDir() error {
	return createDir(logDir())
}

func logFile() string {
	return filepath.Join(logDir(), "lycurgus.log")
}

func createLogFile(logFile string) (*os.File, error) {
	if err := createLogDir(); err != nil {
		return nil, err
	}
	return os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}
