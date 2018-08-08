package utils

import (
	"testing"
)

func Test_FileExists_NotOk(t *testing.T) {
	var badFilename string = "this-file-does-not-exist.test"
	if FileExists(badFilename) {
		t.Errorf("File (%s) does not exist but was detected as existing.", badFilename)
	}
}

func Test_FileExists_Ok(t *testing.T) {
	var goodFilename string = "./utils_test.go"
	if !FileExists(goodFilename) {
		t.Errorf("File (%s) does not exist but should.", goodFilename)
	}
}

func Test_UserHome(t *testing.T) {
	if UserHome() == "" {
		t.Error("Failed to determine User Home Pather.")
	}
}
