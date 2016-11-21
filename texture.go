package main

import "github.com/qeedquan/go-media/sdl"

type Texture struct {
	filename string
	texture  *sdl.Texture
	ref      int
}
