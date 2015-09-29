package main

import (
	"flag"
	"fmt"
	"github.com/gonutz/binpacker"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

var outputImageFile = flag.String("o", "./all.png", "Atlas output file name, make this a PNG.")
var binSize = flag.Int("s", 1024, "Atlas image size, the image will be s by s pixels in size.")
var tableFile = flag.String("t", "./table.txt", "Output table file containing the mappings.")

func main() {
	flag.Parse()

	if err := os.Remove(*outputImageFile); err != nil {
		fmt.Println(err)
	}

	var imagePaths []string
	walk := func(path string, _ os.FileInfo, _ error) error {
		_, filename := filepath.Split(path)
		if strings.HasSuffix(filename, ".png") {
			imagePaths = append(imagePaths, path)
		}
		return nil
	}
	err := filepath.Walk(".", walk)
	if err != nil {
		panic(err)
	}

	err = pack(imagePaths)
	if err != nil {
		panic(err)
	}
}

func pack(paths []string) error {
	packer := binpacker.New(*binSize, *binSize)
	bin := image.NewRGBA(image.Rect(0, 0, *binSize, *binSize))
	boundsTable := make(map[string]binpacker.Rect)

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return err
		}

		rect, err := packer.Insert(img.Bounds().Dx(), img.Bounds().Dy())
		if err != nil {
			return err
		}

		draw.Draw(bin, bounds(rect), img, img.Bounds().Min, draw.Src)
		boundsTable[id(path)] = rect
	}

	if err := saveTable(boundsTable, *tableFile); err != nil {
		return err
	}

	return saveImage(bin, *outputImageFile)
}

func bounds(r binpacker.Rect) image.Rectangle {
	return image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
}

func id(path string) string {
	_, filename := filepath.Split(path)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func saveTable(table map[string]binpacker.Rect, path string) error {
	if err := os.Remove(path); err != nil {
		fmt.Println(err)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for id, rect := range table {
		_, err = file.WriteString(fmt.Sprintf("%s %v %v %v %v\n",
			id, rect.X, rect.Y, rect.Width, rect.Height))
		if err != nil {
			return err
		}
	}

	return nil
}

func saveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
