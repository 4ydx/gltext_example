package main

import (
	"fmt"
	"github.com/4ydx/gltext"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"golang.org/x/image/math/fixed"
	"os"
	"runtime"
)

var useStrictCoreProfile = (runtime.GOOS == "darwin")

func main() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic("glfw error")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	if useStrictCoreProfile {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	}
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("Opengl version", version)

	fd, err := os.Open("font/luximr.ttf")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	scale := fixed.Int26_6(32)
	font, err := gltext.LoadTruetype(fd, scale, 32, 127)
	if err != nil {
		panic(err)
	}
	width, height := window.GetSize()

	font.ResizeWindow(float32(width), float32(height))

	text := gltext.LoadText(font)
	str := "ABCDEFG"
	text.SetString(str)

	xPos := float32(-width)
	flow := float32(1)
	color := float32(0)
	gl.ClearColor(0.4, 0.4, 0.4, 0.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		color += flow * 0.01
		if color > 1.0 {
			color = 1
			flow = -1
		}
		if color < 0.0 {
			color = 0
			flow = +1
		}
		xPos = flow * float32(width) * color
		text.SetPosition(xPos, 0)
		text.SetColor(color, color, color)
		text.SetScale(color + 0.5)
		text.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
	text.Release()
	font.Release()
}
