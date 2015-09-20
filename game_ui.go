package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gonutz/settlers/game"
)

func NewGameUI(window Window) (*gameUI, error) {
	graphics, err := newGraphics()
	if err != nil {
		return nil, err
	}
	ui := &gameUI{
		game:     game.New([]game.Color{game.Red, game.Blue, game.White, game.Orange}, 0),
		window:   window,
		camera:   newCamera(),
		graphics: graphics,
	}
	if err := ui.init(); err != nil {
		return nil, err
	}
	return ui, nil
}

type gameUI struct {
	game           *game.Game
	window         Window
	camera         *camera
	graphics       *graphics
	mouseX, mouseY float64
}

type Window interface {
	Close()
}

func (ui *gameUI) init() error {
	return ui.graphics.createGameBackground(ui.game)
}

func (ui *gameUI) KeyDown(key glfw.Key) {
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
	if ui.game.State == game.BuildingFirstSettlement {
		gameX, gameY := ui.camera.windowToGame(ui.mouseX, ui.mouseY)
		corner, hit := screenToCorner(gameX, gameY)
		if hit && ui.game.CanBuildSettlementAt(corner) {
			ui.game.BuildSettlement(corner)
		}
	}
}

func (ui *gameUI) MouseEntered()             {}
func (ui *gameUI) MouseExited()              { ui.mouseX, ui.mouseY = -10000, -10000 }
func (ui *gameUI) MouseMovedTo(x, y float64) { ui.mouseX, ui.mouseY = x, y }

func (ui *gameUI) RuneTyped(r rune) {}

func (ui *gameUI) WindowSizeChangedTo(width, height int) {
	ui.camera.windowSizeChangedTo(width, height)
}

func (ui *gameUI) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	ui.graphics.drawBackground()

	// TODO draw all buildings and robber

	player := ui.game.GetCurrentPlayer()
	if ui.game.State == game.BuildingFirstSettlement {
		ui.graphics.showInstruction("Build your first settlement", player.Color)
		gameX, gameY := ui.camera.windowToGame(ui.mouseX, ui.mouseY)
		corner, hit := screenToCorner(gameX, gameY)
		canBuild := hit && ui.game.CanBuildSettlementAt(corner)
		if canBuild {
			x, y := cornerToScreen(corner)
			ui.graphics.drawSettlementAt(x, y, player.Color)
		} else {
			ui.graphics.drawHoveringSettlementAt(gameX, gameY, player.Color)
		}
	} else if ui.game.State == game.BuildingFirstRoad {
		ui.graphics.showInstruction("Build your first road", player.Color)
	}

	ui.graphics.drawSettlementAt(0, 0, game.Blue)
}
