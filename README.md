# 🎮 gh-games

Terminal games as a GitHub CLI extension. Play right from your command line!

## Install

```sh
gh extension install maxbeizer/gh-games
```

## Games

### 🟩 Guess

Guess the hidden 5-letter word in 6 tries. After each guess, letters are colored:
- 🟩 **Green** — correct letter, correct position
- 🟨 **Yellow** — correct letter, wrong position  
- ⬜ **Gray** — letter not in the word

```sh
gh games guess           # Daily word (same for everyone today)
gh games guess --random  # Fresh random word
gh games guess --hard    # Guesses must be real words
```

### 🔗 Group

Find four groups of four related words among sixteen. Groups are color-coded by difficulty:
- 🟨 Yellow (easy) → 🟩 Green → 🟦 Blue → 🟪 Purple (expert)

4 mistakes and it's game over!

```sh
gh games group
```

## Development

```sh
make build         # Build binary
make test          # Run tests
make ci            # Build + vet + test with race detector
make install-local # Install from local checkout
```

## License

MIT

This project is not affiliated with the New York Times or any other company ⚖️
