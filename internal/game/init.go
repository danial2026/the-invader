package game

import (
	"fmt"
	"image"
	"math/rand"

	"invaders/internal/ai"
	"invaders/internal/entity"
	"invaders/internal/input"
	"invaders/internal/logger"
	"invaders/internal/renderer"
)

// New creates a new game instance.
func New(cfg Config) (*Game, error) {
	r, err := renderer.New(cfg.WindowWidth, cfg.WindowHeight)
	if err != nil {
		return nil, err
	}

	// Apply Matrix effect to background
	r.ApplyMatrixEffect()

	// Load bios immediately
	logger.Info("LOAD: Alien bios - Reading data/aliens.json")
	allData, err := loadBios("data/aliens.json")
	if err != nil {
		logger.Error("Failed to load bios: %v", err)
		fmt.Printf("Warning: Failed to load bios: %v\n", err)
		allData = []AlienData{{Name: "Unknown", Bio: "???", Personality: "Unknown"}}
	} else {
		logger.Info("LOAD: Alien bios - Loaded %d aliens", len(allData))
	}

	// Initialize game
	g := &Game{
		config:         cfg,
		state:          StateMenu,
		mode:           ModeSinglePlayer,
		renderer:       r,
		input:          input.New(),
		aliens:         make([]*entity.Sprite, 0),
		bombs:          make([]*entity.Sprite, 0),
		beams:          make([]*entity.Sprite, 0),
		alienDirection: 1,
		bot:            ai.NewBot(),
		allAlienData:   allData,
		selectedCount:  7, // Default
		alienData:      make([]AlienData, 0),
		infoCycles:     make([]int, 0),
		alienAI:        make(map[*entity.Sprite]*ai.EnemyAI),
	}

	logger.Info("Game instance created")
	return g, nil
}

// initEntities sets up all game sprites.
func (g *Game) initEntities(count int) {
	logger.Info("CREATE: Entities - Initializing %d aliens", count)
	alienData := g.allAlienData

	if count > len(alienData) {
		count = len(alienData)
	}
	if count < 1 {
		count = 1
	}

	cannonX := g.config.GameWidth/2 - 20
	cannonY := g.config.WindowHeight - 50
	g.cannon = entity.NewCannon(cannonX, cannonY, CannonSprite, CannonExplode)

	g.beams = make([]*entity.Sprite, 0)
	g.aliens = make([]*entity.Sprite, 0)
	g.shields = 3

	// Copy data for UI randomly
	g.alienData = make([]AlienData, count)
	g.infoCycles = make([]int, count)
	perm := rand.Perm(len(alienData))

	for i := 0; i < count; i++ {
		idx := perm[i]
		g.alienData[i] = alienData[idx]
	}

	gap := 45
	cols := 8
	gridWidth := min(count, cols) * gap
	startX := (g.config.GameWidth - gridWidth) / 2
	if startX < 50 {
		startX = 50
	}

	sprites := []image.Rectangle{Alien1Sprite, Alien2Sprite, Alien3Sprite}
	altSprites := []image.Rectangle{Alien1aSprite, Alien2aSprite, Alien3aSprite}

	for i := 0; i < count; i++ {
		data := g.alienData[i]
		spriteType := i % 3
		col := i % cols
		row := i / cols

		x := startX + col*gap
		y := 30 + row*40

		alien := entity.NewAlien(x, y, sprites[spriteType], altSprites[spriteType], AlienExplode, 100)
		alien.Name = data.Name
		g.aliens = append(g.aliens, alien)

		// Create independent AI for each alien
		g.alienAI[alien] = ai.NewEnemyAI(data.Name, mapPersonality(data.Personality))
	}

	g.storyTimer = make([]int, count)
	g.storyAlpha = make([]float64, count)
	g.storyState = make([]int, count)

	for i := 0; i < count; i++ {
		g.storyAlpha[i] = 1.0
		g.storyTimer[i] = i * 100
	}

	logger.Info("CREATE: Entities - Complete (%d aliens with controllers)", count)
	g.state = StatePlaying
}
