package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"runtime"
	"time"
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

	window.SetKeyCallback(keyCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetCursorEnterCallback(cursorEnterCallback)
	window.SetCursorPosCallback(cursorPositionCallback)
	window.SetCharCallback(charCallback)
	window.SetSizeCallback(sizeCallback)

	ui, err = NewGameUI(uiWindow{window})
	if err != nil {
		fmt.Println("NewGameUI():", err)
		return
	}

	window.SetSize(640, 480)
	lastUpdate := time.Now().Add(-time.Hour)
	const frameTimeInSeconds = 1.0 / 60
	for !window.ShouldClose() {
		glfw.PollEvents()
		now := time.Now()
		if now.Sub(lastUpdate).Seconds() > frameTimeInSeconds {
			lastUpdate = now
			ui.Draw()
			window.SwapBuffers()
		}
	}
}

type uiWindow struct {
	*glfw.Window
}

func (w uiWindow) Close() {
	w.Window.SetShouldClose(true)
}

var ui *gameUI

func keyCallback(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	if action == glfw.Press {
		ui.KeyDown(key)
	}
}

func mouseButtonCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	if action == glfw.Press {
		ui.MouseButtonDown(button)
	}
}

func cursorEnterCallback(_ *glfw.Window, entered bool) {
	if entered {
		ui.MouseExited()
	} else {
		ui.MouseEntered()
	}
}

func cursorPositionCallback(_ *glfw.Window, xpos float64, ypos float64) {
	ui.MouseMovedTo(xpos, ypos)
}

func charCallback(_ *glfw.Window, char rune) {
	ui.RuneTyped(char)
}

func sizeCallback(_ *glfw.Window, width, height int) {
	ui.WindowSizeChangedTo(width, height)
}
