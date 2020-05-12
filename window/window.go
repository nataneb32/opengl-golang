package window

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window interface {
	Terminate()
	SwapBuffers()
	ShouldClose() bool
	SetKeyCallback(c glfw.KeyCallback)
}

type GLFWWindow struct {
	window *glfw.Window
}

func (w *GLFWWindow) Terminate() {
	glfw.Terminate()
}

func (w *GLFWWindow) SwapBuffers() {
	w.window.SwapBuffers()
}

func (w *GLFWWindow) ShouldClose() bool {
	return w.window.ShouldClose()
}

func (w *GLFWWindow) SetKeyCallback(c glfw.KeyCallback) {
	w.window.SetKeyCallback(c)
}

func CreateGLFWWindow(width, height int, title string) Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	return &GLFWWindow{window}
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action,
	mods glfw.ModifierKey) {
	if action == glfw.Press {

	}
}
