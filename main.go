package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// 全局变量

var gWindow *sdl.Window
var gScreenSurface *sdl.Surface
var gKeyPressSurface [5]*sdl.Surface
var gCurrentSurface *sdl.Surface

const (
	KEY_PRESS_SURFACE_DEFAULT = iota
	KEY_PRESS_SURFACE_UP
	KEY_PRESS_SURFACE_DOWN
	KEY_PRESS_SURFACE_LEFT
	KEY_PRESS_SURFACE_RIGHT
	KEY_PRESS_SURFACE_TOTAL
)

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
	gWindow.UpdateSurface()

	return nil
}

// 加载媒体
func loadMedia() (err error) {

	gKeyPressSurface[KEY_PRESS_SURFACE_DEFAULT] = loadSurface("assets/press.bmp")
	gKeyPressSurface[KEY_PRESS_SURFACE_DOWN] = loadSurface("assets/down.bmp")
	gKeyPressSurface[KEY_PRESS_SURFACE_LEFT] = loadSurface("assets/left.bmp")
	gKeyPressSurface[KEY_PRESS_SURFACE_RIGHT] = loadSurface("assets/right.bmp")
	gKeyPressSurface[KEY_PRESS_SURFACE_UP] = loadSurface("assets/up.bmp")

	return nil
}

func loadSurface(src string) *sdl.Surface {
	var loadImg *sdl.Surface
	loadImg, err := sdl.LoadBMP(src)

	if err != nil {
		fmt.Println("加载bmp资源错误，Error：", err)
	}
	// defer loadImg.Free()

	return loadImg
}

// 监听
func listen() {
	var event sdl.Event
	var running bool

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				// 结束事件
				running = false
			case *sdl.KeyDownEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				// 键盘事件
				switch t.Keysym.Sym {
				case sdl.K_UP:
					gCurrentSurface = gKeyPressSurface[KEY_PRESS_SURFACE_UP]
				case sdl.K_DOWN:
					gCurrentSurface = gKeyPressSurface[KEY_PRESS_SURFACE_DOWN]
				case sdl.K_LEFT:
					gCurrentSurface = gKeyPressSurface[KEY_PRESS_SURFACE_LEFT]
				case sdl.K_RIGHT:
					gCurrentSurface = gKeyPressSurface[KEY_PRESS_SURFACE_RIGHT]
				default:
					gCurrentSurface = gKeyPressSurface[KEY_PRESS_SURFACE_DEFAULT]
				}

				src := sdl.Rect{X: 0, Y: 0, W: 640, H: 480}
				dst := sdl.Rect{X: 0, Y: 0, W: 640, H: 480}

				// 填充刷新
				gCurrentSurface.Blit(&src, gScreenSurface, &dst)
				gWindow.UpdateSurface()
			}
		}

	}
}

// 资源注销
func close() {
	// 资源释放
	for _, s := range gKeyPressSurface {
		s.Free()
	}

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

	listen()

	// sdl.Delay(5000)
	close()
}
