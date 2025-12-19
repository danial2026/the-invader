package input

import (
	termbox "github.com/nsf/termbox-go"
)

type Key int

const (
	KeyNone Key = iota
	KeyLeft
	KeyRight
	KeyUp
	KeyDown
	KeyShootLeft
	KeyShootRight
	KeyQuit
	KeyPause
	KeyEnter
	KeySpace
)

type Handler struct {
	events chan termbox.Event
}

func New() *Handler {
	h := &Handler{
		events: make(chan termbox.Event, 1000),
	}

	go func() {
		for {
			h.events <- termbox.PollEvent()
		}
	}()

	return h
}

func (h *Handler) Poll() Key {
	select {
	case ev := <-h.events:
		if ev.Type == termbox.EventKey {
			return h.translateKey(ev)
		}
	default:
	}
	return KeyNone
}

func (h *Handler) translateKey(ev termbox.Event) Key {
	switch ev.Key {
	case termbox.KeyArrowLeft:
		return KeyLeft
	case termbox.KeyArrowRight:
		return KeyRight
	case termbox.KeyArrowUp:
		return KeyUp
	case termbox.KeyArrowDown:
		return KeyDown
	case termbox.KeyEsc:
		return KeyQuit
	case termbox.KeyEnter:
		return KeyEnter
	case termbox.KeySpace:
		return KeySpace
	}

	switch ev.Ch {
	case '\r', '\n':
		return KeyEnter
	case ' ':
		return KeySpace
	case 's', 'S':
		return KeyDown
	case 'w', 'W':
		return KeyUp
	case 'a', 'A':
		return KeyLeft
	case 'd', 'D':
		return KeyRight
	case 'q', 'Q':
		return KeyShootLeft
	case 'e', 'E':
		return KeyShootRight
	case 'c', 'C':
		return KeyQuit
	case 'p', 'P':
		return KeyPause
	}

	return KeyNone
}


