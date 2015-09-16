package main

import "image"

type subImage struct {
	image.Image
	bounds image.Rectangle
}

func (i subImage) Bounds() image.Rectangle {
	return i.bounds
}
