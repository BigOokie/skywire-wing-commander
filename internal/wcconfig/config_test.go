// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.
package wcconfig

import (
	"testing"
	"time"

	"github.com/go-test/deep"
)

func Test_ConfigString(t *testing.T) {

	expectstr := "[WingCommander]\n" +
		"  twofactorenabled = false\n" +
		"  analyticsenabled = false\n" +
		"[AppAnalytics]\n" +
		"  clientuuid = \n" +
		"  userid = \n" +
		"[SkyManager]\n" +
		"  address = \"127.0.0.1:8000\"\n" +
		"  discoveryaddress = \"discovery.skycoin.net:8001\"\n" +
		"[Telegram]\n" +
		"  apikey = \"ABC123\"\n" +
		"  chatid = 123456789\n" +
		"  admin  = \"@TESTUSER\"\n" +
		"  debug  = false\n" +
		"[Monitor]\n" +
		"  intervalsec = 10s\n" +
		"  heartbeatintmin = 2h0m0s\n" +
		"  discoverymonitorintmin = 2h0m0s\n"

	var config Config
	config.WingCommander.TwoFactorEnabled = false
	config.SkyManager.Address = "127.0.0.1:8000"
	config.SkyManager.DiscoveryAddress = "discovery.skycoin.net:8001"
	config.Telegram.APIKey = "ABC123"
	config.Telegram.ChatID = 123456789
	config.Telegram.Admin = "@TESTUSER"
	config.Telegram.Debug = false
	config.Monitor.IntervalSec = 10 * time.Second
	config.Monitor.HeartbeatIntMin = 120 * time.Minute
	config.Monitor.DiscoveryMonitorIntMin = 120 * time.Minute

	if diff := deep.Equal(config.String(), expectstr); diff != nil {
		t.Error(diff)
	}
}

func Test_Config_IsEmpty(t *testing.T) {
	emptyConfig := Config{}
	if !IsEmpty(emptyConfig) {
		t.Error("Expected: Config should be empty")
	}
}

func Test_LoadConfigParameters_BadFileName(t *testing.T) {
	// Load configuration
	config, err := LoadConfigParameters("file-does-not-exist", ".", map[string]interface{}{
		"telegram.debug":                 false,
		"monitor.intervalsec":            10,
		"monitor.heartbeatintmin":        120,
		"monitor.discoverymonitorintmin": 120,
		"skymanager.address":             "127.0.0.1:8000",
		"skymanager.discoveryaddress":    "discovery.skycoin.net:8001",
	})

	if err == nil {
		t.Error(err)
	}

	if !IsEmpty(config) {
		t.Error("Expected: Config should be empty")
	}
}

func Test_LoadConfigParameters_AllParams(t *testing.T) {
	// Load configuration
	config, err := LoadConfigParameters("configtest-allparams", "./testdata", map[string]interface{}{
		"telegram.debug":                 false,
		"monitor.intervalsec":            10,
		"monitor.heartbeatintmin":        120,
		"monitor.discoverymonitorintmin": 120,
		"skymanager.address":             "127.0.0.1:8000",
		"skymanager.discoveryaddress":    "discovery.skycoin.net:8001",
	})

	if err != nil {
		t.Error(err)
	}

	if IsEmpty(config) {
		t.Error("Expected: Config should be populated")
	}
}

func Test_LoadConfigParameters_NoDefaultParams(t *testing.T) {
	// Load configuration
	config, err := LoadConfigParameters("configtest-nodefaults", "./testdata", map[string]interface{}{
		"telegram.debug":                 false,
		"monitor.intervalsec":            10,
		"monitor.heartbeatintmin":        120,
		"monitor.discoverymonitorintmin": 120,
		"skymanager.address":             "127.0.0.1:8000",
		"skymanager.discoveryaddress":    "discovery.skycoin.net:8001",
	})

	if err != nil {
		t.Error(err)
	}

	if IsEmpty(config) {
		t.Error("Expected: Config should be populated")
	}
}

func Test_LoadConfigParameters_BadParamData(t *testing.T) {
	// Load configuration
	config, err := LoadConfigParameters("configtest-badparamdata", "./testdata", map[string]interface{}{
		"telegram.debug":                 false,
		"monitor.intervalsec":            10,
		"monitor.heartbeatintmin":        120,
		"monitor.discoverymonitorintmin": 120,
		"skymanager.address":             "127.0.0.1:8000",
		"skymanager.discoveryaddress":    "discovery.skycoin.net:8001",
	})

	if err == nil {
		t.Error(err)
	}

	if !IsEmpty(config) {
		t.Error("Expected: Config should be empty")
	}
}
