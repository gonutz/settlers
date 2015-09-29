package main

type rect struct{ x, y, w, h int }

func (r rect) moveBy(dx, dy int) rect {
	return rect{r.x + dx, r.y + dy, r.w, r.h}
}

func (r rect) contains(x, y int) bool {
	return x >= r.x && y >= r.y && x < r.x+r.w && y < r.y+r.h
}

func (r rect) center() (int, int) {
	return r.x + r.w/2, r.y + r.h/2
}

func (r rect) expandBy(margin int) rect {
	return rect{r.x - margin, r.y - margin, r.w + 2*margin, r.h + 2*margin}
}
