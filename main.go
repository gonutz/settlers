package main

import (
	"github.com/gonutz/prototype/draw"
	"github.com/gonutz/settlers/game"
	"math/rand"
	"strconv"
)

func main() {
	game := game.New([]game.Color{game.Red, game.Blue, game.White}, rand.Int)

	draw.RunWindow("Settlers", 1400, 1100, draw.Resizable, func(window draw.Window) {
		if window.WasKeyPressed("escape") {
			window.Close()
		}

		for i := 2; i <= 12; i++ {
			// TODO what about 7?
			if window.WasKeyPressed(strconv.Itoa(i)) {
				game.DealResources(i)
			}
		}

		drawGame(game, window)
	})
}

const (
	tileW       = 200
	tileH       = 212
	tileYOffset = 162
)

func drawGame(g game.Game, window draw.Window) {
	w, h := window.Size()
	xOffset, yOffset := 0, 0
	window.FillRect(0, 0, w, h, draw.DarkBlue)
	mouseX, mouseY := window.MousePosition()
	mouseX -= xOffset
	mouseY -= yOffset
	for _, tile := range g.GetTiles() {
		x := xOffset + tile.Position.X*tileW/2
		y := yOffset + tile.Position.Y*tileYOffset
		visible := window.IsMouseDown(draw.LeftButton) ||
			!hexContains(x, y, mouseX, mouseY)
		if visible {
			var file string
			switch tile.Terrain {
			case game.Forest:
				file = "./forest.png"
			case game.Field:
				file = "./field.png"
			case game.Mountains:
				file = "./mountains.png"
			case game.Pasture:
				file = "./pasture.png"
			case game.Desert:
				file = "./desert.png"
			case game.Hills:
				file = "./hills.png"
			case game.Water:
				file = "./water.png"
			}
			window.DrawImageFile(file, x, y)
		}
		if tile.Number != 0 || tile.Terrain == game.Desert {
			number := strconv.Itoa(tile.Number)
			const scale = 4
			w, h := window.GetScaledTextSize(number, scale)
			x, y := x+(tileW-w)/2, y+(tileH-h)/2
			plateSize := max(w, h) + 10
			window.FillEllipse(x+w/2-plateSize/2, y+h/2-plateSize/2,
				plateSize, plateSize, draw.White)
			numberColor := draw.Black
			if tile.Number == 6 || tile.Number == 8 {
				numberColor = draw.Red
			}
			window.DrawScaledText(number, x+3, y+3, scale, numberColor)
			if tile.HasRobber {
				window.FillEllipse(x+w/2-plateSize/2, y+h/2-plateSize/2,
					plateSize, plateSize, draw.DarkGray)
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func hexContains(x, y, px, py int) bool {
	w, h := tileW, tileH
	yOffset := h - tileYOffset
	if px < x || px >= x+w || py < y || py >= y+h {
		return false
	}
	// now it is definitely in bouding rectangle
	relX, relY := px-x, py-y
	if relY < yOffset {
		if relX < w/2 {
			// upper-left quarter
			return relX > (yOffset-relY)*2
		} else {
			// upper-right quarter
			return relX-w/2 <= relY*2
		}
	}
	if relY > tileYOffset {
		if relX < w/2 {
			// bottom-left quarter
			return relX > (relY-tileYOffset)*2
		} else {
			// bottom-right quarter
			return w-relX > (relY-tileYOffset)*2
		}
	}
	return true
}
