package ai

import (
	"math"
	"math/rand"

	"invaders/internal/entity"
	"invaders/internal/logger"
)

type Action int

const (
	ActionNone Action = iota
	ActionMoveLeft
	ActionMoveRight
	ActionShoot
)

type Controller interface {
	Decide(cannon *entity.Sprite, aliens []*entity.Sprite, bombs []*entity.Sprite, beams []*entity.Sprite) Action
}

type Bot struct {
	lastShot    int
	minShotGap  int
}

func NewBot() *Bot {
	logger.Info("INIT: AI - Player bot created")
	return &Bot{
		lastShot:   0,
		minShotGap: 5,
	}
}

func (b *Bot) Decide(cannon *entity.Sprite, aliens []*entity.Sprite, bombs []*entity.Sprite, beams []*entity.Sprite) Action {
	if cannon == nil || !cannon.IsAlive() {
		return ActionNone
	}

	cannonX := cannon.Position.X
	cannonY := cannon.Position.Y

	if action := b.dodgeBombs(cannonX, cannonY, bombs); action != ActionNone {
		return action
	}

	target := b.findNearestAlien(cannonX, aliens)
	if target == nil {
		return ActionNone
	}

	targetX := target.Position.X + target.Size.Dx()/2
	cannonCenterX := cannonX + cannon.Size.Dx()/2

	if abs(targetX-cannonCenterX) < 8 {
		if len(beams) < 3 {
			b.lastShot = 0
			return ActionShoot
		}
	}

	if targetX < cannonCenterX {
		return ActionMoveLeft
	} else if targetX > cannonCenterX {
		return ActionMoveRight
	}

	return ActionNone
}

func (b *Bot) dodgeBombs(cannonX, cannonY int, bombs []*entity.Sprite) Action {
	const dangerZone = 30
	const verticalThreat = 80

	for _, bomb := range bombs {
		if bomb == nil || !bomb.IsAlive() {
			continue
		}
		if cannonY-bomb.Position.Y > verticalThreat {
			continue
		}
		dx := bomb.Position.X - cannonX
		if abs(dx) < dangerZone {
			if dx > 0 {
				return ActionMoveLeft
			}
			return ActionMoveRight
		}
	}
	return ActionNone
}

func (b *Bot) findNearestAlien(cannonX int, aliens []*entity.Sprite) *entity.Sprite {
	var nearest *entity.Sprite
	minDist := math.MaxFloat64
	for _, alien := range aliens {
		if alien == nil || !alien.IsAlive() {
			continue
		}
		dist := math.Abs(float64(alien.Position.X - cannonX))
		if dist < minDist {
			minDist = dist
			nearest = alien
		}
	}
	return nearest
}

type EnemyPersonality int

const (
	PersonalityAggressive EnemyPersonality = iota
	PersonalityDefensive
	PersonalityChaotic
	PersonalityHunter
	PersonalityStalker
	PersonalityBerserker
	PersonalitySniper
)

type EnemyAI struct {
	Name        string
	Personality EnemyPersonality
	Speed       int
	ShootRate   float64
	Target      *entity.Sprite
	Direction     int
	MoveTimer     int
	VerticalSpeed int
}

func NewEnemyAI(name string, personality EnemyPersonality) *EnemyAI {
	speed := 3
	shootRate := 0.005

	var personalityName string
	switch personality {
	case PersonalityAggressive:
		speed = 5
		shootRate = 0.03
		personalityName = "Aggressive"
	case PersonalityDefensive:
		speed = 2
		shootRate = 0.009
		personalityName = "Defensive"
	case PersonalityChaotic:
		speed = 4
		shootRate = 0.02
		personalityName = "Chaotic"
	case PersonalityHunter:
		speed = 4
		shootRate = 0.018
		personalityName = "Hunter"
	case PersonalityStalker:
		speed = 1
		shootRate = 0.02
		personalityName = "Stalker"
	case PersonalityBerserker:
		speed = 6
		shootRate = 0.04
		personalityName = "Berserker"
	case PersonalitySniper:
		speed = 1
		shootRate = 0.01
		personalityName = "Sniper"
	default:
		personalityName = "Unknown"
	}

	logger.Info("CREATE: AI Controller - %s (Personality: %s)", name, personalityName)

	return &EnemyAI{
		Name:        name,
		Personality: personality,
		Speed:       speed,
		ShootRate:   shootRate,
		Direction:   1,
		MoveTimer:   rand.Intn(30) + 10,
		VerticalSpeed: 0,
	}
}

func (e *EnemyAI) Update(enemy *entity.Sprite, cannon *entity.Sprite, beams []*entity.Sprite, bounds struct{ MinX, MaxX, MinY, MaxY int }) bool {
	if enemy == nil || !enemy.IsAlive() {
		return false
	}

	e.MoveTimer--

	switch e.Personality {
	case PersonalityAggressive:
		e.moveAggressive(enemy, cannon, bounds)
	case PersonalityDefensive:
		e.moveDefensive(enemy, cannon, beams, bounds)
	case PersonalityChaotic:
		e.moveChaotic(enemy, bounds)
	case PersonalityHunter:
		e.moveHunter(enemy, cannon, bounds)
	case PersonalityStalker:
		e.moveStalker(enemy, cannon, bounds)
	case PersonalityBerserker:
		e.moveBerserker(enemy, bounds)
	case PersonalitySniper:
		e.moveSniper(enemy, bounds)
	}

	return rand.Float64() < e.ShootRate
}

func (e *EnemyAI) moveAggressive(enemy *entity.Sprite, cannon *entity.Sprite, bounds struct{ MinX, MaxX, MinY, MaxY int }) {
	if cannon != nil {
		dx := cannon.Position.X - enemy.Position.X
		step := e.Speed
		
		if abs(dx) < 20 {
			step = 1
		}
		
		if dx > 10 {
			enemy.Move(step, 0)
		} else if dx < -10 {
			enemy.Move(-step, 0)
		}
	}
	if e.MoveTimer <= 0 && enemy.Position.Y < 400 {
		enemy.Move(0, 15)
		e.MoveTimer = 40
	}
	e.checkBounds(enemy, bounds)
}

func (e *EnemyAI) moveDefensive(enemy *entity.Sprite, cannon *entity.Sprite, beams []*entity.Sprite, bounds struct{ MinX, MaxX, MinY, MaxY int }) {
	for _, beam := range beams {
		if beam == nil || !beam.IsAlive() { continue }
		if beam.Position.Y > enemy.Position.Y {
			dx := beam.Position.X - (enemy.Position.X + enemy.Size.Dx()/2)
			if abs(dx) < 25 {
				if dx > 0 {
					enemy.Move(-e.Speed*3, 0)
				} else {
					enemy.Move(e.Speed*3, 0)
				}
				e.checkBounds(enemy, bounds)
				return
			}
		}
	}

	if cannon != nil {
		dist := abs(enemy.Position.X - cannon.Position.X)
		if dist < 40 {
			if enemy.Position.X < cannon.Position.X {
				enemy.Move(-e.Speed*2, 0)
			} else {
				enemy.Move(e.Speed*2, 0)
			}
		} else {
			center := (bounds.MaxX - bounds.MinX) / 2
			if enemy.Position.X < center {
				enemy.Move(1, 0)
			} else {
				enemy.Move(-1, 0)
			}
		}
	}
	e.checkBounds(enemy, bounds)
}

func (e *EnemyAI) moveChaotic(enemy *entity.Sprite, bounds struct{ MinX, MaxX, MinY, MaxY int }) {
	if e.MoveTimer <= 0 {
		e.Direction = []int{-1, 1}[rand.Intn(2)]
		e.VerticalSpeed = []int{-2, 0, 2}[rand.Intn(3)]
		e.MoveTimer = rand.Intn(20) + 5
		if rand.Float64() < 0.3 {
			enemy.Move(0, []int{-15, 15}[rand.Intn(2)])
		}
	}
	enemy.Move(e.Speed*e.Direction, e.VerticalSpeed)
	e.checkBounds(enemy, bounds)
}

func (e *EnemyAI) moveHunter(enemy *entity.Sprite, cannon *entity.Sprite, bounds struct{ MinX, MaxX, MinY, MaxY int }) {
	if cannon != nil {
		targetX := cannon.Position.X
		dx := targetX - enemy.Position.X
		
		if abs(dx) > 10 {
			if dx > 0 {
				enemy.Move(e.Speed, 0)
			} else {
				enemy.Move(-e.Speed, 0)
			}
		}
	}
	e.checkBounds(enemy, bounds)
}

func (e *EnemyAI) moveStalker(enemy *entity.Sprite, cannon *entity.Sprite, bounds struct{ MinX, MaxX, MinY, MaxY int }) {
	if cannon != nil {
		dx := cannon.Position.X - enemy.Position.X
		if abs(dx) < 10 {
			e.MoveTimer = 20
		} else if e.MoveTimer <= 0 {
			if dx > 0 {
				enemy.Move(1, 0)
			} else {
				enemy.Move(-1, 0)
			}
		}
	}
}

func (e *EnemyAI) moveBerserker(enemy *entity.Sprite, bounds struct{ MinX, MaxX, MinY, MaxY int }) {
	enemy.Move(e.Speed*e.Direction, int(math.Sin(float64(enemy.Position.X)/50.0)*5.0))
	if enemy.Position.X <= bounds.MinX || enemy.Position.X >= bounds.MaxX-enemy.Size.Dx() {
		e.Direction *= -1
		enemy.Move(0, 20)
	}
}

func (e *EnemyAI) moveSniper(enemy *entity.Sprite, bounds struct{ MinX, MaxX, MinY, MaxY int }) {
	if e.MoveTimer <= 0 {
		targetX := []int{50, 300, 550}[rand.Intn(3)]
		dx := targetX - enemy.Position.X
		if abs(dx) > 5 {
			dir := 1
			if dx < 0 { dir = -1 }
			enemy.Move(dir*3, 0)
		} else {
			e.MoveTimer = 100
		}
	}
}

func (e *EnemyAI) checkBounds(enemy *entity.Sprite, bounds struct{ MinX, MaxX, MinY, MaxY int }) {
	if enemy.Position.X < bounds.MinX {
		enemy.Position.X = bounds.MinX
		e.Direction = 1
	}
	if enemy.Position.X > bounds.MaxX-20 {
		enemy.Position.X = bounds.MaxX - 20
		e.Direction = -1
	}
	if enemy.Position.Y < bounds.MinY {
		enemy.Position.Y = bounds.MinY
		e.VerticalSpeed = 2
	}
	if enemy.Position.Y > bounds.MaxY-enemy.Size.Dy() {
		enemy.Position.Y = bounds.MaxY - enemy.Size.Dy()
		e.VerticalSpeed = -2
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

