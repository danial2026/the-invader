package game

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"
)

func (g *Game) renderPaused() {
	g.renderer.NewFrame()
	midX := g.config.WindowWidth / 2
	midY := g.config.WindowHeight / 2
	g.renderer.DrawText(midX-30, midY, "PAUSED", color.White)
	g.renderer.Present()
}

func (g *Game) generateEndingText(win bool) string {
	if win {
		lines := []string{
			"VICTORY... BUT AT WHAT COST?",
			"",
			"The skies are clear. and their defense has fallen.",
			"But somewhere, far away, a family waits for a ship that will never return.",
		}
		
		if len(g.alienData) > 0 {
			idx := rand.Intn(len(g.alienData))
			data := g.alienData[idx]
			name := data.Name
			
			lines = append(lines, "")
			lines = append(lines, fmt.Sprintf("%s never made it home.", name))
			
			if len(data.Family) > 0 {
				rel := data.Family[0]
				lines = append(lines, fmt.Sprintf("%s (%s) will keep looking at the sky,", rel.Name, rel.Relation, "."))
				lines = append(lines, "")
				lines = append(lines, "\"We are orphans and fatherless, our mothers are as widows.\" - Lamentations 5:3")
			} else {
				lines = append(lines, "Their story ends here, in cold silence.")
			}
		}
		return strings.Join(lines, "\n")
	} else {
		lines := []string{
			"DEFEAT.",
			"",
			"The fighters continue their journey.",
		}
		
		var survivors []int
		for i, a := range g.aliens {
			if a.IsAlive() {
				survivors = append(survivors, i)
			}
		}
		
		if len(survivors) > 0 {
			idx := survivors[rand.Intn(len(survivors))]
			alien := g.aliens[idx]
			
			lines = append(lines, "")
			lines = append(lines, fmt.Sprintf("%s survives.", alien.Name))
			lines = append(lines, "The war continues for generations.")
			lines = append(lines, "Only survivors return to what remains.")
		}
		return strings.Join(lines, "\n")
	}
}

func (g *Game) renderGameOver() {
	g.renderer.NewFrame()
	
	midX := g.config.WindowWidth / 2
	midY := g.config.WindowHeight / 2
	
	g.renderer.DrawRect(midX-300, midY-150, 600, 300, color.RGBA{20, 20, 20, 255})
	g.renderer.DrawRect(midX-300, midY-150, 600, 2, color.White)
	g.renderer.DrawRect(midX-300, midY+150, 600, 2, color.White)
	
	lines := strings.Split(g.gameOverText, "\n")
	startY := midY - 60
	
	for i, line := range lines {
		width := len(line) * 6
		textX := midX - (width / 2)
		
		var col color.Color = color.White
		if i == 0 { 
			col = color.RGBA{255, 0, 0, 255}
		}
		g.renderer.DrawText(textX, startY + (i*20), line, col)
	}
	
	g.renderer.DrawText(midX-60, midY+120, "PRESS 'C' TO QUIT", color.RGBA{100, 100, 100, 255})
	g.renderer.Present()
}

func (g *Game) render() {
	g.renderer.NewFrame()

	for _, alien := range g.aliens {
		if !alien.IsAlive() {
			if alien.FilterD != nil {
				g.renderer.DrawSprite(alien.FilterD, alien.Position.X, alien.Position.Y)
			}
			continue
		}
		
		filter := alien.Filter
		if g.loop%40 < 20 && alien.FilterA != nil {
			filter = alien.FilterA
		}
		g.renderer.DrawSprite(filter, alien.Position.X, alien.Position.Y)
		
		g.renderer.DrawText(alien.Position.X, alien.Position.Y-10, alien.Name, color.RGBA{0, 255, 0, 255})
		g.renderer.DrawText(alien.Position.X, alien.Position.Y+15, fmt.Sprintf("[%d]", alien.Ammo), color.RGBA{255, 0, 0, 255})
	}

	for _, bomb := range g.bombs {
		if bomb.IsAlive() {
			g.renderer.DrawSprite(bomb.Filter, bomb.Position.X, bomb.Position.Y)
		}
	}

	if g.cannon.IsAlive() {
		g.renderer.DrawSprite(g.cannon.Filter, g.cannon.Position.X, g.cannon.Position.Y)
		
		if g.shields > 0 {
			pulse := (g.loop % 30)
			alpha := 100 + pulse*3
			shieldColor := color.RGBA{0, 255, 255, uint8(alpha)}
			g.renderer.DrawRect(g.cannon.Position.X-2, g.cannon.Position.Y-2, g.cannon.Size.Dx()+4, 2, shieldColor)
			g.renderer.DrawRect(g.cannon.Position.X-2, g.cannon.Position.Y+g.cannon.Size.Dy(), g.cannon.Size.Dx()+4, 2, shieldColor)
			g.renderer.DrawRect(g.cannon.Position.X-2, g.cannon.Position.Y-2, 2, g.cannon.Size.Dy()+4, shieldColor)
			g.renderer.DrawRect(g.cannon.Position.X+g.cannon.Size.Dx(), g.cannon.Position.Y-2, 2, g.cannon.Size.Dy()+4, shieldColor)
		}

		g.renderer.DrawText(g.cannon.Position.X, g.cannon.Position.Y+15, fmt.Sprintf("AMMO: %d", g.cannon.Ammo), color.RGBA{0, 255, 255, 255})
	} else {
		g.renderer.DrawSprite(g.cannon.FilterE, g.cannon.Position.X, g.cannon.Position.Y)
	}

	for _, beam := range g.beams {
		if beam.IsAlive() {
			g.renderer.DrawRect(beam.Position.X, beam.Position.Y, 2, 8, color.RGBA{100, 255, 100, 255})
		}
	}

	panelX := g.config.GameWidth
	g.renderer.DrawRect(panelX, 0, g.config.UIWidth, g.config.WindowHeight, color.RGBA{0, 20, 0, 255})
	g.renderer.DrawRect(panelX, 0, 2, g.config.WindowHeight, color.RGBA{0, 255, 0, 255})
	
	g.renderer.DrawText(panelX+10, 20, fmt.Sprintf("SCORE: %05d", g.score), color.White)
	g.renderer.DrawText(panelX+200, 20, "SHIELDS:", color.White)
	
	shieldBoxX := panelX + 270
	shieldBoxY := 22
	boxW, boxH, spacing := 15, 8, 5
	
	for s := 0; s < 3; s++ {
		boxX := shieldBoxX + s*(boxW+spacing)
		g.renderer.DrawRect(boxX, shieldBoxY, boxW, boxH, color.RGBA{40, 40, 40, 255})
		
		if s < g.shields {
			var shieldCol color.RGBA
			switch g.shields {
			case 3: shieldCol = color.RGBA{0, 255, 255, 255}
			case 2: shieldCol = color.RGBA{255, 255, 0, 255}
			case 1: shieldCol = color.RGBA{255, 0, 0, 255}
			}
			g.renderer.DrawRect(boxX, shieldBoxY, boxW, boxH, shieldCol)
		}
	}

	modeStr := "1P"
	if g.mode == ModeAIBattle { modeStr = "AI" }
	g.renderer.DrawText(panelX+10, 40, fmt.Sprintf("MODE:  %s", modeStr), color.White)
	g.renderer.DrawText(panelX+10, 70, "TARGET INFO:", color.RGBA{0, 255, 0, 255})
	
	col2X := panelX + 220
	g.renderer.DrawRect(col2X-10, 70, 2, g.config.WindowHeight-70, color.RGBA{0, 100, 0, 255})
	g.renderer.DrawText(col2X, 70, "STORIES:", color.RGBA{0, 255, 255, 255})

	startY := 100
	for i := range g.aliens {
		if i < len(g.aliens) {
			yPos := startY + (i * 70)
			alien := g.aliens[i]
			data := g.alienData[i]
			
			nameColor := color.RGBA{255, 255, 0, 255}
			traitColor := color.RGBA{200, 200, 200, 255}
			var bioColor color.Color = color.White

			if !alien.IsAlive() {
				nameColor = color.RGBA{150, 150, 150, 255} 
				traitColor = color.RGBA{100, 100, 100, 255}
				bioColor = color.RGBA{100, 100, 100, 255}
			}

			g.renderer.DrawText(panelX+10, yPos, alien.Name, nameColor)
			g.renderer.DrawText(panelX+10, yPos+15, "Trait: "+data.Personality, traitColor)
			
			bioShort := data.Bio
			if len(bioShort) > 25 { bioShort = bioShort[:25] + "..." }
			g.renderer.DrawText(panelX+10, yPos+30, bioShort, bioColor)

			g.storyTimer[i]++
			const StoryDuration = 900
			
			switch g.storyState[i] {
			case 0:
				if g.storyTimer[i] > StoryDuration { g.storyState[i] = 1 }
				g.storyAlpha[i] = 1.0
			case 1:
				g.storyAlpha[i] -= 0.05
				if g.storyAlpha[i] <= 0 {
					g.storyAlpha[i] = 0
					g.storyState[i] = 2
					g.storyTimer[i] = 0
					totalOptions := len(data.Stories) + len(data.Family) + len(data.Friends)
					if totalOptions > 0 { g.infoCycles[i] = rand.Intn(totalOptions) }
				}
			case 2:
				g.storyAlpha[i] += 0.05
				if g.storyAlpha[i] >= 1.0 {
					g.storyAlpha[i] = 1.0
					g.storyState[i] = 0
				}
			}
			
			applyAlpha := func(c color.RGBA, a float64) color.RGBA {
				return color.RGBA{
					R: uint8(float64(c.R) * a),
					G: uint8(float64(c.G) * a),
					B: uint8(float64(c.B) * a),
					A: 255,
				}
			}
				
			cycle := g.infoCycles[i]
			var title, person, storyText string
			var titleColor color.RGBA
			
			if cycle < len(data.Stories) {
				title, person, storyText = "PERSONAL MEMORY", alien.Name, data.Stories[cycle]
				titleColor = color.RGBA{100, 255, 100, 255}
			} else {
				rem := cycle - len(data.Stories)
				if rem < len(data.Family) {
					rel := data.Family[rem]
					title, person = fmt.Sprintf("FAMILY: %s", rel.Relation), rel.Name
					if len(rel.Stories) > 0 { storyText = rel.Stories[i % len(rel.Stories)] }
					titleColor = color.RGBA{255, 100, 255, 255}
				} else {
					rem -= len(data.Family)
					if rem < len(data.Friends) {
						rel := data.Friends[rem]
						title, person = fmt.Sprintf("FRIEND: %s", rel.Relation), rel.Name
						if len(rel.Stories) > 0 { storyText = rel.Stories[i % len(rel.Stories)] }
						titleColor = color.RGBA{100, 200, 255, 255}
					}
				}
			}
			
			if storyText != "" {
				tCol, pCol, sCol := titleColor, color.RGBA{255, 255, 255, 255}, color.RGBA{200, 200, 200, 255}
				if !alien.IsAlive() {
					tCol, pCol, sCol = color.RGBA{100, 100, 100, 255}, color.RGBA{100, 100, 100, 255}, color.RGBA{80, 80, 80, 255}
				}
				g.renderer.DrawText(col2X, yPos, title, applyAlpha(tCol, g.storyAlpha[i]))
				g.renderer.DrawText(col2X, yPos+10, person, applyAlpha(pCol, g.storyAlpha[i]))
				if len(storyText) > 0 {
					textColor := applyAlpha(sCol, g.storyAlpha[i])
					g.renderer.DrawText(col2X, yPos+25, storyText[0:min(len(storyText), 35)], textColor)
					if len(storyText) > 35 { g.renderer.DrawText(col2X, yPos+35, storyText[35:min(len(storyText), 70)], textColor) }
				}
			}
		}
	}
	
	if g.state == StateConfirmQuit {
		cx, cy := g.config.WindowWidth/2-200, g.config.WindowHeight/2-50
		g.renderer.DrawRect(cx, cy, 400, 100, color.RGBA{50, 0, 0, 255})
		g.renderer.DrawRect(cx, cy, 400, 2, color.RGBA{255, 0, 0, 255})
		g.renderer.DrawRect(cx, cy+100, 400, 2, color.RGBA{255, 0, 0, 255})
		g.renderer.DrawRect(cx, cy, 2, 100, color.RGBA{255, 0, 0, 255})
		g.renderer.DrawRect(cx+400, cy, 2, 100, color.RGBA{255, 0, 0, 255})
		g.renderer.DrawText(cx+20, cy+30, "ARE YOU SURE YOU WANT TO QUIT?", color.White)
		g.renderer.DrawText(cx+40, cy+60, "PRESS 'C' AGAIN TO LEAVE", color.RGBA{255, 0, 0, 255})
	}

	g.renderer.Present()
}
