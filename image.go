package main

import "github.com/qeedquan/go-media/sdl"

type Image struct {
	filename string
	texture  *sdl.Texture

	dest sdl.Rect
	clip sdl.Rect

	renderEngine *RenderEngine
	width        int
	height       int
}

func newImage() Image {
	return Image{}
}

func (i *Image) init(renderEngine *RenderEngine) {
	i.renderEngine = renderEngine
	_, _, i.width, i.height, _ = i.texture.Query()
	i.dest.W, i.clip.W = int32(i.width), int32(i.width)
	i.dest.H, i.clip.H = int32(i.height), int32(i.height)
}

func (i *Image) load(renderEngine *RenderEngine, filename string) {
	if renderEngine != nil {
		renderEngine.loadImage(i, filename)
	}
}

func (i *Image) getWidth() int {
	return i.width
}

func (i *Image) getHeight() int {
	return i.height
}

func (i *Image) setPos(x, y int) {
	i.dest.X = int32(x)
	i.dest.Y = int32(y)
}

func (i *Image) setClip(x, y, w, h int) {
	i.clip = sdl.Rect{int32(x), int32(y), int32(w), int32(h)}
	i.dest.W, i.dest.H = int32(w), int32(h)
}

func (i *Image) render() {
	if i.renderEngine != nil {
		i.renderEngine.renderImage(i)
	}
}

func (i *Image) ref() {
	if i.renderEngine != nil {
		i.renderEngine.cacheLookup(i.filename)
	}
}
