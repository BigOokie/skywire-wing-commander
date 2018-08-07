package utils

import (
	"testing"
)

func TestFileExistsNotOk(t *testing.T) {
	var badFilename string = "this-file-does-not-exist.test"
	if FileExists(badFilename) {
		t.Errorf("File (%s) does not exist but was detected as existing.", badFilename)
	}
}

func TestFileExistsOk(t *testing.T) {
	var goodFilename string = "./utils_test.go"
	if !FileExists(goodFilename) {
		t.Errorf("File (%s) does not exist but should.", goodFilename)
	}
}

func TestUserHome(t *testing.T) {
	if UserHome() == "" {
		t.Error("Failed to determine User Home Pather.")
	}
}
