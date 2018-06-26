package main

// BotConfig struct is used to store run-time configuration
// information for the bot application.
type BotConfig struct {
	BotToken      string `json:"bot_token"`
	ChatID        int64  `json:"chat_id"`
	BotDebug      bool   `json:"botdebug"`
	ClientMonitor *ClientMonitorConfig
}

// ClientMonitorConfig struct manages configuration parameters
// for the Client (file) Monitor
type ClientMonitorConfig struct {
	ClientFile     string `json:"clientfile"`
	MonitorRunning bool   `json:"monitorrunning"`
}
