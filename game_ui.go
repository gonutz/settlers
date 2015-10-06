package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gonutz/settlers/game"
	"github.com/gonutz/settlers/lang"
	"github.com/gonutz/settlers/settings"
	"math/rand"
	"strconv"
	"time"
)

const (
	NewGameOption = iota
	JoinRemoteGameOption
	ChooseLanguageOption
	QuitOption
	LanguageOKOption
	ThreePlayersOption
	FourPlayersOption
	StartGameOption
	NewGameBackOption

	LanguageOptionOffset = 1000
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGameUI(win Window) (*gameUI, error) {
	if err := settings.Settings.Load(); err != nil {
		fmt.Println("cannot load last settings:", err)
	}

	graphics, err := newGraphics()
	if err != nil {
		return nil, err
	}
	g := game.New([]game.Color{game.Red, game.Blue, game.White}, 1)

	// main menu
	size := func(w, h int) rect { return rect{w: w, h: h} }
	newGame := newButton(lang.NewGame, size(500, 80), NewGameOption)
	joinRemoteGame := newButton(lang.JoinRemoteGame, size(500, 80), JoinRemoteGameOption)
	chooseLanguage := newButton(lang.LanguageWord, size(500, 80), ChooseLanguageOption)
	quit := newButton(lang.Quit, size(500, 80), QuitOption)
	mainMenu := newWindow(
		rect{0, 0, gameW, gameH},
		newVerticalFlowLayout(20),
		newGame,
		joinRemoteGame,
		chooseLanguage,
		quit,
	)

	// language menu
	var langBoxes []*checkBox
	for language := lang.Language(0); language < lang.LastLanguage; language++ {
		action := LanguageOptionOffset + int(language)
		cb := newCheckBox(lang.Item(language), size(300, 80), action)
		cb.setChecked(settings.Settings.Language == int(language))
		langBoxes = append(langBoxes, cb)
	}
	languageMenu := newWindow(
		rect{0, 0, gameW, gameH},
		newVerticalFlowLayout(0),
		newCheckBoxGroup(langBoxes...),
		newSpacer(size(0, 20)),
		newButton(lang.OK, size(300, 80), LanguageOKOption),
	)
	languageMenu.setVisible(false)

	// new game menu
	threePlayers := newCheckBox(lang.ThreePlayers, size(350, 80), ThreePlayersOption)
	threePlayers.checked = settings.Settings.PlayerCount == 3
	fourPlayers := newCheckBox(lang.FourPlayers, size(350, 80), FourPlayersOption)
	fourPlayers.checked = settings.Settings.PlayerCount == 4
	var playerMenus [4]*window
	for i := range playerMenus {
		playerIndex := i // need to copy this for use in closures
		nameText := newTextBox(lang.Name, rect{0, 0, 500, 80}, graphics.font)
		nameText.text = settings.Settings.PlayerNames[i]
		nameText.onTextChange(func(text string) {
			settings.Settings.PlayerNames[playerIndex] = text
		})
		playHere := newCheckBox(lang.PlayHere, rect{0, 0, 500, 80}, -1)
		playHere.onCheckChange(func(checked bool) {
			if checked {
				settings.Settings.PlayerTypes[playerIndex] = settings.Human
			}
		})
		playHere.checked = settings.Settings.PlayerTypes[i] == settings.Human
		playAI := newCheckBox(lang.AIPlayer, rect{0, 0, 500, 80}, -1)
		playAI.onCheckChange(func(checked bool) {
			if checked {
				settings.Settings.PlayerTypes[playerIndex] = settings.AI
			}
		})
		playAI.checked = settings.Settings.PlayerTypes[i] == settings.AI
		ipText := newTextBox(lang.IP, rect{0, 0, 500, 80}, graphics.font)
		ipText.text = settings.Settings.IPs[i]
		ipText.onTextChange(func(text string) {
			settings.Settings.IPs[playerIndex] = text
		})
		portText := newTextBox(lang.Port, rect{0, 0, 500, 80}, graphics.font)
		portText.text = settings.Settings.Ports[i]
		portText.onTextChange(func(text string) {
			settings.Settings.Ports[playerIndex] = text
		})
		connectButton := newButton(lang.Connect, size(500, 80), -1)
		playNetwork := newCheckBox(lang.NetworkPlayer, rect{0, 0, 500, 80}, -1)
		playNetwork.onCheckChange(func(checked bool) {
			ipText.setEnabled(checked)
			portText.setEnabled(checked)
			connectButton.setEnabled(checked)
			if checked {
				settings.Settings.PlayerTypes[playerIndex] = settings.NetworkPlayer
			}
		})
		playNetwork.setChecked(settings.Settings.PlayerTypes[i] == settings.NetworkPlayer)

		playerMenus[i] = newWindow(
			rect{},
			newVerticalFlowLayout(0),
			newSpacer(rect{0, 0, 620, 30}),
			nameText,
			newCheckBoxGroup(
				playHere,
				playAI,
				playNetwork,
			),
			ipText,
			portText,
			connectButton,
			newSpacer(rect{0, 0, 0, 30}),
		)
	}
	var playerTabs [4]*tab
	for i := range playerTabs {
		col := game.Color(i)
		playerTabs[i] = newTab(fullPlayerColor(col), playerMenus[i], i < settings.Settings.PlayerCount)
	}
	playersSheet := newTabSheet(60, playerTabs[:]...)

	newGameMenu := newWindow(
		rect{0, 0, gameW, gameH},
		newVerticalFlowLayout(20),
		newCheckBoxGroup(threePlayers, fourPlayers),
		playersSheet,
		newButton(lang.StartGame, rect{0, 0, 400, 80}, StartGameOption),
		newButton(lang.Back, rect{0, 0, 400, 80}, NewGameBackOption),
	)
	newGameMenu.setVisible(false)

	ui := &gameUI{
		game:           g,
		window:         win,
		camera:         newCamera(),
		graphics:       graphics,
		mainMenu:       mainMenu,
		newGameMenu:    newGameMenu,
		languageMenu:   languageMenu,
		playerTabSheet: playersSheet,
		lastPlayerTab:  playerTabs[3],
	}
	ui.buyMenu = newBuyMenu(graphics, ui)
	ui.gui = newComposite(ui.mainMenu, ui.newGameMenu, ui.languageMenu)
	if err := ui.init(); err != nil {
		return nil, err
	}
	ui.setLanguage(lang.Language(settings.Settings.Language))
	//ui.DEBUG_initGame()
	return ui, nil
}

type gameUI struct {
	game           *game.Game
	window         Window
	camera         *camera
	graphics       *graphics
	mouseX, mouseY float64
	buyMenu        *buyMenu
	gui            *composite
	mainMenu       *window
	languageMenu   *window
	newGameMenu    *window
	playerTabSheet *tabSheet
	lastPlayerTab  *tab
	quitting       bool
}

type Window interface {
	Close()
	SetTitle(title string)
}

func (ui *gameUI) Game() *game.Game { return ui.game }

func (ui *gameUI) setLanguage(id lang.Language) {
	lang.CurrentLanguage = id
	settings.Settings.Language = int(id)
	ui.window.SetTitle(lang.Get(lang.Title))
}

func (ui *gameUI) init() error {
	return ui.graphics.createGameBackground(ui.game)
}

func (ui *gameUI) KeyDown(key glfw.Key) {
	ui.gui.keyPressed(key)
	if key == glfw.KeyEscape {
		ui.window.Close()
	}
	if key == glfw.Key1 {
		ui.game.CurrentPlayer = 0
	}
	if key == glfw.Key2 {
		ui.game.CurrentPlayer = 1
	}
	if key == glfw.Key3 {
		ui.game.CurrentPlayer = 2
	}
	if key == glfw.Key4 {
		ui.game.CurrentPlayer = 3
	}
}

func (ui *gameUI) MouseButtonDown(button glfw.MouseButton) {
	if button != glfw.MouseButtonLeft {
		return // TODO handle right click when building to undo both buying and build
	}

	gameX, gameY := ui.camera.windowToGame(ui.mouseX, ui.mouseY)

	if ui.game.State == game.NotStarted {
		if action := ui.gui.click(gameX, gameY); action != -1 {
			switch action {
			case NewGameOption:
				ui.mainMenu.visible = false
				ui.newGameMenu.visible = true
			case ChooseLanguageOption:
				ui.mainMenu.visible = false
				ui.languageMenu.visible = true
			case QuitOption:
				ui.window.Close()
			case LanguageOKOption:
				ui.languageMenu.visible = false
				ui.mainMenu.visible = true
			case NewGameBackOption:
				ui.newGameMenu.visible = false
				ui.mainMenu.visible = true
			case StartGameOption:
				ui.game = game.New([]game.Color{game.Red, game.Blue, game.White}, rand.Int())
				ui.game = game.New([]game.Color{game.Red, game.Blue, game.White}, 0) // TODO remove
				ui.init()
				ui.game.Start()
			case ThreePlayersOption:
				ui.lastPlayerTab.visible = false
				ui.playerTabSheet.relayout()
				settings.Settings.PlayerCount = 3
			case FourPlayersOption:
				ui.lastPlayerTab.visible = true
				ui.playerTabSheet.relayout()
				settings.Settings.PlayerCount = 4
			}
			if action >= LanguageOptionOffset {
				language := lang.Language(action - LanguageOptionOffset)
				ui.setLanguage(language)
			}
		}
	} else if ui.game.State == game.ChoosingNextAction {
		ui.buyMenu.click(gameX, gameY)
	} else if ui.game.State == game.BuildingFirstSettlement ||
		ui.game.State == game.BuildingSecondSettlement ||
		ui.game.State == game.BuildingNewSettlement {
		corner, hit := screenToCorner(gameX, gameY)
		if hit && ui.game.CanBuildSettlementAt(corner) {
			ui.game.BuildSettlement(corner)
		}
	} else if ui.game.State == game.BuildingFirstRoad ||
		ui.game.State == game.BuildingSecondRoad ||
		ui.game.State == game.BuildingNewRoad {
		edge, hit := screenToEdge(gameX, gameY)
		if hit && ui.game.CanBuildRoadAt(edge) {
			ui.game.BuildRoad(edge)
		}
	} else if ui.game.State == game.BuildingNewCity {
		corner, hit := screenToCorner(gameX, gameY)
		if hit && ui.game.CanBuildCityAt(corner) {
			ui.game.BuildCity(corner)
		}
	} else if ui.game.State == game.RollingDice {
		center := rect{gameW/2 - 100, gameH/2 - 50, 200, 100}
		if center.contains(gameX, gameY) {
			ui.game.RollTheDice()
		}
	}
}

func (ui *gameUI) MouseEntered() {}
func (ui *gameUI) MouseExited()  { ui.mouseX, ui.mouseY = -10000, -10000 }

func (ui *gameUI) MouseMovedTo(x, y float64) {
	gameX, gameY := ui.camera.windowToGame(ui.mouseX, ui.mouseY)
	ui.gui.mouseMovedTo(gameX, gameY)
	ui.mouseX, ui.mouseY = x, y
}

func (ui *gameUI) RuneTyped(r rune) {
	ui.gui.runeTyped(r)
}

func (ui *gameUI) WindowSizeChangedTo(width, height int) {
	ui.camera.windowSizeChangedTo(width, height)
}

func (ui *gameUI) Draw() {
	ui.drawBaseGame()

	if ui.game.State > game.BuildingSecondRoad {
		ui.buyMenu.update()
		ui.buyMenu.draw()
	}

	player := ui.game.GetCurrentPlayer()
	ui.graphics.drawResources(ui.game.GetCurrentPlayer().Resources, playerColor(player.Color))
	color := player.Color
	if ui.game.State == game.NotStarted {
		color = game.White
	}
	ui.graphics.showInstruction(ui.stateInstruction(), color)

	if ui.game.State != game.RollingDice && ui.game.State > game.BuildingSecondRoad {
		ui.graphics.drawDice(ui.game.Dice)
	}

	gameX, gameY := ui.camera.windowToGame(ui.mouseX, ui.mouseY)

	if ui.game.State == game.NotStarted {
		ui.gui.draw(ui.graphics)
	} else if ui.game.State == game.ChoosingNextAction {
	} else if ui.game.State == game.BuildingFirstSettlement ||
		ui.game.State == game.BuildingSecondSettlement ||
		ui.game.State == game.BuildingNewSettlement {
		corner, hit := screenToCorner(gameX, gameY)
		canBuild := hit && ui.game.CanBuildSettlementAt(corner)
		if canBuild {
			x, y := cornerToScreen(corner)
			ui.graphics.drawSettlementAt(x, y, player.Color)
		} else {
			ui.graphics.drawHoveringSettlementAt(gameX, gameY, player.Color)
		}
	} else if ui.game.State == game.BuildingFirstRoad ||
		ui.game.State == game.BuildingSecondRoad ||
		ui.game.State == game.BuildingNewRoad {
		edge, hit := screenToEdge(gameX, gameY)
		canBuild := hit && ui.game.CanBuildRoadAt(edge)
		if canBuild {
			x, y := edgeToScreen(edge)
			ui.graphics.drawRoadAt(x, y, edge, player.Color)
		} else {
			ui.graphics.drawHoveringRoadAt(gameX, gameY, player.Color)
		}
	} else if ui.game.State == game.BuildingNewCity {
		corner, hit := screenToCorner(gameX, gameY)
		canBuild := hit && ui.game.CanBuildCityAt(corner)
		if canBuild {
			x, y := cornerToScreen(corner)
			ui.graphics.drawCityAt(x, y, player.Color)
		} else {
			ui.graphics.drawHoveringCityAt(gameX, gameY, player.Color)
		}
	} else if ui.game.State == game.RollingDice {
		const d = 100
		ui.graphics.rect(gameW/2-2*d, gameH/2-d, 4*d, 2*d, [4]float32{1, 1, 1, 0.8})
		ui.graphics.drawImageCenteredAt("dice", gameW/2, gameH/2)
	}
}

func (ui *gameUI) stateInstruction() string {
	switch ui.game.State {
	case game.NotStarted:
		return lang.Get(lang.Menu)
	case game.BuildingFirstSettlement:
		return lang.Get(lang.BuildFirstSettlement)
	case game.BuildingFirstRoad:
		return lang.Get(lang.BuildFirstRoad)
	case game.BuildingSecondSettlement:
		return lang.Get(lang.BuildSecondSettlement)
	case game.BuildingSecondRoad:
		return lang.Get(lang.BuildSecondRoad)
	case game.BuildingNewRoad:
		return lang.Get(lang.BuildRoad)
	case game.BuildingNewSettlement:
		return lang.Get(lang.BuildSettlement)
	case game.BuildingNewCity:
		return lang.Get(lang.BuildCity)
	case game.ChoosingNextAction:
		return lang.Get(lang.ChooseNextAction)
	case game.RollingDice:
		return lang.Get(lang.RollDice)
	}
	return "Unknown State: " + strconv.Itoa(int(ui.game.State))
}

func (ui *gameUI) drawBaseGame() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	ui.graphics.drawBackground()

	// draw roads first
	for _, p := range ui.game.GetPlayers() {
		for _, r := range p.GetBuiltRoads() {
			x, y := edgeToScreen(r.Position)
			ui.graphics.drawRoadAt(x, y, r.Position, p.Color)
		}
	}
	// draw buildings above the roads
	for _, p := range ui.game.GetPlayers() {
		for _, s := range p.GetBuiltSettlements() {
			x, y := cornerToScreen(s.Position)
			ui.graphics.drawSettlementAt(x, y, p.Color)
		}
		for _, c := range p.GetBuiltCities() {
			x, y := cornerToScreen(c.Position)
			ui.graphics.drawCityAt(x, y, p.Color)
		}
	}

	ui.graphics.drawRobber(tileToScreen(ui.game.Robber.Position))
}

func (ui *gameUI) Finish() {
	if err := settings.Settings.Save(); err != nil {
		fmt.Println("cannot save settings:", err)
	}
}

func (ui *gameUI) DEBUG_initGame() {
	ui.game.Players[0].Settlements[0].Position = game.TileCorner{4, 2}
	ui.game.Players[0].Settlements[1].Position = game.TileCorner{6, 2}
	ui.game.Players[0].Roads[0].Position = game.TileEdge{9, 2}
	ui.game.Players[0].Roads[1].Position = game.TileEdge{11, 2}

	ui.game.Players[1].Settlements[0].Position = game.TileCorner{4, 3}
	ui.game.Players[1].Settlements[1].Position = game.TileCorner{6, 3}
	ui.game.Players[1].Roads[0].Position = game.TileEdge{9, 3}
	ui.game.Players[1].Roads[1].Position = game.TileEdge{11, 3}

	ui.game.Players[2].Settlements[0].Position = game.TileCorner{4, 5}
	ui.game.Players[2].Settlements[1].Position = game.TileCorner{6, 5}
	ui.game.Players[2].Roads[0].Position = game.TileEdge{9, 5}
	ui.game.Players[2].Roads[1].Position = game.TileEdge{11, 5}

	ui.game.State = game.RollingDice

	ui.game.Players[0].Resources[game.Wool] = 100
	ui.game.Players[0].Resources[game.Grain] = 100
	ui.game.Players[0].Resources[game.Ore] = 100
	ui.game.Players[0].Resources[game.Lumber] = 100
	ui.game.Players[0].Resources[game.Brick] = 100
}
