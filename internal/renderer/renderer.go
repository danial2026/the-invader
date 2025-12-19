package renderer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"invaders/internal/logger"
	"github.com/disintegration/gift"
	termbox "github.com/nsf/termbox-go"
)

type Renderer struct {
	SpriteSheet image.Image // Main sprite sheet
	Background  image.Image // Background image
	buffer      *image.RGBA // Current frame buffer
	width       int
	height      int
	termType    TerminalType // Detected terminal type
}

func New(width, height int) (*Renderer, error) {
	logger.Info("INIT: Renderer - Starting initialization")
	
	// Detect terminal type
	termType := DetectTerminal()
	logger.Info("INIT: Terminal - Detected type: %s", termType.String())
	
	// Initialize termbox for input handling
	logger.Info("INIT: Termbox - Initializing...")
	if err := termbox.Init(); err != nil {
		return nil, fmt.Errorf("failed to init termbox: %w", err)
	}
	logger.Info("INIT: Termbox - Initialized successfully")
	
	logger.Info("LOAD: Sprite sheet - Reading imgs/sprites.png")
	sprites, err := loadImage("imgs/sprites.png")
	if err != nil {
		return nil, fmt.Errorf("failed to load sprites: %w", err)
	}
	logger.Info("LOAD: Sprite sheet - Loaded successfully")

	logger.Info("LOAD: Background - Reading imgs/bg.png")
	bg, err := loadImage("imgs/bg.png")
	if err != nil {
		return nil, fmt.Errorf("failed to load background: %w", err)
	}
	logger.Info("LOAD: Background - Loaded successfully")

	// Resize background to cover the entire screen
	logger.Info("INIT: Renderer - Resizing background to %dx%d", width, height)
	g := gift.New(
		gift.Resize(width, height, gift.LanczosResampling),
	)
	dst := image.NewRGBA(g.Bounds(bg.Bounds()))
	g.Draw(dst, bg)
	bg = dst

	logger.Info("INIT: Renderer - Complete (Width: %d, Height: %d)", width, height)
	return &Renderer{
		SpriteSheet: sprites,
		Background:  bg,
		buffer:      image.NewRGBA(image.Rect(0, 0, width, height)),
		width:       width,
		height:      height,
		termType:    termType,
	}, nil
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (r *Renderer) NewFrame() {
	gift.New().Draw(r.buffer, r.Background)
}

func (r *Renderer) DrawSprite(filter *gift.GIFT, x, y int) {
	if filter == nil {
		return
	}
	filter.DrawAt(r.buffer, r.SpriteSheet, image.Pt(x, y), gift.OverOperator)
}

func (r *Renderer) Present() {
	switch r.termType {
	case TerminalITerm2:
		r.presentITerm2()
	case TerminalSixel:
		r.presentSixel()
	case TerminalANSI:
		r.presentANSI()
	default:
		// Fallback to ANSI for unknown terminals
		r.presentANSI()
	}
}

func (r *Renderer) presentITerm2() {
	var buf bytes.Buffer
	png.Encode(&buf, r.buffer)
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	// iTerm2 escape sequence for inline images
	// \x1b[2;0H moves cursor to row 2, column 0
	// \x1b]1337;File=inline=1:<base64>\a displays the image
	fmt.Printf("\x1b[2;0H\x1b]1337;File=inline=1:%s\a", encoded)
}

func (r *Renderer) Buffer() *image.RGBA {
	return r.buffer
}

func (r *Renderer) ApplyMatrixEffect() {
	logger.Info("INIT: Renderer - Applying Matrix effect to background")
	// Create a green-tinted version of sprites
	greenFilter := gift.New(
		gift.ColorFunc(func(r0, g0, b0, a0 float32) (r, g, b, a float32) {
			// Convert to grayscale, then tint green
			gray := 0.299*r0 + 0.587*g0 + 0.114*b0
			return gray * 0.2, gray * 1.0, gray * 0.3, a0
		}),
	)

	// Apply to background
	dst := image.NewRGBA(r.Background.Bounds())
	greenFilter.Draw(dst, r.Background)
	r.Background = dst
	logger.Info("INIT: Renderer - Matrix effect applied")

	// Apply to sprites (optional - keep original colors for variety)
	// Uncomment if you want full Matrix effect on sprites too:
	// spriteDst := image.NewRGBA(r.SpriteSheet.Bounds())
	// greenFilter.Draw(spriteDst, r.SpriteSheet)
	// r.SpriteSheet = spriteDst
}

func (r *Renderer) DrawText(x, y int, text string, col color.Color) {
	DrawText(r.buffer, x, y, text, col)
}

func (r *Renderer) DrawTextScaled(x, y int, text string, col color.Color, scale float64) {
	DrawTextScaled(r.buffer, x, y, text, col, scale)
}

func (r *Renderer) DrawRect(x, y, w, h int, col color.Color) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			r.buffer.Set(x+dx, y+dy, col)
		}
	}
}

func (r *Renderer) Close() {
	logger.Info("Renderer closing")
	termbox.Close()
}
