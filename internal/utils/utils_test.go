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
