package lang

var CurrentLanguage Language

type Language int

const (
	English Language = iota
	German
	LanguageCount // TODO this has to always come last
)

func Get(id Item) string {
	return languages[CurrentLanguage][id]
}

type Item int

const (
	EnglishName Item = iota
	GermanName
	Title
	Menu
	NewGame
	JoinRemoteGame
	Quit
	LanguageWord
	ThreePlayers
	FourPlayers
	Back
	StartGame
	OK
	PlayHere
	AIPlayer
	NetworkPlayer
	Name
	IP
	Port
)

var languages = [][]string{
	// English
	[]string{
		"English",
		"Deutsch",
		"Settlers",
		"Menu",
		"New Game",
		"Join Remote Game",
		"Quit",
		"Language",
		"3 Players",
		"4 Players",
		"Back",
		"Start Game",
		"OK",
		"Play on this PC",
		"Computer AI",
		"Network Player",
		"Name:",
		"IP:",
		"Port:",
	},

	// German
	[]string{
		"English",
		"Deutsch",
		"Siedler",
		"Menü",
		"Neues Spiel",
		"Netzwerkspiel beitreten",
		"Beenden",
		"Sprache",
		"3 Spieler",
		"4 Spieler",
		"Zurück",
		"Spiel starten",
		"OK",
		"Spielt hier",
		"Computerspieler",
		"Netzwerkspieler",
		"Name:",
		"IP:",
		"Port:",
	},
}
