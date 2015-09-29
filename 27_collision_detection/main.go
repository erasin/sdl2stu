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

// 限定帧数
const SCREEN_FPS uint32 = 60
const SCREEN_TICK_PER_FRAME uint32 = 1000 / SCREEN_FPS

// 标题设定
var windowTitle = "SDL2 Tutorial tu 23 "

// 字体
var gFont *sdlTtf.Font

// 字体纹理
var gFPSTextTexture LTexture

// 字体色彩
var textColor sdl.Color = sdl.Color{R: 0, G: 0, B: 0, A: 255}

// 渲染点
var gDotTexture LTexture

var wall sdl.Rect

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
	// 初始化
	gFPSTextTexture = newLTexture()
	gDotTexture = newLTexture()

	wall = sdl.Rect{X: 300, Y: 40, W: 40, H: 400}

	err = gDotTexture.LoadFromFile("dot.bmp")
	if err != nil {
		fmt.Println("加载dot.bmp失败", err)
		return
	}

	// 打开字体文件
	gFont, err = sdlTtf.OpenFont("../NotoSans-Regular.ttf", 28)
	if err != nil {
		fmt.Println("加载字体失败", err)
		return
	}

	return nil
}

// 资源注销
func close() {
	gFPSTextTexture.Free()
	gDotTexture.Free()

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
	var fpsTimer LTimer // FPS 计时器
	var capTimer LTimer // 片段时间

	// 计时器
	fpsTimer = newLTimer()
	capTimer = newLTimer()

	// 开始统计每秒多少帧
	var countedFrames float32 = 0
	fpsTimer.start()

	var dot Dot = newDot()

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				// 结束事件
				running = false
			}
			// dot处理
			dot.handleEvent(event)
		}

		var avgFPS float32 = countedFrames / (float32(fpsTimer.getTicks()) / 1000)
		if avgFPS > 2000000 {
			avgFPS = 0
		}

		updateRender(avgFPS, &dot)
		countedFrames = countedFrames + 1

		// 如果帧数超过（提早完成）则延时
		frameTicks := capTimer.getTicks()
		if frameTicks < SCREEN_TICK_PER_FRAME {
			sdl.Delay(SCREEN_TICK_PER_FRAME - frameTicks)
		}
	}
}

// 更新处理
func updateRender(avgFPS float32, dot *Dot) {

	// 清空屏幕
	gRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	var err error = gRenderer.Clear()
	if err != nil {
		fmt.Println("clear err:", err)
	}

	timeText := fmt.Sprintf("FPS: %.2f", avgFPS)
	err = gFPSTextTexture.LoadFromText(timeText, textColor)
	if err != nil {
		fmt.Println("时间文字无法加载！Error：", err)
	}

	gFPSTextTexture.Render(0, 0, nil)

	gRenderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	gRenderer.DrawRect(&wall)

	// 渲染dot
	dot.move(&wall)
	dot.render()

	// 更新渲染器
	gRenderer.Present()
	gFPSTextTexture.Free()

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
