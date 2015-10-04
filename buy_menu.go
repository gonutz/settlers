package main

import "github.com/gonutz/settlers/game"

const (
	cellW = 70
	cellH = 70
)

func newBuyMenu(g *graphics, gamer gamer) *buyMenu {
	symbols := []string{"lumber", "brick", "wool", "ore", "grain"}
	maxHeight := 0
	for _, symbol := range symbols {
		_, h := g.imageSize(symbol + "_symbol")
		if h > maxHeight {
			maxHeight = h
		}
	}

	mainBounds := rect{0, gameH - 3*tileH/2, 500, 4 * cellH}
	iconBounds := rect{mainBounds.w, 0, 100, tileH / 2}
	iconBounds.y = mainBounds.y + mainBounds.h - iconBounds.h
	return &buyMenu{
		graphics:   g,
		gamer:      gamer,
		mainBounds: mainBounds,
		iconBounds: iconBounds,
		xOffset:    -leftBorder - mainBounds.w,
		left:       -leftBorder - mainBounds.w,
		right:      -leftBorder,
		state:      closed,
	}
}

// The menu comes out of the left screen border, at first it is only an icon
// and when you click that, the whole menu expands to the right.
//
//main bounds
//|       |
//v       v
//--------+
// the    |__
// menu   |  |
//-----------+
//        ^  ^
//        |  |
//     icon bounds
//
type buyMenu struct {
	graphics    *graphics
	gamer       gamer
	mainBounds  rect
	iconBounds  rect
	left, right int
	xOffset     int
	state       menuState
}

type gamer interface {
	Game() *game.Game
}

type menuState int

const (
	closed menuState = iota
	opening
	opened
	closing
)

func (m *buyMenu) draw() {
	// draw background and icon
	m.drawRect(m.mainBounds)
	m.drawRect(m.iconBounds)
	x, y := m.iconBounds.center()
	m.graphics.drawImageCenteredAt("hammer_icon", x+m.xOffset, y)

	// TODO have the arrow?
	//arrow := "right_arrow"
	//if m.state == closed || m.state == closing {
	//arrow = "left_arrow"
	//}
	//m.graphics.drawImageCenteredAt(arrow, x+m.xOffset, m.iconBounds.y+m.iconBounds.h-15)

	// draw costs
	costs := [][]string{
		[]string{"lumber_symbol", "brick_symbol"},
		[]string{"lumber_symbol", "brick_symbol", "grain_symbol", "wool_symbol"},
		[]string{"grain_symbol", "grain_symbol", "ore_symbol", "ore_symbol", "ore_symbol"},
		[]string{"grain_symbol", "wool_symbol", "ore_symbol"},
	}
	for line, cost := range costs {
		for column, resource := range cost {
			x := m.mainBounds.x + m.xOffset + column*cellW + cellW/2
			y := m.mainBounds.y + line*cellH + cellH/2
			m.graphics.drawImageCenteredAt(resource, x, y)
		}
	}
	color := colorToString(m.gamer.Game().GetCurrentPlayer().Color)
	symbols := []string{"road_" + color + "_up", "settlement_" + color, "city_" + color, "card_symbol"}
	for line, symbol := range symbols {
		x := m.xOffset + m.mainBounds.w - cellW
		y := m.mainBounds.y + line*cellH + cellH/2
		color := [4]float32{1, 1, 1, 1}
		if !m.canBuyItem(line) {
			color[3] = 0.3
		}
		m.graphics.drawColoredImageCenteredAt(symbol, x, y, color)
	}
}

func (m *buyMenu) canBuyItem(index int) bool {
	game := m.gamer.Game()
	switch index {
	case 0:
		return game.CanBuyRoad()
	case 1:
		return game.CanBuySettlement()
	case 2:
		return game.CanBuyCity()
	case 3:
		return game.CanBuyDevelopmentCard()
	}
	panic("illegal index")
}

func (m *buyMenu) drawRect(r rect) {
	m.graphics.rect(r.x+m.xOffset, r.y, r.w, r.h, [4]float32{0.6, 0.4, 0.1, 0.8})
}

func (m *buyMenu) update() {
	const speed = 20
	if m.state == opening {
		m.xOffset += speed
		if m.xOffset >= m.right {
			m.xOffset = m.right
			m.state = opened
		}
	}
	if m.state == closing {
		m.xOffset -= speed
		if m.xOffset <= m.left {
			m.xOffset = m.left
			m.state = closed
		}
	}
}

func (m *buyMenu) click(x, y int) {
	if m.iconBounds.moveBy(m.xOffset, 0).contains(x, y) {
		if m.state == opened || m.state == opening {
			m.state = closing
		} else {
			m.state = opening
		}
		return
	}

	areaW := 2 * cellW
	left := m.xOffset + m.mainBounds.w - areaW
	for line := 0; line < 4; line++ {
		top := m.mainBounds.y + line*cellH
		area := rect{left, top, areaW, cellH}
		if area.contains(x, y) && m.canBuyItem(line) {
			m.state = closing
			m.buyItem(line)
			return
		}
	}
}

func (m *buyMenu) buyItem(index int) {
	game := m.gamer.Game()
	switch index {
	case 0:
		game.BuyRoad()
	case 1:
		game.BuySettlement()
	case 2:
		game.BuyCity()
	case 3:
		game.BuyDevelopmentCard()
	default:
		panic("illegal index")
	}
}
