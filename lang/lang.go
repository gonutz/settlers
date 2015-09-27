package lang

var CurrentLanguage Language

type Language int

const (
	English Language = iota
	German
	LanguageCount // TODO this has to always come last
)

func Get(id ItemID) string {
	return languages[CurrentLanguage][id]
}

func GetLanguageName(id Language) string {
	return languages[id][LanguageName]
}

type ItemID int

const (
	LanguageName ItemID = iota
	Title
	Menu
	NewGame
	JoinRemoteGame
	Quit
	LanguageWord
)

var languages = [][]string{
	// English
	[]string{
		"English",
		"Settlers",
		"Menu",
		"New Game",
		"Join Remote Game",
		"Quit",
		"Language",
	},

	// German
	[]string{
		"Deutsch",
		"Siedler",
		"Menü",
		"Neues Spiel",
		"Netzwerkspiel beitreten",
		"Beenden",
		"Sprache",
	},
}
