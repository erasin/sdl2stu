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
var windowTitle = "SDL2 Tutorial tu 17 mouse "

// 图片纹理
var gArrowTexture LTexture

// 全局紫婷
var gFont *sdlTtf.Font

// 字体纹理
var gTextTexture LTexture

// 按钮大小
var buttonWidth, buttonHeight int32 = 300, 200

// 按钮数量
const totalButtons int = 4

// button clips
var gSpriteClips [totalButtons]sdl.Rect
var gButtonSpriteSheetTexture LTexture

//------------

// 纹理对象
type LTexture struct {
	mTexture *sdl.Texture
	mWidth   int32
	mHeight  int32
}

// 默认空纹理
func newLTexture() (l LTexture) {
	return LTexture{
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

type LButtonSprite int

const (
	BUTTON_SPRITE_MOUSE_OUT = iota
	BUTTON_SPRITE_MOUSE_OVER_MOTION
	BUTTON_SPRITE_MOUSE_DOWN
	BUTTON_SPRITE_MOUSE_UP
	BUTTON_SPRITE_TOTAL
)

var gButton [totalButtons]LButton

// 鼠标按钮对象
// b := LButton{ mPostion: sdl.Point{X:0,Y:0}, mCurrentSprite: BUTTON_SPRITE_MOUSE_OUT }
type LButton struct {
	mPostion       sdl.Point
	mCurrentSprite int
}

func newLButton(x, y int32) (b LButton) {
	return LButton{
		mPostion:       sdl.Point{X: x, Y: y},
		mCurrentSprite: BUTTON_SPRITE_MOUSE_OUT,
	}
}

// 设定 top,left 坐标
func (b *LButton) setPosition(x, y int32) {
	b.mPostion = sdl.Point{X: x, Y: y}
}

// 鼠标移动事件触发
func (b *LButton) handleEvent(x, y int32, t uint32) {

	var inside bool = true
	// 鼠标在左边

	switch {
	case x < b.mPostion.X:
		inside = false
		break
	case x > (b.mPostion.X + buttonWidth):
		inside = false
		break
	case y < b.mPostion.Y:
		inside = false
		break
	case y > (b.mPostion.Y + buttonHeight):
		inside = false
		break
	}

	if !inside {
		b.mCurrentSprite = BUTTON_SPRITE_MOUSE_OUT
	} else {
		switch t {
		case sdl.MOUSEMOTION:
			b.mCurrentSprite = BUTTON_SPRITE_MOUSE_OVER_MOTION
			break
		case sdl.MOUSEBUTTONDOWN:
			b.mCurrentSprite = BUTTON_SPRITE_MOUSE_DOWN
			break
		case sdl.MOUSEBUTTONUP:
			b.mCurrentSprite = BUTTON_SPRITE_MOUSE_UP
			break
		}
	}

}

func (b *LButton) render() {
	gButtonSpriteSheetTexture.Render(b.mPostion.X, b.mPostion.Y, &gSpriteClips[b.mCurrentSprite])
}

// 初始化
func sdlinit() (err error) {
	// 初始化
	if err = sdl.Init(sdl.INIT_AUDIO); err != nil {
		fmt.Println("初始化失败 !Error:", err)
		return err
	}

	// 设置线性纹理优先级
	if !sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1") {
		fmt.Println()
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
	gButtonSpriteSheetTexture = newLTexture()
	err = gButtonSpriteSheetTexture.LoadFromFile("button.png")
	if err != nil {
		fmt.Println("加载按钮失败！")
		return
	}

	// 切片设定
	for i := 0; i < BUTTON_SPRITE_TOTAL; i++ {
		gSpriteClips[i] = sdl.Rect{X: 0, Y: int32(i * 200), W: buttonWidth, H: buttonHeight}
	}

	// 设定位置
	gButton[0] = newLButton(0, 0)
	gButton[1] = newLButton(screenWidth-buttonWidth, 0)
	gButton[2] = newLButton(0, screenHeight-buttonHeight)
	gButton[3] = newLButton(screenWidth-buttonWidth, screenHeight-buttonHeight)

	fmt.Println(gButton)
	return nil
}

// 资源注销
func close() {
	gButtonSpriteSheetTexture.Free()

	gRenderer.Destroy()
	gWindow.Destroy()

	sdlImg.Quit()
	sdl.Quit()
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
			case *sdl.MouseMotionEvent:
				for i := 0; i < totalButtons; i++ {
					gButton[i].handleEvent(t.X, t.Y, t.Type)
				}
				break
			case *sdl.MouseButtonEvent:
				for i := 0; i < totalButtons; i++ {
					gButton[i].handleEvent(t.X, t.Y, t.Type)
				}
				break
			}

		}

		updateRender()

	}
}

// 更新处理
func updateRender() {
	// 清空屏幕
	gRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	gRenderer.Clear()

	for i := 0; i < totalButtons; i++ {
		gButton[i].render()
	}

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
