package main

import (
	"image"
	"image/draw"
	"runtime"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/kbinani/screenshot"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	bounds := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}

	h := float32(img.Bounds().Size().Y)
	w := float32(img.Bounds().Size().X)

	destImage := imaging.Blur(img, 5)
	err = imaging.Save(destImage, "out.jpg")

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Decorated, glfw.False)
	window, err := glfw.CreateWindow(int(w), int(h), "Gopherlay", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.SetKeyCallback(glfw.KeyCallback(handleClose))

	gl.ClearColor(0, 0, 0, 0)
	gl.Enable(gl.TEXTURE_2D)
	texture := loadTexture(destImage)
	defer gl.DeleteTextures(1, &texture)

	for !window.ShouldClose() {
		gl.Viewport(0, 0, int32(w), int32(h))
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		gl.Ortho(0, float64(w), float64(h), 0, 0, 100)
		gl.MatrixMode(gl.PROJECTION)

		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.LoadIdentity()

		gl.Begin(gl.QUADS)
		gl.Color3f(1, 1, 1)
		gl.TexCoord2i(0, 0)
		gl.Vertex2i(0, 0)
		gl.TexCoord2i(0, 1)
		gl.Vertex2i(0, int32(h))
		gl.TexCoord2i(1, 1)
		gl.Vertex2i(int32(w), int32(h))
		gl.TexCoord2i(1, 0)
		gl.Vertex2i(int32(w), 0)
		gl.End()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func loadTexture(img image.Image) uint32 {
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}

func handleClose(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.Destroy()
	}
}
