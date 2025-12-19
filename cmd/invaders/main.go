package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"invaders/internal/game"
	"invaders/internal/logger"
)

func main() {
	aiMode := flag.Bool("ai", false, "Run in AI Battle mode (spectator)")
	flag.Parse()

	if err := logger.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
	}
	defer logger.Close()

	logger.Info("Starting The Invader")

	cfg := game.DefaultConfig()
	g, err := game.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize game: %v\n", err)
		os.Exit(1)
	}
	defer g.Close()

	if *aiMode {
		g.SetMode(game.ModeAIBattle)
		logger.Info("Starting in AI Battle Mode")
		fmt.Println("Starting AI Battle mode...")
	} else {
		g.SetMode(game.ModeSinglePlayer)
		logger.Info("Starting in Single Player Mode")
	}

	// Panic Recovery
	defer func() {
		if r := recover(); r != nil {
			f, _ := os.Create("crash.log")
			if f != nil {
				fmt.Fprintf(f, "Panic: %v\nStack: %s\n", r, debug.Stack())
				f.Close()
			}
			logger.Error("PANIC: %v\nStack: %s", r, debug.Stack())
			fmt.Fprintf(os.Stderr, "Game crashed! See crash.log for details.\n")
			
			if g != nil {
				g.Close()
			} else {
				logger.Close()
			}
			os.Exit(1)
		}
	}()

	g.Run()
}
