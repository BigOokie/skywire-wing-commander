# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- TOML based configuration. A valid `config.toml` file is required for the application to start. An example reference configuration file is provided (`config.toml.example`)
- Added a formal changelog (based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)), and document version history.
- Added `/licences` command to allow users to query the Bot for licence attributions.

### Changed
- Changed versioning. This project now adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html). Note that versions prior to this do not follow formal Semantic Versioning.
- Renamed application and project to Skywire-Wing-Commander (from skywire-telegram-notify-bot).
- Restructured repository layout (folders).

### Deprecated
### Removed
- Support for commande line parameters. All runtime configuration must now be sourced from `config.toml`
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

[Unreleased]: https://github.com/BigOokie/Skywire-Wing-Commander/compare/master...dev
[v0.0.3-alpha]: https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.0.2-alpha...v0.0.3-alpha
[v0.0.2-alpha]: https://github.com/BigOokie/Skywire-Wing-Commander/compare/v0.0.1-alpha...v0.0.2-alpha
[v0.0.1-alpha]: https://github.com/BigOokie/Skywire-Wing-Commander/commit/70153f0777a3d71bdc15bb4509c0b36ce45e096b



