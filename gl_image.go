package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"os"
)

type glImage struct {
	id                       uint32
	Width, Height            int
	left, top, right, bottom float32
}

func LoadGLImageFromFile(path string) (*glImage, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewGLImageFromReader(file)
}

func NewGLImageFromReader(r io.Reader) (*glImage, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return NewGLImageFromImage(img)
}

func NewGLImageFromImage(img image.Image) (*glImage, error) {
	var rgba *image.RGBA
	if asRGBA, ok := img.(*image.RGBA); ok {
		rgba = asRGBA
	} else {
		rgba = image.NewRGBA(img.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			return nil, errors.New("unsupported stride")
		}
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	}

	var tex uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Bounds().Dx()),
		int32(rgba.Bounds().Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	return &glImage{
		tex,
		img.Bounds().Dx(),
		img.Bounds().Dy(),
		0.0,
		0.0,
		1.0,
		1.0,
	}, nil
}

func (img *glImage) SubImage(x, y, w, h int) (*glImage, error) {
	if x < 0 || y < 0 || x+w > img.Width || y+h > img.Height {
		return nil, errors.New(fmt.Sprintf(
			"sub image %v,%v,%v,%v lies outside image",
			x, y, w, h))
	}
	// TODO consider deeper nesting
	fw, fh := 1.0/float32(img.Width), 1.0/float32(img.Height)
	return &glImage{
		img.id,
		w,
		h,
		float32(x) * fw,
		float32(y) * fh,
		float32(x+w) * fw,
		float32(y+h) * fh,
	}, nil
}

func (img *glImage) DrawAtXY(x, y int) {
	img.DrawColoredAtXY(x, y, [4]float32{1, 1, 1, 1})
}

func (img *glImage) DrawColoredAtXY(x, y int, color [4]float32) {
	// TODO have the state knwon globally somewhere so this does not need to be
	// called all the time
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, img.id)

	gl.Begin(gl.QUADS)

	gl.Color4f(color[0], color[1], color[2], color[3])
	gl.TexCoord2f(img.left, img.top)
	gl.Vertex2i(int32(x), int32(y))

	gl.Color4f(color[0], color[1], color[2], color[3])
	gl.TexCoord2f(img.right, img.top)
	gl.Vertex2i(int32(x+img.Width), int32(y))

	gl.Color4f(color[0], color[1], color[2], color[3])
	gl.TexCoord2f(img.right, img.bottom)
	gl.Vertex2i(int32(x+img.Width), int32(y+img.Height))

	gl.Color4f(color[0], color[1], color[2], color[3])
	gl.TexCoord2f(img.left, img.bottom)
	gl.Vertex2i(int32(x), int32(y+img.Height))

	gl.End()
}
