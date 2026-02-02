package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseScene struct {
	loaded bool
}

func NewPauseScene() *PauseScene {
	return &PauseScene{
		loaded: false,
	}
}

func (p *PauseScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 255, 0, 255})
	ebitenutil.DebugPrint(screen, "Press Enter to unpause")
}

func (p *PauseScene) FirstLoad() {
	p.loaded = true
}

func (p *PauseScene) IsLoaded() bool {
	return p.loaded
}

func (p *PauseScene) OnEnter() {
}

func (p *PauseScene) OnExit() {
}

func (p *PauseScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	return PauseSceneId
}

var _ Scene = (*PauseScene)(nil)
