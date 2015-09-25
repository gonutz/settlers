package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gonutz/settlers/game"
	"math/rand"
	"time"
)

const (
	NewGameOption = iota
	JoinRemoteGameOption
	QuitOption
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGameUI(window Window) (*gameUI, error) {
	graphics, err := newGraphics()
	if err != nil {
		return nil, err
	}
	g := game.New([]game.Color{game.Red, game.Blue, game.White}, 1)
	ui := &gameUI{
		game:     g,
		window:   window,
		camera:   newCamera(),
		graphics: graphics,
		buyMenu:  newBuyMenu(graphics, g),
		mainMenu: newMenu(
			50,
			newButton("New Game", rect{500, 400, 450, 80}, NewGameOption),
			newButton("Join Remote Game", rect{500, 500, 450, 80}, JoinRemoteGameOption),
			newButton("Quit", rect{500, 600, 450, 80}, QuitOption),
			newTextField(rect{500, 700, 450, 80}, graphics.font, -1),
		),
		joinGameMenu: newMenu(
			50,
			newButton("Yes", rect{500, 400, 450, 80}, -1),
			newButton("OK", rect{500, 500, 450, 80}, -1),
			newButton("Well...", rect{500, 600, 450, 80}, -1),
		),
	}
	ui.mainMenu.show()
	ui.menus = []*menu{ui.mainMenu, ui.joinGameMenu}
	if err := ui.init(); err != nil {
		return nil, err
	}
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
	menus          []*menu
	mainMenu       *menu
	joinGameMenu   *menu
	quitting       bool
}

type Window interface {
	Close()
}

func (ui *gameUI) init() error {
	return ui.graphics.createGameBackground(ui.game)
}

func (ui *gameUI) KeyDown(key glfw.Key) {
	for _, menu := range ui.menus {
		menu.keyPressed(key)
	}
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
		return
	}

	gameX, gameY := ui.camera.windowToGame(ui.mouseX, ui.mouseY)

	if ui.game.State == game.NotStarted {
		for _, menu := range ui.menus {
			if action := menu.click(gameX, gameY); action != -1 {
				switch action {
				case NewGameOption:
					ui.game = game.New([]game.Color{game.Red, game.Blue, game.White}, rand.Int())
					ui.init()
				case JoinRemoteGameOption:
					ui.mainMenu.hide()
					ui.joinGameMenu.show()
				case QuitOption:
					ui.window.Close()
				}
				return // ignore other event handling if any button was pressed
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
	for _, menu := range ui.menus {
		menu.mouseMovedTo(gameX, gameY)
	}
	ui.mouseX, ui.mouseY = x, y
}

func (ui *gameUI) RuneTyped(r rune) {
	for _, menu := range ui.menus {
		menu.runeTyped(r)
	}
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
		for _, menu := range ui.menus {
			menu.draw(ui.graphics)
		}
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
		return "Menu"
	case game.BuildingFirstSettlement:
		return "Build your first Settlement"
	case game.BuildingFirstRoad:
		return "Build your first Road"
	case game.BuildingSecondSettlement:
		return "Build your second Settlement"
	case game.BuildingSecondRoad:
		return "Build your second Road"
	case game.BuildingNewRoad:
		return "Build your Road"
	case game.BuildingNewSettlement:
		return "Build your Settlement"
	case game.BuildingNewCity:
		return "Build your City"
	case game.ChoosingNextAction:
		return "Choose your next Action"
	case game.RollingDice:
		return "Let's roll the Dice"
	}
	return "Unknown State"
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
