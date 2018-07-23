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

// SkyManagerMonitor is used to monitor a Sky Manager and provide messages to the
// main process when specific events are detected.
type SkyManagerMonitor struct {
	ManagerAddress string
	Running        bool
}

// NewMonitor creates a SkyManagerMonitor which will monitor the provided managerip.
func NewMonitor(manageraddress string) *SkyManagerMonitor {
	return &SkyManagerMonitor{
		ManagerAddress: manageraddress,
		Running:        false,
	}
}

// Run starts the SkyManagerMonitor.
// If `ctx` is not nil, the monitor will listen to ctx.Done() and stop monitoring
// when it recieves the signal.
func (m *SkyManagerMonitor) Run(ctx context.Context) {
	var done <-chan struct{}
	if ctx != nil {
		done = ctx.Done()
	}

	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticker.C:
			//msgText := m.getGetAllNodes()
			//sendBotMsg(m, msgText, false)
			return
		case <-done:
			return
		}
	}
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

// Defines an in-memory slice (dynamic array) based on the connectedNode struct
type connectedNodeSlice []connectedNode

func (m *SkyManagerMonitor) getGetAllNodes() string {
	var respstr string

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

func (m *SkyManagerMonitor) getGetAllNodesList() (cns map[string]connectedNodeSlice, err error) {
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
