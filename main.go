package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/FunctionPointerXDD/Trader/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
}

func (g *Game) Update() error {

	// react to key presses
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2
	}

	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer {
			if sprite.X < g.player.X {
				sprite.X += 0.5
			} else if sprite.X > g.player.X {
				sprite.X -= 0.5
			}
			if sprite.Y < g.player.Y {
				sprite.Y += 0.5
			} else if sprite.Y > g.player.Y {
				sprite.Y -= 0.5
			}
		}
	}

	for _, potion := range g.potions {
		// 플레이어와 포션 이미지가 겹치는지 확인 (16x16 크기 기준)
		if !potion.isUsed && g.player.X > potion.X-16.0 && g.player.X < potion.X+16.0 &&
			g.player.Y > potion.Y-16.0 && g.player.Y < potion.Y+16.0 {
			g.player.Health += potion.AmtHeal
			log.Printf("Picked up potion!. Health: %d\n", g.player.Health)
			potion.isUsed = true
			// 한 번 먹은 포션은 제거하거나 비활성화 해야 함. 일단 로그 확인용.
		}
	}

	return nil
}

// 그리기는 절차적으로 수행된다. 따라서 그릴 때 마다 이전에 그린 부분에 덮어 씌워진다.
func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255}) // blue background

	opts := ebiten.DrawImageOptions{}

	//loop over the layers
	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {

			//catch display position
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			//catch tileset position (id-> tileset's id)
			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	opts.GeoM.Translate(g.player.X, g.player.Y)

	// draw our player
	screen.DrawImage(
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}

	opts.GeoM.Reset()

	for _, sprite := range g.potions {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}
	opts.GeoM.Reset()

	for _, sprite := range g.potions {
		if sprite.isUsed {
			sprite.Img.Clear()
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//return 320, 240
	return ebiten.WindowSize() // 창 크기 만큼 레이아웃 사이즈 조절
}

func main() {
	// 로그 파일 설정: WSL이나 일부 환경에서 터미널 출력이 안 보일 경우를 대비해 파일로 남김
	f, err := os.OpenFile("game.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f) // 모든 log.Print... 함수가 이 파일에 쓰도록 설정

	ebiten.SetWindowSize(640, 480) // 기본 창 사이즈
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode((ebiten.WindowResizingModeEnabled)) // 전체 창 모드

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/ninga.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/LifePot.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/tileset/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &Player{
			Sprite: &Sprite{
				Img: playerImg,
				X:   50.0,
				Y:   50.0,
			},
			Health: 3,
		},
		enemies: []*Enemy{
			{
				&Sprite{
					Img: skeletonImg,
					X:   100.0,
					Y:   100.0,
				},
				true,
			},
			{
				&Sprite{
					Img: skeletonImg,
					X:   150.0,
					Y:   150.0,
				},
				true,
			},
		},
		potions: []*Potion{
			{
				&Sprite{
					Img: potionImg,
					X:   210.0,
					Y:   100.0,
				},
				1.0,
				false,
			},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
	}

	if err := ebiten.RunGame(&game); err != nil { // Game이라는 구조체(이름은 상관없음) 하나를 정의해서 Update, Draw, Layout에 인터페이스 역할을 수행한다.
		log.Fatal(err)
	}
}
