package camera

import "math"

type Camera struct {
	X, Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

// Camera인스턴스에서 사용가능한 리시버 함수 (go 언어 기능)
func (c *Camera) FollowTarget(targetX, targetY, screenWidth, screenHeight float64) {
	/** 카메라는 실제로 존재하지 않는다.
	배경자체를 카메라 오프셋 만큼 타켓이 이동하는 방향의 반대 방향으로 움직이면 카메라가 이동하는 효과를 얻을 수 있다. */
	c.X = -targetX + screenWidth/2.0
	c.Y = -targetY + screenHeight/2.0
}

/* 카메라가 배경 밖으로 벗어나지 않게 해주는 함수*/
func (c *Camera) Constrain(tilemapWidthPixels, tilemapHeightPixels, screenWidth, screenHeight float64) {
	c.X = math.Min(c.X, 0.0)
	c.Y = math.Min(c.Y, 0.0)

	c.X = math.Max(c.X, screenWidth-tilemapWidthPixels)
	c.Y = math.Max(c.Y, screenHeight-tilemapHeightPixels)
}
