package renderer

import "os"

type TerminalType int

const (
	TerminalUnknown TerminalType = iota
	TerminalITerm2
	TerminalSixel
	TerminalANSI
)

func (t TerminalType) String() string {
	switch t {
	case TerminalITerm2:
		return "iTerm2"
	case TerminalSixel:
		return "Sixel"
	case TerminalANSI:
		return "ANSI"
	default:
		return "Unknown"
	}
}

func DetectTerminal() TerminalType {
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")
	
	if termProgram == "iTerm.app" {
		return TerminalITerm2
	}
	
	if term == "xterm" || term == "xterm-256color" || 
	   term == "mlterm" || term == "foot" || term == "wezterm" {
		return TerminalSixel
	}
	
	return TerminalANSI
}

