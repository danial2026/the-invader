package game

import (
	"image"
)

type AlienData struct {
	Name        string     `json:"name"`
	Personality string     `json:"personality"`
	Bio         string     `json:"bio"`
	Stories     []string   `json:"stories"`
	Family      []Relative `json:"family"`
	Friends     []Relative `json:"friends"`
}

type Relative struct {
	Name     string   `json:"name"`
	Relation string   `json:"relation"`
	Stories  []string `json:"stories"`
}

type GameMode int

const (
	ModeSinglePlayer GameMode = iota
	ModeAIBattle
)

type State int

const (
	StateMenu State = iota
	StatePlaying
	StatePaused
	StateSelectCount
	StateConfirmQuit
	StateGameOver
)

// SpriteRects defines sprite positions in the sprite sheet.
var (
	CannonSprite  = image.Rect(20, 47, 38, 59)
	CannonExplode = image.Rect(0, 47, 16, 57)
	Alien1Sprite  = image.Rect(0, 0, 20, 14)
	Alien1aSprite = image.Rect(20, 0, 40, 14)
	Alien2Sprite  = image.Rect(0, 14, 20, 26)
	Alien2aSprite = image.Rect(20, 14, 40, 26)
	Alien3Sprite  = image.Rect(0, 27, 20, 40)
	Alien3aSprite = image.Rect(20, 27, 40, 40)
	AlienExplode  = image.Rect(0, 60, 16, 68)
	BeamSprite    = image.Rect(20, 60, 22, 65)
	BombSprite    = image.Rect(0, 70, 10, 79)
)

// Config holds game configuration.
type Config struct {
	WindowWidth  int
	WindowHeight int
	GameWidth    int // Width of playable area
	UIWidth      int // Width of side panel
	AliensPerRow int
	BombProb     float64
	BombSpeed    int
}

// DefaultConfig returns sensible default configuration.
func DefaultConfig() Config {
	return Config{
		// 1100x600 total, 600 game, 500 UI
		WindowWidth:  1100,
		WindowHeight: 600,
		GameWidth:    600,
		UIWidth:      500,
		AliensPerRow: 8,
		BombProb:     0.005,
		BombSpeed:    10,
	}
}
