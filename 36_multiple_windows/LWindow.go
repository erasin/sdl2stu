package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  = 640
	SCREEN_HEIGHT = 480
)

type LWindow struct {
	mWindow   *sdl.Window
	mRenderer *sdl.Renderer
	mWindowID int

	// window dimensions
	mWidth  int
	mHeight int

	// window focus
	mMouseFocus    bool
	mKeyboardFocus bool
	mFullScreen    bool
	mMinimized     bool
	mShown         bool
}

func newLWindow() LWindow {
	return LWindow{
		mWindow:   nil,
		mRenderer: nil,
		mWindowID: -1,

		mWidth:  0,
		mHeight: 0,

		mMouseFocus:    false,
		mKeyboardFocus: false,
		mFullScreen:    false,
		mShown:         false,
		// mMinimized:     false,
	}
}

// 创建
func (w *LWindow) init() (err error) {

	// 窗口创建
	w.mWindow, err = sdl.CreateWindow("SDL WINDOWS", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_UNDEFINED, SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return
	}

	w.mMouseFocus = true
	w.mKeyboardFocus = true
	w.mWidth = SCREEN_WIDTH
	w.mHeight = SCREEN_HEIGHT

	// 纹理器
	w.mRenderer, err = sdl.CreateRenderer(w.mWindow, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		w.mWindow.Destroy()
		return
	}

	w.mRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	w.mWindowID = int(w.mWindow.GetID())

	return nil
}

func (w *LWindow) handleEvent(e sdl.Event) {
	switch t := e.(type) {
	case *sdl.WindowEvent:
		// Caption 解说词更新
		updateCaption := false
		switch t.Event {
		case sdl.WINDOWEVENT_SHOWN:
			w.mShown = true
			break
		case sdl.WINDOWEVENT_HIDDEN:
			w.mShown = false
			break
		case sdl.WINDOWEVENT_SIZE_CHANGED:
			w.mWidth = int(t.Data1)
			w.mHeight = int(t.Data2)
			w.mRenderer.Present()
			break
		case sdl.WINDOWEVENT_EXPOSED:
			w.mRenderer.Present()
			break
		case sdl.WINDOWEVENT_ENTER:
			w.mMouseFocus = true
			updateCaption = true
			break
		case sdl.WINDOWEVENT_LEAVE:
			w.mMouseFocus = false
			updateCaption = true
			break
		case sdl.WINDOWEVENT_FOCUS_GAINED:
			w.mKeyboardFocus = true
			updateCaption = true
			break
		case sdl.WINDOWEVENT_FOCUS_LOST:
			w.mKeyboardFocus = false
			updateCaption = true
			break
		case sdl.WINDOWEVENT_MINIMIZED:
			w.mMinimized = true
			break
		case sdl.WINDOWEVENT_MAXIMIZED:
			w.mMinimized = false
			break
		case sdl.WINDOWEVENT_RESTORED:
			w.mMinimized = false
			break
		case sdl.WINDOWEVENT_CLOSE:
			w.mWindow.Hide()
			break
		}
		if updateCaption {
			fmt.Printf("SDL WINDOWS:\n\t ID : %d\n\t MouseFocus:%t\n\t KeyboardFocus:%t", w.mWindowID, w.mMouseFocus, w.mKeyboardFocus)
			title := fmt.Sprintf("SDL WINDOWS - ID : %d  MouseFocus:%t  KeyboardFocus:%t", w.mWindowID, w.mMouseFocus, w.mKeyboardFocus)
			w.mWindow.SetTitle(title)
		}
	}
}

func (w *LWindow) focus() {
	if !w.mShown {
		w.mWindow.Show()
	}
	// 移动窗口到当前
	w.mWindow.Raise()
}

func (w *LWindow) render() {
	if !w.mMinimized {
		w.mRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
		w.mRenderer.Clear()
		w.mRenderer.Present()
	}
}

func (w *LWindow) free() {
	if w.mWindow != nil {
		w.mWindow.Destroy()
	}
	w.mMouseFocus = false
	w.mKeyboardFocus = false
	w.mWidth = 0
	w.mHeight = 0
}

func (w *LWindow) getWidth() int {
	return w.mWidth
}

func (w *LWindow) getHeight() int {
	return w.mHeight
}

func (w *LWindow) hasMouseFocus() bool {
	return w.mMouseFocus
}

func (w *LWindow) hasKeyboardFocus() bool {
	return w.mKeyboardFocus
}

func (w *LWindow) isMinimized() bool {
	return w.mMinimized
}

func (w *LWindow) isShown() bool {
	return w.mShown
}
