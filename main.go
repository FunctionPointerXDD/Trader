package main

import (
	"image"
	"log"

	"github.com/FunctionPointerXDD/Trader/entities"
	"github.com/hajimehoshi/ebiten/v2"
)

func CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16.0,
				int(sprite.Y)+16.0,
			),
		) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - 16.0
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16.0,
				int(sprite.Y)+16.0,
			),
		) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - 16.0
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}

func main() {

	ebiten.SetWindowSize(640, 480) // 기본 창 사이즈
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode((ebiten.WindowResizingModeEnabled)) // 전체 창 모드

	var game *Game
	game = NewGame()
	if err := ebiten.RunGame(game); err != nil { // Game이라는 구조체(이름은 상관없음) 하나를 정의해서 Update, Draw, Layout에 인터페이스 역할을 수행한다.
		log.Fatal(err)
	}
}
