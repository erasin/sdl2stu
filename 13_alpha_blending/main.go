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
var gRender *sdl.Renderer

// 设定窗口
var screenWidth, screenHeight int32 = 640, 480

// 标题设定
var windowTitle = "SDL2 Tutorial tu11 "

var gModulatedTexture *LTexture
var gBackgroundTexture *LTexture

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
	newTexture, err = gRender.CreateTextureFromSurface(loadedSurface)
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
	gRender.Copy(l.mTexture, clip, &dst)
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
	gRender, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Println("无法创建渲染器 ", err)
		return err
	}

	gRender.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)

	// 检测是否支持PNG
	if i := sdlImg.Init(sdlImg.INIT_PNG); i < 0 {
		fmt.Println("图片加载器PNG失败! Error", sdlImg.GetError())
		return err
	}

	return nil
}

// 加载媒体
func loadMedia() (err error) {
	gModulatedTexture, err = newLTexture("fadeout.png")
	if err != nil {
		fmt.Println("加载前景失败")
		return
	}

	// 设定混合模式
	gModulatedTexture.SetBlendMode(sdl.BLENDMODE_BLEND)

	gBackgroundTexture, err = newLTexture("fadein.png")
	if err != nil {
		fmt.Println("加载背景失败")
		return
	}

	return nil
}

// 资源注销
func close() {
	gModulatedTexture.Free()
	gBackgroundTexture.Free()

	gRender.Destroy()
	gWindow.Destroy()

	sdlImg.Quit()
	sdl.Quit()
}

// 监听
func listen() {
	var event sdl.Event
	var running bool

	// 色差
	var r, g, b uint8 = 255, 255, 255
	// 透明度
	var alpha uint8 = 255

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				// 结束事件
				running = false
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_q:
					r += 32
				case sdl.K_w:
					g += 32
				case sdl.K_e:
					b += 32
				case sdl.K_a:
					r -= 32
				case sdl.K_s:
					g -= 32
				case sdl.K_UP:
					if (alpha + 32) > 255 {
						alpha = 255
					} else {
						alpha += 32
					}
				case sdl.K_DOWN:
					if (alpha - 32) < 0 {
						alpha = 0
					} else {
						alpha -= 32
					}
				}
			}
		}

		updateRender(r, g, b, alpha)

	}
}

// 更新处理
func updateRender(r, g, b, alpha uint8) {
	// 清空屏幕
	gRender.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	gRender.Clear()

	// 渲染背景
	gBackgroundTexture.Render(0, 0, nil)

	// 透明设定
	gModulatedTexture.SetAlpha(alpha)
	// 色调设定
	gModulatedTexture.SetColor(r, g, b)
	// 渲染前景
	gModulatedTexture.Render(0, 0, nil)

	// 更新渲染器
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

	listen()

	// sdl.Delay(5000)
	close()
}
