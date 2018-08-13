package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BigOokie/skywire-wing-commander/src/wcconst"
)

//var scriptPath = "/src/github.com/BigOokie/skywire-wing-commander/src/scripts/"

// CheckUpdate executes the shell script `wc-updatecheck.sh` to determine if a new version
// of the app is available or not.
func CheckUpdate() (result []byte, err error) {
	var cmd *exec.Cmd
	var gopath = os.Getenv("GOPATH")

	// output, err := exec.Command("/bin/bash", "/home/myname/MyProj/Server/src/temp.sh").CombinedOutput()
	cmd = exec.Command("/bin/bash", filepath.Join(gopath, fmt.Sprintf("%s%s", wcconst.ScriptPath, "wc-checkupdate.sh")))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	result = out
	return
}

// UpdateApp executes the shell script `wc-update.sh` to perform an update of the app to the new version
func UpdateApp() (result []byte, err error) {
	var cmd *exec.Cmd
	var gopath = os.Getenv("GOPATH")
	cmd = exec.Command(filepath.Join(gopath, fmt.Sprintf("%s%s", wcconst.ScriptPath, "wc-update.sh")))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	result = out
	return
}
