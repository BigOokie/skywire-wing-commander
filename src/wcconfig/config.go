// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package wcconfig

import (
	"fmt"
	"strings"
	"time"

	viper "github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
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
	Address string `mapstructure:"address"`
}

// MonitorParameters struct defines the configuration parameters that
// are used by the Monitor which polls the SkyManager
type MonitorParameters struct {
	IntervalSec     time.Duration `mapstructure:"intervalsec"`
	HeartbeatIntMin time.Duration `mapstructure:"heartbeatintmin"`
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

// DebugLogConfig will log debug information for the passed Config structure
func (c *Config) DebugLogConfig() {
	log.Debugf("Wing Commander Configuration:\n%s", c.String())
}

// LoadConfigParameters loads config parameters from the provided config file.
// The default file is config.toml in the .wingcommander folder in the users Home
// JSON, toml or yaml file can be used (toml preferred and default).
func LoadConfigParameters(configName string) (Config, error) {
	confNameSplit := strings.Split(configName, ".")
	fileType := confNameSplit[len(confNameSplit)-1]
	switch fileType {
	case "toml", "json", "yaml", "yml":
		viper.SetConfigType(confNameSplit[len(confNameSplit)-1])
	default:
		return Config{}, fmt.Errorf("invalid blockchain config file type: %s", fileType)
	}

	configName = configName[:len(configName)-(len(fileType)+1)]
	viper.SetConfigName(configName)
	viper.AddConfigPath("$HOME/.wingcommander/")

	// set app config defaults
	setConfigDefaults()

	config := Config{}

	if err := viper.ReadInConfig(); err != nil {
		log.Debug("ReadInConfig")
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Debug("Unmarshal")
		return config, err
	}

	config.Monitor.IntervalSec = config.Monitor.IntervalSec * time.Second
	config.Monitor.HeartbeatIntMin = config.Monitor.HeartbeatIntMin * time.Minute

	// Check if the Admin user is prefixed with `@`
	if !strings.HasPrefix(config.Telegram.Admin, "@") {
		// Add an "@" as the first character
		config.Telegram.Admin = "@" + config.Telegram.Admin
		log.Warnf("ReadConfig: admin username configuration is not prefixed `@`. Runtime config updated to prevent errors.")
	}

	config.DebugLogConfig()
	return config, nil
}

func setConfigDefaults() {

	// Telegram defaults
	viper.SetDefault("telegram.debug", false)

	// Monitor defaults
	viper.SetDefault("monitor.intervalsec", 10)
	viper.SetDefault("monitor.heartbeatintmin", 120)

	// Skyminer defaults
	viper.SetDefault("skymanager.address", "127.0.0.1:8000")
}
