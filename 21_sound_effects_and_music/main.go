package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	sdlImg "github.com/veandco/go-sdl2/sdl_image"
	sdlmixer "github.com/veandco/go-sdl2/sdl_mixer"
	sdlTtf "github.com/veandco/go-sdl2/sdl_ttf"
)

// 全局变量
var gWindow *sdl.Window

// 渲染器
var gRenderer *sdl.Renderer

// 设定窗口
var screenWidth, screenHeight int32 = 640, 480

// 标题设定
var windowTitle = "SDL2 Tutorial tu 21 sound "

// 场景纹理
var gPromptTexture LTexture

// 音乐对象
var gMusic *sdlmixer.Music

// 音乐效果
var gScratch *sdlmixer.Chunk
var gHigh *sdlmixer.Chunk
var gMedium *sdlmixer.Chunk
var gLow *sdlmixer.Chunk

// 全局字体
var gFont *sdlTtf.Font

// 字体纹理
var gTextTexture LTexture

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

// 初始化
func sdlinit() (err error) {
	// 初始化 音频视频
	if err = sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO); err != nil {
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
	gRenderer, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Println("无法创建渲染器 ", err)
		return err
	}

	// 设定渲染器默认色彩
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

	if err = sdlmixer.OpenAudio(44100, sdlmixer.DEFAULT_FORMAT, 2, 2048); err != nil {
		fmt.Println("mixer 初始化失败，ERROR：", err)
		return err
	}

	return nil
}

// 加载媒体
func loadMedia() (err error) {

	// 加载纹理
	gPromptTexture = newLTexture()
	gPromptTexture.LoadFromFile("prompt.png")

	gMusic, err = sdlmixer.LoadMUS("beat.wav")
	if err != nil {
		fmt.Println("加载音频 beat.wav 失败！ERROR:", err)
		return err
	}

	gScratch, err = sdlmixer.LoadWAV("scratch.wav")
	if err != nil {
		fmt.Println("加载效果音 scratch 失败，Error:", err)
		return err
	}
	gHigh, err = sdlmixer.LoadWAV("high.wav")
	if err != nil {
		fmt.Println("加载效果音 high 失败，Error:", err)
		return err
	}
	gMedium, err = sdlmixer.LoadWAV("medium.wav")
	if err != nil {
		fmt.Println("加载效果音 medium 失败，Error:", err)
		return err
	}
	gLow, err = sdlmixer.LoadWAV("low.wav")

	if err != nil {
		fmt.Println("加载效果音 low 失败，Error:", err)
		return err
	}

	return nil
}

// 资源注销
func close() {
	gPromptTexture.Free()

	// 释放音频
	gMusic.Free()
	gScratch.Free()
	gHigh.Free()
	gMedium.Free()
	gLow.Free()

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
		// 轮询未完成的事件
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				// 结束事件
				fmt.Println("结束")
				running = false
				goto L
			case *sdl.KeyDownEvent:
				// 按键后的触发
				fmt.Printf("key:%v,state:%v\n", t.Keysym.Sym, t.State)
				switch t.Keysym.Sym {
				case sdl.K_1:
					// 1 播放特效
					gHigh.Play(-1, 0)
					break
				case sdl.K_2:
					gMedium.Play(-1, 0)
					break
				case sdl.K_3:
					gLow.Play(-1, 0)
					break
				case sdl.K_4:
					gScratch.Play(-1, 0)
					break
				case sdl.K_9:
					// 如果没有播放则播放gMusic
					if !sdlmixer.PlayingMusic() {
						// -1 无限播放 0 单播放
						gMusic.Play(0)
					} else {
						// 如果暂停中则恢复,否则暂停
						if sdlmixer.Paused(-1) == 1 {
							sdlmixer.Resume(-1)
						} else {
							sdlmixer.Pause(-1)
						}
					}
					break
				case sdl.K_0:
					// 停止所有
					// sdlmixer.HaltChannel(-1)
					sdlmixer.HaltMusic()
				}
			}
		}

		updateRender()

	}
L:
	fmt.Println("结束。。。")
}

// 更新处理
func updateRender() {

	// 清空屏幕
	gRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	gRenderer.Clear()

	gPromptTexture.Render(0, 0, nil)

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
