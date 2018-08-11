package wcconfig

import (
	"fmt"
	"strings"
	"time"

	"github.com/BigOokie/skywire-wing-commander/src/utils"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

// Config structure models the TOML configuration file
type Config struct {
	WingCommander struct {
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
		IntervalSec     time.Duration
		HeartbeatIntMin time.Duration
	}
}

// String is the stringer function for the Config struct
func (c *Config) String() string {
	resultstr := "[WingCommander]\n" +
		"  twofactorenabled = %v\n" +
		"[SkyManager]\n" +
		"  address = \"%s\"\n" +
		"[Telegram]\n" +
		"  apikey = \"%s\"\n" +
		"  chatid = %v\n" +
		"  admin  = \"%s\"\n" +
		"  debug  = %v\n" +
		"[Monitor]\n" +
		"  intervalsec = %v\n" +
		"  heartbeatintmin = %v\n"

	return fmt.Sprintf(resultstr, c.WingCommander.TwoFactorEnabled, c.SkyManager.Address, c.Telegram.APIKey, c.Telegram.ChatID, c.Telegram.Admin, c.Telegram.Debug, c.Monitor.IntervalSec, c.Monitor.HeartbeatIntMin)
}

// ReadConfig will read configuration from the provided TOML config file
func ReadConfig(filename string) (*Config, error) {
	if !utils.FileExists(filename) {
		log.Error("Unable to find config file.")
	}

	var conf Config
	_, err := toml.DecodeFile(filename, &conf)
	if err != nil {
		return nil, err
	}

	log.Debugln("Reading config.")

	// Adjust time durations for interval configurations
	conf.Monitor.IntervalSec = time.Second * conf.Monitor.IntervalSec
	conf.Monitor.HeartbeatIntMin = time.Minute * conf.Monitor.HeartbeatIntMin

	// Check if the Admin user is prefixed with `@`
	if !strings.HasPrefix(conf.Telegram.Admin, "@") {
		// Add an "@" as the first character
		conf.Telegram.Admin = "@" + conf.Telegram.Admin
		log.Warnf("ReadConfig: admin username configuration is not prefixed `@`. Runtime config updated to prevent errors.")
	}

	DebugLogConfig(&conf)

	return &conf, nil
}

// DebugLogConfig will log debug information for the passed Config structure
func DebugLogConfig(conf *Config) {
	log.Debugf("Wing Commander Configuration:\n%s", conf.String())
}
