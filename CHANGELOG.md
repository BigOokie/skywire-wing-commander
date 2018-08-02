# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

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

[v0.1.0-alpha.1]:
https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.0.3-alpha...v0.1.0-alpha.1
[v0.0.3-alpha]: https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.0.2-alpha...v0.0.3-alpha
[v0.0.2-alpha]: https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.0.1-alpha...v0.0.2-alpha
[v0.0.1-alpha]: https://github.com/BigOokie/Skywire-Wing-Commander/commit/70153f0777a3d71bdc15bb4509c0b36ce45e096b



