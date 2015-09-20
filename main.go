package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// 全局变量

var gWindow *sdl.Window
var gScreenSurface *sdl.Surface
var gHelloWorld *sdl.Surface

// 设定窗口
var screenWidth, screenHeight int = 640, 480

// 标题设定
var windowTitle = "SDL2 Tutorial tu2"

// 初始化
func sdlinit() (err error) {
	// 初始化
	if err = sdl.Init(sdl.INIT_AUDIO); err != nil {
		fmt.Println("初始化失败 !Error:", err)
		return err
	}

	gWindow, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("window 创建失败！Error:", err)
		return err
	}
	// defer window.Destroy()

	gScreenSurface, _ = gWindow.GetSurface()
	gScreenSurface.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(screenWidth), H: int32(screenHeight)}, sdl.MapRGB(gScreenSurface.Format, 173, 22, 196))
	return nil
}

// 加载媒体
func loadMedia() (err error) {

	gHelloWorld, err = sdl.LoadBMP("assets/test.bmp")
	if err != nil {
		fmt.Println("加载bmp资源错误，Error：", err)
		return err
	}

	return nil
}

// 资源注销
func close() {
	gHelloWorld.Free()
	gWindow.Destroy()
	sdl.Quit()
}

func main() {
	if sdlinit() != nil {
		fmt.Println("初始化失败！")
		os.Exit(0)
	}

	if loadMedia() != nil {
		fmt.Println("加载媒体失败！")
		os.Exit(1)
	}

	src := sdl.Rect{X: 0, Y: 0, W: 640, H: 480}
	dst := sdl.Rect{X: 0, Y: 0, W: 640, H: 480}

	// 刷新
	gHelloWorld.Blit(&src, gScreenSurface, &dst)

	gWindow.UpdateSurface()

	sdl.Delay(5000)
	close()
}
