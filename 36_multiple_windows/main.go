package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var total_windows int = 3

var gWindows []LWindow

func sdlinit() (err error) {
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Println(err)
		return
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	if err = gWindows[0].init(); err != nil {
		return
	}

	return nil
}

func close() {
	for i := 0; i < total_windows; i++ {
		gWindows[i].free()
	}
	sdl.Quit()
}

//
func main() {

	gWindows = make([]LWindow, total_windows)
	err := sdlinit()
	if err != nil {
		os.Exit(0)
	}

	for i := 1; i < total_windows; i++ {
		gWindows[i].init()
	}

	var event sdl.Event
	var running bool = true

	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				// 结束事件
				running = false
			}

			for i := 0; i < total_windows; i++ {
				gWindows[i].handleEvent(event)
				gWindows[i].render()
				fmt.Println(i)
			}

		}
	}

	close()

}
