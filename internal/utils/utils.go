// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"os"
	"runtime"

	"github.com/marcsauter/single"
	log "github.com/sirupsen/logrus"
	latest "github.com/tcnksm/go-latest"
)

// UserHome returns the current user home path
func UserHome() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}

// FileExists checks if the provided file exists or not
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Errorf("File (%s) does not exist. Error: %s", filename, err)
		return false
	}

	log.Debugf("File (%s) exists.", filename)
	return true
}

// UpdateAvailable will perform a check against the specified
// repository to determine if the passed in version tag is the latest or not
func UpdateAvailable(ownername, reponame, versiontag string) (result bool, updateMsg string) {
	log.Debugf("UpdateAvailable: Owner: %s, Repo: %s, Version: %s",
		ownername, reponame, versiontag)
	result = false
	updateMsg = "An error occurred checking for updates."
	githubTag := &latest.GithubTag{
		Owner:      ownername,
		Repository: reponame,
	}

	res, _ := latest.Check(githubTag, versiontag)
	if res.Outdated {
		result = true
		updateMsg = fmt.Sprintf("%s is not latest, you should upgrade to v%s", versiontag, res.Current)
	} else if res.New {
		result = false
		updateMsg = fmt.Sprintf("%s is newer than the latest version on GitHub (v%s).", versiontag, res.Current)
	} else if res.Latest {
		result = false
		updateMsg = fmt.Sprintf("v%s is the latest version.", res.Current)
	}
	log.Infof("UpdateAvailable: %s", updateMsg)
	return
}

// InitAppInstance will attempt to initalise an instance of the application based on the provided value of appID.
// A FATAL error will occur causing the application to exit if another instance
// of the application is detected as already running.
func InitAppInstance(appID string) (s *single.Single) {
	s = single.New(appID)
	if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
		msgAppInstErr := "Another instance of Wing Commander has been detected running on this system.\n\n" +
			"To identify and terminate (kill) ALL instances of Wing Commander on this system, run:\n\n" +
			"   pgrep wcbot | xargs kill\n\n" +
			"Exiting\n"
		log.Fatal(msgAppInstErr)
	} else if err != nil {
		// Another error occurred, might be worth handling it as well
		log.Fatalf("Failed to acquire exclusive app lock: %v", err)
	}
	return
}
