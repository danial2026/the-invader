package game

import (
	"encoding/json"
	"io/ioutil"

	"strings"

	"invaders/internal/ai"
	"invaders/internal/entity"
)

func mapPersonality(p string) ai.EnemyPersonality {
	p = strings.ToLower(p)
	if strings.Contains(p, "aggressive") || strings.Contains(p, "angry") || strings.Contains(p, "fighter") {
		return ai.PersonalityAggressive
	}
	if strings.Contains(p, "defensive") || strings.Contains(p, "shy") || strings.Contains(p, "scared") || strings.Contains(p, "anxious") {
		return ai.PersonalityDefensive
	}
	if strings.Contains(p, "chaotic") || strings.Contains(p, "crazy") || strings.Contains(p, "wild") || strings.Contains(p, "coffee") {
		return ai.PersonalityChaotic
	}
	if strings.Contains(p, "hunter") || strings.Contains(p, "tracker") || strings.Contains(p, "rival") {
		return ai.PersonalityHunter
	}
	if strings.Contains(p, "stalker") || strings.Contains(p, "creepy") || strings.Contains(p, "quiet") {
		return ai.PersonalityStalker
	}
	if strings.Contains(p, "berserker") || strings.Contains(p, "trucker") || strings.Contains(p, "strong") {
		return ai.PersonalityBerserker
	}
	if strings.Contains(p, "sniper") || strings.Contains(p, "precise") || strings.Contains(p, "calm") {
		return ai.PersonalitySniper
	}
	return ai.PersonalityChaotic
}

func loadBios(path string) ([]AlienData, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data []AlienData
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func collide(s1, s2 *entity.Sprite) bool {
	if s1 == nil || s2 == nil {
		return false
	}
	a := s1.Bounds()
	b := s2.Bounds()
	return a.Min.X < b.Max.X && a.Max.X > b.Min.X &&
		a.Min.Y < b.Max.Y && a.Max.Y > b.Min.Y
}
