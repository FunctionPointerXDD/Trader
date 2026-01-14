package main

import (
	"image"
	"image/color"
	"log"

	"github.com/FunctionPointerXDD/Trader/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapJSON *TilemapJSON
	tilesets    []Tileset
	tilemapImg  *ebiten.Image
	cam         *Camera
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
		if !potion.IsUsed && g.player.X > potion.X-16.0 && g.player.X < potion.X+16.0 &&
			g.player.Y > potion.Y-16.0 && g.player.Y < potion.Y+16.0 {
			g.player.Health += potion.AmtHeal
			log.Printf("Picked up potion!. Health: %d\n", g.player.Health)
			potion.IsUsed = true
		}
	}

	g.cam.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.cam.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*16.0,
		float64(g.tilemapJSON.Layers[0].Height)*16.0,
		320,
		240,
	)

	return nil
}

// 그리기는 절차적으로 수행된다. 따라서 그릴 때 마다 이전에 그린 부분에 덮어 씌워진다.
func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255}) // blue background

	opts := ebiten.DrawImageOptions{}

	//loop over the layers
	for layerIndex, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {

			if id == 0 {
				continue
			}
			//catch display position
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			img := g.tilesets[layerIndex].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))
			opts.GeoM.Translate(g.cam.X, g.cam.Y)
			screen.DrawImage(img, &opts)
			opts.GeoM.Reset()
		}
	}

	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

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
		opts.GeoM.Translate(g.cam.X, g.cam.Y)
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
		if sprite.IsUsed {
			continue
		}
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)
		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
	//return ebiten.WindowSize() // 창 크기 만큼 레이아웃 사이즈 조절
}

func main() {

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

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   50.0,
				Y:   50.0,
			},
			Health: 3,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   100.0,
					Y:   100.0,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   150.0,
					Y:   150.0,
				},
				FollowsPlayer: true,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   210.0,
					Y:   100.0,
				},
				AmtHeal: 1.0,
				IsUsed:  false,
			},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
		tilesets:    tilesets,
		cam:         NewCamera(0.0, 0.0),
	}

	if err := ebiten.RunGame(&game); err != nil { // Game이라는 구조체(이름은 상관없음) 하나를 정의해서 Update, Draw, Layout에 인터페이스 역할을 수행한다.
		log.Fatal(err)
	}
}
