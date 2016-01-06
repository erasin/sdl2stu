package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	sdlImg "github.com/veandco/go-sdl2/sdl_image"
)

// 全局变量
var gWindow *sdl.Window

// 渲染器
var gRenderer *sdl.Renderer

// 设定窗口
var screenWidth, screenHeight int32 = 640, 480

// 标题设定
var windowTitle = "SDL2 Tutorial tu 19 joy"

const JOYSTICK_DEAD_ZONE int = 8000

var gGameController *sdl.Joystick

var gArrowTexture *LTexture

// 纹理对象
type LTexture struct {
	mTexture *sdl.Texture
	mWidth   int32
	mHeight  int32
}

// 加载纹理
func newLTexture(src string) (l *LTexture, err error) {
	// l.Free()

	var loadedSurface *sdl.Surface
	var newTexture *sdl.Texture

	// PNG 图片资源加载
	loadedSurface, err = sdlImg.Load(src)
	if err != nil {
		fmt.Println("加载PNG资源错误，Error：", err)
	}

	// Color key image 设置透明元素
	loadedSurface.SetColorKey(1, sdl.MapRGB(loadedSurface.Format, 0, 255, 255))
	newTexture, err = gRenderer.CreateTextureFromSurface(loadedSurface)
	if err != nil {
		fmt.Println("纹理渲染失败，Error:", err)
	}
	l = &LTexture{
		mTexture: newTexture,
		mHeight:  loadedSurface.H,
		mWidth:   loadedSurface.W,
	}
	loadedSurface.Free()
	return l, err
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
	if err = sdl.Init(sdl.INIT_AUDIO | sdl.INIT_JOYSTICK); err != nil {
		fmt.Println("初始化失败 !Error:", err)
		return err
	}
	// joy

	// 设定纹理线性过滤
	if sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1") {
		fmt.Println("警告：线性纹理过滤无法设定！")
	}

	if sdl.NumJoysticks() < 1 {
		fmt.Println("警告：没有任何 joysticks手柄连接！ ")
	} else {
		// 加载手柄
		gGameController = sdl.JoystickOpen(0)
		if gGameController == nil {
			fmt.Printf("警告！无法打开手柄控制，SDL ERROR:%s\n", sdl.GetError())
		}
	}

	gWindow, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int(screenWidth), int(screenHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("window 创建失败！Error:", err)
		return err
	}
	// defer window.Destroy()

	// 渲染器
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

	return nil
}

// 加载媒体
func loadMedia() (err error) {
	gArrowTexture, err = newLTexture("arrow.png")
	if err != nil {
		fmt.Println("加载箭头失败")
		return
	}

	return nil
}

// 资源注销
func close() {
	gArrowTexture.Free()

	// 释放手柄
	gGameController.Close()

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
			case *sdl.JoyAxisEvent:
				fmt.Println(t.Which, t.Axis, t.Value, t.Type)
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

	gArrowTexture.RenderEx((screenWidth-gArrowTexture.mWidth)/2,
		(screenHeight-gArrowTexture.mHeight)/2, nil, degrees, nil, flipType)

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
