package main

import "github.com/veandco/go-sdl2/sdl"

// 计时器
type Dot struct {
	//The dimensions of the dot
	width  int
	height int
	// Maximum axis velocity of the dot

	vel int

	mPosX, mPosY int
	mVelX, mVelY int
}

func newDot() Dot {
	return Dot{
		width:  20,
		height: 20,
		vel:    10,
	}
}

// 键盘事件 调整速度
func (dot *Dot) handleEvent(e sdl.Event) {
	switch t := e.(type) {
	case *sdl.KeyDownEvent:
		if t.Repeat == 0 {
			switch t.Keysym.Sym {
			case sdl.K_UP:
				dot.mVelY -= dot.vel
				break
			case sdl.K_DOWN:
				dot.mVelY += dot.vel
				break
			case sdl.K_LEFT:
				dot.mVelX -= dot.vel
				break
			case sdl.K_RIGHT:
				dot.mVelX += dot.vel
				break
			}
		}
		break
	case *sdl.KeyUpEvent:
		if t.Repeat == 0 {
			switch t.Keysym.Sym {
			case sdl.K_UP:
				dot.mVelY += dot.vel
				break
			case sdl.K_DOWN:
				dot.mVelY -= dot.vel
				break
			case sdl.K_LEFT:
				dot.mVelX += dot.vel
				break
			case sdl.K_RIGHT:
				dot.mVelX -= dot.vel
				break
			}
		}
	}
}

// 移动
func (dot *Dot) move() {

	// 移动坐标
	dot.mPosX += dot.mVelX

	// 越界处理
	if dot.mPosX < 0 || dot.mPosX+dot.width > int(screenWidth) {
		// 返回
		dot.mPosX -= dot.mVelX
	}

	// 移动
	dot.mPosY += dot.mVelY

	if dot.mPosY < 0 || dot.mPosY+dot.height > int(screenHeight) {
		dot.mPosY -= dot.mVelY
	}

}

// 渲染
func (dot *Dot) render() {
	gDotTexture.Render(int32(dot.mPosX), int32(dot.mPosY), nil)
}
