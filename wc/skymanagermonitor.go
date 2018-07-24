package wingcommander

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
	log.Debugf("SkyManagerMonitor Run: Start (Interval: %v)", time.Second*pollInt)
	ticker := time.NewTicker(time.Second * pollInt)

	var oldresp, newresp string

	for {
		select {
		case <-ticker.C:
			newresp = m.getAllNodes()
			if newresp != oldresp {
				log.Debugln(newresp)
				oldresp = newresp
			} else {
				log.Debugln("SkyManagerMonitor: No change.")
			}
			//return
		case <-runctx.Done(): //<-done:
			//m.Running = false
			log.Debugln("SkyManagerMonitor - Done Event.")
			return
		}
	}
	log.Debugln("SkyManagerMonitor Run: End")
}

// IsRunning determines if the SkyMgrMonitor is running or not.
// This is assessed based on the assignment of the context cancel function (one is assigned if it is running).
func (m *SkyManagerMonitor) IsRunning() bool {
	return m.CancelFunc != nil
}

// connectedNode structure stores JSON response from /conn/getAll API
type connectedNode struct {
	Key         string `json:"key"`
	Conntype    string `json:"type"`
	SendBytes   int    `json:"send_bytes"`
	RecvBytes   int    `json:"recv_bytes"`
	LastAckTime int    `json:"last_ack_time"`
	StartTime   int    `json:"start_time"`
}

func (m *SkyManagerMonitor) getAllNodes() string {
	var respstr string
	//log.Debugln("SkyManagerMonitor.getAllNodes")
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

func (m *SkyManagerMonitor) getAllNodesList() (cns map[string]connectedNodeSlice, err error) {
	log.Debugln("SkyManagerMonitor.getAllNodesList")
	apiURL := fmt.Sprintf("http://%s/%s", m.ManagerAddress, managerAPIGetAllConnectedNodes)
	resp, err := http.Get(apiURL)
	if err != nil {
		cns = nil
		log.Error(err)
		return
	}

	defer resp.Body.Close()
	respbuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cns = nil
		log.Error(err)
		return
	}

	cns = make(map[string]connectedNodeSlice)
	err = json.Unmarshal(respbuf, &cns)
	if err != nil {
		return
	}
	return
}
