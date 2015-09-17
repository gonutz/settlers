package main

import "github.com/gonutz/settlers/game"

const (
	tileW           = 200
	tileH           = 212
	tileSlopeHeight = 50
	tileYOffset     = tileH - tileSlopeHeight
)

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

func tileToScreen(p game.TilePosition) (x, y, w, h int) {
	return p.X * tileW / 2, p.Y * tileYOffset, tileW, tileH
}

func screenToCorner(x, y int) (corner game.TileCorner, hit bool) {
	const maxDist = 45
	tileX := (x + tileW/4) / (tileW / 2)
	if abs(x-tileX*tileW/2) > maxDist {
		return
	}
	topH := tileH - tileSlopeHeight
	tileY := (y - tileSlopeHeight/2 + topH/2) / topH
	nextScreenY := tileY * topH
	if (tileX%2 == 0 && tileY%2 == 1) || (tileX%2 == 1 && tileY%2 == 0) {
		nextScreenY += tileSlopeHeight
	}
	if abs(y-nextScreenY) > maxDist {
		return
	}
	return game.TileCorner{tileX, tileY}, true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func screenToEdge(x, y int) (edge game.TileEdge, hit bool) {
	edge.X = (x + tileW/8) / (tileW / 4)
	sectionH := tileH - tileSlopeHeight
	edge.Y = y / sectionH
	hit = isValidEdge(edge.X, edge.Y)
	return
}

func isValidEdge(x, y int) bool {
	if y%2 == 0 {
		return x%4 != 0
	}
	return (x-2)%4 != 0
}
