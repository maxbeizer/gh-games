package trivia

// Question represents a single trivia question with four choices.
type Question struct {
	Text     string
	Choices  [4]string
	Answer   int // 0-3 index into Choices
	Category string
}

// AllQuestions is the full bank of trivia questions.
var AllQuestions = []Question{
	// ─── Programming ──────────────────────────────────────────
	{
		Text:     "What does CSS stand for?",
		Choices:  [4]string{"Computer Style Sheets", "Cascading Style Sheets", "Creative Style System", "Colorful Style Sheets"},
		Answer:   1,
		Category: "Programming",
	},
	{
		Text:     "Which language was created by Guido van Rossum?",
		Choices:  [4]string{"Ruby", "Java", "Python", "Perl"},
		Answer:   2,
		Category: "Programming",
	},
	{
		Text:     "What year was Go first released?",
		Choices:  [4]string{"2007", "2009", "2011", "2012"},
		Answer:   1,
		Category: "Programming",
	},
	{
		Text:     "What does HTML stand for?",
		Choices:  [4]string{"Hyper Text Markup Language", "High Tech Modern Language", "Hyper Transfer Markup Language", "Home Tool Markup Language"},
		Answer:   0,
		Category: "Programming",
	},
	{
		Text:     "Which company developed the Java programming language?",
		Choices:  [4]string{"Microsoft", "Apple", "Sun Microsystems", "IBM"},
		Answer:   2,
		Category: "Programming",
	},
	{
		Text:     "What symbol is used for single-line comments in Python?",
		Choices:  [4]string{"//", "#", "--", "/*"},
		Answer:   1,
		Category: "Programming",
	},
	{
		Text:     "What does SQL stand for?",
		Choices:  [4]string{"Structured Query Language", "Simple Query Logic", "Standard Question Language", "Sequential Query Language"},
		Answer:   0,
		Category: "Programming",
	},
	{
		Text:     "Which of these is NOT a JavaScript framework?",
		Choices:  [4]string{"React", "Angular", "Django", "Vue"},
		Answer:   2,
		Category: "Programming",
	},
	{
		Text:     "What is the output of 'print(type([]))' in Python?",
		Choices:  [4]string{"<class 'tuple'>", "<class 'list'>", "<class 'dict'>", "<class 'set'>"},
		Answer:   1,
		Category: "Programming",
	},
	{
		Text:     "Which language is known as the 'language of the web'?",
		Choices:  [4]string{"Python", "Java", "JavaScript", "C++"},
		Answer:   2,
		Category: "Programming",
	},
	{
		Text:     "What does 'defer' do in Go?",
		Choices:  [4]string{"Delays execution until function returns", "Creates a goroutine", "Handles errors", "Imports a package"},
		Answer:   0,
		Category: "Programming",
	},
	{
		Text:     "What is the mascot of the Go programming language?",
		Choices:  [4]string{"A snake", "A gopher", "A ferret", "A penguin"},
		Answer:   1,
		Category: "Programming",
	},
	{
		Text:     "Which keyword declares a constant in JavaScript?",
		Choices:  [4]string{"var", "let", "const", "static"},
		Answer:   2,
		Category: "Programming",
	},
	{
		Text:     "In Git, what does 'HEAD' refer to?",
		Choices:  [4]string{"The first commit", "The current branch tip", "The remote repository", "The staging area"},
		Answer:   1,
		Category: "Programming",
	},
	{
		Text:     "What is 'nil' in Go?",
		Choices:  [4]string{"An integer zero", "A boolean false", "The zero value for pointers, interfaces, maps, slices, channels, and functions", "An empty string"},
		Answer:   2,
		Category: "Programming",
	},
	{
		Text:     "Which sorting algorithm has the best average-case time complexity?",
		Choices:  [4]string{"Bubble Sort — O(n²)", "Merge Sort — O(n log n)", "Selection Sort — O(n²)", "Insertion Sort — O(n²)"},
		Answer:   1,
		Category: "Programming",
	},

	// ─── Tech ─────────────────────────────────────────────────
	{
		Text:     "What company created the iPhone?",
		Choices:  [4]string{"Google", "Samsung", "Apple", "Microsoft"},
		Answer:   2,
		Category: "Tech",
	},
	{
		Text:     "What does API stand for?",
		Choices:  [4]string{"Application Programming Interface", "Applied Program Interaction", "Automated Protocol Integration", "Application Process Integration"},
		Answer:   0,
		Category: "Tech",
	},
	{
		Text:     "What is the default port for HTTPS?",
		Choices:  [4]string{"80", "8080", "443", "22"},
		Answer:   2,
		Category: "Tech",
	},
	{
		Text:     "What does CPU stand for?",
		Choices:  [4]string{"Central Processing Unit", "Computer Personal Unit", "Central Program Utility", "Core Processing Unit"},
		Answer:   0,
		Category: "Tech",
	},
	{
		Text:     "Which company owns GitHub?",
		Choices:  [4]string{"Google", "Amazon", "Microsoft", "Meta"},
		Answer:   2,
		Category: "Tech",
	},
	{
		Text:     "What does DNS stand for?",
		Choices:  [4]string{"Digital Network Service", "Domain Name System", "Data Node Server", "Dynamic Name Service"},
		Answer:   1,
		Category: "Tech",
	},
	{
		Text:     "What year was the World Wide Web invented?",
		Choices:  [4]string{"1985", "1989", "1993", "1995"},
		Answer:   1,
		Category: "Tech",
	},
	{
		Text:     "What does SSD stand for?",
		Choices:  [4]string{"Super Speed Disk", "Solid State Drive", "System Storage Device", "Shared Storage Disk"},
		Answer:   1,
		Category: "Tech",
	},
	{
		Text:     "Which protocol is used for sending email?",
		Choices:  [4]string{"HTTP", "FTP", "SMTP", "SSH"},
		Answer:   2,
		Category: "Tech",
	},
	{
		Text:     "What does RAM stand for?",
		Choices:  [4]string{"Read Access Memory", "Random Access Memory", "Rapid Application Memory", "Runtime Allocated Memory"},
		Answer:   1,
		Category: "Tech",
	},
	{
		Text:     "What is the most popular version control system?",
		Choices:  [4]string{"SVN", "Mercurial", "Git", "Perforce"},
		Answer:   2,
		Category: "Tech",
	},
	{
		Text:     "What does SSH stand for?",
		Choices:  [4]string{"Secure Shell", "System Shell Host", "Safe Socket Handler", "Secure System Hub"},
		Answer:   0,
		Category: "Tech",
	},
	{
		Text:     "What is Docker primarily used for?",
		Choices:  [4]string{"Database management", "Containerization", "Code editing", "Version control"},
		Answer:   1,
		Category: "Tech",
	},
	{
		Text:     "What does URL stand for?",
		Choices:  [4]string{"Universal Resource Locator", "Uniform Resource Locator", "Universal Reference Link", "Unified Resource Locator"},
		Answer:   1,
		Category: "Tech",
	},
	{
		Text:     "Which cloud provider has the largest market share?",
		Choices:  [4]string{"Google Cloud", "Microsoft Azure", "Amazon Web Services", "IBM Cloud"},
		Answer:   2,
		Category: "Tech",
	},
	{
		Text:     "What does JSON stand for?",
		Choices:  [4]string{"JavaScript Object Notation", "Java Standard Object Naming", "JavaScript Ordered Nodes", "JSON Standard Object Notation"},
		Answer:   0,
		Category: "Tech",
	},

	// ─── Science ──────────────────────────────────────────────
	{
		Text:     "What planet is closest to the Sun?",
		Choices:  [4]string{"Venus", "Earth", "Mercury", "Mars"},
		Answer:   2,
		Category: "Science",
	},
	{
		Text:     "What is the chemical symbol for gold?",
		Choices:  [4]string{"Go", "Gd", "Au", "Ag"},
		Answer:   2,
		Category: "Science",
	},
	{
		Text:     "How many bits are in a byte?",
		Choices:  [4]string{"4", "6", "8", "16"},
		Answer:   2,
		Category: "Science",
	},
	{
		Text:     "What is the speed of light in a vacuum (approx)?",
		Choices:  [4]string{"300,000 km/s", "150,000 km/s", "1,000,000 km/s", "30,000 km/s"},
		Answer:   0,
		Category: "Science",
	},
	{
		Text:     "What is the binary representation of the decimal number 10?",
		Choices:  [4]string{"1100", "1010", "1001", "1110"},
		Answer:   1,
		Category: "Science",
	},
	{
		Text:     "What does the 'http' in a URL stand for?",
		Choices:  [4]string{"HyperText Transfer Protocol", "High Tech Transfer Process", "Hyper Transfer Text Protocol", "Home Tool Transfer Protocol"},
		Answer:   0,
		Category: "Science",
	},
	{
		Text:     "What gas makes up most of Earth's atmosphere?",
		Choices:  [4]string{"Oxygen", "Carbon Dioxide", "Nitrogen", "Hydrogen"},
		Answer:   2,
		Category: "Science",
	},
	{
		Text:     "How many planets are in our solar system?",
		Choices:  [4]string{"7", "8", "9", "10"},
		Answer:   1,
		Category: "Science",
	},
	{
		Text:     "What is the largest organ in the human body?",
		Choices:  [4]string{"Liver", "Brain", "Skin", "Heart"},
		Answer:   2,
		Category: "Science",
	},
	{
		Text:     "What is the powerhouse of the cell?",
		Choices:  [4]string{"Nucleus", "Ribosome", "Mitochondria", "Golgi Apparatus"},
		Answer:   2,
		Category: "Science",
	},
	{
		Text:     "What is absolute zero in Celsius?",
		Choices:  [4]string{"-273.15°C", "-100°C", "-459.67°C", "0°C"},
		Answer:   0,
		Category: "Science",
	},
	{
		Text:     "What element has the atomic number 1?",
		Choices:  [4]string{"Helium", "Lithium", "Hydrogen", "Carbon"},
		Answer:   2,
		Category: "Science",
	},
	{
		Text:     "What is the hardest natural substance on Earth?",
		Choices:  [4]string{"Gold", "Iron", "Diamond", "Quartz"},
		Answer:   2,
		Category: "Science",
	},
	{
		Text:     "How many bones are in the adult human body?",
		Choices:  [4]string{"186", "206", "226", "256"},
		Answer:   1,
		Category: "Science",
	},
	{
		Text:     "What is the chemical formula for water?",
		Choices:  [4]string{"CO2", "H2O", "NaCl", "O2"},
		Answer:   1,
		Category: "Science",
	},
	{
		Text:     "What planet is known as the Red Planet?",
		Choices:  [4]string{"Jupiter", "Saturn", "Mars", "Venus"},
		Answer:   2,
		Category: "Science",
	},

	// ─── Git ──────────────────────────────────────────────────
	{
		Text:     "What command creates a new Git branch?",
		Choices:  [4]string{"git new", "git branch", "git create", "git init"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What does 'git stash' do?",
		Choices:  [4]string{"Deletes changes permanently", "Saves changes temporarily without committing", "Pushes to remote", "Creates a new branch"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What does 'git clone' do?",
		Choices:  [4]string{"Creates a new branch", "Copies a remote repository locally", "Merges two branches", "Deletes a repository"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What is a 'merge conflict' in Git?",
		Choices:  [4]string{"When Git can't automatically combine changes", "When a branch is deleted", "When a commit is reverted", "When the remote is unreachable"},
		Answer:   0,
		Category: "Git",
	},
	{
		Text:     "What does 'git log' show?",
		Choices:  [4]string{"Uncommitted changes", "The commit history", "Remote branches", "Stashed changes"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What does 'git rebase' do?",
		Choices:  [4]string{"Deletes commits", "Replays commits on top of another base", "Creates a new repository", "Pushes changes to remote"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What file tells Git to ignore certain files?",
		Choices:  [4]string{".gitconfig", ".gitignore", ".gitmodules", ".gitattributes"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What does 'git diff' show?",
		Choices:  [4]string{"The commit log", "Differences between file versions", "Branch list", "Remote URLs"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What is a 'pull request'?",
		Choices:  [4]string{"A command to download code", "A request to merge changes into a branch", "A way to delete a branch", "A type of Git hook"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What does 'git cherry-pick' do?",
		Choices:  [4]string{"Applies a specific commit to the current branch", "Deletes selected files", "Picks the best merge strategy", "Selects files to stage"},
		Answer:   0,
		Category: "Git",
	},
	{
		Text:     "What is 'git bisect' used for?",
		Choices:  [4]string{"Splitting a repository", "Finding the commit that introduced a bug", "Merging two branches", "Creating a tag"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What does 'git fetch' do differently from 'git pull'?",
		Choices:  [4]string{"It deletes local changes", "It downloads changes without merging", "It pushes to remote", "Nothing, they're the same"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "How do you undo the last Git commit (keeping changes)?",
		Choices:  [4]string{"git revert HEAD", "git reset --soft HEAD~1", "git undo", "git rollback"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What does 'git remote -v' show?",
		Choices:  [4]string{"Verbose commit log", "Remote repository URLs", "Version information", "Branch differences"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What is a Git 'tag' used for?",
		Choices:  [4]string{"Labeling a branch", "Marking a specific commit (like a release)", "Adding comments to files", "Categorizing issues"},
		Answer:   1,
		Category: "Git",
	},
	{
		Text:     "What does 'git blame' show?",
		Choices:  [4]string{"Who wrote each line of a file", "Files with errors", "Merge conflicts", "Deleted branches"},
		Answer:   0,
		Category: "Git",
	},

	// ─── Fun/Nerd ─────────────────────────────────────────────
	{
		Text:     "In what year was the first tweet sent?",
		Choices:  [4]string{"2004", "2005", "2006", "2007"},
		Answer:   2,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What was the first programming language?",
		Choices:  [4]string{"FORTRAN", "COBOL", "Assembly", "Plankalkül"},
		Answer:   3,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What does the 'S' in 'HTTPS' stand for?",
		Choices:  [4]string{"Simple", "Secure", "Standard", "System"},
		Answer:   1,
		Category: "Fun/Nerd",
	},
	{
		Text:     "Who co-founded Apple with Steve Jobs?",
		Choices:  [4]string{"Bill Gates", "Steve Wozniak", "Tim Cook", "Elon Musk"},
		Answer:   1,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What is the answer to life, the universe, and everything?",
		Choices:  [4]string{"7", "13", "42", "100"},
		Answer:   2,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What year was the original iPhone released?",
		Choices:  [4]string{"2005", "2006", "2007", "2008"},
		Answer:   2,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What color is the default GitHub contribution graph?",
		Choices:  [4]string{"Blue", "Green", "Purple", "Orange"},
		Answer:   1,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What is the name of the Linux mascot?",
		Choices:  [4]string{"Tux", "Penguin Pete", "Linux Larry", "Beastie"},
		Answer:   0,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What was Google's original name?",
		Choices:  [4]string{"Alphabet", "BackRub", "PageRank", "SearchBot"},
		Answer:   1,
		Category: "Fun/Nerd",
	},
	{
		Text:     "In what decade was email invented?",
		Choices:  [4]string{"1960s", "1970s", "1980s", "1990s"},
		Answer:   1,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What fictional AI became self-aware in the Terminator movies?",
		Choices:  [4]string{"HAL 9000", "WOPR", "Skynet", "Ultron"},
		Answer:   2,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What does the 'F' in FIFO stand for?",
		Choices:  [4]string{"Fast", "First", "Final", "Full"},
		Answer:   1,
		Category: "Fun/Nerd",
	},
	{
		Text:     "How many keys are on a standard full-size keyboard?",
		Choices:  [4]string{"88", "101", "104", "110"},
		Answer:   2,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What is the most popular programming language on GitHub (2024)?",
		Choices:  [4]string{"Python", "JavaScript", "TypeScript", "Java"},
		Answer:   1,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What movie features the quote 'I'm in'?",
		Choices:  [4]string{"The Matrix", "Hackers", "Every hacking movie ever", "WarGames"},
		Answer:   2,
		Category: "Fun/Nerd",
	},
	{
		Text:     "What is the Konami Code?",
		Choices:  [4]string{"↑↑↓↓←→←→BA", "↑↓←→ABAB", "←→↑↓BABA", "ABAB↑↓←→"},
		Answer:   0,
		Category: "Fun/Nerd",
	},
}
