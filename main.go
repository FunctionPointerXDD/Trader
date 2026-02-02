package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

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
