package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// clientConnection is a structure that represents the JSON file structure
// of the clients.json file (from Skywire project)
type clientConnection struct {
	Label   string `json:"label"`
	NodeKey string `json:"nodeKey"`
	AppKey  string `json:"appKey"`
	Count   int    `json:"count"`
}

// selectClientFile tests a number of known locations for the Client.json file
func selectClientFile() string {
	// Default to the Users home folder - but lets check
	clientfile := filepath.Join(userHome(), ".skywire", "manager", "clients.json")
	log.Debugf("[selectClientFile] Checking if file exists: %s", clientfile)

	// If we cant find the file in the users home folder - check the $GOPATH
	if fileExists(clientfile) == false {
		gopath := os.Getenv("GOPATH")
		clientfile = filepath.Join(gopath, "bin", ".skywire", "manager", "clients.json")
		log.Debugf("[selectClientFile] Checking if file exists: %s", clientfile)
		// Can't find what we are looking for
		if fileExists(clientfile) == false {
			log.Panicln("Unable to detect location of client file.")
		}
	}

	log.Debugf("[selectClientFile] Selected file: %s", clientfile)
	return clientfile
}

// Defines an in-memory slice (dynamic array) based on the ClientConnection struct
type clientConnectionSlice []clientConnection

// Determines if the specified ClientConnection exists within the clientConnectionSlice
func (c clientConnectionSlice) Exist(rf clientConnection) bool {
	for _, v := range c {
		if v.AppKey == rf.AppKey && v.NodeKey == rf.NodeKey {
			return true
		}
	}
	return false
}

// Reads the physical Skywire Clients.JSON file into an in-memory structure
func readClientConnectionConfig() (cfs map[string]clientConnectionSlice, err error) {
	fb, err := ioutil.ReadFile(oldconfig.ClientFile)
	if err != nil {
		if os.IsNotExist(err) {
			cfs = nil
			err = nil
			return
		}
		return
	}
	cfs = make(map[string]clientConnectionSlice)
	err = json.Unmarshal(fb, &cfs)
	if err != nil {
		return
	}
	return
}

// getClientConnectionListString will iterate over the ClientConnectionConfig JSON
// file and return a formatted string for all Clients and their Nodes
func getClientConnectionListString() string {
	var clientsb strings.Builder
	var condir string

	// Read the Client Connection Config (JSON) into ccc
	ccc, err := readClientConnectionConfig()
	if err == nil {
		// Iterate ccc reading the Keys (k)
		for k := range ccc {
			// Output to our string builder the current Client Type (from K)
			// Add an newline if this isnt the first itteration
			if clientsb.String() != "" {
				clientsb.WriteString("\n")
			}

			switch k {
			case "socket":
				condir = "Outbound"
			case "socksc":
				condir = "Inbound"
			default:
				condir = "??"
			}

			clientsb.WriteString(fmt.Sprintf("ClientType: [%s](%s)\n", k, condir))
			// Iterate all Nodes in the current client type (ccc[k])
			for _, b := range ccc[k] {
				// Output the details for each node of this client type
				clientsb.WriteString(fmt.Sprintf("Label:   %s\n", b.Label))
				clientsb.WriteString(fmt.Sprintf("NodeKey: %s\n", b.NodeKey))
				clientsb.WriteString(fmt.Sprintf("AppKey:  %s\n", b.AppKey))
				clientsb.WriteString("\n")
			}
		}
	}
	// Return the built string
	return clientsb.String()
}

// getClientConnectionCountString will iterate over the ClientConnectionConfig JSON
// file and return a formatted string containing the count of Clients connected
// Both Outbound and Inbound
func getClientConnectionCountString() string {
	var clientsb strings.Builder
	var condir string

	// Read the Client Connection Config (JSON) into ccc
	ccc, err := readClientConnectionConfig()
	if err == nil {
		// Iterate ccc reading the Keys (k)
		for k := range ccc {
			// Output to our string builder the current Client Type (from K)
			// Add an newline if this isnt the first itteration
			//if clientsb.String() != "" {
			//	clientsb.WriteString("\n")
			//}

			switch k {
			case "socket":
				condir = "Outbound"
			case "socksc":
				condir = "Inbound"
			default:
				condir = "??"
			}

			clientsb.WriteString(fmt.Sprintf("ClientType: [%s](%s)  Count:%v\n", k, condir, len(ccc[k])))
		}
	}
	// Return the built string
	return clientsb.String()
}
