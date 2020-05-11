package main

import (
	"log"
	"path/filepath"

	"github.com/ProtonMail/go-appdir"
	"github.com/natefinch/lumberjack"
)

func logDir() string {
	dirs := appdir.New(appName)
	return dirs.UserLogs()
}

func logFile() string {
	return filepath.Join(logDir(), "lycurgus.log")
}

func initLog(config *Config) {
	if config.LogEnabled && config.LogPath != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   config.LogPath,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		})
	}
}
