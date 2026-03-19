package group

import "math/rand"

// Difficulty represents how hard a category is to identify.
type Difficulty int

const (
	Easy   Difficulty = iota // Yellow
	Medium                   // Green
	Hard                     // Blue
	Expert                   // Purple
)

// Category is a group of 4 words that share a connection.
type Category struct {
	Name       string
	Words      []string // exactly 4 words
	Difficulty Difficulty
}

// AllCategories contains the full bank of available categories.
var AllCategories = []Category{
	// ─── Easy (Yellow) ───────────────────────────────────────
	{Name: "Primary Colors", Words: []string{"RED", "BLUE", "YELLOW", "GREEN"}, Difficulty: Easy},
	{Name: "Seasons", Words: []string{"SPRING", "SUMMER", "FALL", "WINTER"}, Difficulty: Easy},
	{Name: "Planets", Words: []string{"MARS", "VENUS", "SATURN", "JUPITER"}, Difficulty: Easy},
	{Name: "Fruits", Words: []string{"APPLE", "BANANA", "GRAPE", "MANGO"}, Difficulty: Easy},
	{Name: "Card Suits", Words: []string{"HEARTS", "CLUBS", "SPADES", "DIAMONDS"}, Difficulty: Easy},
	{Name: "Compass Directions", Words: []string{"NORTH", "SOUTH", "EAST", "WEST"}, Difficulty: Easy},
	{Name: "Ocean Animals", Words: []string{"SHARK", "WHALE", "DOLPHIN", "OCTOPUS"}, Difficulty: Easy},
	{Name: "Musical Instruments", Words: []string{"GUITAR", "PIANO", "DRUMS", "VIOLIN"}, Difficulty: Easy},
	{Name: "School Subjects", Words: []string{"MATH", "SCIENCE", "HISTORY", "ENGLISH"}, Difficulty: Easy},
	{Name: "Fast Food Chains", Words: []string{"WENDYS", "SUBWAY", "ARBYS", "POPEYES"}, Difficulty: Easy},
	{Name: "Dog Breeds", Words: []string{"POODLE", "BEAGLE", "BOXER", "COLLIE"}, Difficulty: Easy},
	{Name: "Vegetables", Words: []string{"CARROT", "CELERY", "ONION", "PEPPER"}, Difficulty: Easy},
	{Name: "Weather", Words: []string{"RAIN", "SNOW", "HAIL", "SLEET"}, Difficulty: Easy},
	{Name: "Coins", Words: []string{"PENNY", "NICKEL", "DIME", "QUARTER"}, Difficulty: Easy},
	{Name: "Breakfast Foods", Words: []string{"BACON", "WAFFLE", "PANCAKE", "CEREAL"}, Difficulty: Easy},
	{Name: "Body Parts", Words: []string{"ELBOW", "KNEE", "ANKLE", "WRIST"}, Difficulty: Easy},

	// ─── Medium (Green) ──────────────────────────────────────
	{Name: "Programming Languages", Words: []string{"PYTHON", "RUST", "SWIFT", "RUBY"}, Difficulty: Medium},
	{Name: "Greek Letters", Words: []string{"ALPHA", "BETA", "GAMMA", "DELTA"}, Difficulty: Medium},
	{Name: "Chess Pieces", Words: []string{"ROOK", "BISHOP", "KNIGHT", "QUEEN"}, Difficulty: Medium},
	{Name: "Pasta Shapes", Words: []string{"PENNE", "FUSILLI", "ORZO", "ZITI"}, Difficulty: Medium},
	{Name: "Dances", Words: []string{"SALSA", "TANGO", "WALTZ", "RUMBA"}, Difficulty: Medium},
	{Name: "Fabrics", Words: []string{"SILK", "DENIM", "COTTON", "VELVET"}, Difficulty: Medium},
	{Name: "Gemstones", Words: []string{"RUBY", "TOPAZ", "OPAL", "JADE"}, Difficulty: Medium},
	{Name: "Elements", Words: []string{"IRON", "GOLD", "NEON", "ARGON"}, Difficulty: Medium},
	{Name: "Types of Bread", Words: []string{"RYE", "PUMPER", "SOURDOUGH", "BRIOCHE"}, Difficulty: Medium},
	{Name: "Martial Arts", Words: []string{"KARATE", "JUDO", "AIKIDO", "KENDO"}, Difficulty: Medium},
	{Name: "Hat Types", Words: []string{"FEDORA", "BERET", "BEANIE", "STETSON"}, Difficulty: Medium},
	{Name: "Currencies", Words: []string{"DOLLAR", "EURO", "YEN", "POUND"}, Difficulty: Medium},
	{Name: "Constellations", Words: []string{"ORION", "LYRA", "DRACO", "GEMINI"}, Difficulty: Medium},
	{Name: "Spices", Words: []string{"CUMIN", "THYME", "SAGE", "CLOVE"}, Difficulty: Medium},
	{Name: "Sharks", Words: []string{"HAMMER", "BULL", "TIGER", "MAKO"}, Difficulty: Medium},
	{Name: "Mythical Creatures", Words: []string{"DRAGON", "PHOENIX", "GRIFFIN", "HYDRA"}, Difficulty: Medium},

	// ─── Hard (Blue) ─────────────────────────────────────────
	{Name: "Types of Clouds", Words: []string{"CIRRUS", "STRATUS", "CUMULUS", "NIMBUS"}, Difficulty: Hard},
	{Name: "Shakespeare Plays", Words: []string{"HAMLET", "OTHELLO", "MACBETH", "TEMPEST"}, Difficulty: Hard},
	{Name: "Poker Terms", Words: []string{"FLOP", "RIVER", "TURN", "BLIND"}, Difficulty: Hard},
	{Name: "Architectural Styles", Words: []string{"GOTHIC", "BAROQUE", "TUDOR", "DECO"}, Difficulty: Hard},
	{Name: "Jazz Musicians", Words: []string{"MONK", "MILES", "MINGUS", "BIRD"}, Difficulty: Hard},
	{Name: "Rhetorical Devices", Words: []string{"IRONY", "SIMILE", "LITOTES", "ZEUGMA"}, Difficulty: Hard},
	{Name: "Types of Whiskey", Words: []string{"BOURBON", "SCOTCH", "RYE", "MALT"}, Difficulty: Hard},
	{Name: "Fencing Terms", Words: []string{"PARRY", "RIPOSTE", "LUNGE", "FLECHE"}, Difficulty: Hard},
	{Name: "Knot Types", Words: []string{"BOWLINE", "CLOVE", "REEF", "HITCH"}, Difficulty: Hard},
	{Name: "Philosophy Schools", Words: []string{"STOIC", "CYNIC", "SKEPTIC", "EPICURE"}, Difficulty: Hard},
	{Name: "Geologic Eras", Words: []string{"JURASSIC", "TRIASSIC", "PERMIAN", "CAMBRIAN"}, Difficulty: Hard},
	{Name: "Cheese Types", Words: []string{"BRIE", "GOUDA", "HAVARTI", "FONTINA"}, Difficulty: Hard},
	{Name: "Film Noir Terms", Words: []string{"FEMME", "HEIST", "ALIBI", "SLEUTH"}, Difficulty: Hard},
	{Name: "Sailing Terms", Words: []string{"TACK", "JIB", "HALYARD", "CLEAT"}, Difficulty: Hard},
	{Name: "Opera Terms", Words: []string{"ARIA", "LIBRETTO", "SOPRANO", "OVERTURE"}, Difficulty: Hard},
	{Name: "Woodworking Joints", Words: []string{"DOVETAIL", "MORTISE", "TENON", "RABBET"}, Difficulty: Hard},

	// ─── Expert (Purple) ─────────────────────────────────────
	{Name: "___ Run", Words: []string{"HOME", "DRY", "BULL", "TRIAL"}, Difficulty: Expert},
	{Name: "___ Board", Words: []string{"CARD", "DART", "SKATE", "CHALK"}, Difficulty: Expert},
	{Name: "Fire ___", Words: []string{"PLACE", "TRUCK", "WORKS", "FLY"}, Difficulty: Expert},
	{Name: "___ Light", Words: []string{"FLASH", "MOON", "SPOT", "HIGH"}, Difficulty: Expert},
	{Name: "Words in NATO Alphabet", Words: []string{"TANGO", "FOXTROT", "LIMA", "OSCAR"}, Difficulty: Expert},
	{Name: "___ Stone", Words: []string{"LIME", "KEY", "MILE", "GRAVE"}, Difficulty: Expert},
	{Name: "Anagram of a Color", Words: []string{"LUBE", "ANGER", "REIGN", "DALE"}, Difficulty: Expert},
	{Name: "Double Letters", Words: []string{"BALLOON", "BUFFALO", "MILLION", "RACCOON"}, Difficulty: Expert},
	{Name: "___ Fish", Words: []string{"SWORD", "BLOW", "CAT", "ANGEL"}, Difficulty: Expert},
	{Name: "Water ___", Words: []string{"FALL", "MARK", "PROOF", "MELON"}, Difficulty: Expert},
	{Name: "Also a Car Brand", Words: []string{"JAGUAR", "ACCORD", "MUSTANG", "ECLIPSE"}, Difficulty: Expert},
	{Name: "___ Back", Words: []string{"SET", "DRAW", "PAPER", "QUARTER"}, Difficulty: Expert},
	{Name: "Head ___", Words: []string{"BAND", "LINE", "PHONE", "HUNTER"}, Difficulty: Expert},
	{Name: "Black ___", Words: []string{"BERRY", "SMITH", "BIRD", "MAIL"}, Difficulty: Expert},
	{Name: "Hidden Body Part", Words: []string{"SHINGLE", "ARMED", "ALCHEMY", "ANTHEM"}, Difficulty: Expert},
	{Name: "___ House", Words: []string{"WARE", "GREEN", "POWER", "FIRE"}, Difficulty: Expert},
}

// GeneratePuzzle picks 4 categories (one per difficulty) randomly and returns
// them along with a shuffled flat list of all 16 words.
func GeneratePuzzle() (categories [4]Category, shuffledWords []string) {
	byDifficulty := map[Difficulty][]Category{}
	for _, c := range AllCategories {
		byDifficulty[c.Difficulty] = append(byDifficulty[c.Difficulty], c)
	}

	for d := Easy; d <= Expert; d++ {
		pool := byDifficulty[d]
		categories[d] = pool[rand.Intn(len(pool))]
	}

	// Collect all 16 words and verify no overlap
	seen := map[string]bool{}
	for i := range categories {
		for _, w := range categories[i].Words {
			if seen[w] {
				// Collision detected — retry with a fresh pick
				return GeneratePuzzle()
			}
			seen[w] = true
		}
	}

	shuffledWords = make([]string, 0, 16)
	for _, c := range categories {
		shuffledWords = append(shuffledWords, c.Words...)
	}
	rand.Shuffle(len(shuffledWords), func(i, j int) {
		shuffledWords[i], shuffledWords[j] = shuffledWords[j], shuffledWords[i]
	})

	return categories, shuffledWords
}
