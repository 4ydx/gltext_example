package main

import (
	"fmt"
	gltext "github.com/4ydx/gltext"
	glfw "github.com/go-gl/glfw3"
	gl32 "github.com/go-gl/glow/gl-core/3.2/gl"
	"github.com/go-gl/glow/gl-core/3.3/gl"
	"os"
	"runtime"
)

var useStrictCoreProfile = (runtime.GOOS == "darwin")

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func findCenter(windowWidth int, windowHeight int, x1, x2 gltext.Point) (lowerLeft gltext.Point) {
	widthHalf := windowWidth / 2
	heightHalf := windowHeight / 2

	lineWidthHalf := (x2.X - x1.X) / 2
	lineHeightHalf := (x2.Y - x1.Y) / 2

	lowerLeft.X = float32(widthHalf) - lineWidthHalf
	lowerLeft.Y = float32(heightHalf) - lineHeightHalf
	return
}

func main() {
	runtime.LockOSThread()

	glfw.SetErrorCallback(errorCallback)
	if !glfw.Init() {
		panic("glfw error")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	if useStrictCoreProfile {
		glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	}
	glfw.WindowHint(glfw.OpenglDebugContext, glfw.True)

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	if err := gl32.Init(); err != nil {
		fmt.Println("could not initialize GL 3.2")
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("Opengl version", version)

	fd, err := os.Open("font/luximr.ttf")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	scale := int32(32)
	line, err := gltext.LoadTruetype(fd, scale, 32, 127)
	if err != nil {
		panic(err)
	}
	width, height := window.GetSize()

	line.ResizeWindow(float32(width), float32(height))

	str := ">-- I am Batman --<"
	x1, x2 := line.SetString(str)

	// find the center of the string based on the bounding box
	fmt.Printf("bounding box %v %v\n", x1, x2)
	lowerLeft := findCenter(width, height, x1, x2)
	line.SetPosition(lowerLeft.X, lowerLeft.Y)
	line.SetTextLowerBound(0.5)
	// alpha == 1 means no opacity
	line.SetColor(1, 1, 1, 1)

	//gl.ClearColor(1, 1, 1, 0.0)
	flow := float32(1)
	color := float32(0.0)
	gl.ClearColor(0, 1, 1, 0.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		color += flow * 0.01
		if color > 1.0 {
			color = 1
			flow = -1
		}
		if color < 0.0 {
			color = 0
			flow = +1
		}
		line.SetColor(color, color, color, 1)
		line.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
	line.Release()
}
