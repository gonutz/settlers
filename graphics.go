package main

import "github.com/gonutz/settlers/game"

func cornerToScreen(c game.TileCorner) (x, y int) {
	y = c.Y * (tileH - tileSlopeHeight)
	if c.X%2 != c.Y%2 {
		y += tileSlopeHeight
	}
	return c.X * tileW / 2, y
}

func edgeToScreen(e game.TileEdge) (x, y int) {
	x = e.X * tileW / 4
	y = tileSlopeHeight/2 + e.Y*(tileH-tileSlopeHeight)
	if isEven(e.X) {
		y = y + tileSlopeHeight/2 + (tileH-2*tileSlopeHeight)/2
	}
	return
}

func isEven(x int) bool {
	return x%2 == 0
}
func isEdgeVertical(e game.TileEdge) bool {
	return isEven(e.X)
}

func isEdgeGoingDown(e game.TileEdge) bool {
	if isEven(e.Y) {
		return (e.X-1)%4 == 0
	}
	return (e.X-3)%4 == 0
}

const (
	tileW           = 200
	tileH           = 212
	tileSlopeHeight = 50
	tileYOffset     = tileH - tileSlopeHeight
)

func tileToScreen(p game.TilePosition) (x, y, w, h int) {
	return p.X * tileW / 2, p.Y * tileYOffset, tileW, tileH
}
