package game

import (
	"fmt"
	"image/color"
	"time"

	"invaders/internal/input"
	"invaders/internal/logger"
)

func (g *Game) showStartScreen() bool {
	for {
		g.renderer.NewFrame()
		midX := g.config.WindowWidth / 2
		g.renderer.DrawTextScaled(midX-240, 100, "THE INVADER", color.RGBA{0, 255, 0, 255}, 7.0)
		
		v1 := "Deliver me, O LORD, from the evil man: preserve me from the violent man;"
		v2 := "Which imagine mischiefs in their heart; continually are they gathered together for war."
		
		g.renderer.DrawTextScaled(midX-int(float64(len(v1))*4.5), 180, v1, color.RGBA{255, 0, 0, 255}, 1.5)
		g.renderer.DrawTextScaled(midX-int(float64(len(v2))*4.5), 200, v2, color.RGBA{255, 100, 100, 255}, 1.5)

		g.renderer.DrawText(midX-120, 240, "PRESS [ENTER] OR [SPACE] TO START", color.White)
		g.renderer.DrawText(midX-120, 260, "PRESS 'C' TO QUIT", color.White)
		
		g.renderer.DrawText(midX-120, 300, "CONTROLS:", color.RGBA{0, 255, 0, 255})
		g.renderer.DrawText(midX-120, 320, "WASD : MOVE", color.White)
		g.renderer.DrawText(midX-120, 340, "Q/E  : SHOOT LEFT/RIGHT", color.White)

		g.renderer.Present()
		
		key := g.input.Poll()
		if key == input.KeyEnter || key == input.KeySpace { return true }
		if key == input.KeyQuit { return false }
		time.Sleep(50 * time.Millisecond)
	}
}

func (g *Game) showSelectionScreen() bool {
	g.selectedCount = 3
	maxCount := len(g.allAlienData)
	minCount := 1

	for {
		g.renderer.NewFrame()
		midX := g.config.WindowWidth / 2
		centerY := g.config.WindowHeight / 2
		g.renderer.DrawTextScaled(midX-160, centerY-60, "SELECT ENEMY COUNT", color.RGBA{0, 255, 0, 255}, 2.0)
		g.renderer.DrawText(midX-140, centerY, "<", color.White)
		g.renderer.DrawText(midX+140, centerY, ">", color.White)
		
		col := color.RGBA{0, 255, 0, 255}
		if g.selectedCount == maxCount { col = color.RGBA{255, 100, 100, 255} }
		g.renderer.DrawText(midX-20, centerY, fmt.Sprintf("%02d", g.selectedCount), col)
		g.renderer.DrawText(midX-120, centerY+60, "USE ARROW KEYS TO CHANGE", color.White)
		g.renderer.DrawText(midX-100, centerY+80, "PRESS [ENTER] to CONFIRM", color.White)

		g.renderer.Present()

		key := g.input.Poll()
		if key == input.KeyEnter || key == input.KeySpace { return true }
		if key == input.KeyQuit { return false }
		if key == input.KeyLeft || key == input.KeyDown {
			g.selectedCount--
			if g.selectedCount < minCount { g.selectedCount = maxCount }
		}
		if key == input.KeyRight || key == input.KeyUp {
			g.selectedCount++
			if g.selectedCount > maxCount { g.selectedCount = minCount }
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (g *Game) showLoadingScreen() {
	logger.Info("Loading screen: Starting initialization sequence")
	
	steps := []struct {
		name   string
		action func()
	}{
		{"Preparing Game Engine", func() { time.Sleep(200 * time.Millisecond) }},
		{"Initializing AI Systems", func() { time.Sleep(200 * time.Millisecond) }},
		{"Loading Alien Data", func() { time.Sleep(200 * time.Millisecond) }},
		{"Creating AI Controllers", func() { time.Sleep(200 * time.Millisecond) }},
	}

	midX := g.config.WindowWidth / 2
	midY := g.config.WindowHeight / 2
	greenColor := color.RGBA{0, 255, 0, 255}
	whiteColor := color.RGBA{255, 255, 255, 255}
	yellowColor := color.RGBA{255, 255, 0, 255}
	dimColor := color.RGBA{100, 100, 100, 255}

	for i, step := range steps {
		g.renderer.NewFrame()
		g.renderer.DrawTextScaled(midX-160, midY-100, "INITIALIZING SYSTEMS...", greenColor, 2.0)
		
		barWidth, barHeight := 400, 20
		barX, barY := midX - barWidth/2, midY - 50
		g.renderer.DrawRect(barX, barY, barWidth, barHeight, color.RGBA{30, 30, 30, 255})
		
		progress := float64(i+1) / float64(len(steps))
		fillWidth := int(float64(barWidth) * progress)
		g.renderer.DrawRect(barX, barY, fillWidth, barHeight, greenColor)
		g.renderer.DrawText(midX-15, barY+5, fmt.Sprintf("%d%%", int(progress * 100)), whiteColor)
		
		yOffset := midY + 20
		for j := 0; j < i; j++ { g.renderer.DrawText(midX-150, yOffset+j*25, "[✓] "+steps[j].name, dimColor) }
		
		dots := []string{"", ".", "..", "..."}
		g.renderer.DrawText(midX-150, yOffset+i*25, "[...] " + step.name + dots[i%4], yellowColor)
		
		for j := i + 1; j < len(steps); j++ { g.renderer.DrawText(midX-150, yOffset+j*25, "[ ] "+steps[j].name, dimColor) }
		
		g.renderer.Present()
		logger.Info("Loading step: %s", step.name)
		step.action()
	}
	
	g.renderer.NewFrame()
	g.renderer.DrawTextScaled(midX-160, midY-100, "INITIALIZING SYSTEMS...", greenColor, 2.0)
	barX, barY := midX - 200, midY - 50
	g.renderer.DrawRect(barX, barY, 400, 20, greenColor)
	g.renderer.DrawText(midX-15, barY+5, "100%", whiteColor)
	
	yOffset := midY + 20
	for j, step := range steps { g.renderer.DrawText(midX-150, yOffset+j*25, "[✓] "+step.name, whiteColor) }
	
	g.renderer.DrawText(midX-40, midY+140, "READY", greenColor)
	g.renderer.Present()
	logger.Info("Loading screen: Complete")
	time.Sleep(500 * time.Millisecond)
}
