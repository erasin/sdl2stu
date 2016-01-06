package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	runtime.UnlockOSThread()

	var winTtitle string = "Go-sdl2 opengl2"
	var winWidth, winHeight int = 800, 600
	var window *sdl.Window
	var context sdl.GLContext
	var event sdl.Event
	var running bool
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	window, err = sdl.CreateWindow(winTtitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	context, err = sdl.GL_CreateContext(window)
	if err != nil {
		panic(err)
	}
	defer sdl.GL_DeleteContext(context)

	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Viewport(0, 0, int32(winWidth), int32(winHeight))

	running = true

	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Printf("[%d ms] MouseMotion \t", t.Timestamp, t.Which)

			}
		}
		drawgl(winWidth, winHeight)
		sdl.GL_SwapWindow(window)
	}

}

func drawgl(w, h int) {

	ratio := w / h

	gl.Viewport(0, 0, int32(w), int32(h))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(float64(-ratio), float64(ratio), -1, 1, 1, -1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Rotatef(float32(sdl.GetTicks()*50), 0., 0., 1.)
	gl.Begin(gl.TRIANGLES)
	gl.Color3f(1., 0., 0.)
	gl.Vertex3f(-0.6, -0.4, 0.)
	gl.Color3f(0., 1., 0.)
	gl.Vertex3f(0.6, -0.4, 0.)
	gl.Color3f(0., 0., 1.)
	gl.Vertex3f(0., 0.6, 0.)
	gl.End()

}
