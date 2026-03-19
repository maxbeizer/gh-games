package cross

// Cell represents a single cell in the crossword grid.
type Cell struct {
	Letter rune // the answer letter, 0 for black cells
	Black  bool
}

// Clue represents a single crossword clue.
type Clue struct {
	Number    int
	Direction string // "across" or "down"
	Text      string
}

// Puzzle represents a complete mini crossword puzzle.
type Puzzle struct {
	Grid  [5][5]Cell
	Clues []Clue
}

// L is a shorthand to create a letter cell.
func L(r rune) Cell { return Cell{Letter: r} }

// B creates a black (blocked) cell.
var B = Cell{Black: true}

// Puzzles is the collection of pre-built mini crossword puzzles.
var Puzzles = []Puzzle{
	// Puzzle 0
	// S W A M P
	// H O R S E
	// A ■ E ■ A
	// R I S K S
	// P E A R L
	{
		Grid: [5][5]Cell{
			{L('S'), L('W'), L('A'), L('M'), L('P')},
			{L('H'), L('O'), L('R'), L('S'), L('E')},
			{L('A'), B, L('E'), B, L('A')},
			{L('R'), L('I'), L('S'), L('K'), L('S')},
			{L('P'), L('E'), L('A'), L('R'), L('L')},
		},
		Clues: []Clue{
			{1, "across", "Boggy wetland"},
			{6, "across", "Equine animal"},
			{7, "across", "Zone or region"},
			{8, "across", "Danger or hazard"},
			{10, "across", "Gem from an oyster"},
			{1, "down", "Pointed, like a blade"},
			{2, "down", "Cable alternative, once"},
			{3, "down", "Regions or zones"},
			{4, "down", "Mrs. or Mr."},
			{5, "down", "Request earnestly"},
		},
	},
	// Puzzle 1
	// B R A V E
	// L O N E R
	// A ■ C ■ A
	// S T E E L
	// T R E A T
	{
		Grid: [5][5]Cell{
			{L('B'), L('R'), L('A'), L('V'), L('E')},
			{L('L'), L('O'), L('N'), L('E'), L('R')},
			{L('A'), B, L('C'), B, L('A')},
			{L('S'), L('T'), L('E'), L('E'), L('L')},
			{L('T'), L('R'), L('E'), L('A'), L('T')},
		},
		Clues: []Clue{
			{1, "across", "Courageous"},
			{6, "across", "Solitary person"},
			{7, "across", "Happen just the once"},
			{8, "across", "Iron alloy"},
			{10, "across", "Halloween ___"},
			{1, "down", "Explosion"},
			{2, "down", "Umpire's call after three strikes"},
			{3, "down", "Once more"},
			{4, "down", "Flat bread from India"},
			{5, "down", "Genuine"},
		},
	},
	// Puzzle 2
	// C L A S P
	// R A N G E
	// A ■ T ■ N
	// F L I E S
	// T R A I N
	{
		Grid: [5][5]Cell{
			{L('C'), L('L'), L('A'), L('N'), L('P')},
			{L('R'), L('A'), L('N'), L('G'), L('E')},
			{L('A'), B, L('T'), B, L('N')},
			{L('F'), L('L'), L('I'), L('E'), L('S')},
			{L('T'), L('R'), L('A'), L('I'), L('N')},
		},
		Clues: []Clue{
			{1, "across", "Jewelry fastener"},
			{6, "across", "Mountain ___"},
			{7, "across", "Top-of-the-line, rhymes with tent"},
			{8, "across", "Buzzes around food"},
			{10, "across", "Locomotive"},
			{1, "down", "Skillful work"},
			{2, "down", "Comparable to"},
			{3, "down", "Satin or silk"},
			{4, "down", "Gene variant"},
			{5, "down", "Writing instrument"},
		},
	},
	// Puzzle 3
	// S T A R S
	// T I L E D
	// O ■ L ■ G
	// R A I S E
	// E N D E D
	{
		Grid: [5][5]Cell{
			{L('S'), L('T'), L('A'), L('R'), L('S')},
			{L('T'), L('I'), L('L'), L('E'), L('D')},
			{L('O'), B, L('L'), B, L('G')},
			{L('R'), L('A'), L('I'), L('S'), L('E')},
			{L('E'), L('N'), L('D'), L('E'), L('D')},
		},
		Clues: []Clue{
			{1, "across", "Celestial twinklers"},
			{6, "across", "Like a mosaic floor"},
			{7, "across", "Not feeling well"},
			{8, "across", "Increase, as wages"},
			{10, "across", "Finished, concluded"},
			{1, "down", "Shop or boutique"},
			{2, "down", "Like some walls"},
			{3, "down", "Partner to \"void\""},
			{4, "down", "Go up"},
			{5, "down", "Sword's sharp side"},
		},
	},
	// Puzzle 4
	// G R A S P
	// L I N E R
	// A ■ K ■ I
	// S P O R T
	// S T O N E
	{
		Grid: [5][5]Cell{
			{L('G'), L('R'), L('A'), L('S'), L('P')},
			{L('L'), L('I'), L('N'), L('E'), L('R')},
			{L('A'), B, L('K'), B, L('I')},
			{L('S'), L('P'), L('O'), L('R'), L('T')},
			{L('S'), L('T'), L('O'), L('N'), L('E')},
		},
		Clues: []Clue{
			{1, "across", "Understand, get a grip on"},
			{6, "across", "Cruise ship"},
			{7, "across", "Twist or bend in a hose"},
			{8, "across", "Athletic activity"},
			{10, "across", "Rock, as in Rolling ___s"},
			{1, "down", "Eyewear material"},
			{2, "down", "Lean, as a tower"},
			{3, "down", "Call off, cancel"},
			{4, "down", "Painful, tender"},
			{5, "down", "Baked dessert"},
		},
	},
	// Puzzle 5
	// F L A M E
	// L O C A L
	// A ■ K ■ I
	// T R E N D
	// S T E E P
	{
		Grid: [5][5]Cell{
			{L('F'), L('L'), L('A'), L('M'), L('E')},
			{L('L'), L('O'), L('C'), L('A'), L('L')},
			{L('A'), B, L('K'), B, L('I')},
			{L('T'), L('R'), L('E'), L('N'), L('D')},
			{L('S'), L('T'), L('E'), L('E'), L('P')},
		},
		Clues: []Clue{
			{1, "across", "Fire's glow"},
			{6, "across", "Neighborhood bar"},
			{7, "across", "Flaw or defect"},
			{8, "across", "Fashion direction"},
			{10, "across", "Very expensive, pricewise"},
			{1, "down", "Apartment or condo"},
			{2, "down", "Lucky charm, like a four-leaf clover"},
			{3, "down", "Things to acknowledge"},
			{4, "down", "Tidy up"},
			{5, "down", "Journey, as on a cruise"},
		},
	},
	// Puzzle 6
	// P L A N T
	// R I V E R
	// I ■ S ■ A
	// D O N O R
	// E L B O W
	{
		Grid: [5][5]Cell{
			{L('P'), L('L'), L('A'), L('N'), L('T')},
			{L('R'), L('I'), L('V'), L('E'), L('R')},
			{L('I'), B, L('S'), B, L('A')},
			{L('D'), L('O'), L('N'), L('O'), L('R')},
			{L('E'), L('L'), L('B'), L('O'), L('W')},
		},
		Clues: []Clue{
			{1, "across", "Fern or ficus"},
			{6, "across", "Flowing waterway"},
			{7, "across", "Compared to"},
			{8, "across", "Blood bank contributor"},
			{10, "across", "Joint below the shoulder"},
			{1, "down", "Self-satisfaction"},
			{2, "down", "Petroleum product"},
			{3, "down", "Notice or catch sight of"},
			{4, "down", "Not one person"},
			{5, "down", "Uncooked"},
		},
	},
	// Puzzle 7
	// D R A F T
	// R O U S E
	// I ■ N ■ A
	// V E N U E
	// E D G E S
	{
		Grid: [5][5]Cell{
			{L('D'), L('R'), L('A'), L('F'), L('T')},
			{L('R'), L('O'), L('U'), L('S'), L('E')},
			{L('I'), B, L('N'), B, L('A')},
			{L('V'), L('E'), L('N'), L('U'), L('E')},
			{L('E'), L('D'), L('G'), L('E'), L('S')},
		},
		Clues: []Clue{
			{1, "across", "First version of a paper"},
			{6, "across", "Wake from sleep"},
			{7, "across", "Sister or brother"},
			{8, "across", "Concert location"},
			{10, "across", "Borders or boundaries"},
			{1, "down", "Steer, like a car"},
			{2, "down", "Had a debt"},
			{3, "down", "Firearm"},
			{4, "down", "Blend together"},
			{5, "down", "Ocean, poetically"},
		},
	},
	// Puzzle 8
	// C H E S T
	// L O V E R
	// A ■ A ■ I
	// S K A T E
	// P A N E L
	{
		Grid: [5][5]Cell{
			{L('C'), L('H'), L('E'), L('S'), L('T')},
			{L('L'), L('O'), L('V'), L('E'), L('R')},
			{L('A'), B, L('A'), B, L('I')},
			{L('S'), L('K'), L('A'), L('T'), L('E')},
			{L('P'), L('A'), L('N'), L('E'), L('L')},
		},
		Clues: []Clue{
			{1, "across", "Treasure ___"},
			{6, "across", "Romeo, e.g."},
			{7, "across", "Large bird that doesn't fly (abbr. sound)"},
			{8, "across", "Ice rink activity"},
			{10, "across", "Discussion group"},
			{1, "down", "Hand slap, casually"},
			{2, "down", "Tall Japanese tree"},
			{3, "down", "Lava rock island"},
			{4, "down", "Fixed, as a game"},
			{5, "down", "Path or road"},
		},
	},
	// Puzzle 9
	// S H A R K
	// T O N I C
	// O ■ G ■ A
	// R E A D S
	// E N D O W
	{
		Grid: [5][5]Cell{
			{L('S'), L('H'), L('A'), L('R'), L('K')},
			{L('T'), L('O'), L('N'), L('I'), L('C')},
			{L('O'), B, L('G'), B, L('A')},
			{L('R'), L('E'), L('A'), L('D'), L('S')},
			{L('E'), L('N'), L('D'), L('O'), L('W')},
		},
		Clues: []Clue{
			{1, "across", "Ocean predator"},
			{6, "across", "Gin and ___"},
			{7, "across", "Small amount, like a smidge"},
			{8, "across", "Enjoys a book"},
			{10, "across", "Fund, as a chair at a university"},
			{1, "down", "Stock up on"},
			{2, "down", "Chick's sound"},
			{3, "down", "Plan or schedule"},
			{4, "down", "Wireless Internet"},
			{5, "down", "What a bird does"},
		},
	},
	// Puzzle 10
	// M O U S E
	// A G I L E
	// P ■ R ■ A
	// L E A R N
	// E D G E S
	{
		Grid: [5][5]Cell{
			{L('M'), L('O'), L('U'), L('S'), L('E')},
			{L('A'), L('G'), L('I'), L('L'), L('E')},
			{L('P'), B, L('R'), B, L('A')},
			{L('L'), L('E'), L('A'), L('R'), L('N')},
			{L('E'), L('D'), L('G'), L('E'), L('S')},
		},
		Clues: []Clue{
			{1, "across", "Computer pointing device"},
			{6, "across", "Nimble and quick"},
			{7, "across", "Two of a kind"},
			{8, "across", "Acquire knowledge"},
			{10, "across", "Cliff sides"},
			{1, "down", "Tree syrup source"},
			{2, "down", "Vintage, like wine"},
			{3, "down", "Hitchhiker's Guide prefix"},
			{4, "down", "Moves like a snake"},
			{5, "down", "More than enough"},
		},
	},
	// Puzzle 11
	// T A B L E
	// R A I S E
	// A ■ L ■ V
	// C O L O R
	// E N D E D
	{
		Grid: [5][5]Cell{
			{L('T'), L('A'), L('B'), L('L'), L('E')},
			{L('R'), L('A'), L('I'), L('S'), L('E')},
			{L('A'), B, L('L'), B, L('V')},
			{L('C'), L('O'), L('L'), L('O'), L('R')},
			{L('E'), L('N'), L('D'), L('E'), L('D')},
		},
		Clues: []Clue{
			{1, "across", "Furniture for dining"},
			{6, "across", "Lift up or increase"},
			{7, "across", "Sick, under the weather"},
			{8, "across", "Red, blue, or green"},
			{10, "across", "All finished"},
			{1, "down", "Follow a path"},
			{2, "down", "Starting from"},
			{3, "down", "Invoiced"},
			{4, "down", "By oneself"},
			{5, "down", "In the evening"},
		},
	},
	// Puzzle 12
	// Q U I L T
	// U N I O N
	// A ■ C ■ I
	// K N E E L
	// E L D E R
	{
		Grid: [5][5]Cell{
			{L('Q'), L('U'), L('I'), L('L'), L('T')},
			{L('U'), L('N'), L('I'), L('O'), L('N')},
			{L('A'), B, L('C'), B, L('I')},
			{L('K'), L('N'), L('E'), L('E'), L('L')},
			{L('E'), L('D'), L('G'), L('E'), L('R')},
		},
		Clues: []Clue{
			{1, "across", "Bed covering, often handmade"},
			{6, "across", "Labor group"},
			{7, "across", "Frozen water, in London"},
			{8, "across", "Genuflect"},
			{10, "across", "Accounting book"},
			{1, "down", "Shiver, as in fear"},
			{2, "down", "Bare, unclothed"},
			{3, "down", "Chunk of frozen water"},
			{4, "down", "Poem of praise"},
			{5, "down", "Grading system"},
		},
	},
	// Puzzle 13
	// C R A N E
	// O U T E R
	// B ■ A ■ A
	// R I D E R
	// A L E R T
	{
		Grid: [5][5]Cell{
			{L('C'), L('R'), L('A'), L('N'), L('E')},
			{L('O'), L('U'), L('T'), L('E'), L('R')},
			{L('B'), B, L('A'), B, L('A')},
			{L('R'), L('I'), L('D'), L('E'), L('R')},
			{L('A'), L('L'), L('E'), L('R'), L('T')},
		},
		Clues: []Clue{
			{1, "across", "Construction machine or wading bird"},
			{6, "across", "External, not inner"},
			{7, "across", "Greek letter, a to z"},
			{8, "across", "Horseback person"},
			{10, "across", "Vigilant, on guard"},
			{1, "down", "King cobra, for one"},
			{2, "down", "Governing body"},
			{3, "down", "Go with the flow"},
			{4, "down", "Get closer"},
			{5, "down", "Space or region"},
		},
	},
	// Puzzle 14
	// W H A L E
	// R I P E N
	// A ■ E ■ D
	// P L E A S
	// S E N S E
	{
		Grid: [5][5]Cell{
			{L('W'), L('H'), L('A'), L('L'), L('E')},
			{L('R'), L('I'), L('P'), L('E'), L('N')},
			{L('A'), B, L('E'), B, L('D')},
			{L('P'), L('L'), L('E'), L('A'), L('S')},
			{L('S'), L('E'), L('N'), L('S'), L('E')},
		},
		Clues: []Clue{
			{1, "across", "Enormous sea mammal"},
			{6, "across", "Become ready to eat, as fruit"},
			{7, "across", "Quick look"},
			{8, "across", "Court requests"},
			{10, "across", "Common ___ (logic)"},
			{1, "down", "Gift ___"},
			{2, "down", "Greek letter, a to o"},
			{3, "down", "Gap or crack"},
			{4, "down", "Rent out"},
			{5, "down", "Conclusion"},
		},
	},
}
