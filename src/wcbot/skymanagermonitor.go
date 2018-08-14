// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/BigOokie/skywire-wing-commander/src/skynode"
	"github.com/BigOokie/skywire-wing-commander/src/wcconst"
	log "github.com/sirupsen/logrus"
)

const (
	managerAPIGetAllConnectedNodes = "/conn/getAll"
)

// SkyManagerMonitor is used to monitor a Sky Manager and provide messages to the
// main process when specific events are detected.
type SkyManagerMonitor struct {
	ManagerAddress       string
	CancelFunc           func()
	monitorStatusMsgChan chan string
	connectedNodes       skynode.NodeInfoMap
}

// NewMonitor creates a SkyManagerMonitor which will monitor the provided managerip.
func NewMonitor(manageraddress string) *SkyManagerMonitor {
	return &SkyManagerMonitor{
		ManagerAddress:       manageraddress,
		CancelFunc:           nil,
		monitorStatusMsgChan: nil,
		connectedNodes:       make(skynode.NodeInfoMap),
	}
}

// Run starts the SkyManagerMonitor.
// If `ctx` is not nil, the monitor will listen to ctx.Done() and stop monitoring
// when it recieves the signal.
func (m *SkyManagerMonitor) Run(runctx context.Context, statusMsgChan chan<- string, pollInt time.Duration) {
	log.Debugf("SkyManagerMonitor Run: Start (Interval: %v)", pollInt)
	defer log.Debugln("SkyManagerMonitor Run: End")

	ticker := time.NewTicker(pollInt)

	for {
		select {
		case <-ticker.C:
			newcns, err := m.getAllNodesList()
			if err != nil {
				log.Error(err)
			} else {
				// Maintain the list of connected nodes
				m.maintainConnectedNodesList(newcns, statusMsgChan)
			}
		case <-runctx.Done():
			log.Debugln("SkyManagerMonitor - Done Event.")
			return
		}
	}
}

// IsRunning determines if the SkyMgrMonitor is running or not.
// This is assessed based on the assignment of the context cancel function (one is assigned if it is running).
func (m *SkyManagerMonitor) IsRunning() bool {
	return m.CancelFunc != nil
}

// getAllNodesStr requests the list of connected Nodes from the Manager and returns the raw JSON response as a string
func (m *SkyManagerMonitor) getAllNodesStr() string {
	var respstr string
	log.Debugln("SkyManagerMonitor.getAllNodesStr")
	apiURL := fmt.Sprintf("http://%s/%s", m.ManagerAddress, managerAPIGetAllConnectedNodes)
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Error(err)
	} else {
		defer resp.Body.Close()
		respbuf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		respstr = string(respbuf)
	}
	return respstr
}

// getAllNodesList requests the list of connected Nodes from the Manager and returns an array (slice) of connectedNode
func (m *SkyManagerMonitor) getAllNodesList() (cns skynode.NodeInfoSlice, err error) {
	log.Debugln("SkyManagerMonitor.getAllNodesList")
	apiURL := fmt.Sprintf("http://%s/%s", m.ManagerAddress, managerAPIGetAllConnectedNodes)
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Error(err)
		return
	}

	defer resp.Body.Close()

	respbuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	err = json.Unmarshal(respbuf, &cns)
	if err != nil {
		log.Error(err)
	}
	return
}

// maintainConnectedNodeList is responsible for maintaining (adding, updating and deleting) Nodes from the
// Monitors internal connectedNodeList.
// TODO: Support use of channels for signaling of change events
func (m *SkyManagerMonitor) maintainConnectedNodesList(newcns skynode.NodeInfoSlice, statusMsgChan chan<- string) {
	// Make sure the newcns structure is not nil, and return if it is (do nothing)
	if newcns == nil {
		log.Error("SkyManagerMonitor.maintainConnectedNodesList: newcns is nil.")
		return
	}

	// Compare the new connected node list (newcns) against the current list.
	// If they are not different we dont need to do anything
	for _, v := range newcns {
		_, hasKey := m.connectedNodes[v.Key]
		if hasKey {
			// Node key found
			// Until I can figure out a better way - lets replace the existing entry with the new data
			// Delete and then add the new instance
			delete(m.connectedNodes, v.Key)
			m.connectedNodes[v.Key] = v
		} else {
			// Add new NodeInfo
			m.connectedNodes[v.Key] = v
			msg := fmt.Sprintf(wcconst.MsgNodeConnected, v.Key, len(m.connectedNodes))
			log.Debugln(msg)
			statusMsgChan <- msg
		}
	}

	// If the number of Nodes in the connectedNodes list greater than
	// the number of Nodes returned from the last request, we need to
	// prune the connectedNodes list (i.e. some Nodes have been disconnected)
	if len(m.connectedNodes) > len(newcns) {
		niMap := skynode.NodeInfoSliceToMap(newcns)
		// Iterate the connectedNodes and delete any that are not found
		// in the newly returned connected Node list (niMap)
		for _, v := range m.connectedNodes {
			_, hasKey := niMap[v.Key]
			if !hasKey {
				// Node Key not found
				// Delete the Node from the Connected Node List
				log.Debugf("SkyManagerMonitor.maintainConnectedNodesList: Node Removed:\n%s\n", v.FmtString())
				delete(m.connectedNodes, v.Key)
				msg := fmt.Sprintf(wcconst.MsgNodeDisconnected, v.Key, len(m.connectedNodes))
				log.Debugln(msg)
				statusMsgChan <- msg
			}
		}
	}
	return
}

// GetConnectedNodeCount will return the count of Nodes within the connectedNodes structure
// If the structure is nil (not yet assigned), 0 will be returned
func (m *SkyManagerMonitor) GetConnectedNodeCount() int {
	if m.connectedNodes == nil {
		return 0
	}
	return len(m.connectedNodes)
}
