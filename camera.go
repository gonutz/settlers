package main

import (
	"github.com/go-gl/gl/v2.1/gl"
)

const (
	leftBorder   = 50
	rightBorder  = 50
	bottomBorder = 100
	topBorder    = 100
	gameW        = 7.0 * 200
	gameH        = 7.0*tileYOffset + tileSlopeHeight
)

func newCamera() *camera { return &camera{} }

type camera struct {
	WindowWidth, WindowHeight int
	Left, Right, Top, Bottom  float64
}

func (c *camera) windowToGame(x, y float64) (int, int) {
	xPercent := x / float64(c.WindowWidth)
	yPercent := y / float64(c.WindowHeight)
	relX := c.Left + xPercent*(c.Right-c.Left)
	relY := c.Top + yPercent*(c.Bottom-c.Top)
	return int(relX + 0.5), int(relY + 0.5)
}

func (c *camera) gameToWindow(x, y int) (float64, float64) {
	xPercent := (float64(x) - c.Left) / (c.Right - c.Left)
	yPercent := (float64(y) - c.Top) / (c.Bottom - c.Top)
	return xPercent * float64(c.WindowWidth), yPercent * float64(c.WindowHeight)
}

func (cam *camera) windowSizeChangedTo(width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	cam.WindowWidth, cam.WindowHeight = width, height
	cam.recalcOrthoBorders()
}

func (cam *camera) recalcOrthoBorders() {
	const totalW = gameW + leftBorder + rightBorder
	const totalH = gameH + topBorder + bottomBorder
	const totalRatio = float64(totalW) / totalH

	windowRatio := float64(cam.WindowWidth) / float64(cam.WindowHeight)

	var horizontalBorder, verticalBorder float64

	if windowRatio > totalRatio {
		// window is wider than game => borders left and right
		horizontalBorder = (windowRatio*totalH - totalW) / 2
	} else {
		// window is higher than game => borders on top and bottom
		verticalBorder = (totalW/windowRatio - totalH) / 2
	}

	cam.Left = -leftBorder - horizontalBorder
	cam.Right = gameW + rightBorder + horizontalBorder
	cam.Top = -topBorder - verticalBorder
	cam.Bottom = gameH + bottomBorder + verticalBorder

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(cam.Left, cam.Right, cam.Bottom, cam.Top, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
}
