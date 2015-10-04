package lang

var CurrentLanguage Language

type Language int

const (
	English Language = iota
	German
	LastLanguage // NOTE this has to always come last
)

const LanguageCount = int(LastLanguage)

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
	Connect
	BuildFirstSettlement
	BuildSecondSettlement
	BuildFirstRoad
	BuildSecondRoad
	BuildRoad
	BuildSettlement
	BuildCity
	ChooseNextAction
	RollDice
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
		"Connect",
		"Build your first Settlement",
		"Build your second Settlement",
		"Build your first Road",
		"Build your second Road",
		"Build your Road",
		"Build your Settlement",
		"Build your City",
		"Choose your next Action",
		"Roll the Dice",
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
		"Verbinden",
		"Baue deine erste Siedlung",
		"Baue deine zweite Siedlung",
		"Baue deine erste Straße",
		"Baue deine zweite Straße",
		"Baue deine Straße",
		"Baue deine Siedlung",
		"Baue deine Stadt",
		"Wähle deine nächste Aktion",
		"Würfle",
	},
}
