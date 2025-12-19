package renderer

import (
	"fmt"
)

func (r *Renderer) presentANSI() {
	fmt.Print("\x1b[2J\x1b[H")
	
	bounds := r.buffer.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	scaleX := 2
	scaleY := 4
	
	termWidth := width / scaleX
	termHeight := height / scaleY
	
	for ty := 0; ty < termHeight; ty++ {
		for tx := 0; tx < termWidth; tx++ {
			px := tx * scaleX
			py := ty * scaleY
			
			var rSum, gSum, bSum uint32
			count := 0
			for dy := 0; dy < scaleY && py+dy < height; dy++ {
				for dx := 0; dx < scaleX && px+dx < width; dx++ {
					c := r.buffer.At(px+dx, py+dy)
					r, g, b, _ := c.RGBA()
					rSum += r >> 8
					gSum += g >> 8
					bSum += b >> 8
					count++
				}
			}
			
			if count > 0 {
				avgR := uint8(rSum / uint32(count))
				avgG := uint8(gSum / uint32(count))
				avgB := uint8(bSum / uint32(count))
				
				colorCode := rgbToANSI256(avgR, avgG, avgB)
				fmt.Printf("\x1b[48;5;%dm ", colorCode)
			} else {
				fmt.Print("\x1b[0m ")
			}
		}
		fmt.Print("\x1b[0m\n")
	}
}

func rgbToANSI256(r, g, b uint8) int {
	if r < 8 && g < 8 && b < 8 {
		return 16
	}
	
	avg := (int(r) + int(g) + int(b)) / 3
	rDiff := abs(int(r) - avg)
	gDiff := abs(int(g) - avg)
	bDiff := abs(int(b) - avg)
	
	if rDiff < 10 && gDiff < 10 && bDiff < 10 {
		gray := avg * 23 / 255
		return 232 + gray
	}
	
	r6 := int(r) * 5 / 255
	g6 := int(g) * 5 / 255
	b6 := int(b) * 5 / 255
	
	return 16 + 36*r6 + 6*g6 + b6
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

