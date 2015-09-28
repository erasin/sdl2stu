package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	sdlImg "github.com/veandco/go-sdl2/sdl_image"
	sdlTtf "github.com/veandco/go-sdl2/sdl_ttf"
)

// 全局变量
var gWindow *sdl.Window

// 渲染器
var gRenderer *sdl.Renderer

// 设定窗口
var screenWidth, screenHeight int32 = 640, 480

// 标题设定
var windowTitle = "SDL2 Tutorial tu 23 "

// 字体
var gFont *sdlTtf.Font

// 字体纹理
var gTimeTextTexture LTexture
var gPausePromptTexture LTexture
var gStartPromptTexture LTexture

// 初始化
func sdlinit() (err error) {
	// 初始化
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Println("初始化失败 !Error:", err)
		return err
	}

	if !sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1") {
		fmt.Println("警告: Linear texture filtering not enabled!")
	}

	gWindow, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int(screenWidth), int(screenHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("window 创建失败！Error:", err)
		return err
	}
	// defer window.Destroy()

	// 渲染器 vsynced
	gRenderer, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Println("无法创建渲染器 ", err)
		return err
	}

	gRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)

	// 检测是否支持PNG
	if i := sdlImg.Init(sdlImg.INIT_PNG); i < 0 {
		fmt.Println("图片加载器PNG失败! Error", sdlImg.GetError())
		return err
	}

	// 字体支持
	if err = sdlTtf.Init(); err != nil {
		fmt.Println("字体无法初始化！Error:", err)
		return err
	}

	return nil
}

// 加载媒体
func loadMedia() (err error) {
	gTimeTextTexture = newLTexture()
	gStartPromptTexture = newLTexture()
	gPausePromptTexture = newLTexture()

	// 打开字体文件
	gFont, err = sdlTtf.OpenFont("../NotoSans-Regular.ttf", 28)
	if err != nil {
		fmt.Println("加载字体失败", err)
		return
	}

	// 字体色彩
	textColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}

	err = gStartPromptTexture.LoadFromText("Press S to Reset Start or Stop Timer.", textColor)
	if err != nil {
		return
	}

	err = gPausePromptTexture.LoadFromText("Press P to Pause or Unpause the Timer", textColor)
	if err != nil {
		return
	}

	return nil
}

// 资源注销
func close() {
	gTimeTextTexture.Free()
	gStartPromptTexture.Free()
	gPausePromptTexture.Free()

	gFont.Close()

	gRenderer.Destroy()
	gWindow.Destroy()

	sdlTtf.Quit()
	sdlImg.Quit()
	sdl.Quit()
}

// 监听
func listen() {
	var event sdl.Event
	var running bool
	var timer LTimer

	// 字体色彩
	textColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}

	// 计时器
	timer = newLTimer()

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				// 结束事件
				running = false
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_s:
					if timer.isStarted() {
						timer.stop()
					} else {
						timer.start()
					}
					break
				case sdl.K_p:
					if timer.isPaused() {
						timer.unpause()
					} else {
						timer.pause()
					}
				}
			}
		}

		updateRender(timer, textColor)

	}
}

// 更新处理
func updateRender(timer LTimer, textColor sdl.Color) {

	// 清空屏幕
	gRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	gRenderer.Clear()

	timeText := fmt.Sprintf("Seconds since start time %d", timer.getTicks())
	err := gTimeTextTexture.LoadFromText(timeText, textColor)
	if err != nil {
		fmt.Println("时间文字无法加载！Error：", err)
	}

	gStartPromptTexture.Render((screenWidth-gStartPromptTexture.mWidth)/2, 0, nil)
	gPausePromptTexture.Render((screenWidth-gPausePromptTexture.mWidth)/2, gPausePromptTexture.mHeight, nil)

	gTimeTextTexture.Render((screenWidth-gTimeTextTexture.mWidth)/2, (screenHeight-gTimeTextTexture.mHeight)/2, nil)

	// 更新渲染器
	gRenderer.Present()
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
