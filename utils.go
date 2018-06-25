package main

import (
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// userHome returns the current user home path
func userHome() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}

// fileExists checks if the provided file exists or not
func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Warnf("File does not exist: %s", filename)
		return false
	} else {
		log.Debugf("File exists: %s", filename)
		return true
	}
}
