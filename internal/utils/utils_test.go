// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.
package utils

import (
	"testing"

	"github.com/BigOokie/skywire-wing-commander/internal/wcconst"
)

func Test_FileExists_NotOk(t *testing.T) {
	badFilename := "this-file-does-not-exist.test"
	if FileExists(badFilename) {
		t.Errorf("File (%s) does not exist but was detected as existing.", badFilename)
	}
}

func Test_FileExists_Ok(t *testing.T) {
	goodFilename := "./utils_test.go"
	if !FileExists(goodFilename) {
		t.Errorf("File (%s) does not exist but should.", goodFilename)
	}
}

func Test_UserHome(t *testing.T) {
	if UserHome() == "" {
		t.Error("Failed to determine User Home Path")
	}
}

func Test_InitAppInstance_Ok(t *testing.T) {
	appInst := InitAppInstance(wcconst.AppInstanceID + "Test_InitAppInstance_Ok")
	defer appInst.Unlock()
	if appInst == nil {
		t.Error("Failed to obtain application instance.")
	}
}

func Test_ReleaseAppInstance_Ok(t *testing.T) {
	appInst := InitAppInstance(wcconst.AppInstanceID + "Test_ReleaseAppInstance_Ok")
	if appInst == nil {
		t.Error("Failed to obtain application instance.")
	}

	ReleaseAppInstance(appInst)
}

func Test_ReleaseAppInstance_Nil_Ok(t *testing.T) {
	ReleaseAppInstance(nil)
}

func Test_ReleaseAppInstance_Double_Unlock(t *testing.T) {
	appInst := InitAppInstance(wcconst.AppInstanceID + "Test_ReleaseAppInstance_Double_Unlock")
	if appInst == nil {
		t.Error("Failed to obtain application instance.")
	}
	err := appInst.TryUnlock()
	if err != nil {
		t.Error("TryUnlock failded")
	}

	ReleaseAppInstance(appInst)
}

func Test_UpgradeAvailable_BadRepo(t *testing.T) {
	result, msg := UpdateAvailable("BigOokie", "ThisRepoDoesntExist", "1.0")
	if result {
		t.Errorf("An upgrade should not be available as the Repo does not exist. %s", msg)
	}
}
func Test_UpgradeAvailable_GoodRepo_BadVersion(t *testing.T) {
	result, msg := UpdateAvailable("BigOokie", "skywire-wing-commander", "vBAD-99999")
	if result {
		t.Errorf("An upgrade should not be available. The version tag should not exist. %s", msg)
	}
}

func Test_UpgradeAvailable_GoodRepo_VersionOK_Older(t *testing.T) {
	result, msg := UpdateAvailable("BigOokie", "skywire-wing-commander", "0.0.1")
	if !result {
		t.Errorf("An upgrade should be available. The newer version tag should exist. %s", msg)
	}
}

func Test_UpgradeAvailable_GoodRepo_VersionOK_Newer(t *testing.T) {
	result, msg := UpdateAvailable("BigOokie", "skywire-wing-commander", "1000.0.0")
	if result {
		t.Errorf("%s", msg)
	}
}

func Test_UpgradeAvailable_GoodRepo_VersionOK_Latest(t *testing.T) {
	result, msg := UpdateAvailable("BigOokie", "skywire-wing-commander", wcconst.BotVersion)
	if result {
		t.Errorf("%s", msg)
	}
}
