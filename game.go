package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/FunctionPointerXDD/Trader/animations"
	"github.com/FunctionPointerXDD/Trader/components"
	"github.com/FunctionPointerXDD/Trader/constants"
	"github.com/FunctionPointerXDD/Trader/entities"
	"github.com/FunctionPointerXDD/Trader/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	tilemapJSON       *TilemapJSON
	tilesets          []Tileset
	tilemapImg        *ebiten.Image
	cam               *Camera
	colliders         []image.Rectangle
}

func (g *Game) Update() error {

	// react to key presses

	g.player.Dx = 0.0
	g.player.Dy = 0.0
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Dx = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Dy = 2
	}

	g.player.X += g.player.Dx
	CheckCollisionHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy
	CheckCollisionVertical(g.player.Sprite, g.colliders)

	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		activeAnim.Update()
	}

	for _, enemy := range g.enemies {

		enemy.Dx = 0.0
		enemy.Dy = 0.0
		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.Dx = 0.5
			} else if enemy.X > g.player.X {
				enemy.Dx = -0.5
			}
			if enemy.Y < g.player.Y {
				enemy.Dy = 0.5
			} else if enemy.Y > g.player.Y {
				enemy.Dy = -0.5
			}
		}
		enemy.X += enemy.Dx
		CheckCollisionHorizontal(enemy.Sprite, g.colliders)
		enemy.Y += enemy.Dy
		CheckCollisionVertical(enemy.Sprite, g.colliders)
	}

	for _, potion := range g.potions {
		if !potion.IsUsed && g.player.X > potion.X-16.0 && g.player.X < potion.X+16.0 &&
			g.player.Y > potion.Y-16.0 && g.player.Y < potion.Y+16.0 {
			g.player.Health += potion.AmtHeal
			fmt.Printf("Picked up potion!. Health: %d\n", g.player.Health)
			potion.IsUsed = true
		}
	}

	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	cX, cY := ebiten.CursorPosition()
	cX -= int(g.cam.X)
	cY -= int(g.cam.Y)
	g.player.CombatComp.Update()
	pRect := image.Rect(
		int(g.player.X),
		int(g.player.Y),
		int(g.player.X)+constants.Tilesize,
		int(g.player.Y)+constants.Tilesize,
	)

	deadEnemies := make(map[int]struct{})
	for index, enemy := range g.enemies {
		enemy.CombatComp.Update()
		rect := image.Rect(
			int(enemy.X),
			int(enemy.Y),
			int(enemy.X)+constants.Tilesize,
			int(enemy.Y)+constants.Tilesize,
		)

		if rect.Overlaps(pRect) {
			if enemy.CombatComp.Attack() {
				g.player.CombatComp.Damage(enemy.CombatComp.AttackPower())
				fmt.Println(
					fmt.Sprintf("player damaged. health: %d\n", g.player.CombatComp.Health()),
				)
				if g.player.CombatComp.Health() <= 0 {
					fmt.Println("player has died!")
				}
			}
		}

		//is cursor in rect?
		if cX > rect.Min.X && cX < rect.Max.X && cY > rect.Min.Y && cY < rect.Max.Y {
			if clicked &&
				//플레이어의 공격(클릭)이 플레이어 중심으로 5칸 이내 범위(원)에 속하면 공격 허용
				math.Sqrt(
					math.Pow(
						float64(cX)-g.player.X+(constants.Tilesize/2),
						2,
					)+math.Pow(
						float64(cY)-g.player.Y+(constants.Tilesize/2),
						2,
					),
				) < constants.Tilesize*5 {
				fmt.Println("damagind enemy")
				enemy.CombatComp.Damage(g.player.CombatComp.AttackPower())

				if enemy.CombatComp.Health() <= 0 {
					deadEnemies[index] = struct{}{} //빈 구조체 타입값 적용
					fmt.Println("enemy has been eliminated.")
				}
			}
		}
	}
	if len(deadEnemies) > 0 {
		newEnemies := make([]*entities.Enemy, 0)
		for index, enemy := range g.enemies {
			if _, isDead := deadEnemies[index]; !isDead {
				newEnemies = append(newEnemies, enemy)
			}
		}
		g.enemies = newEnemies
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

	playerFrame := 0
	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		playerFrame = activeAnim.Frame()
	}
	// draw our player
	screen.DrawImage(
		g.player.Img.SubImage(
			g.playerSpriteSheet.Rect(playerFrame), // if activeAnim is nil, then playFrame is Zero(0), So crop 0 index rect image.
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

	for _, collider := range g.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.cam.X),
			float32(collider.Min.Y)+float32(g.cam.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			true,
		)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
	//return ebiten.WindowSize() // 창 크기 만큼 레이아웃 사이즈 조절
}

func NewGame() *Game {

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

	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)

	return &Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   50.0,
				Y:   50.0,
			},
			Health: 3,
			Animations: map[entities.PlayerState]*animations.Animation{
				entities.Up:    animations.NewAnimation(5, 13, 4, 20),
				entities.Down:  animations.NewAnimation(4, 12, 4, 20),
				entities.Left:  animations.NewAnimation(6, 14, 4, 20),
				entities.Right: animations.NewAnimation(7, 15, 4, 20),
			},
			CombatComp: components.NewBasicCombat(3, 1),
		},
		playerSpriteSheet: playerSpriteSheet,
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   100.0,
					Y:   100.0,
				},
				FollowsPlayer: true,
				CombatComp:    components.NewEnemyCombat(3, 1, 30),
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   150.0,
					Y:   150.0,
				},
				FollowsPlayer: false,
				CombatComp:    components.NewEnemyCombat(3, 1, 30),
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
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}
}
