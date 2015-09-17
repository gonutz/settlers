package main

type camera struct {
	WindowWidth, WindowHeight int
	Left, Right, Top, Bottom  float64
}

func (c *camera) screenToGame(x, y float64) (int, int) {
	xPercent := x / float64(c.WindowWidth)
	yPercent := y / float64(c.WindowHeight)
	relX := c.Left + xPercent*(c.Right-c.Left)
	relY := c.Top + yPercent*(c.Bottom-c.Top)
	return int(relX + 0.5), int(relY + 0.5)
}
