package main

import (
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
	"runtime"
	"time"
)

func init() {
	runtime.LockOSThread()
}

var windowW = 640
var windowH = 480

const (
	tileW           = 200
	tileH           = 212
	tileSlopeHeight = 50
	tileYOffset     = tileH - tileSlopeHeight
)

func main() {
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

	stash := fontstash.New(512, 512)
	fontID, err := stash.AddFont("./MorrisRoman-Black.ttf")
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

	g := game.New([]game.Color{game.Red, game.White, game.Blue}, rand.Int)

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

	var background image.Image
	var backImg *glImage
	go func() {
		tileImageFile, err := os.Open("./terrain_tiles.png")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer tileImageFile.Close()
		tilesImage, _, err := image.Decode(tileImageFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		gameW, gameH := 7*200, 7*tileYOffset+tileSlopeHeight
		back := image.NewRGBA(image.Rect(0, 0, gameW, gameH))
		clearToTransparent(back)
		drawGameIntoImage(back, g, tilesImage)
		background = back
	}()

	drawGame := func() {
		if backImg == nil && background != nil {
			backImg, _ = NewGLImageFromImage(background)
		}
		if backImg != nil {
			backImg.DrawAtXY(0, 0)
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
		gameW, gameH := 7.0*200, 7.0*tileYOffset+tileSlopeHeight
		gameRatio := gameW / gameH
		windowRatio := float64(windowW) / float64(windowH)
		if windowRatio > gameRatio {
			// window is wider than game => borders left and right
			border := (windowRatio*gameH - gameW) / 2
			gl.Ortho(-border, border+gameW, gameH, 0, -1, 1)
		} else {
			// window is higher than game => borders on top and bottom
			border := (gameW/windowRatio - gameH) / 2
			gl.Ortho(0, gameW, gameH+border, -border, -1, 1)
		}
		gl.ClearColor(0, 0, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		drawGame()

		if len(lines) > 0 {
			stash.BeginDraw()
			const fontSize = 35
			const cursorText = "_"
			cursor := ""
			if showCursor > 500 {
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

		showCursor = (showCursor + 1) % 1000
		frames++
		if time.Now().Sub(start).Seconds() >= 1.0 {
			fmt.Println(frames, "fps")
			frames = 0
			start = time.Now()
		}
	}
}

func removeLastRune(s string) string {
	b := 0
	for i := range s {
		b = i
	}
	return s[:b]
}

func tileToScreen(p game.TilePosition) (x, y int) {
	return p.X * tileW / 2, p.Y * tileYOffset
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
	draw.Draw(img, img.Bounds(), image.NewUniform(col.RGBA{0, 0, 0, 0}), image.ZP, draw.Src)
}

func drawGameIntoImage(dest draw.Image, g *game.Game, tileImage image.Image) {
	forest := subImage{tileImage, bounds(1, 429, 200, 212)}
	field := subImage{tileImage, bounds(1, 215, 200, 212)}
	moutains := subImage{tileImage, bounds(405, 1, 200, 212)}
	pasture := subImage{tileImage, bounds(203, 215, 200, 212)}
	hills := subImage{tileImage, bounds(203, 1, 200, 212)}
	desert := subImage{tileImage, bounds(405, 215, 200, 212)}
	water := subImage{tileImage, bounds(1, 1, 200, 212)}

	numberPlate := subImage{tileImage, bounds(607, 1, 70, 70)}
	var numbers [13]image.Image
	for i, n := range []int{2, 3, 4, 5, 6, 8, 9, 10, 11, 12} {
		numbers[n] = subImage{tileImage, bounds(i*72+1, 643, 70, 70)}
	}

	for _, tile := range g.GetTiles() {
		x := tile.Position.X * tileW / 2
		y := tile.Position.Y * tileYOffset
		var img image.Image
		switch tile.Terrain {
		case game.Forest:
			img = forest
		case game.Field:
			img = field
		case game.Mountains:
			img = moutains
		case game.Pasture:
			img = pasture
		case game.Desert:
			img = desert
		case game.Hills:
			img = hills
		case game.Water:
			img = water
		}
		draw.Draw(dest, img.Bounds().Sub(img.Bounds().Min).Add(image.Pt(x, y)),
			img, img.Bounds().Min, draw.Over)
		if tile.Number != 0 {
			numberImg := numbers[tile.Number]
			x, y := tileToScreen(tile.Position)
			x += (tileW - numberImg.Bounds().Dx()) / 2
			y += (tileH - numberImg.Bounds().Dy()) / 2
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

type subImage struct {
	image.Image
	bounds image.Rectangle
}

func (i subImage) Bounds() image.Rectangle {
	return i.bounds
}
