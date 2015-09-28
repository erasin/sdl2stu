package main

import "github.com/veandco/go-sdl2/sdl"

// 计时器
type LTimer struct {
	mStartTicks  uint32 // 开始时间
	mPausedTicks uint32 // 暂停时间
	mPaused      bool   // 暂停状态
	mStarted     bool   // 开始状态
}

func newLTimer() LTimer {
	return LTimer{
		mStartTicks:  0,
		mPausedTicks: 0,
		mPaused:      false,
		mStarted:     false,
	}
}

func (t *LTimer) start() {
	// 开始计时
	t.mStarted = true
	t.mPaused = false
	// 获取sdl时钟时间
	t.mStartTicks = sdl.GetTicks()
	t.mPausedTicks = 0
}

func (t *LTimer) stop() {
	t.mStarted = false
	t.mPaused = true
	t.mStartTicks = 0
	t.mPausedTicks = 0
}

func (t *LTimer) pause() {
	// 如果进行中 并 未暂停
	if t.mStarted && !t.mPaused {
		t.mPaused = true
		// 计算暂停时间
		t.mPausedTicks = sdl.GetTicks() - t.mStartTicks
		t.mStartTicks = 0
	}
}

func (t *LTimer) unpause() {
	// 如果进行中并暂停
	if t.mStarted && t.mPaused {
		t.mPaused = false
		// 重置开始时间
		t.mStartTicks = sdl.GetTicks() - t.mPausedTicks
		t.mPausedTicks = 0
	}
}

func (t *LTimer) getTicks() (time uint32) {
	if t.mStarted {
		if t.mPaused {
			time = t.mPausedTicks
		} else {
			time = sdl.GetTicks() - t.mStartTicks
		}
	}
	return
}

func (t *LTimer) isStarted() bool {
	return t.mStarted
}

func (t *LTimer) isPaused() bool {
	return t.mPaused
}
