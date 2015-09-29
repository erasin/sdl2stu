package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	DOT_WIDTH  int = 20
	DOT_HEIGHT int = 20
)

// 计时器
type Dot struct {
	//The dimensions of the dot
	width  int
	height int
	// Maximum axis velocity of the dot

	vel int

	mPosX, mPosY int
	mVelX, mVelY int

	// 碰撞矩形
	mCollider sdl.Rect
}

func newDot() Dot {
	return Dot{
		width:  20,
		height: 20,
		vel:    10,

		mPosX: 0,
		mPosY: 0,
		mVelX: 0,
		mVelY: 0,

		mCollider: sdl.Rect{X: 0, Y: 0, W: int32(DOT_WIDTH), H: int32(DOT_HEIGHT)},
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
func (dot *Dot) move(wall *sdl.Rect) {

	// 移动坐标
	dot.mPosX += dot.mVelX
	dot.mCollider.X = int32(dot.mPosX)

	// 越界处理
	if dot.mPosX < 0 || dot.mPosX+dot.width > int(screenWidth) || checkCollision(dot.mCollider, wall) {
		// 返回
		dot.mPosX -= dot.mVelX
		dot.mCollider.X = int32(dot.mPosX)
	}

	// 移动
	dot.mPosY += dot.mVelY
	dot.mCollider.Y = int32(dot.mPosY)

	if dot.mPosY < 0 || dot.mPosY+dot.height > int(screenHeight) || checkCollision(dot.mCollider, wall) {
		dot.mPosY -= dot.mVelY
		dot.mCollider.Y = int32(dot.mPosY)
	}

}

// 渲染
func (dot *Dot) render() {
	gDotTexture.Render(int32(dot.mPosX), int32(dot.mPosY), nil)
}

// 碰撞探测
//Box collision detector
func checkCollision(a sdl.Rect, b *sdl.Rect) (ok bool) {
	// 设定矩形边
	var leftA, leftB, rightA, rightB, topA, topB, bottomA, bottomB int32

	leftA, rightA, topA, bottomA = a.X, a.X+a.W, a.Y, a.Y+a.H
	leftB, rightB, topB, bottomB = b.X, b.X+b.W, b.Y, b.Y+b.H
	// fmt.Println(leftA, rightA, topA, bottomA)
	// fmt.Println(leftB, rightB, topB, bottomB)

	// 如果A边在B内部
	if bottomA <= topB {
		fmt.Printf("bottomA: %d, topB: %d\n", bottomA, topB)
		return false
	}

	if topA >= bottomB {
		fmt.Printf("topA: %d, bottomB: %d\n", topA, bottomB)
		return false
	}

	if rightA <= leftB {
		fmt.Printf("rightA: %d, leftB: %d\n", rightA, leftB)
		return false
	}
	if leftA >= rightB {
		fmt.Printf("leftA: %d, rightB: %d\n", leftA, rightB)
		return false
	}

	return true
}
