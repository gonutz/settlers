package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gonutz/settlers/game"
	"math/rand"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		fmt.Println("glfw.Init():", err)
		return
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 1)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(640, 480, "Settlers", nil, nil)
	if err != nil {
		fmt.Println("glfw.CreateWindow():", err)
		return
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		fmt.Println("gl.Init():", err)
		return
	}

	game.New([]game.Color{game.Red, game.White, game.Blue}, rand.Int)

	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.ClearColor(1, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		window.SwapBuffers()
	}
}
