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

	// code from here
	gltext.IsDebug = false

	fd, err := os.Open("font/font_1_honokamin.ttf")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	// Japanese character ranges
	// http://www.rikai.com/library/kanjitables/kanji_codes.unicode.shtml
	// http://www.binaryhexconverter.com/hex-to-decimal-converter
	// 3000 - 3030 -> 12288 - 12336
	// 3040 - 309f -> 12352 - 12447
	// 30a0 - 30ff -> 12448 - 12543
	// 4e00 - 9faf -> 19968 - 40879
	// ff00 - ffef -> 65280 - 65519
	scale := fixed.Int26_6(32)
	runesPerRow := fixed.Int26_6(128)
	runeRanges := make(gltext.RuneRanges, 0)
	runeRange := gltext.RuneRange{Low: 12288, High: 12336}
	runeRanges = append(runeRanges, runeRange)
	runeRange = gltext.RuneRange{Low: 12352, High: 12447}
	runeRanges = append(runeRanges, runeRange)
	runeRange = gltext.RuneRange{Low: 12448, High: 12543}
	runeRanges = append(runeRanges, runeRange)
	runeRange = gltext.RuneRange{Low: 19968, High: 40879}
	runeRanges = append(runeRanges, runeRange)

	font, err := gltext.LoadTruetype("fontconfigs")
	if err == nil {
		fmt.Println("Font loaded from disk...")
	} else {
		font, err = gltext.NewTruetype(fd, scale, runeRanges, runesPerRow)
		if err != nil {
			panic(err)
		}
		font.Config.Save("fontconfigs")
	}

	width, height := window.GetSize()
	font.ResizeWindow(float32(width), float32(height))

	scaleMin, scaleMax := float32(1.0), float32(1.1)
	text := gltext.NewText(font, scaleMin, scaleMax)
	str := "梅干しが大好き。ウメボシガダイスキ。"
	for _, s := range str {
		fmt.Printf("%c: %d\n", s, rune(s))
	}
	text.SetString(str)

	color := float32(0)
	gl.ClearColor(0.4, 0.4, 0.4, 0.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		color += color * 0.01
		if color > 1.0 {
			color = 1
		}
		if color < 0.0 {
			color = 0
		}
		text.SetColor(color, color, color)
		text.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
	text.Release()
	font.Release()
}
