package game

import (
	"math/rand"
	"time"

	"invaders/internal/ai"
	"invaders/internal/entity"
	"invaders/internal/input"
	"invaders/internal/logger"
)

// Run starts the game main loop.
func (g *Game) Run() {
	rand.Seed(time.Now().UTC().UnixNano())

	if !g.showStartScreen() {
		return
	}

	if !g.showSelectionScreen() {
		return
	}

	g.showLoadingScreen()
	g.initEntities(g.selectedCount)

	for {
		if g.state == StateGameOver {
			g.renderGameOver()
			key := g.input.Poll()
			if key == input.KeyQuit || key == input.KeyEnter {
				break
			}
			time.Sleep(50 * time.Millisecond)
			continue
		}

		if g.state == StatePaused {
			g.renderPaused()
			if g.input.Poll() == input.KeyPause {
				g.state = StatePlaying
			}
			time.Sleep(50 * time.Millisecond)
			continue
		}

		g.update()
		g.render()
		time.Sleep(16 * time.Millisecond)
	}
}

func (g *Game) update() {
	if g.mode == ModeAIBattle {
		g.handleAIInput()
	} else {
		g.handleHumanInput()
	}

	if g.state == StatePlaying {
		g.updateAliens()
		g.updateBombs()
		g.updateBeam()
		g.checkEndConditions()
	}

	if g.loop%60 == 0 {
		logger.Info("Heartbeat: Frame %d, State %d, Entities: %d", g.loop, g.state, len(g.aliens))
	}
	g.loop++
}

func (g *Game) handleHumanInput() {
	key := g.input.Poll()

	if key == input.KeyQuit {
		if g.state == StatePlaying || g.state == StatePaused {
			g.state = StateConfirmQuit
		} else if g.state == StateConfirmQuit {
			g.state = StateGameOver
		}
		return
	}

	if g.state == StateConfirmQuit && key != input.KeyNone {
		g.state = StatePlaying
		return
	}

	if key == input.KeyPause {
		if g.state == StatePlaying {
			g.state = StatePaused
		} else if g.state == StatePaused {
			g.state = StatePlaying
		}
		return
	}

	if g.state == StatePlaying {
		switch key {
		case input.KeyLeft:
			g.cannon.Move(-10, 0)
		case input.KeyRight:
			g.cannon.Move(10, 0)
		case input.KeyUp:
			g.cannon.Move(0, -10)
		case input.KeyDown:
			g.cannon.Move(0, 10)
		case input.KeyShootLeft:
			g.fireBeam(2)
		case input.KeyShootRight:
			g.fireBeam(14)
		}
		g.clampCannon()
	}
}

func (g *Game) handleAIInput() {
	key := g.input.Poll()

	if key == input.KeyQuit {
		if g.state == StatePlaying || g.state == StatePaused {
			g.state = StateConfirmQuit
		} else if g.state == StateConfirmQuit {
			g.state = StateGameOver
		}
		return
	}

	if g.state == StateConfirmQuit && key != input.KeyNone {
		g.state = StatePlaying
		return
	}

	action := g.bot.Decide(g.cannon, g.aliens, g.bombs, g.beams)
	switch action {
	case ai.ActionMoveLeft:
		g.cannon.Move(-10, 0)
	case ai.ActionMoveRight:
		g.cannon.Move(10, 0)
	case ai.ActionShoot:
		g.fireBeam(8)
	}
}

func (g *Game) updateAliens() {
	bounds := struct{ MinX, MaxX, MinY, MaxY int }{
		MinX: 0,
		MaxX: g.config.GameWidth,
		MinY: 0,
		MaxY: g.config.WindowHeight,
	}

	for _, alien := range g.aliens {
		if !alien.IsAlive() {
			continue
		}

		if brain, ok := g.alienAI[alien]; ok {
			shouldShoot := brain.Update(alien, g.cannon, g.beams, bounds)
			if shouldShoot {
				g.fireEnemyBomb(alien)
			}
		}

		// Check beam collisions
		for _, beam := range g.beams {
			if beam.IsAlive() && collide(alien, beam) {
				alien.Kill()
				g.score += alien.Points
				beam.Kill()
			}
		}
	}

	// Anti-collision logic
	for i, a1 := range g.aliens {
		if !a1.IsAlive() {
			continue
		}
		for j, a2 := range g.aliens {
			if i == j || !a2.IsAlive() {
				continue
			}

			distX := a1.Position.X - a2.Position.X
			distY := a1.Position.Y - a2.Position.Y
			if abs(distX) < 30 && abs(distY) < 30 {
				if distX > 0 {
					a1.Move(1, 0)
				} else {
					a1.Move(-1, 0)
				}
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *Game) updateBombs() {
	for _, bomb := range g.bombs {
		if !bomb.IsAlive() {
			continue
		}

		bomb.Move(0, g.config.BombSpeed)

		if collide(bomb, g.cannon) {
			if g.shields > 0 {
				g.shields--
				bomb.Kill()
				logger.Info("Shield hit! Shields remaining: %d", g.shields)
			} else {
				g.cannon.Kill()
				g.state = StateGameOver
			}
		}

		if bomb.Position.Y > g.config.WindowHeight {
			bomb.Kill()
		}
	}
}

func (g *Game) fireEnemyBomb(alien *entity.Sprite) {
	if alien.Ammo > 0 {
		bomb := entity.NewProjectile(alien.Position.X+7, alien.Position.Y+15, BombSprite)
		bomb.Activate()
		g.bombs = append(g.bombs, bomb)
		alien.Ammo--
	}
}

func (g *Game) fireBeam(offset int) {
	if len(g.beams) < 10 && g.cannon.Ammo > 0 {
		beamX := g.cannon.Position.X + offset
		beam := entity.NewProjectile(beamX, g.cannon.Position.Y, BeamSprite)
		beam.Activate()
		g.beams = append(g.beams, beam)
		g.cannon.Ammo--
	}
}

func (g *Game) clampCannon() {
	if g.cannon.Position.X < 0 {
		g.cannon.Position.X = 0
	}
	if g.cannon.Position.X > g.config.GameWidth-g.cannon.Size.Dx() {
		g.cannon.Position.X = g.config.GameWidth - g.cannon.Size.Dx()
	}
	if g.cannon.Position.Y < 0 {
		g.cannon.Position.Y = 0
	}
	if g.cannon.Position.Y > g.config.WindowHeight-50 {
		g.cannon.Position.Y = g.config.WindowHeight - 50
	}
}

func (g *Game) updateBeam() {
	activeBeams := make([]*entity.Sprite, 0)
	for _, beam := range g.beams {
		beam.Move(0, -15)

		if beam.Position.Y > 0 && beam.IsAlive() {
			activeBeams = append(activeBeams, beam)
		} else {
			beam.Kill()
		}
	}
	g.beams = activeBeams
}

func (g *Game) checkEndConditions() {
	allDead := true
	for _, a := range g.aliens {
		if a.IsAlive() {
			allDead = false
			break
		}
	}
	if allDead {
		logger.Info("Game Won!")
		g.gameOverText = g.generateEndingText(true)
		g.state = StateGameOver
	}

	if !g.cannon.IsAlive() {
		logger.Info("Game Lost - Cannon Destroyed")
		g.gameOverText = g.generateEndingText(false)
		g.state = StateGameOver
	}
}
