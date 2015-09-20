package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// 设定窗口
var screenWidth, screenHeight int = 640, 480

// 标题设定
var windowTitle = "SDL2 Tutorial"

func main() {
	//
	var window *sdl.Window
	// var renderer *sdl.Renderer
	var screenSurface *sdl.Surface
	var err error
	// var src, dst sdl.Rect
	// var texture *sdl.Texture

	// 初始化
	if err = sdl.Init(sdl.INIT_AUDIO); err != nil {
		fmt.Println("初始化失败 !Error:", err)
		os.Exit(0)
	}

	window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("window 创建失败！Error:", err)
		os.Exit(1)
	}
	defer window.Destroy()

	screenSurface, _ = window.GetSurface()
	maprgb := sdl.MapRGB(screenSurface.Format, 173, 22, 196)
	screenSurface.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(screenWidth), H: int32(screenHeight)}, maprgb)

	window.UpdateSurface()
	//
	// // 渲染器
	// renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_PRESENTVSYNC)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "无法创建渲染器 %s\n", err)
	// 	os.Exit(2)
	// }
	// defer renderer.Destroy()
	//
	// var rmask uint32 = 0
	// var gmask uint32 = 32
	// var bmask uint32 = 32
	// var amask uint32 = 32
	// var x, y int32 = 0, 0
	//
	// screenSurface, _ = sdl.CreateRGBSurface(0, int32(32), int32(screenHeight), 32, rmask, gmask, bmask, amask)
	//
	// texture, err = renderer.CreateTextureFromSurface(screenSurface)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "加载失败texture: %s\n", err)
	// }
	// defer texture.Destroy()
	//
	// // 全屏加载
	// src = sdl.Rect{X: 0, Y: 0, W: int32(screenWidth), H: int32(screenHeight)}
	// dst = sdl.Rect{X: x, Y: y, W: int32(screenWidth), H: int32(screenHeight)}
	//
	// renderer.Clear()
	// // renderer.SetDrawColor(255, 0, 0, 0)
	// // renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(screenWidth), H: int32(screenHeight)})
	// renderer.Copy(texture, &src, &dst)
	// renderer.Present()

	sdl.Delay(5000)
}
