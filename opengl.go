package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"./event"

	"./camera"
	"./window"
	mgl "github.com/go-gl/mathgl/mgl32"

	// OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 500
	height = 500
)
const (
	vertexShaderSource = `
	#version 410
	in vec3 vp;
	uniform mat4 MVP;
	out vec4 vertexColor;

	void main() {
			gl_Position = MVP * vec4(vp, 1.0);
			vertexColor = vec4(sin(vp.x), sin(vp.y), sin(vp.z), 1);
	}
` + "\x00"

	fragmentShaderSource = `
	#version 410
	out vec4 frag_colour;
	in vec4 vertexColor;
	void main() {
			frag_colour = vertexColor;
	}
` + "\x00"
)

var (
	triangle = []float32{
		-1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		0.0, 1.0, 0.0,
	}
)

type Scene struct {
	v mgl.Vec3
}

func (s *Scene) Notify(e event.Event) {
	if ke, ok := e.(event.KeyEvent); ok {
		if ke.Key == glfw.KeyUp {
			if ke.Action == 1 {
				s.v = mgl.Vec3{0, 1, 0}
			} else if ke.Action == 0 {
				s.v = mgl.Vec3{0, 0, 0}
			}
		}
		if ke.Key == glfw.KeyDown {
			if ke.Action == 1 {
				s.v = mgl.Vec3{0, -1, 0}
			} else if ke.Action == 0 {
				s.v = mgl.Vec3{0, 0, 0}
			}
		}
		if ke.Key == glfw.KeyRight {
			if ke.Action == 1 {
				s.v = mgl.Vec3{1, 0, 0}
			} else if ke.Action == 0 {
				s.v = mgl.Vec3{0, 0, 0}
			}
		}
		if ke.Key == glfw.KeyLeft {
			if ke.Action == 1 {
				s.v = mgl.Vec3{-1, 0, 0}
			} else if ke.Action == 0 {
				s.v = mgl.Vec3{0, 0, 0}
			}
		}
	}
}

func main() {
	runtime.LockOSThread()

	w := window.CreateGLFWWindow(width, height, "Titulo")
	defer w.Terminate()
	program := initOpenGL()

	vao := makeVao(triangle)

	camera := camera.CreateCamera3D(mgl.Vec3{-3, 0, -3}, mgl.Vec3{3, 0, 3})

	mvpCStr, free := gl.Strs("MVP")
	defer free()
	matrixID := gl.GetUniformLocation(program, *mvpCStr)

	eventH := event.CreateEventHandler()
	mainScene := &Scene{}
	w.SetKeyCallback(eventH.KeyCallback)
	previousTime := glfw.GetTime()
	for !w.ShouldClose() {

		time := glfw.GetTime()
		deltaT := float32(time - previousTime)
		previousTime = time

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		gl.Enable(gl.DEPTH_TEST)
		mvp := camera.GetProjectionMatrix(width, height)
		gl.UniformMatrix4fv(matrixID, 1, false, &mvp[0])

		eventH.Subscribe(mainScene)
		camera.Move(mainScene.v.Mul(deltaT))

		// camera.Rotate(1*deltaT, 0*deltaT)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

		glfw.PollEvents()
		w.SwapBuffers()
	}

}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func draw(vao uint32, w window.Window, program uint32, matrixID int32, c camera.Camera3D) {
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
