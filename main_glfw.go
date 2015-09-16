package main

import (
	"bufio"
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gonutz/fontstash.go/fontstash"
	"github.com/gonutz/settlers/game"
	"image"
	col "image/color"
	"image/draw"
	_ "image/png"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func init() {
	runtime.LockOSThread()
}

var windowW = 800
var windowH = 600

var images = make(map[string]image.Image)
var glImages = make(map[string]*glImage)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := glfw.Init(); err != nil {
		fmt.Println("glfw.Init():", err)
		return
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Decorated, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 1)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(10, 10, "Settlers", nil, nil)
	if err != nil {
		fmt.Println("glfw.CreateWindow():", err)
		return
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		fmt.Println("gl.Init():", err)
		return
	}

	var maxTextureSize int32
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxTextureSize)
	fmt.Println(maxTextureSize)

	stash := fontstash.New(512, 512)
	fontID, err := stash.AddFont(resourcePath("MorrisRoman-Black.ttf"))
	if err != nil {
		fmt.Println(err)
		return
	}
	stash.SetYInverted(true)
	font := &font{stash, fontID, 35}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	g := game.New([]game.Color{game.Red, game.White, game.Blue, game.Orange}, rand.Int)
	// NOTE for testing purposes, initialize some players
	g.Players[0].Settlements[0].Position = game.TileCorner{5, 2}
	g.Players[0].Roads[0].Position = game.TileEdge{11, 2}
	g.Players[0].Roads[1].Position = game.TileEdge{10, 2}
	g.Players[0].Roads[2].Position = game.TileEdge{11, 3}
	g.Players[1].Cities[0].Position = game.TileCorner{7, 2}
	g.Players[1].Roads[0].Position = game.TileEdge{13, 2}
	g.Players[2].Cities[0].Position = game.TileCorner{8, 1}
	g.Players[2].Roads[0].Position = game.TileEdge{16, 1}
	g.Players[2].Roads[1].Position = game.TileEdge{17, 2}
	g.Players[3].Settlements[0].Position = game.TileCorner{5, 4}
	g.Players[3].Roads[0].Position = game.TileEdge{11, 4}
	g.Players[3].Roads[1].Position = game.TileEdge{13, 4}
	g.Players[3].Roads[2].Position = game.TileEdge{15, 4}
	g.Players[3].Roads[3].Position = game.TileEdge{17, 4}
	g.Players[3].Roads[4].Position = game.TileEdge{19, 4}
	g.Players[3].Roads[5].Position = game.TileEdge{21, 4}
	g.Robber.Position = game.TilePosition{5, 4}

	var lines []string
	window.SetCharCallback(func(_ *glfw.Window, r rune) {
		if len(lines) == 0 {
			lines = []string{""}
		}
		lines[len(lines)-1] += string(r)
	})
	window.SetCharModsCallback(func(_ *glfw.Window, r rune, _ glfw.ModifierKey) {
	})
	window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
		}
		if len(lines) > 0 && key == glfw.KeyBackspace && (action == glfw.Press || action == glfw.Repeat) {
			if len(lines[len(lines)-1]) == 0 {
				lines = lines[:len(lines)-1]
			} else {
				lines[len(lines)-1] = removeLastRune(lines[len(lines)-1])
			}
		}
		if (action == glfw.Press || action == glfw.Repeat) &&
			(key == glfw.KeyEnter || key == glfw.KeyKPEnter) {
			lines = append(lines, "")
		}
	})
	window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	window.SetSizeCallback(func(_ *glfw.Window, w, h int) {
		windowW, windowH = w, h
		gl.Viewport(0, 0, int32(w), int32(h))
	})

	if err := loadImages(); err != nil {
		fmt.Println(err)
		return
	}

	var background image.Image
	var backImg *glImage
	go func() {
		gameW, gameH := 7*200, 7*tileYOffset+tileSlopeHeight
		back := image.NewRGBA(image.Rect(0, 0, gameW, gameH))
		clearToTransparent(back)
		drawGameBackgroundIntoImage(back, g)
		background = back
	}()

	gameColorToString := func(c game.Color) string {
		switch c {
		case game.Red:
			return "red"
		case game.Orange:
			return "orange"
		case game.Blue:
			return "blue"
		default:
			return "white"
		}
	}

	roadImage := func(pos game.TileEdge, c game.Color) *glImage {
		color := gameColorToString(c)
		dir := "up"
		if isEdgeVertical(pos) {
			dir = "vertical"
		}
		if isEdgeGoingDown(pos) {
			dir = "down"
		}
		return glImages["road_"+color+"_"+dir]
	}

	settlementImage := func(pos game.TileCorner, c game.Color) *glImage {
		return glImages["settlement_"+gameColorToString(c)]
	}

	cityImage := func(pos game.TileCorner, c game.Color) *glImage {
		return glImages["city_"+gameColorToString(c)]
	}

	drawGame := func() {
		if backImg == nil && background != nil {
			backImg, _ = NewGLImageFromImage(background)
		}
		if backImg != nil {
			backImg.DrawAtXY(0, 0)

			for _, p := range g.GetPlayers() {
				for _, r := range p.GetBuiltRoads() {
					x, y := edgeToScreen(r.Position)
					img := roadImage(r.Position, p.Color)
					img.DrawAtXY(x-img.Width/2, y-img.Height/2)
				}
				for _, s := range p.GetBuiltSettlements() {
					x, y := cornerToScreen(s.Position)
					img := settlementImage(s.Position, p.Color)
					img.DrawAtXY(x-img.Width/2, y-(5*img.Height/8))
				}
				for _, c := range p.GetBuiltCities() {
					x, y := cornerToScreen(c.Position)
					img := cityImage(c.Position, p.Color)
					img.DrawAtXY(x-img.Width/2, y-(5*img.Height/8))
				}
			}

			x, y, w, h := tileToScreen(g.Robber.Position)
			robber := glImages["robber"]
			robber.DrawAtXY(x+(w-robber.Width)/2, y+(h-robber.Height)/2)
		}
	}

	buildMenu := &menu{color{0.5, 0.4, 0.8, 0.8}, rect{0, 500, 800, 250}}

	showCursor := 0
	start := time.Now()
	frames := 0

	window.SetSize(windowW, windowH)

	for !window.ShouldClose() {
		glfw.PollEvents()

		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		const controlsHeight = 0 // TODO reserve are for stats and menus
		gameW, gameH := 7.0*200, 7.0*tileYOffset+tileSlopeHeight+controlsHeight
		gameRatio := gameW / gameH
		windowRatio := float64(windowW) / float64(windowH)
		var left, right, bottom, top float64
		if windowRatio > gameRatio {
			// window is wider than game => borders left and right
			border := (windowRatio*gameH - gameW) / 2
			left, right = -border, border+gameW
			bottom, top = gameH, 0
		} else {
			// window is higher than game => borders on top and bottom
			border := (gameW/windowRatio - gameH) / 2
			left, right = 0, gameW
			bottom, top = gameH+border, -border
		}
		gl.Ortho(left, right, bottom, top, -1, 1)
		gl.ClearColor(0, 0, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		drawGame()

		if len(lines) > 0 {
			stash.BeginDraw()
			const fontSize = 35
			const cursorText = "_"
			cursor := ""
			if showCursor > 60 {
				cursor = cursorText
			}
			for i, line := range lines {
				output := line
				if i == len(lines)-1 {
					output += cursor
				}
				font.Write(output, 0, float64(i+1)*fontSize)
			}
			if len(lines) == 0 {
				font.Write(cursor, 0, fontSize)
			}
			stash.EndDraw()
		}

		buildMenu.draw()

		window.SwapBuffers()

		showCursor = (showCursor + 1) % 120
		frames++
		if time.Now().Sub(start).Seconds() >= 1.0 {
			fmt.Println(frames, "fps")
			frames = 0
			start = time.Now()
		}
	}
}

func loadImages() error {
	imgFile, err := os.Open(filepath.Join("images", "all.png"))
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

	tableFile, err := os.Open(filepath.Join("images", "table.txt"))
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
		images[id] = subImage{img, bounds(x, y, w, h)}
		glImages[id], err = glImage.SubImage(x, y, w, h)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeLastRune(s string) string {
	b := 0
	for i := range s {
		b = i
	}
	return s[:b]
}

type menu struct {
	color color
	pos   rect
}

func (m *menu) draw() {
	return // TODO
	gl.Disable(gl.TEXTURE_2D)
	gl.Begin(gl.QUADS)
	gl.Color4f(m.color.r, m.color.g, m.color.b, m.color.a)
	gl.Vertex2f(m.pos.x, m.pos.y)
	gl.Color4f(m.color.r, m.color.g, m.color.b, m.color.a)
	gl.Vertex2f(m.pos.x+m.pos.w, m.pos.y)
	gl.Color4f(m.color.r, m.color.g, m.color.b, m.color.a)
	gl.Vertex2f(m.pos.x+m.pos.w, m.pos.y+m.pos.h)
	gl.Color4f(m.color.r, m.color.g, m.color.b, m.color.a)
	gl.Vertex2f(m.pos.x, m.pos.y+m.pos.h)
	gl.End()
}

type color struct{ r, g, b, a float32 }

type rect struct{ x, y, w, h float32 }

type font struct {
	stash *fontstash.Stash
	id    int
	Size  float64
}

func (f *font) Write(text string, x, y float64) {
	f.stash.DrawText(f.id, f.Size, x, y, text)
}

func clearToTransparent(img draw.Image) {
	draw.Draw(img, img.Bounds(), image.NewUniform(col.RGBA{0, 0, 0, 0}),
		image.ZP, draw.Src)
}

func drawGameBackgroundIntoImage(dest draw.Image, g *game.Game) {
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
}

func bounds(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}

func resourcePath(filename string) string {
	return filepath.Join(os.Getenv("GOPATH"),
		"src", "github.com", "gonutz", "settlers", filename)
}
