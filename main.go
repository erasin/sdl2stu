package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/veandco/go-sdl2/sdl"
	sdlImg "github.com/veandco/go-sdl2/sdl_image"
)

// 全局变量

var gWindow *sdl.Window

var gScreenSurface *sdl.Surface
var gKeyPressSurface [5]*sdl.Surface
var gCurrentSurface *sdl.Surface

// 渲染器
var gRender *sdl.Renderer

// 纹理渲染
var gTexture *sdl.Texture

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

	// 检测是否支持PNG
	if e := sdlImg.Init(sdlImg.INIT_PNG); e < 0 {
		fmt.Println("图片加载器，吃书啊PNG失败! Error", e)
		return err
	}

	// 渲染器
	gRender, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Println("无法创建渲染器 ", err)
		return err
	}

	// gScreenSurface, _ = gWindow.GetSurface()
	// gScreenSurface.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(screenWidth), H: int32(screenHeight)}, sdl.MapRGB(gScreenSurface.Format, 173, 22, 196))
	// // gWindow.UpdateSurface()

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

// 加载贴图
func loadSurface(src string) *sdl.Surface {
	var loadedSurface *sdl.Surface
	// var optimizedSurface *sdl.Surface
	var err error

	regIsPng := regexp.MustCompile(".(png|PNG)$")

	// 判定资源
	if regIsPng.MatchString(src) {
		// PNG 图片资源加载
		loadedSurface, err = sdlImg.Load(src)
		if err != nil {
			fmt.Println("加载PNG资源错误，Error：", err)
		}
	} else {
		// BMP资源加载
		loadedSurface, err = sdl.LoadBMP(src)
		if err != nil {
			fmt.Println("加载BMP资源错误，Error：", err)
		}
	}

	// 释放旧数据
	// defer loadedSurface.Free()

	// 转换
	// optimizedSurface, err = loadedSurface.Convert(gScreenSurface.Format, 0)
	//
	// if err != nil {
	// 	fmt.Println("转换图片资源错误，Error：", err)
	// }

	return loadedSurface
}

// 加载纹理
func loadTexture(s *sdl.Surface) *sdl.Texture {
	var newTexture *sdl.Texture
	var err error
	newTexture, err = gRender.CreateTextureFromSurface(s)
	if err != nil {
		fmt.Println("纹理渲染失败，Error:", err)
	}
	return newTexture
}

// 刷新渲染器
func updateRender(s *sdl.Surface) {
	gRender.Clear()

	src := sdl.Rect{X: 0, Y: 0, W: int32(screenWidth), H: int32(screenHeight)}
	dst := sdl.Rect{X: 0, Y: 0, W: int32(screenWidth), H: int32(screenHeight)}

	// 贴图
	gRender.Copy(loadTexture(s), &src, &dst)
	gRender.Present()
}

// 刷新窗口
func updateWindow(s *sdl.Surface) {
	src := sdl.Rect{X: 0, Y: 0, W: s.W, H: s.H}
	dst := sdl.Rect{X: 0, Y: 0, W: int32(screenWidth), H: int32(screenHeight)}
	// 填充刷新
	// 缩放处理
	if s.W < gScreenSurface.W {
		s.BlitScaled(&src, gScreenSurface, &dst)
	} else {
		s.Blit(&src, gScreenSurface, &dst)
	}
	gWindow.UpdateSurface()
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
				updateRender(gCurrentSurface)
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
	gTexture.Destroy()
	gRender.Destroy()
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

	// 加载图片
	gCurrentSurface = loadSurface("assets/06_loaded.png")
	// updateWindow(gCurrentSurface)
	updateRender(gCurrentSurface)

	// 加载初始色彩
	// var rmask uint32 = 54
	// var gmask uint32 = 196
	// var bmask uint32 = 36
	// var amask uint32 = 100

	// gScreenSurface, _ = sdl.CreateRGBSurface(0, int32(screenWidth), int32(screenHeight), 32, rmask, gmask, bmask, amask)
	// updateRender(gScreenSurface)
	// defer gScreenSurface.Free()

	listen()

	// sdl.Delay(5000)
	close()
}
