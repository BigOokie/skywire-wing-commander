# Wing Commander Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security

## [v0.2.0-beta.4] - 2018-08-19
### Added
- Added `/checkupdate` command. Requests the Bot to check GitHub to determine if a new version is available. The Bot will report back its findings. Note: This command does not perform the upgrade.
### Changed
### Deprecated
### Removed
### Fixed
### Security

## [v0.2.0-beta.3] - 2018-08-18
### Added
- Added monitoring of the Skywire Discovery Server to ensure local Nodes are connected. The connection of all locally connected Nodes are now checked as part of the routine Heartbeat cycle. The Heartbeat status message will report both the number of locally connected Nodes (Nodes connected to the local Manager) as well as the number of these Nodes that are currently registered with the Discovery Server. This feature is intended to help the Skyfleet community preserve their Node up-times and remain elegible for monthly Testnet rewards.
- Added error reporting notifications to report when the Bot is unable to connect to the local manager. This error notification will be reported each time the local Manager Node is polled (default 10sec).
### Changed
### Deprecated
### Removed
### Fixed
### Security

## [v0.2.0-beta.2] - 2018-08-12
### Added
- Added autorestart shell script (`wcstart.sh`). Thanks to @Cryptovinnie
- Added config file generator shell script (`wcbuildconfig.sh`). Thanks to @Cryptovinnie
- Added `/showconfig` command. This will ask the Bot to display its current runtime configuration, as stored in the `config.toml` file.
- Added `/checkupdate` command. This allows the Bot to check for new releases on GitHub and report back the status.
### Changed
- `/status` command will now respond differently depending on if the Monitor is currently runing or not (i.e. `/start`). When the Monitor is not running, the responce will indicate this. When the Monitor is running, the responce will indicate the current number of connected Nodes.
### Deprecated
### Removed
### Fixed
### Security

## [v0.2.0-beta.1] - 2018-08-10
### Added
- Code tests.
- Added list of Contributors [CONTRIBUTORS.md](CONTRIBUTORS.md)
- Added ability for Bot to compensate for missing `@` on `admin` user in `config.toml`. Previosuly is the admin user was not prefixed with an `@` the Bot would refuse to respond to messages from your user on Telegram. Now the Bot will compensate (the `config.toml` file is unchanged) by updating its runtime configuration to add the missing `@` if required. The Bot will log a warning to notify you.
### Changed
- Minor updates to messages generated from Wing Commander.
### Deprecated
### Removed
### Fixed
- Minor fixes to code and doco based on community feedback during ALPAH release stage.
### Security

## [v0.1.1-alpha.1] - 2018-08-03
### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
- Prevent the **Wing Commander** bot from responding to any user who is not configured as the Admin (in the `config.toml`). Messages or commander from non-admin users are simply dropped and ignored. Previously they could be processed - meaning any Telegram user could initiate a private chat with an instance of the **Wing Commander** bot that they do not own. This has now been resolved.  All users should upgrade to this version. v0.1.0-alpha.1 should not be used.

## [v0.1.0-alpha.1] - 2018-08-02
### Added
- TOML based configuration. A valid `config.toml` file is required in the users `$HOME/.wingcommander` folder. The `$HOME` folder may be different depending on your OS. An example reference configuration file is provided (`config.toml.example`).
- Added a formal changelog (based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)), and document version history.

### Changed
- Changed versioning. This project now adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html). Note that versions prior to this do not follow formal Semantic Versioning.
- Renamed application and project to Skywire-Wing-Commander (from skywire-telegram-notify-bot).
- Restructured repository layout (folders).
- Notifications for Node `connect` and `disconnect` events are sent via Telegram. Monitoring is initiated by the `/start` command, and terminated using the `/stop` command. Monitoring interval (seconds) is set within the `config.toml`. Each time a connect or disconnect notification is provided, the Bot will also provide a count of the number of Nodes currently connected to the Manager.
- The `/start` command now also initiates a heartbeat based on a configurable interval (`config.toml`).

### Deprecated
### Removed
- Support for commande line parameters. All runtime configuration must now be sourced from `config.toml`
- The `/heartbeat` command has been removed. A heartbeat is now managed by the `/start` command.
### Fixed
### Security

## [v0.0.3-alpha] - 2018-06-26
### Added
- Support for TravisCI
### Changed
- `Manager` connection status changes now provide a structured Telegram message. Previously a raw JSON dump was sent as the message.
### Deprecated
### Removed
### Fixed
- Pevent `Manager` connection monitor from running more than one instance.
### Security

## [v0.0.2-alpha] - 2018-06-24
### Added
- Support for Official Sky Miner (running the official OrangePi images). The Official Miner uses different file locations for `Manager` configuration files (when compared to DIY Miners). Specifically the app was monitoring for changes in the `clients.json` file.
- Telegram status notifications now include count of connections by type (and directions - incomming or outgoing).
### Changed
- Major refactor of code base to redesign approach.

### Deprecated
### Removed
### Fixed
### Security

## [v0.0.1-alpha] - 2018-06-23
### Added
- Support `bottoken` command line parameter. Allows the Telegram API Key to be passed to the application. The application will refuse to start if this is not provided. The API Key will be provided by the Telegram `@BotFather`.
- Support `botdebug` command line parameter. Allows debugging of the Telegram Bot API interface to be enabled (off by default).
- Support for `Manager` connection status monitoring (via `clients.json` file) on SkyCoin Miners running the non-official DIY platforms (such as Raspberry Pi). Official Miners are not supported at this stage.
- Bot will send raw JSON dump of the `clients.json` file from the SkyMiner `Manager` when ever this file changes.
- Added `CREDITS.md` to record licencing and attribution for licenced and derivative works.
- Added Icon: `SKY_RANK_WingCommander.png` (not currently used). Credit Noun Project for derivative work.
- Added Icon: `WingCommanderLogo.png`. Credit SkyCoin project for derivative work.

### Changed
### Deprecated
### Removed
### Fixed
### Security

[Unreleased]: 
https://github.com/BigOokie/Skywire-Wing-Commander/compare/master...dev
[v0.2.0-beta.4]: 
https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.2.0-beta.3...v0.2.0-beta.4
[v0.2.0-beta.3]: 
https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.2.0-beta.2...v0.2.0-beta.3
[v0.2.0-beta.2]: 
https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.2.0-beta.1...v0.2.0-beta.2
[v0.2.0-beta.1]: 
https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.1.1-alpha.1...v0.2.0-beta.1
[v0.1.1-alpha.1]:
https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.1.0-alpha.1...v0.1.1-alpha.1
[v0.1.0-alpha.1]:
https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.0.3-alpha...v0.1.0-alpha.1
[v0.0.3-alpha]:
https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.0.2-alpha...v0.0.3-alpha
[v0.0.2-alpha]:
https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.0.1-alpha...v0.0.2-alpha
[v0.0.1-alpha]:
https://github.com/BigOokie/Skywire-Wing-Commander/commit/70153f0777a3d71bdc15bb4509c0b36ce45e096b
