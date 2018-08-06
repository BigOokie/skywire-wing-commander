package main

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
