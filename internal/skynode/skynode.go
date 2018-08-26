// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package skynode

import (
	"fmt"
)

// NodeInfo structure stores of JSON response from /conn/getAll API
type NodeInfo struct {
	Key         string `json:"key"`
	Conntype    string `json:"type"`
	SendBytes   int    `json:"send_bytes"`
	RecvBytes   int    `json:"recv_bytes"`
	LastAckTime int    `json:"last_ack_time"`
	StartTime   int    `json:"start_time"`
}

// NodeInfoSlice defines an in-memory (dynamic) array of NodeInfo structures
type NodeInfoSlice []NodeInfo

// NodeInfoMap defines a string key based map of NodeInfo structs
type NodeInfoMap map[string]NodeInfo

// NodeInfoSliceToMap convers a provided NodeInfoSlice to a NodeInfoMap
func NodeInfoSliceToMap(nis NodeInfoSlice) NodeInfoMap {
	niMap := make(map[string]NodeInfo)
	for _, ni := range nis {
		niMap[ni.Key] = ni
	}
	return niMap
}

// NodesAreEqual determines if two instances of a NodeInfo structure (a and b) represent the same Node based on their Keys.
// Not that this is not an equality check of the structure - but simply that the two structures represent the same Node.
// Other elements of the strucutre may be different.
func NodesAreEqual(a, b NodeInfo) bool {
	// Check both NodeInfo structures are for the same
	return a.Key == b.Key
}

// NodeInfoSliceEqual determines if two instances of a NodeInfoSlice (a and b) are equal - on the basis that
// they contain the same list of NodeInfo structures.
// TODO: May need to review this function and how it determines equality. Very simplistice ATM.
func NodeInfoSliceEqual(a, b NodeInfoSlice) bool {
	// Check both contain the same number of connectedNode elements.
	if len(a) != len(b) {
		return false
	}

	// Iterate a and check the Node Key for each connected Node against b.
	for i, v := range a {
		if !NodesAreEqual(v, b[i]) {
			// If a different Node is found in the list,
			// consider the two inequal.
			return false
		}
	}
	return true
}

// String satisfies the fmt.Stringer interface for the NodeInfo type
func (ni NodeInfo) String() string {
	msg := "Key: %s, Type: %s, SendBytes: %v, RecvBytes: %v, LastAckTime: %vs, StartTime: %vs"
	return fmt.Sprintf(msg, ni.Key, ni.Conntype, ni.SendBytes, ni.RecvBytes, ni.LastAckTime, ni.StartTime)
}

// FmtString returns a rich formatted string representing the data of the NodeInfo struct
func (ni NodeInfo) FmtString() string {
	msg := "Node Information:\n" +
		"Key        : %s\n" +
		"Type       : %s\n" +
		"SendBytes  : %v\n" +
		"RecvBytes  : %v\n" +
		"LastAckTime: %vs\n" +
		"StartTime  : %vs"

	return fmt.Sprintf(msg, ni.Key, ni.Conntype, ni.SendBytes, ni.RecvBytes, ni.LastAckTime, ni.StartTime)
}

/*
type RWMap struct {
	sync.RWMutex
	m map[string]int
}

// Get is a wrapper for getting the value from the underlying map
func (r *RWMap) Get(key string) int {
	r.RLock()
	defer r.RUnlock()
	return r.m[key]
}

// Set is a wrapper for setting the value of a key in the underlying map
func (r *RWMap) Set(key string, val int) {
	r.Lock()
	defer r.Unlock()
	r.m[key] = val
}
*/
