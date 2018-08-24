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
	"sync"
	"time"

	"github.com/BigOokie/skywire-wing-commander/src/skynode"
	"github.com/BigOokie/skywire-wing-commander/src/wcconst"
	log "github.com/sirupsen/logrus"
)

const (
	managerAPIGetAllConnectedNodes = "conn/getAll"
)

// SkyManagerMonitor is used to monitor a Sky Manager and provide messages to the
// main process when specific events are detected.
type SkyManagerMonitor struct {
	ManagerAddress       string
	DiscoveryAddress     string
	cancelFunc           func()
	monitorStatusMsgChan chan string
	connectedNodes       skynode.NodeInfoMap
	discConnNodeCount    int
	m                    sync.Mutex
}

// SetCancelFunc is a thread-safe function for setting the cancelFunc
// on the SkyManagerMonitor struct
func (smm *SkyManagerMonitor) SetCancelFunc(cf func()) {
	smm.m.Lock()
	defer smm.m.Unlock()
	smm.cancelFunc = cf
}

// GetCancelFunc is a thread-safe function for accessing (getting) the
// value of cancelFunc on the SkyManagerMonitor struct
func (smm *SkyManagerMonitor) GetCancelFunc() func() {
	smm.m.Lock()
	defer smm.m.Unlock()
	return smm.cancelFunc
}

// DoCancelFunc is a thread-safe function for calling the
// cancelFunc on the SkyManagerMonitor struct
func (smm *SkyManagerMonitor) DoCancelFunc() {
	smm.m.Lock()
	defer smm.m.Unlock()
	if smm.cancelFunc != nil {
		smm.cancelFunc()
	}
}

// NewMonitor creates a SkyManagerMonitor which will monitor the provided managerip.
func NewMonitor(manageraddress, discoveryaddress string) *SkyManagerMonitor {
	return &SkyManagerMonitor{
		ManagerAddress:       manageraddress,
		DiscoveryAddress:     discoveryaddress,
		cancelFunc:           nil,
		monitorStatusMsgChan: nil,
		connectedNodes:       make(skynode.NodeInfoMap),
		discConnNodeCount:    0,
	}
}

// RunManagerMonitor starts the SkyManagerMonitor monitoring of the local Manager Node.
// If `ctx` is not nil, the monitor will listen to ctx.Done() and stop monitoring
// when it recieves the signal.
func (smm *SkyManagerMonitor) RunManagerMonitor(runctx context.Context, statusMsgChan chan<- string, pollInt time.Duration) {
	log.Debugf("SkyManagerMonitor::RunManagerMonitor: Start (Interval: %v)", pollInt)
	defer log.Debugln("SkyManagerMonitor::RunManagerMonitor: End")

	ticker := time.NewTicker(pollInt)

	for {
		select {
		case <-ticker.C:
			newcns, err := getAllNodesList(smm.ManagerAddress)
			if err != nil {
				log.Error(err)
				statusMsgChan <- wcconst.MsgErrorGetNodes
			} else {
				// Maintain the list of connected nodes
				smm.maintainConnectedNodesList(newcns, statusMsgChan)
			}
		case <-runctx.Done():
			log.Debugln("SkyManagerMonitor::RunManagerMonitor: Done Event.")
			return
		}
	}
}

/*
// RunDiscoveryMonitor starts the SkyManagerMonitor monitoring of the Skywire Discovery Node.
// If `ctx` is not nil, the monitor will listen to ctx.Done() and stop monitoring
// when it recieves the signal.
func (smm *SkyManagerMonitor) RunDiscoveryMonitor(runctx context.Context, statusMsgChan chan<- string, pollInt time.Duration) {
	log.Debugf("SkyManagerMonitor::RunDiscoveryMonitor: Start (Interval: %v)", pollInt)
	defer log.Debugln("SkyManagerMonitor::RunDiscoveryMonitor: End")

	ticker := time.NewTicker(pollInt)

	for {
		select {
		case <-ticker.C:
			discNodes, err := getAllNodesList(smm.ManagerAddress)
			if err != nil {
				log.Error(err)
				statusMsgChan <- wcconst.MsgErrorGetDiscNodes
			} else {
				// Check the local Nodes are connected to Discovery Node
				smm.checkNodeDiscoveryConnection(discNodes, statusMsgChan)
			}
		case <-runctx.Done():
			log.Debugln("SkyManagerMonitor::RunDiscoveryMonitor: Done Event.")
			return
		}
	}
}
*/

// ConnectedDiscNodeCount returns a count the locally Managed Nodes that are connected to the
// Discovery Node
func (smm *SkyManagerMonitor) ConnectedDiscNodeCount() (int, error) {
	log.Debug("SkyManagerMonitor::RefreshDiscoveryConnectionCount: Start")
	defer log.Debugln("SkyManagerMonitor::RefreshDiscoveryConnectionCount: End")
	discConnNodeCount := 0

	discNodes, err := getAllNodesList(smm.DiscoveryAddress)
	if err != nil {
		log.Errorf("SkyManagerMonitor.RefreshDiscoveryConnectionCount: Error contacting Discovery Server: %v", err)
		return discConnNodeCount, err

	} else {
		// Check the local Nodes are connected to Discovery Node
		//if len(smm.connectedNodes) == 0 {
		if smm.GetConnectedNodeCount() == 0 {
			log.Debug("SkyManagerMonitor.RefreshDiscoveryConnectionCount: Connected Node list is empty. No work to do.")
			return discConnNodeCount, nil
		}

		smm.m.Lock()
		defer smm.m.Unlock()

		// Compare the list of Nodes connected to the Discovery Node (disccns) against the
		// current list of locally conected nodes.
		// If our local Nodes are not listed as connected to the Discovery Node we need to raise an alert
		discNodeMap := skynode.NodeInfoSliceToMap(discNodes)
		for _, v := range smm.connectedNodes {
			_, hasKey := discNodeMap[v.Key]
			if hasKey {
				// Node Key found in the Discover Node Map
				log.Debugf("SkyManagerMonitor.RefreshDiscoveryConnectionCount: Node Connected:\n%s\n", v.FmtString())
				discConnNodeCount++
			} else {
				log.Debugf("SkyManagerMonitor.RefreshDiscoveryConnectionCount: Node Not Connected:\n%s\n", v.FmtString())
			}
		}

		log.Debugf("%d Nodes Connected to Discovery", discConnNodeCount)
	}
	return discConnNodeCount, nil
}

// IsRunning determines if the SkyMgrMonitor is running or not.
// This is assessed based on the assignment of the context cancel function (one is assigned if it is running).
func (smm *SkyManagerMonitor) IsRunning() bool {
	return smm.GetCancelFunc() != nil
}

// getAllNodesStr requests the list of connected Nodes from the Manager and returns the raw JSON response as a string
func (smm *SkyManagerMonitor) getAllNodesStr() string {
	var respstr string
	log.Debugln("SkyManagerMonitor.getAllNodesStr")
	apiURL := fmt.Sprintf("http://%s/%s", smm.ManagerAddress, managerAPIGetAllConnectedNodes)

	//http.Header.Add
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
func getAllNodesList(managerAddr string) (cns skynode.NodeInfoSlice, err error) {
	log.Debugln("SkyManagerMonitor.getAllNodesList")
	userAgent := "Wing Commander Telegram Bot " + wcconst.BotVersion

	client := &http.Client{}
	client.Timeout = time.Second * 30
	apiURL := fmt.Sprintf("http://%s/%s", managerAddr, managerAPIGetAllConnectedNodes)
	log.Debugln(apiURL)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		log.Errorf("http.NewRequest() failed with '%s'\n", err)
		return
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("client.Do() failed with '%s'\n", err)
		return
	}

	defer resp.Body.Close()
	respbuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("ioutil.ReadAll() failed with '%s'\n", err)
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
func (smm *SkyManagerMonitor) maintainConnectedNodesList(newcns skynode.NodeInfoSlice, statusMsgChan chan<- string) {
	smm.m.Lock()
	defer smm.m.Unlock()

	// Make sure the newcns structure is not nil, and return if it is (do nothing)
	if newcns == nil {
		log.Error("SkyManagerMonitor.maintainConnectedNodesList: newcns is nil.")
		return
	}

	// Compare the new connected node list (newcns) against the current list.
	// If they are not different we dont need to do anything
	for _, v := range newcns {
		_, hasKey := smm.connectedNodes[v.Key]
		if hasKey {
			// Node key found
			// Until I can figure out a better way - lets replace the existing entry with the new data
			// Delete and then add the new instance
			delete(smm.connectedNodes, v.Key)
			smm.connectedNodes[v.Key] = v
		} else {
			// Add new NodeInfo
			smm.connectedNodes[v.Key] = v
			msg := fmt.Sprintf(wcconst.MsgNodeConnected, v.Key, len(smm.connectedNodes))
			log.Debugln(msg)
			statusMsgChan <- msg
		}
	}

	// If the number of Nodes in the connectedNodes list greater than
	// the number of Nodes returned from the last request, we need to
	// prune the connectedNodes list (i.e. some Nodes have been disconnected)
	if len(smm.connectedNodes) > len(newcns) {
		niMap := skynode.NodeInfoSliceToMap(newcns)
		// Iterate the connectedNodes and delete any that are not found
		// in the newly returned connected Node list (niMap)
		for _, v := range smm.connectedNodes {
			_, hasKey := niMap[v.Key]
			if !hasKey {
				// Node Key not found
				// Delete the Node from the Connected Node List
				log.Debugf("SkyManagerMonitor.maintainConnectedNodesList: Node Removed:\n%s\n", v.FmtString())
				delete(smm.connectedNodes, v.Key)
				msg := fmt.Sprintf(wcconst.MsgNodeDisconnected, v.Key, len(smm.connectedNodes))
				log.Debugln(msg)
				statusMsgChan <- msg
			}
		}
	}
	return
}

// checkNodeDiscoveryConnection is responsible for checking the list of Nodes currently connected to the local Manager
// against the list of Nodes reported as connected to the Skywire Discovery Node. If our local Nodes are not reported
// as connected to the Discovery Node, we need to raise an alert using the provided statusMsgChan
//TODO: Refactor this. We are doing the same thing as in other functions essentially. We need to restucture this (but later)
func (smm *SkyManagerMonitor) checkNodeDiscoveryConnection(disccns skynode.NodeInfoSlice, statusMsgChan chan<- string) {
	if smm.GetConnectedNodeCount() == 0 {
		log.Debug("SkyManagerMonitor.checkNodeDiscoveryConnection: Connected Node list is empty. No work to do.")
		return
	}

	smm.m.Lock()
	defer smm.m.Unlock()

	// Make sure the disccns structure is not nil, and return if it is (do nothing)
	if disccns == nil {
		log.Error("SkyManagerMonitor.checkNodeDiscoveryConnection: disccns is nil.")
		return
	}

	discConnNodeCount := 0

	// Compare the list of Nodes connected to the Discovery Node (disccns) against the
	// current list of locally conected nodes.
	// If our local Nodes are not listed as connected to the Discovery Node we need to raise an alert
	discNodeMap := skynode.NodeInfoSliceToMap(disccns)
	for _, v := range smm.connectedNodes {
		_, hasKey := discNodeMap[v.Key]
		if hasKey {
			// Node Key found in the Discover Node Map
			log.Debugf("SkyManagerMonitor.checkNodeDiscoveryConnection: Node Connected:\n%s\n", v.FmtString())
			discConnNodeCount++
		} else {
			log.Debugf("SkyManagerMonitor.checkNodeDiscoveryConnection: Node Not Connected:\n%s\n", v.FmtString())
			msg := fmt.Sprintf("Discovery Disconnected: Node: %s", v.Key)
			log.Debugln(msg)
			statusMsgChan <- msg
		}
	}

	smm.discConnNodeCount = discConnNodeCount
	msg := fmt.Sprintf("%d Nodes Connected to Discovery", smm.discConnNodeCount)
	log.Debugln(msg)
	//statusMsgChan <- msg
	return
}

// GetConnectedNodeCount will return the count of Nodes within the connectedNodes structure
// If the structure is nil (not yet assigned), 0 will be returned
func (smm *SkyManagerMonitor) GetConnectedNodeCount() int {
	smm.m.Lock()
	defer smm.m.Unlock()
	if smm.connectedNodes == nil {
		return 0
	}
	return len(smm.connectedNodes)
}
