package main

import (
	"fmt"
	"github.com/gonutz/prototype/draw"
	"github.com/gonutz/settlers/game"
	"math/rand"
	"strconv"
)

func main() {
	g := game.New([]game.Color{game.Red, game.Blue, game.White}, rand.Int)
	state := buildingSettlement
	g.Players[0].Roads[0].Position = game.TileEdge{9, 2}

	draw.RunWindow("Settlers", 1400, 1000, draw.Resizable, func(window draw.Window) {
		if window.WasKeyPressed("escape") {
			window.Close()
		}

		if window.WasKeyPressed("s") {
			state = buildingSettlement
		}
		if window.WasKeyPressed("r") {
			state = buildingRoad
		}
		if window.WasKeyPressed("c") {
			state = buildingCity
		}

		for i := 1; i <= 4; i++ {
			if window.WasKeyPressed(strconv.Itoa(i)) {
				fmt.Println(g.CurrentPlayer, "->", i)
				g.CurrentPlayer = i - 1
			}
		}

		if state == buildingSettlement || state == buildingCity || state == buildingRoad {
			canBuild := g.CanPlayerBuildSettlement
			build := g.BuildSettlement
			screenToTile := screenToCorner
			if state == buildingCity {
				canBuild = g.CanPlayerBuildCity
				build = g.BuildCity
			}
			if state == buildingRoad {
				canBuild = g.CanPlayerBuildRoad
				build = g.BuildRoad
				screenToTile = screenToEdge
			}
			if canBuild() {
				for _, click := range window.Clicks() {
					if click.Button == draw.LeftButton {
						corner, ok := screenToTile(click.X, click.Y, 40)
						if ok {
							build(corner)
						} else {
							fmt.Println("miss")
						}
					}
				}
			}
		}

		for i := 2; i <= 12; i++ {
			// TODO what about 7?
			if window.WasKeyPressed(strconv.Itoa(i)) {
				g.DealResources(i)
			}
		}

		drawGame(g, window)
		if state == buildingSettlement || state == buildingCity {
			mx, my := window.MousePosition()
			corner, hit := screenToCorner(mx, my, 40)
			if hit {
				x, y := cornerToScreen(corner)
				color := currentPlayerColor(g)
				color.A = 0.75
				window.FillEllipse(x-40, y-40, 80, 80, color)
				if state == buildingSettlement && !g.CanBuildSettlementAt(game.TileCorner(corner)) {
					color.A = 1
					window.DrawLine(x-50, y-50, x+50, y+50, color)
					window.DrawLine(x-50, y+50, x+50, y-50, color)
				}
			}
		}
		if state == buildingRoad {
			mx, my := window.MousePosition()
			edge, hit := screenToEdge(mx, my, 40)
			if hit {
				x, y := edgeToScreen(game.TileEdge(edge))
				color := currentPlayerColor(g)
				color.A = 0.75
				window.FillEllipse(x-40, y-40, 80, 80, color)
				if !g.CanBuildRoadAt(game.TileEdge(edge)) {
					color.A = 1
					window.DrawLine(x-50, y-50, x+50, y+50, color)
					window.DrawLine(x-50, y+50, x+50, y-50, color)
				}
			}
		}
	})
}

// TODO keep the current state in game.Game?
type gameState int

const (
	buildingSettlement gameState = iota
	buildingCity
	buildingRoad
)

func currentPlayerColor(g *game.Game) draw.Color {
	if g.CurrentPlayer < 0 || g.CurrentPlayer >= g.PlayerCount {
		panic("invalid current player")
	}
	switch g.Players[g.CurrentPlayer].Color {
	case game.Red:
		return draw.Red
	case game.Blue:
		return draw.Blue
	case game.Orange:
		return draw.LightBrown
	}
	return draw.White
}

const (
	tileW            = 200
	tileH            = 212
	tileSlopeHeight  = 50
	tileYOffset      = tileH - tileSlopeHeight
	tileMiddleHeight = tileH - 2*(tileH-tileYOffset)
)

func drawGame(g *game.Game, window draw.Window) {
	// draw tiles and numbers
	w, h := window.Size()
	window.FillRect(0, 0, w, h, draw.DarkBlue)
	for _, tile := range g.GetTiles() {
		x := tile.Position.X * tileW / 2
		y := tile.Position.Y * tileYOffset
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
		if tile.Number != 0 || tile.Terrain == game.Desert {
			number := strconv.Itoa(tile.Number)
			const scale = 4
			w, h := window.GetScaledTextSize(number, scale)
			x, y := x+(tileW-w)/2, y+(tileH-h)/2
			plateSize := max(w, h) + 10
			window.FillEllipse(x+w/2-plateSize/2, y+h/2-plateSize/2,
				plateSize, plateSize, draw.RGBA(1, 1, 1, 0.5))
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

	// draw buildings
	for _, player := range g.GetPlayers() {
		for _, s := range player.GetBuiltSettlements() {
			var file string
			switch player.Color {
			case game.Red:
				file = "./settlement_red.png"
			case game.Blue:
				file = "./settlement_blue.png"
			case game.White:
				file = "./settlement_white.png"
			case game.Orange:
				file = "./settlement_orange.png"
			}
			x, y := cornerToScreen(game.Point(s.Position))
			window.DrawImageFile(file, x-30/2, y-42/2)
		}
		for _, s := range player.GetBuiltCities() {
			var file string
			switch player.Color {
			case game.Red:
				file = "./city_red.png"
			case game.Blue:
				file = "./city_blue.png"
			case game.White:
				file = "./city_white.png"
			case game.Orange:
				file = "./city_orange.png"
			}
			x, y := cornerToScreen(game.Point(s.Position))
			window.DrawImageFile(file, x-65/2, y-42/2)
		}
		for _, s := range player.GetBuiltRoads() {
			var color draw.Color
			switch player.Color {
			case game.Red:
				color = draw.Red
			case game.Blue:
				color = draw.Blue
			case game.White:
				color = draw.White
			case game.Orange:
				color = draw.LightBrown // TODO have orange in prototype/draw
			}
			x, y := edgeToScreen(s.Position)
			window.FillRect(x-20, y-20, 40, 40, color)
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

func hexagonAtTile(tileX, tileY int) polygon {
	x, y := tileToScreen(tileX, tileY)
	return hexagonAt(x, y)
}

func tileToScreen(tileX, tileY int) (x, y int) {
	x = tileX * tileW / 2
	y = tileY * tileYOffset
	return
}

func hexagonAt(x, y int) polygon {
	return polygon{
		{x + tileW/2, y},
		{x + tileW - 1, y + tileSlopeHeight},
		{x + tileW - 1, y + tileH - tileSlopeHeight},
		{x + tileW/2, y + tileH - 1},
		{x, y + tileH - tileSlopeHeight},
		{x, y + tileSlopeHeight},
	}
}

func (p polygon) draw(window draw.Window, color draw.Color) {
	for i := range p {
		j := (i + 1) % len(p)
		window.DrawLine(p[i].x, p[i].y, p[j].x, p[j].y, color)
	}
}

func hexContains(x, y, px, py int) bool {
	// TODO use polygon containment here?
	return isPointInPolygon(point{px, py}, hexagonAt(x, y))

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

func highlightCorner(game *game.Game, window draw.Window) {
	mx, my := window.MousePosition()
	w, h := window.Size()
	var position point
	const radius = 40
	const squareRadius = radius * radius

	yOffset := 0
	for y := 0; y < h; y += tileH - tileSlopeHeight {
		for x := 0; x < w; x += tileW / 2 {
			if squareDist(x, y+yOffset, mx, my) < squareRadius {
				position = point{x, y + yOffset}
			}
			yOffset = tileSlopeHeight - yOffset
		}
		yOffset = tileSlopeHeight - yOffset
	}

	if position.x != 0 || position.y != 0 {
		window.FillEllipse(
			position.x-radius, position.y-radius,
			2*radius, 2*radius,
			draw.RGBA(1, 0, 0.5, 0.75))
	}
}

func squareDist(x1, y1, x2, y2 int) int {
	dx, dy := x1-x2, y1-y2
	return dx*dx + dy*dy
}

func boundingRect(p game.TilePosition) rect {
	x := p.X * tileW / 2
	y := p.Y * tileYOffset
	return rect{x, y, tileW, tileH}
}

type rect struct{ x, y, w, h int }

func (r rect) contains(x, y int) bool {
	return x >= r.x && x < r.x+r.w && y >= r.y && y < r.y+r.h
}

func isPointInPolygon(p point, v polygon) bool {
	contained := false
	for i, j := 0, len(v)-1; i < len(v); i, j = i+1, i {
		if ((v[i].y > p.y) != (v[j].y > p.y)) &&
			(p.x < (v[j].x-v[i].x)*(p.y-v[i].y)/(v[j].y-v[i].y)+v[i].x) {
			contained = !contained
		}
	}
	return contained
}

type point struct{ x, y int }

type polygon []point

func highlightEdge(game *game.Game, window draw.Window) {
	mx, my := window.MousePosition()
	for _, tile := range game.GetTiles() {
		hex := hexagonAtTile(tile.Position.X, tile.Position.Y)

		r := rect{hex[1].x - tileSlopeHeight/2, hex[1].y, tileSlopeHeight, hex[2].y - hex[1].y}
		if r.contains(mx, my) {
			window.FillRect(r.x, r.y, r.w, r.h, draw.RGBA(1.0, 0, 0.5, 0.75))
			return
		}

		r = rect{hex[3].x, hex[2].y, hex[2].x - hex[3].x, hex[3].y - hex[2].y}
		if r.contains(mx, my) {
			window.FillRect(r.x, r.y, r.w, r.h, draw.RGBA(1.0, 0, 0.5, 0.75))
			return
		}

		r = rect{hex[4].x, hex[4].y, hex[3].x - hex[4].x, hex[3].y - hex[4].y}
		if r.contains(mx, my) {
			window.FillRect(r.x, r.y, r.w, r.h, draw.RGBA(1.0, 0, 0.5, 0.75))
			return
		}
	}
}

func cornerToScreen(c game.Point) (x, y int) {
	y = c.Y * (tileH - tileSlopeHeight)
	if c.X%2 != c.Y%2 {
		y += tileSlopeHeight
	}
	return c.X * tileW / 2, y
}

func screenToCorner(x, y, maxDist int) (corner game.Point, hit bool) {
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
	return game.Point{tileX, tileY}, true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func edgeToScreen(e game.TileEdge) (x, y int) {
	x = e.X * tileW / 4
	y = tileSlopeHeight/2 + e.Y*(tileH-tileSlopeHeight)
	if e.X%2 == 0 {
		y = y + tileSlopeHeight/2 + (tileH-2*tileSlopeHeight)/2
	}
	return
}

func screenToEdge(x, y, maxDist int) (edge game.Point, hit bool) {
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
