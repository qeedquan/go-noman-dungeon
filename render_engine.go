package main

import (
	"log"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

const (
	VIEW_W = 224
	VIEW_H = 128
)

type Display struct {
	*sdl.Window
	*sdl.Renderer
}

type RenderEngine struct {
	screen     *Display
	font       *sdlttf.Font
	fullscreen bool
	cache      []Texture
}

func newDisplay(w, h int, wflag sdl.WindowFlags) (*Display, error) {
	window, renderer, err := sdl.CreateWindowAndRenderer(w, h, wflag)
	if err != nil {
		return nil, err
	}
	return &Display{window, renderer}, nil
}

func newRenderEngine() *RenderEngine {
	desktop, err := sdl.GetDesktopDisplayMode(0)
	if err != nil {
		log.Fatal("sdl: ", err)
	}

	width := int(float64(desktop.W)/1.5/VIEW_W) * VIEW_W
	height := int(float64(desktop.H)/1.5/VIEW_H) * VIEW_H

	// linear interpolation and others make jaggy lines
	// probably because the texture map for tiles has
	// the context tiles transparent below the main tiles
	// so it causes to the interpolator to use it
	// might try to fix it by separating the 2 tileset (????)
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "nearest")

	wflag := sdl.WINDOW_RESIZABLE
	screen, err := newDisplay(width, height, wflag)
	if err != nil {
		log.Fatal("sdl: ", err)
	}
	screen.SetTitle("Noman's Dungeon")
	screen.SetLogicalSize(VIEW_W, VIEW_H)

	err = sdlttf.Init()
	if err != nil {
		log.Fatal("sdl: ", err)
	}

	name := dpath("dejavu_sans_mono.ttf")
	font, err := sdlttf.OpenFont(name, 10)
	if err != nil {
		log.Fatal("sdl: ", name, ": ", err)
	}

	sdl.ShowCursor(0)

	return &RenderEngine{
		screen: screen,
		font:   font,
	}
}

func (p *RenderEngine) clear() {
	p.screen.SetDrawColor(sdlcolor.Black)
	p.screen.Clear()
}

func (p *RenderEngine) commitFrame() {
	p.screen.Present()
}

func (p *RenderEngine) toggleFullscreen() {
	if p.fullscreen {
		p.screen.SetFullscreen(0)
	} else {
		p.screen.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
	}
	p.fullscreen = !p.fullscreen
}

func (p *RenderEngine) loadImage(image *Image, filename string) {
	if image == nil {
		return
	}

	texture := p.cacheLookup(filename)
	if texture != nil {
		image.filename = filename
		image.texture = texture
		image.init(p)
		return
	}

	texture, err := sdlimage.LoadTextureFile(p.screen.Renderer, filename)
	if err == nil {
		texture.SetBlendMode(sdl.BLENDMODE_BLEND)
		image.filename = filename
		image.texture = texture
		p.cacheStore(image)
	}
	image.init(p)

	if err != nil {
		log.Fatal("image: ", err)
	}
}

func (p *RenderEngine) renderImage(image *Image) {
	p.screen.Copy(image.texture, &image.clip, &image.dest)
}

func (p *RenderEngine) renderText(image *Image, text string, color sdl.Color) {
	if p.font == nil || image == nil {
		return
	}

	image.filename = ""
	if image.texture != nil {
		image.texture.Destroy()
		image.texture = nil
	}

	surface, err := p.font.RenderUTF8Solid(text, color)
	if err != nil {
		log.Fatal("sdl: ", err)
	}
	image.texture, err = p.screen.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatal("sdl: ", err)
	}
	image.texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	image.init(p)
}

func (p *RenderEngine) cacheLookup(filename string) *sdl.Texture {
	for i := range p.cache {
		c := &p.cache[i]
		if filename == c.filename {
			c.ref++
			return c.texture
		}
	}
	return nil
}

func (p *RenderEngine) cacheStore(image *Image) {
	if image.filename != "" && image.texture != nil {
		tex := Texture{
			filename: image.filename,
			texture:  image.texture,
			ref:      1,
		}
		p.cache = append(p.cache, tex)
	}
}
