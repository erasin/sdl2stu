package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	sdlImg "github.com/veandco/go-sdl2/sdl_image"
	sdlTtf "github.com/veandco/go-sdl2/sdl_ttf"
)

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
