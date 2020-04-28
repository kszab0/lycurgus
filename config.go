package main

import (
	"os"
	"path/filepath"
	"runtime"
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
