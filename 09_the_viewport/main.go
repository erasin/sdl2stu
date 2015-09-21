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

// 渲染器
var gRender *sdl.Renderer

// 纹理渲染
var gTexture *sdl.Texture

// 设定窗口
var screenWidth, screenHeight int32 = 640, 480

// 标题设定
var windowTitle = "SDL2 Tutorial tu9"

// 初始化
func sdlinit() (err error) {
	// 初始化
	if err = sdl.Init(sdl.INIT_AUDIO); err != nil {
		fmt.Println("初始化失败 !Error:", err)
		return err
	}

	gWindow, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int(screenWidth), int(screenHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("window 创建失败！Error:", err)
		return err
	}
	// defer window.Destroy()

	// 渲染器
	gRender, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Println("无法创建渲染器 ", err)
		return err
	}

	gRender.SetDrawColor(255, 255, 255, 100)

	// 检测是否支持PNG

	if i := sdlImg.Init(sdlImg.INIT_PNG); i < 0 {
		fmt.Println("图片加载器PNG失败! Error", sdlImg.GetError())
		return err
	}

	return nil
}

// 加载媒体
func loadMedia() (err error) {
	gTexture = loadTexture("viewport.png")
	return nil
}

// 加载纹理
func loadTexture(src string) *sdl.Texture {
	var newTexture *sdl.Texture
	var err error
	var loadedSurface *sdl.Surface

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

	newTexture, err = gRender.CreateTextureFromSurface(loadedSurface)
	if err != nil {
		fmt.Println("纹理渲染失败，Error:", err)
	}
	loadedSurface.Free()
	return newTexture
}

// 刷新渲染器
func updateRender(t *sdl.Texture) {
	gRender.Clear()

	src := sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight}
	dst := sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight}

	// 贴图
	gRender.Copy(t, &src, &dst)
	gRender.Present()
}

// 监听
func listen() {
	var event sdl.Event
	var running bool

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				// 结束事件
				running = false
			}
		}
	}
}

// 资源注销
func close() {
	gTexture.Destroy()
	gRender.Destroy()
	gWindow.Destroy()

	sdlImg.Quit()
	sdl.Quit()
}

func drawView() {

	gRender.SetDrawColor(117, 125, 124, 100)
	gRender.Clear()

	// 上左
	topLeftViewPort := sdl.Rect{X: 0, Y: 0, W: screenWidth / 2, H: screenHeight / 2}
	gRender.SetViewport(&topLeftViewPort)
	// 渲染
	gRender.Copy(gTexture, nil, nil)

	// 上右
	topRightViewPort := sdl.Rect{X: screenWidth / 2, Y: 0, W: screenWidth / 2, H: screenHeight / 2}
	gRender.SetViewport(&topRightViewPort)
	// 渲染
	gRender.Copy(gTexture, nil, nil)

	// 下方
	bottomViewPort := sdl.Rect{X: 0, Y: screenHeight / 2, W: screenWidth, H: screenHeight / 2}
	gRender.SetViewport(&bottomViewPort)
	// 渲染
	gRender.Copy(gTexture, nil, nil)

	gRender.Present()
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

	drawView()
	listen()

	// sdl.Delay(5000)
	close()
}
