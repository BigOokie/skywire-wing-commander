package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	managerAPIGetAllConnectedNodes = "/conn/getAll"
)

// connectedNode structure stores of JSON response from /conn/getAll API
type connectedNode struct {
	Key         string `json:"key"`
	Conntype    string `json:"type"`
	SendBytes   int    `json:"send_bytes"`
	RecvBytes   int    `json:"recv_bytes"`
	LastAckTime int    `json:"last_ack_time"`
	StartTime   int    `json:"start_time"`
}

// Defines an in-memory slice (dynamic array) based on the connectedNode struct
type connectedNodeSlice []connectedNode

// SkyManagerMonitor is used to monitor a Sky Manager and provide messages to the
// main process when specific events are detected.
type SkyManagerMonitor struct {
	ManagerAddress string
	CancelFunc     func()
}

// NewMonitor creates a SkyManagerMonitor which will monitor the provided managerip.
func NewMonitor(manageraddress string) *SkyManagerMonitor {
	return &SkyManagerMonitor{
		ManagerAddress: manageraddress,
		CancelFunc:     nil,
	}
}

// Run starts the SkyManagerMonitor.
// If `ctx` is not nil, the monitor will listen to ctx.Done() and stop monitoring
// when it recieves the signal.
func (m *SkyManagerMonitor) Run(runctx context.Context, pollInt time.Duration) {
	log.Debugf("SkyManagerMonitor Run: Start (Interval: %v)", pollInt)
	defer log.Debugln("SkyManagerMonitor Run: End")

	ticker := time.NewTicker(pollInt)

	//var oldresp, newresp string
	var oldcns connectedNodeSlice

	for {
		select {
		case <-ticker.C:

			newcns, err := m.getAllNodesList()
			if err != nil {
				log.Debugln(err)
			}

			if len(newcns) != len(oldcns) {
				log.Debugf("SkyManagerMonitor: Connected Node List Changed (Old: %v, New: %v)", len(oldcns), len(newcns))
				log.Debugln(newcns)
				oldcns = newcns
			} else {
				log.Debugln("SkyManagerMonitor: Connected Node List Unchanged.")
			}

			/*
				newresp = m.getAllNodesStr()
				if newresp != oldresp {
					log.Debugln(newresp)
					oldresp = newresp
				} else {
					log.Debugln("SkyManagerMonitor: No change.")
				}
			*/
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

func (m *SkyManagerMonitor) getAllNodesList() (cns connectedNodeSlice, err error) {
	log.Debugln("SkyManagerMonitor.getAllNodesList")
	apiURL := fmt.Sprintf("http://%s/%s", m.ManagerAddress, managerAPIGetAllConnectedNodes)
	//apiURL := fmt.Sprintf("http://%s/%s", "discovery.skycoin.net:8001", managerAPIGetAllConnectedNodes)
	resp, err := http.Get(apiURL)

	defer resp.Body.Close()

	respbuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respbuf, &cns)
	return
}
