# 🎮 gh-games

Terminal games as a GitHub CLI extension. Play right from your command line!

## Install

```sh
gh extension install maxbeizer/gh-games
```

## Games

| Command | Game |
|---------|------|
| `gh games guess` | 🟩 Guess the 5-letter word in 6 tries |
| `gh games group` | 🔗 Find four groups of four related words |
| `gh games hive` | 🐝 Find words from 7 letters, center required |
| `gh games hang` | ☠️ Classic hangman |
| `gh games jumble` | 🔀 Unscramble the jumbled word |
| `gh games ladder` | 🪜 Change one letter at a time to reach the target |
| `gh games trivia` | 🧠 10-question trivia quiz |
| `gh games code` | 🔐 Crack the secret color code |
| `gh games cross` | 📰 5×5 mini crossword |

## Sharing Results

After each game, share your spoiler-free results:

- **Clipboard** — press `C` to copy (always available)
- **Slack** — press `S` to post via [gh-slack](https://github.com/github/gh-slack) (requires setup)

### Slack setup

```sh
gh extension install github/gh-slack  # if not already installed
gh games config                        # interactive setup
```

## Development

```sh
make build         # Build binary
make test          # Run tests
make ci            # Build + vet + test with race detector
make install-local # Install from local checkout
```

## License

MIT · Not affiliated with the New York Times or any other company ⚖️
