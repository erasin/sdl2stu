package main

import "github.com/veandco/go-sdl2/sdl"

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
	mColliders []*sdl.Rect
}

func newDot(x, y int) (d Dot) {

	mColliders := make([]*sdl.Rect, 11)

	mColliders[0] = &sdl.Rect{X: 0, Y: 0, W: 6, H: 1}
	mColliders[1] = &sdl.Rect{X: 0, Y: 0, W: 10, H: 1}
	mColliders[2] = &sdl.Rect{X: 0, Y: 0, W: 14, H: 1}
	mColliders[3] = &sdl.Rect{X: 0, Y: 0, W: 16, H: 2}
	mColliders[4] = &sdl.Rect{X: 0, Y: 0, W: 18, H: 2}
	mColliders[5] = &sdl.Rect{X: 0, Y: 0, W: 20, H: 6}
	mColliders[6] = &sdl.Rect{X: 0, Y: 0, W: 18, H: 2}
	mColliders[7] = &sdl.Rect{X: 0, Y: 0, W: 16, H: 2}
	mColliders[8] = &sdl.Rect{X: 0, Y: 0, W: 14, H: 1}
	mColliders[9] = &sdl.Rect{X: 0, Y: 0, W: 10, H: 1}
	mColliders[10] = &sdl.Rect{X: 0, Y: 0, W: 6, H: 1}

	d = Dot{
		width:  20,
		height: 20,
		vel:    10,

		mPosX: x,
		mPosY: y,
		mVelX: 0,
		mVelY: 0,

		mColliders: mColliders,
	}

	d.shiftColliders()
	return d
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
func (dot *Dot) move(otherColliders []*sdl.Rect) {

	// 移动坐标
	dot.mPosX += dot.mVelX
	dot.shiftColliders()

	// 越界处理
	if dot.mPosX < 0 || dot.mPosX+dot.width > int(screenWidth) || checkCollision(dot.mColliders, otherColliders) {
		// 返回
		dot.mPosX -= dot.mVelX
		dot.shiftColliders()
	}

	// 移动
	dot.mPosY += dot.mVelY
	dot.shiftColliders()

	if dot.mPosY < 0 || dot.mPosY+dot.height > int(screenHeight) || checkCollision(dot.mColliders, otherColliders) {
		dot.mPosY -= dot.mVelY
		dot.shiftColliders()
	}

}

// 渲染
func (dot *Dot) render() {
	gDotTexture.Render(int32(dot.mPosX), int32(dot.mPosY), nil)
}

// 相对点偏移碰撞对象
func (dot *Dot) shiftColliders() {
	// 行偏移
	var r int32 = 0

	// Go through the dot's collision boxe
	for i := 0; i < len(dot.mColliders); i += 1 {
		// center the collision box
		dot.mColliders[i].X = int32(dot.mPosX) + (int32(DOT_WIDTH)-dot.mColliders[i].W)/2
		// set the collision box at its row offet
		dot.mColliders[i].Y = int32(dot.mPosY) + r
		// move the row offset down the height of the collision box
		r += dot.mColliders[i].H
	}
}

func (dot *Dot) getColliders() []*sdl.Rect {
	return dot.mColliders
}

// 碰撞探测
//Box collision detector
func checkCollision(a []*sdl.Rect, b []*sdl.Rect) (ok bool) {
	// 设定矩形边
	var leftA, leftB, rightA, rightB, topA, topB, bottomA, bottomB int32

	for Abox := 0; Abox < len(a); Abox++ {

		leftA, rightA, topA, bottomA = a[Abox].X, a[Abox].X+a[Abox].W, a[Abox].Y, a[Abox].Y+a[Abox].H

		for Bbox := 0; Bbox < len(b); Bbox++ {
			leftB, rightB, topB, bottomB = b[Bbox].X, b[Bbox].X+b[Bbox].W, b[Bbox].Y, b[Bbox].Y+b[Bbox].H
		}

		if !(bottomA <= topB || topA >= bottomB || rightA <= leftB || leftA >= rightB) {
			return true
		}
	}

	//If neither set of collision boxes touche
	return false

}
