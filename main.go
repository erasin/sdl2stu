package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// 全局变量

var gWindow *sdl.Window
var gScreenSurface *sdl.Surface

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
func loadMedia(i2 int) (err error) {

	var loadImg *sdl.Surface

	switch i2 {
	case 1:
		loadImg, err = sdl.LoadBMP("assets/test.bmp")
	case 2:
		loadImg, err = sdl.LoadBMP("assets/test_close.bmp")
	default:
		loadImg, err = sdl.LoadBMP("assets/test.bmp")
	}

	if err != nil {
		fmt.Println("加载bmp资源错误，Error：", err)
		return err
	}

	src := sdl.Rect{X: 0, Y: 0, W: 640, H: 480}
	dst := sdl.Rect{X: 0, Y: 0, W: 640, H: 480}

	// 刷新
	loadImg.Blit(&src, gScreenSurface, &dst)

	gWindow.UpdateSurface()

	defer loadImg.Free()

	return nil
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
				loadMedia(2)
				sdl.Delay(2000)
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			case *sdl.MouseButtonEvent:
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			case *sdl.MouseWheelEvent:
				fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y)
			case *sdl.KeyUpEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				// 键盘事件
				switch t.Keysym.Sym {
				case sdl.K_q:
					running = false
				}
			}
		}
	}
}

// 资源注销
func close() {
	gWindow.Destroy()
	sdl.Quit()
}

func main() {
	if sdlinit() != nil {
		fmt.Println("初始化失败！")
		os.Exit(0)
	}

	if loadMedia(1) != nil {
		fmt.Println("加载媒体失败！")
		os.Exit(1)
	}

	listen()

	// sdl.Delay(5000)
	close()
}
