package wingcommander

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

// Config structure models the TOML configuration file
type Config struct {
	WingCommander struct {
		MonitorRunning   bool
		Heartbeat        bool
		TwoFactorEnabled bool
	}
	SkyManager struct {
		Address string
	}
	Telegram struct {
		APIKey string
		ChatID int64
		Admin  string
		Debug  bool
	}
	Monitor struct {
		Interval int
	}
}

// ReadConfig will read configuration from the provided TOML config file
func ReadConfig(filename string) (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile(filename, &conf)
	if err != nil {
		return nil, err
	}

	log.Debugln("ReadConfig")
	DebugLogConfig(&conf)
	return &conf, nil
}

// DebugLogConfig will log debug information for the passed Config structure
func DebugLogConfig(conf *Config) {
	log.Debugln("WingCommander Configs:")
	log.Debugf("  monitorrunning = %v", conf.WingCommander.MonitorRunning)
	log.Debugf("  heartbeat = %v", conf.WingCommander.Heartbeat)
	log.Debugf("  twofactorenabled = %v", conf.WingCommander.TwoFactorEnabled)

	log.Debugln("Node Manager Configs:")
	log.Debugf("  address = %s", conf.SkyManager.Address)

	log.Debugln("Telegram Configs:")
	log.Debugf("  apikey = %s", conf.Telegram.APIKey)
	log.Debugf("  chatid = %v", conf.Telegram.ChatID)
	log.Debugf("  admin  = %s", conf.Telegram.Admin)
	log.Debugf("  debug  = %v", conf.Telegram.Debug)

	log.Debugln("Monitor Configs:")
	log.Debugf("  interval = %v", conf.Monitor.Interval)
}
