package main

import (
	"bufio"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/gonutz/fontstash.go/fontstash"
	"github.com/gonutz/settlers/game"
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func newGraphics() (*graphics, error) {
	g := &graphics{
		images:   make(map[string]image.Image),
		glImages: make(map[string]*glImage),
	}
	if err := g.init(); err != nil {
		return nil, err
	}
	return g, nil
}

type graphics struct {
	images     map[string]image.Image
	glImages   map[string]*glImage
	background *glImage
	fontStash  *fontstash.Stash
	font       *glFont
}

func (g *graphics) init() error {
	if err := g.loadImages(); err != nil {
		return err
	}

	gl.ClearColor(0, 0, 0.7, 1)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	g.fontStash = fontstash.New(512, 512)
	fontID, err := g.fontStash.AddFont(resourcePath("MorrisRoman-Black.ttf"))
	if err != nil {
		return err
	}
	g.fontStash.SetYInverted(true)
	g.font = NewGLFont(g.fontStash, fontID, 45, [4]float32{0, 0, 0, 1})

	return nil
}

func (g *graphics) loadImages() error {
	imgFile, err := os.Open(resourcePath("images", "all.png"))
	if err != nil {
		return err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}

	glImage, err := NewGLImageFromImage(img)
	if err != nil {
		return err
	}

	tableFile, err := os.Open(resourcePath("images", "table.txt"))
	if err != nil {
		return err
	}
	defer tableFile.Close()
	scanner := bufio.NewScanner(tableFile)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) != 5 {
			continue
		}
		id := parts[0]
		x, _ := strconv.Atoi(parts[1])
		y, _ := strconv.Atoi(parts[2])
		w, _ := strconv.Atoi(parts[3])
		h, _ := strconv.Atoi(parts[4])
		x++
		y++
		w -= 2
		h -= 2
		g.images[id] = subImage{img, image.Rect(x, y, x+w, y+h)}
		g.glImages[id], err = glImage.SubImage(x, y, w, h)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourcePath(foldersAndFile ...string) string {
	parts := []string{
		os.Getenv("GOPATH"),
		"src", "github.com", "gonutz", "settlers",
	}
	for _, f := range foldersAndFile {
		parts = append(parts, f)
	}
	return filepath.Join(parts...)
}

func (gr *graphics) createGameBackground(g *game.Game) error {
	dest := image.NewRGBA(image.Rect(0, 0, gameW, gameH))

	images := gr.images

	var numbers [13]image.Image
	for _, n := range []int{2, 3, 4, 5, 6, 8, 9, 10, 11, 12} {
		numbers[n] = images[strconv.Itoa(n)]
	}

	for _, tile := range g.Tiles {
		x, y, w, h := tileToScreen(tile.Position)
		var img image.Image
		switch tile.Terrain {
		case game.Forest:
			img = images["tile_forest"]
		case game.Field:
			img = images["tile_field"]
		case game.Mountains:
			img = images["tile_mountains"]
		case game.Pasture:
			img = images["tile_pasture"]
		case game.Desert:
			img = images["tile_desert"]
		case game.Hills:
			img = images["tile_hills"]
		case game.Water:
			img = images["tile_water"]
		}
		draw.Draw(dest, img.Bounds().Sub(img.Bounds().Min).Add(image.Pt(x, y)),
			img, img.Bounds().Min, draw.Over)

		if tile.Terrain == game.Water && tile.Harbor.Kind != game.NoHarbor {
			id := "harbor_"
			switch tile.Harbor.Direction {
			case game.Right:
				id += "right"
			case game.TopRight:
				id += "top_right"
			case game.TopLeft:
				id += "top_left"
			case game.Left:
				id += "left"
			case game.BottomLeft:
				id += "bottom_left"
			case game.BottomRight:
				id += "bottom_right"
			}
			img := images[id]
			draw.Draw(dest, img.Bounds().Sub(img.Bounds().Min).Add(image.Pt(x, y)),
				img, img.Bounds().Min, draw.Over)

			id = "harbor_"
			switch tile.Harbor.Kind {
			case game.WoolHarbor:
				id += "wool"
			case game.LumberHarbor:
				id += "lumber"
			case game.BrickHarbor:
				id += "brick"
			case game.OreHarbor:
				id += "ore"
			case game.GrainHarbor:
				id += "grain"
			case game.ThreeToOneHarbor:
				id += "3_1"
			}
			img = images[id]
			harborX := x + (w-img.Bounds().Dx())/2
			harborY := y + (h-img.Bounds().Dy())/2
			draw.Draw(dest, img.Bounds().Sub(img.Bounds().Min).Add(image.Pt(harborX, harborY)),
				img, img.Bounds().Min, draw.Over)
		}

		if tile.Number != 0 {
			numberImg := numbers[tile.Number]
			x, y, w, h := tileToScreen(tile.Position)
			x += (w - numberImg.Bounds().Dx()) / 2
			y += (h - numberImg.Bounds().Dy()) / 2
			numberPlate := images["number_plate"]
			draw.Draw(dest,
				numberPlate.Bounds().Sub(numberPlate.Bounds().Min).Add(image.Pt(x, y)),
				numberPlate, numberPlate.Bounds().Min, draw.Over)
			draw.Draw(dest,
				numberImg.Bounds().Sub(numberImg.Bounds().Min).Add(image.Pt(x, y)),
				numberImg, numberImg.Bounds().Min, draw.Over)
		}
	}

	var err error
	gr.background, err = NewGLImageFromImage(dest)

	return err
}

func (g *graphics) drawBackground() {
	g.background.DrawAtXY(0, 0)
}

func (g *graphics) showInstruction(msg string, color game.Color) {
	var x, y, w, h float32
	textWidth, textHeight := g.font.TextSize(msg)
	const border = 25
	w, h = float32(textWidth)+2*border, 90
	x, y = (gameW-w)/2, -topBorder
	glRect(x, y, w, h, 0.5, 0.5, 1, 0.8)
	g.fontStash.BeginDraw()
	g.font.Color = gameColorToFloats(color)
	g.font.Write(msg, float64(x+border), float64(y+float32(100-textHeight)/2)+g.font.Size)
	g.fontStash.EndDraw()
}

func gameColorToFloats(c game.Color) [4]float32 {
	switch c {
	case game.Red:
		return [4]float32{0.7, 0, 0, 1}
	case game.Blue:
		return [4]float32{0, 0, 0.6, 1}
	case game.White:
		return [4]float32{1, 1, 1, 1}
	default: // orange
		return [4]float32{1, 0.5, 0, 1}
	}
}

func glRect(x, y, w, h, r, g, b, a float32) {
	gl.Disable(gl.TEXTURE_2D) // TODO do this once and remember the state
	gl.Begin(gl.QUADS)
	gl.Color4f(r, g, b, a)
	gl.Vertex2f(x, y)
	gl.Color4f(r, g, b, a)
	gl.Vertex2f(x+w, y)
	gl.Color4f(r, g, b, a)
	gl.Vertex2f(x+w, y+h)
	gl.Color4f(r, g, b, a)
	gl.Vertex2f(x, y+h)
	gl.End()
}

func (g *graphics) drawSettlementAt(x, y int, color game.Color) {
	img := g.settlement(color)
	img.DrawAtXY(x-img.Width/2, y-img.Height/2)
}

func (g *graphics) settlement(color game.Color) *glImage {
	switch color {
	case game.White:
		return g.glImages["settlement_white"]
	case game.Blue:
		return g.glImages["settlement_blue"]
	case game.Red:
		return g.glImages["settlement_red"]
	default:
		return g.glImages["settlement_orange"]
	}
}

func (g *graphics) drawHoveringSettlementAt(x, y int, color game.Color) {
	img := g.settlement(color)
	col := playerColor(color)
	col[3] = 0.6
	img.DrawColoredAtXY(x-img.Width/2, y-img.Height/2, col)
}

func playerColor(color game.Color) [4]float32 {
	switch color {
	case game.White:
		return [4]float32{1, 1, 1, 1}
	case game.Blue:
		return [4]float32{0, 0, 1, 1}
	case game.Red:
		return [4]float32{1, 0, 0, 1}
	default:
		return [4]float32{1, 0.5, 0, 1}
	}
}
