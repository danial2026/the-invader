package entity

import (
	"image"

	"github.com/disintegration/gift"
)

type Sprite struct {
	Name     string
	Size     image.Rectangle
	Filter   *gift.GIFT
	FilterA  *gift.GIFT
	FilterE  *gift.GIFT
	FilterD  *gift.GIFT
	Position image.Point
	Status   bool
	Points   int
	Ammo     int
	VX, VY   int
}

func NewAlien(x, y int, sprite, alt, explode image.Rectangle, points int) *Sprite {
	return &Sprite{
		Size:     sprite,
		Filter:   gift.New(gift.Crop(sprite)),
		FilterA:  gift.New(gift.Crop(alt)),
		FilterE:  gift.New(gift.Crop(explode)),
		FilterD:  gift.New(gift.Crop(sprite), gift.Colorize(100, 100, 100), gift.Brightness(-50)),
		Position: image.Pt(x, y),
		Status:   true,
		Points:   points,
		Ammo:     100,
	}
}

func NewCannon(x, y int, sprite, explode image.Rectangle) *Sprite {
	return &Sprite{
		Size:     sprite,
		Filter:   gift.New(gift.Crop(sprite)),
		FilterE:  gift.New(gift.Crop(explode)),
		Position: image.Pt(x, y),
		Status:   true,
		Ammo:     100,
	}
}

func NewProjectile(x, y int, sprite image.Rectangle) *Sprite {
	return &Sprite{
		Size:     sprite,
		Filter:   gift.New(gift.Crop(sprite)),
		Position: image.Pt(x, y),
		Status:   false,
	}
}

func (s *Sprite) Move(dx, dy int) {
	s.Position.X += dx
	s.Position.Y += dy
}

func (s *Sprite) Kill() {
	s.Status = false
}

func (s *Sprite) Activate() {
	s.Status = true
}

func (s *Sprite) IsAlive() bool {
	return s.Status
}

func (s *Sprite) Bounds() image.Rectangle {
	return image.Rect(
		s.Position.X,
		s.Position.Y,
		s.Position.X+s.Size.Dx(),
		s.Position.Y+s.Size.Dy(),
	)
}

