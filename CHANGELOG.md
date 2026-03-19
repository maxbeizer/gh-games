# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com),
and this project adheres to [Semantic Versioning](https://semver.org).

## [0.5.0] - 2026-03-19

### Added
- Share results: clipboard copy + Slack posting via gh-slack
- `gh games config` command for interactive setup (`--show`, `--reset`)
- Slack team/workspace support in config

## [0.4.0] - 2026-03-19

### Added
- Share infrastructure (ShareResult type, clipboard, webhook support)
- Spoiler-free `Summary()` on all 9 game structs

## [0.3.0] - 2026-03-19

### Added
- ☠️ Hangman (`gh games hang`)
- 🔀 Jumble (`gh games jumble`)
- 🪜 Word Ladder (`gh games ladder`)
- 🧠 Trivia (`gh games trivia`)
- 🔐 Mastermind (`gh games code`)
- 📰 Mini Crossword (`gh games cross`)

## [0.2.0] - 2026-03-19

### Added
- 🐝 Hive (`gh games hive`) — word finding from 7 letters

### Fixed
- Hive startup hang (case mismatch in dictionary)
- Hive input rejection (uppercase/lowercase comparison)
- Expanded hive dictionary to 174k words

## [0.1.0] - 2026-03-19

### Added
- 🟩 Guess (`gh games guess`) — 5-letter word guessing
- 🔗 Group (`gh games group`) — find four groups of four words
- Daily and random modes for Guess
- Bubbletea TUI with Lipgloss styling
