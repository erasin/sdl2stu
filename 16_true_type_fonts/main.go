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
var windowTitle = "SDL2 Tutorial tu 15 "

// 图片纹理
var gArrowTexture *LTexture

// 全局紫婷
var gFont *sdlTtf.Font

// 字体纹理
var gTextTexture *LTexture

// 纹理对象
type LTexture struct {
	mTexture *sdl.Texture
	mWidth   int32
	mHeight  int32
}

// 默认空纹理
func newLTexture() (l *LTexture) {
	return &LTexture{
		mTexture: nil,
		mHeight:  0,
		mWidth:   0,
	}
}

func (l *LTexture) LoadFromFile(src string) (err error) {
	l.Free()

	var loadedSurface *sdl.Surface
	var newTexture *sdl.Texture

	// PNG 图片资源加载
	loadedSurface, err = sdlImg.Load(src)
	if err != nil {
		fmt.Println("加载PNG资源错误，Error：", err)
	}

	// Color key image 设置透明元素
	loadedSurface.SetColorKey(1, sdl.MapRGB(loadedSurface.Format, 0, 0xFF, 0xFF))
	newTexture, err = gRenderer.CreateTextureFromSurface(loadedSurface)
	if err != nil {
		fmt.Println("纹理渲染失败，Error:", err)
		return err
	}
	l.mTexture = newTexture
	l.mHeight = loadedSurface.H
	l.mWidth = loadedSurface.W
	loadedSurface.Free()
	return nil
}

// 加载字体纹理
func (l *LTexture) LoadFromText(str string, textColor sdl.Color) (err error) {
	l.Free()

	var textSurface *sdl.Surface
	var newTexture *sdl.Texture

	textSurface, err = gFont.RenderUTF8_Solid(str, textColor)
	if err != nil {
		fmt.Println("无法渲染字体对象！sdl ttf Error:", err, ":", sdlTtf.GetError())
		return err
	}

	newTexture, err = gRenderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		fmt.Println("纹理渲染失败，Error:", err)
		return err
	}
	l.mTexture = newTexture
	l.mWidth = textSurface.W
	l.mHeight = textSurface.H

	textSurface.Free()
	return nil
}

// 释放
func (l *LTexture) Free() {
	if l.mTexture != nil {
		l.mTexture.Destroy()
		l.mWidth = 0
		l.mHeight = 0
	}
}

// 设定调节侧才
func (l *LTexture) SetColor(r, g, b uint8) {
	l.mTexture.SetColorMod(r, g, b)
}

// 调节混合
func (l *LTexture) SetBlendMode(blending sdl.BlendMode) {
	l.mTexture.SetBlendMode(blending)
}

// 透明调节
func (l *LTexture) SetAlpha(alpha uint8) {
	l.mTexture.SetAlphaMod(alpha)
}

// 渲染 切割渲染
func (l *LTexture) Render(x, y int32, clip *sdl.Rect) {
	dst := sdl.Rect{X: x, Y: y, W: l.mWidth, H: l.mHeight}
	if !clip.Empty() {
		dst.W = clip.W
		dst.H = clip.H
	}
	gRenderer.Copy(l.mTexture, clip, &dst)
}

// 渲染 旋转和翻转
func (l *LTexture) RenderEx(x, y int32, clip *sdl.Rect, angle float64, center *sdl.Point, flip sdl.RendererFlip) {
	dst := sdl.Rect{X: x, Y: y, W: l.mWidth, H: l.mHeight}
	if !clip.Empty() {
		dst.W = clip.W
		dst.H = clip.H
	}
	// (l.mTexture, clip, &dst)
	gRenderer.CopyEx(l.mTexture, clip, &dst, angle, center, flip)
}

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
	gRenderer, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_PRESENTVSYNC)
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

	if err = sdlTtf.Init(); err != nil {
		fmt.Println("字体无法初始化！Error:", err)
		return err
	}

	return nil
}

// 加载媒体
func loadMedia() (err error) {
	gArrowTexture = newLTexture()
	err = gArrowTexture.LoadFromFile("arrow.png")
	if err != nil {
		// fmt.Println("加载箭头失败")
		return
	}

	gFont, err = sdlTtf.OpenFont("lazy.ttf", 28)
	if err != nil {
		fmt.Println("加载字体失败", err)
		return
	}

	gTextTexture = newLTexture()
	textColor := sdl.Color{R: 0, G: 0, B: 0}
	err = gTextTexture.LoadFromText("test font file is load?", textColor)
	if err != nil {
		return
	}

	return nil
}

// 资源注销
func close() {
	gArrowTexture.Free()
	gTextTexture.Free()

	gRenderer.Destroy()
	gWindow.Destroy()

	sdlImg.Quit()
	sdl.Quit()
}

// 监听
func listen() {
	var event sdl.Event
	var running bool

	// Angle of rotation 旋转角
	var degrees float64 = 0
	// 翻转设定类型
	var flipType sdl.RendererFlip = sdl.FLIP_NONE

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				// 结束事件
				running = false
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_a:
					degrees -= 60
					break
				case sdl.K_d:
					degrees += 60
					break
				case sdl.K_q:
					flipType = sdl.FLIP_HORIZONTAL
				case sdl.K_w:
					flipType = sdl.FLIP_NONE
				case sdl.K_e:
					flipType = sdl.FLIP_VERTICAL
				}
			}
		}

		updateRender(degrees, flipType)

	}
}

// 更新处理
func updateRender(degrees float64, flipType sdl.RendererFlip) {
	// 清空屏幕
	gRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	gRenderer.Clear()

	// center := &sdl.Point{X: 0, Y: 0}
	gArrowTexture.RenderEx((screenWidth-gArrowTexture.mWidth)/2,
		(screenHeight-gArrowTexture.mHeight)/2, nil, degrees, nil, flipType)

	if degrees != 0 {
		textColor := sdl.Color{R: 0, G: 0, B: 0}
		str := fmt.Sprintf("the row degress is %f", degrees)
		gTextTexture.LoadFromText(str, textColor)
	}

	gTextTexture.Render(0, 0, nil)

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
