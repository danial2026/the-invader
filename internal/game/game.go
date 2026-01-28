package game

import (
	"invaders/internal/ai"
	"invaders/internal/entity"
	"invaders/internal/input"
	"invaders/internal/renderer"
)

type Game struct {
	config   Config
	state    State
	mode     GameMode
	renderer *renderer.Renderer
	input    *input.Handler

	cannon *entity.Sprite
	beams  []*entity.Sprite
	aliens []*entity.Sprite
	bombs  []*entity.Sprite

	score          int
	loop           int
	alienDirection int
	gameOverText   string

	allAlienData  []AlienData
	selectedCount int
	alienData     []AlienData
	infoCycles    []int

	storyTimer []int
	storyAlpha []float64
	storyState []int

	bot *ai.Bot

	alienAI map[*entity.Sprite]*ai.EnemyAI

	shields int
}

func (g *Game) Close() {
	if g.renderer != nil {
		g.renderer.Close()
	}
}

// SetMode sets the game mode (Single Player or AI Battle).
func (g *Game) SetMode(mode GameMode) {
	g.mode = mode
}
