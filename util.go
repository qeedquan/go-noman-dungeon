package main

import (
	"math"
	"path/filepath"

	"github.com/qeedquan/go-media/sdl"
)

func dpath(path ...string) string {
	return filepath.Join(dataDir, filepath.Join(path...))
}

func distance(p1, p2 sdl.Point) float64 {
	x := float64(p2.X - p1.X)
	y := float64(p2.Y - p1.Y)
	return math.Sqrt(x*x + y*y)
}

func clamp(x, a, b int) int {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}
