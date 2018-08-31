// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package wcconfig

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	viper "github.com/spf13/viper"
)

// Config structure models the applications configuration structure
type Config struct {
	WingCommander WingCommanderParameters `mapstructure:"wingcommander"`
	Telegram      TelegramParameters      `mapstructure:"telegram"`
	Monitor       MonitorParameters       `mapstructure:"monitor"`
	SkyManager    SkyManagerParameters    `mapstructure:"skymanager"`
}

// WingCommanderParameters struct defines the configuration parameters that
// are used to manage runtime config for the Wing Commander application
type WingCommanderParameters struct {
	TwoFactorEnabled bool `mapstructure:"twofactorenabled"`
	UpgradeComplete  bool
}

// TelegramParameters struct defines the configuration parameters that
// are used to manage Wing Commander application integrationw it Telegram
type TelegramParameters struct {
	APIKey string `mapstructure:"apikey"`
	ChatID int64  `mapstructure:"chatid"`
	Admin  string `mapstructure:"admin"`
	Debug  bool   `mapstructure:"debug"`
}

// SkyManagerParameters struct defines the configuration parameters that
// are used to manage connectivity with the Skywire Manager
type SkyManagerParameters struct {
	Address          string `mapstructure:"address"`
	DiscoveryAddress string `mapstructure:"discoveryaddress"`
}

// MonitorParameters struct defines the configuration parameters that
// are used by the Monitor which polls the SkyManager
type MonitorParameters struct {
	IntervalSec            time.Duration `mapstructure:"intervalsec"`
	HeartbeatIntMin        time.Duration `mapstructure:"heartbeatintmin"`
	DiscoveryMonitorIntMin time.Duration `mapstructure:"discoverymonitorintmin"`
}

// String is the stringer function for the Config struct
func (c *Config) String() string {
	resultstr := "[WingCommander]\n" +
		"  twofactorenabled = %v\n" +
		"[SkyManager]\n" +
		"  address = %q\n" +
		"  discoveryaddress = %q\n" +
		"[Telegram]\n" +
		"  apikey = %q\n" +
		"  chatid = %v\n" +
		"  admin  = %q\n" +
		"  debug  = %v\n" +
		"[Monitor]\n" +
		"  intervalsec = %v\n" +
		"  heartbeatintmin = %v\n" +
		"  discoverymonitorintmin = %v\n"

	return fmt.Sprintf(resultstr, c.WingCommander.TwoFactorEnabled, c.SkyManager.Address,
		c.SkyManager.DiscoveryAddress, c.Telegram.APIKey, c.Telegram.ChatID,
		c.Telegram.Admin, c.Telegram.Debug, c.Monitor.IntervalSec, c.Monitor.HeartbeatIntMin,
		c.Monitor.DiscoveryMonitorIntMin)
}

// PrintConfig will log debug information for the passed Config structure
func (c *Config) PrintConfig() {
	log.Printf("Wing Commander Configuration:\n%s", c.String())
}

// IsEmpty will compare the current instance of Config against an empty instance
// and return the result of the comparison
func IsEmpty(c Config) bool {
	return (c == (Config{}))
}

// readConfig attempts to read configuration parameters from the provided
// file (`filename`) and utilises the provided set of default values
// If successful it will return a *viper.Viper struct
func readConfig(filename, pathname string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(pathname)
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}

// LoadConfigParameters will load the applications configuration from the
// specified configuration file `filename` (note file extension must not be provided) in the
// specified path `pathname`. The function also provides the ability to specify
// configuration defaults.
// An `error` will be returned if any errors occur.
// A valid `Config` struct will be returned on success.
func LoadConfigParameters(filename, pathname string, defaults map[string]interface{}) (config Config, err error) {
	v1, err := readConfig(filename, pathname, defaults)

	if err != nil {
		return config, err
	}

	if err := v1.Unmarshal(&config); err != nil {
		return config, err
	}

	// Validate and adjust any configuration parameters
	config.Monitor.IntervalSec = config.Monitor.IntervalSec * time.Second
	config.Monitor.HeartbeatIntMin = config.Monitor.HeartbeatIntMin * time.Minute
	config.Monitor.DiscoveryMonitorIntMin = config.Monitor.DiscoveryMonitorIntMin * time.Minute

	// Check if the Admin user is prefixed with `@`
	if !strings.HasPrefix(config.Telegram.Admin, "@") {
		// Add an "@" as the first character
		config.Telegram.Admin = "@" + config.Telegram.Admin
		log.Warnf("ReadConfig: admin username configuration is not prefixed `@`. Runtime config updated to prevent errors.")
	}

	//config.PrintLogConfig()
	return config, nil
}
