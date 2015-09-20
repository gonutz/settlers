package main

import "github.com/gonutz/fontstash.go/fontstash"

func NewGLFont(stash *fontstash.Stash, id int, size float64, color [4]float32) *glFont {
	return &glFont{stash, id, size, color}
}

type glFont struct {
	stash *fontstash.Stash
	id    int
	Size  float64
	Color [4]float32
}

func (f *glFont) Write(text string, x, y float64) {
	f.stash.DrawText(f.id, f.Size, x, y, text, f.Color)
}

func (f *glFont) TextSize(text string) (w, h int) {
	return int(f.stash.GetAdvance(f.id, f.Size, text) + 0.5), int(f.Size + 0.5)
}
