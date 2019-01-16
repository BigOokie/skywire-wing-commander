// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BigOokie/skywire-wing-commander/internal/utils"
	"github.com/BigOokie/skywire-wing-commander/internal/wcconfig"
	"github.com/BigOokie/skywire-wing-commander/internal/wcconst"
	log "github.com/sirupsen/logrus"
)

type cmdlineFlags struct {
	dumpconfig       bool
	version          bool
	help             bool
	about            bool
	upgradecompleted bool
}

type wcBotApp struct {
	config   wcconfig.Config
	cmdFlags cmdlineFlags
}

// loadConfig manages the configuration load specifics
// offloading the detail from the `main()` funct
func (ba *wcBotApp) loadConfig() {
	log.Debugln("wcBotApp.loadConfig: Start")
	defer log.Debugln("wcBotApp.loadConfig: Complete")
	// Load configuration
	c, err := wcconfig.LoadConfigParameters("config", filepath.Join(utils.UserHome(), ".wingcommander"), map[string]interface{}{
		"wingcommander.analyticsenabled": true,
		"telegram.debug":                 false,
		"monitor.intervalsec":            10,
		"monitor.heartbeatintmin":        120,
		"monitor.discoverymonitorintmin": 120,
		"skymanager.address":             "127.0.0.1:8000",
		"skymanager.discoveryaddress":    "discovery.skycoin.net:8001",
	})

	if err != nil {
		log.Fatalf("wcBotApp.loadConfig: Error loading configuration: %s", err)
		return
	}
	ba.config = c
}

func (cf *cmdlineFlags) parseCmdLineFlags() {
	flag.BoolVar(&cf.version, "v", false, "print current version")
	flag.BoolVar(&cf.dumpconfig, "config", false, "print current config")
	flag.BoolVar(&cf.help, "help", false, "print application help")
	flag.BoolVar(&cf.about, "about", false, "print application information")
	flag.BoolVar(&cf.upgradecompleted, "upgradecompleted", false, "signals the application has been restarted following an upgrade")

	flag.Parse()
}

func (cf *cmdlineFlags) handleCmdLineFlags() {
	// if version cmd line flag `-v` then print version info and exit
	if cf.version {
		fmt.Println(wcconst.BotAppVersion)
		fmt.Println("")
		os.Exit(0)
	}

	// if help cmd line flag `-help` then print version info and exit
	if cf.help {
		fmt.Println(wcconst.MsgCmdLineHelp)
		fmt.Println("")
		os.Exit(0)
	}

	// if about cmd line flag `-about` then print version info and exit
	if cf.about {
		fmt.Println(wcconst.MsgAbout)
		fmt.Println("")
		os.Exit(0)
	}

	// if about cmd line flag `-about` then print version info and exit
	if cf.upgradecompleted {
		fmt.Println("Upgrade completed.")
		fmt.Println("")
	}
}

func (ba *wcBotApp) initLogging() {
	// Setup Log Formatter
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
}
