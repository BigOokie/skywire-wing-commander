package wingcommander

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	nodeManagerAPIURL       = "127.0.0.1:8000"
	apiGetAllConnectedNodes = "/conn/getAll"
)

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

func getGetAllNodes() string {
	var respstr string

	apiURL := fmt.Sprintf("http://%s/%s", nodeManagerAPIURL, apiGetAllConnectedNodes)
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

func getGetAllNodesList() (cns map[string]connectedNodeSlice, err error) {
	apiURL := fmt.Sprintf("http://%s/%s", nodeManagerAPIURL, apiGetAllConnectedNodes)
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
