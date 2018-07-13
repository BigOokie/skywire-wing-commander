package wingcommander

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

// TOML Config
type Config struct {
	WingCommander struct {
		MonitorRunning bool
		Heartbeat      bool
	}
	Telegram struct {
		APIKey string
		ChatID int64
		Admin  string
		Debug  bool
	}
}

// ReadConfig will read configuration from the provided TOML config file
func ReadConfig(filename string) (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile(filename, &conf)
	if err != nil {
		return nil, err
	}
	log.Debugln("ReadConfig:")
	log.Debugf("apikey = %s", conf.Telegram.APIKey)
	log.Debugf("chatid = %v", conf.Telegram.ChatID)
	log.Debugf("admin  = %s", conf.Telegram.Admin)
	log.Debugf("debug  = %v", conf.Telegram.Debug)

	return &conf, nil
}
