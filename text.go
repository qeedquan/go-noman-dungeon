package main

import (
	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

type Text struct {
	pos         sdl.Point
	text        string
	image       Image
	imageShadow Image
}

func (t *Text) setText(text string) {
	if t.text == text {
		return
	}

	t.text = text
	renderEngine.renderText(&t.image, text, sdlcolor.White)
	renderEngine.renderText(&t.imageShadow, text, sdlcolor.Black)
	t.setPos(int(t.pos.X), int(t.pos.Y))
}

func (t *Text) setPos(x, y int) {
	t.pos = sdl.Point{int32(x), int32(y)}
	t.image.setPos(x, y)
	t.imageShadow.setPos(x+1, y+1)
}

func (t *Text) render() {
	t.imageShadow.render()
	t.image.render()
}
