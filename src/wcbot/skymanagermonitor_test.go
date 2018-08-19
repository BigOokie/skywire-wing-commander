// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"testing"

	"github.com/BigOokie/skywire-wing-commander/src/skynode"
	"github.com/go-test/deep"
)

func Test_NewMonitor(t *testing.T) {
	expect := &SkyManagerMonitor{
		ManagerAddress:       "0.0.0.0:8000",
		CancelFunc:           nil,
		monitorStatusMsgChan: nil,
		connectedNodes:       make(skynode.NodeInfoMap),
	}

	actual := NewMonitor("0.0.0.0:8000")

	if diff := deep.Equal(expect, actual); diff != nil {
		t.Error(diff)
	}
}

func Test_IsRunning(t *testing.T) {
	expect := &SkyManagerMonitor{
		ManagerAddress:       "0.0.0.0:8000",
		CancelFunc:           nil,
		monitorStatusMsgChan: nil,
		connectedNodes:       make(skynode.NodeInfoMap),
	}

	if expect.IsRunning() {
		t.Fail()
	}

	_, expect.CancelFunc = context.WithCancel(context.Background())
	if !expect.IsRunning() {
		t.Fail()
	}
}

func Test_GetConnectedNodeCount(t *testing.T) {
	monitor := NewMonitor("0.0.0.0:8000")

	if monitor.GetConnectedNodeCount() != 0 {
		t.Fail()
	}

	nodeA := skynode.NodeInfo{
		Key:         "NODE1KEY",
		Conntype:    "TCP",
		SendBytes:   1,
		RecvBytes:   2,
		LastAckTime: 3,
		StartTime:   4}

	nodeB := skynode.NodeInfo{
		Key:         "NODE2KEY",
		Conntype:    "TCP",
		SendBytes:   4,
		RecvBytes:   3,
		LastAckTime: 2,
		StartTime:   1}

	nodeSlice := skynode.NodeInfoSlice{
		nodeA,
		nodeB,
	}

	monitor.connectedNodes = skynode.NodeInfoSliceToMap(nodeSlice)

	if monitor.GetConnectedNodeCount() != 2 {
		t.Fail()
	}
}
