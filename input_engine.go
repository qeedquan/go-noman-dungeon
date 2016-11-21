package main

import "github.com/qeedquan/go-media/sdl"

const (
	ACTION = iota
	EXIT
	LEFT
	RIGHT
	UP
	DOWN
	RESTART
	FULLSCREEN_TOGGLE
	KEY_COUNT
)

type InputEngine struct {
	binding  [KEY_COUNT]sdl.Keycode
	pressing [KEY_COUNT]bool
	lock     [KEY_COUNT]bool
	done     bool
}

func newInputEngine() *InputEngine {
	return &InputEngine{
		binding: [KEY_COUNT]sdl.Keycode{
			ACTION:            sdl.K_RETURN,
			EXIT:              sdl.K_ESCAPE,
			LEFT:              sdl.K_LEFT,
			RIGHT:             sdl.K_RIGHT,
			UP:                sdl.K_UP,
			DOWN:              sdl.K_DOWN,
			RESTART:           sdl.K_r,
			FULLSCREEN_TOGGLE: sdl.K_f,
		},
	}
}

func (p *InputEngine) logic() {
	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}
		switch ev := ev.(type) {
		case sdl.KeyDownEvent:
			for i, key := range p.binding {
				if ev.Sym == key {
					p.pressing[i] = true
				}
			}

		case sdl.KeyUpEvent:
			for i, key := range p.binding {
				if ev.Sym == key {
					p.pressing[i] = false
					p.lock[i] = false
				}
			}

		case sdl.QuitEvent:
			p.done = true
		}
	}
}
