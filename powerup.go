package main

import (
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
)

const (
	POWERUP_ATTACK = iota
	POWERUP_DEFENSE
	POWERUP_POTION
	POWERUP_TREASURE
)

type Powerup struct {
	pos         sdl.Point
	typ         int
	amount      int
	animNormal  Animation
	currentAnim Animation
}

func newPowerup(typ int) *Powerup {
	p := &Powerup{}
	p.typ = typ

	switch p.typ {
	case POWERUP_ATTACK:
		p.animNormal.load(dpath("powerups.png"), 0, 1, 100, 0, TILE_SIZE, TILE_SIZE, 0, 0)
		p.amount = rand.Intn(5) + 1

	case POWERUP_DEFENSE:
		p.animNormal.load(dpath("powerups.png"), 1, 1, 100, 0, TILE_SIZE, TILE_SIZE, 0, 0)
		p.amount = rand.Intn(5) + 1

	case POWERUP_POTION:
		p.animNormal.load(dpath("powerups.png"), 2, 1, 100, 0, TILE_SIZE, TILE_SIZE, 0, 0)
		p.amount = rand.Intn(40) + 10

	case POWERUP_TREASURE:
		p.animNormal.load(dpath("powerups.png"), 3, 1, 100, 0, TILE_SIZE, TILE_SIZE, 0, 0)
		p.amount = 1
	}

	p.currentAnim.setTo(&p.animNormal)
	return p
}

func (p *Powerup) logic() {
	p.currentAnim.logic()
}

func (p *Powerup) render() {
	p.currentAnim.setPos(int(p.pos.X*TILE_SIZE), int(p.pos.Y*TILE_SIZE))
	p.currentAnim.render()
}

func (p *Powerup) setPos(x, y int) {
	p.pos = sdl.Point{int32(x), int32(y)}
}
