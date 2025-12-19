package renderer

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

func (r *Renderer) presentSixel() {
	fmt.Print("\x1b[2J\x1b[H")
	sixelData := encodeSixel(r.buffer)
	fmt.Print(sixelData)
}

func encodeSixel(img *image.RGBA) string {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	var buf bytes.Buffer
	
	buf.WriteString("\x1bPq")
	buf.WriteString("\"1;1;")
	buf.WriteString(fmt.Sprintf("%d;%d", width, height))
	
	colorMap := make(map[color.RGBA]int)
	colorIndex := 0
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.RGBAAt(x, y)
			if _, exists := colorMap[c]; !exists && colorIndex < 256 {
				colorMap[c] = colorIndex
				r := int(float64(c.R) * 100.0 / 255.0)
				g := int(float64(c.G) * 100.0 / 255.0)
				b := int(float64(c.B) * 100.0 / 255.0)
				buf.WriteString(fmt.Sprintf("#%d;2;%d;%d;%d", colorIndex, r, g, b))
				colorIndex++
			}
		}
	}
	
	for y := 0; y < height; y += 6 {
		for c, idx := range colorMap {
			hasPixels := false
			var rowBuf bytes.Buffer
			
			for x := 0; x < width; x++ {
				sixelChar := 0
				for dy := 0; dy < 6 && y+dy < height; dy++ {
					pixel := img.RGBAAt(x, y+dy)
					if pixel == c {
						sixelChar |= (1 << dy)
					}
				}
				
				if sixelChar > 0 {
					hasPixels = true
					rowBuf.WriteByte(byte(sixelChar + 63))
				} else {
					rowBuf.WriteByte('?')
				}
			}
			
			if hasPixels {
				buf.WriteString(fmt.Sprintf("#%d", idx))
				buf.Write(rowBuf.Bytes())
			}
		}
		
		buf.WriteString("$")
		buf.WriteString("-")
	}
	
	buf.WriteString("\x1b\\")
	
	return buf.String()
}

