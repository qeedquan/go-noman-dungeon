package main

import (
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
)

const (
	ENEMY_COOLDOWN = 15
)

const (
	ENEMY_SLIME = iota
)

type Enemy struct {
	Entity
	player      *Player
	actionTicks int
}

func newEnemy(typ int) *Enemy {
	e := &Enemy{}
	e.Entity = newEntity()
	if typ == ENEMY_SLIME {
		e.animNormal.load(dpath("slime.png"), 0, 1, 100, 0, TILE_SIZE, TILE_SIZE, 0, 0)
		e.animHurt.load(dpath("slime.png"), 1, 4, 266, 1, TILE_SIZE, TILE_SIZE, 0, 0)
		e.animDie.load(dpath("slime.png"), 2, 4, 266, 1, TILE_SIZE, TILE_SIZE, 0, 0)
		e.hp = rand.Intn(10) + 10
		e.maxhp = e.hp
		e.attack = rand.Intn(5) + 10
		e.defense = rand.Intn(5) + 1
	}

	e.init()
	return e
}

func (e *Enemy) init() {
	e.isTurn = false
	e.currentAnim.setTo(&e.animNormal)

	e.player = nil
	e.animating = false
}

func (e *Enemy) startTurn() {
	e.isTurn = true
	e.actionTicks = ENEMY_COOLDOWN
}

func (e *Enemy) setPlayer(player *Player) {
	e.player = player
}

func (e *Enemy) actionAttack() {
	if e.player == nil {
		return
	}

	dmg := (rand.Intn(e.attack) + e.attack/2) - (rand.Intn(e.player.defense) + e.player.defense/2)
	if dmg <= 0 {
		dmg = 1
	}

	e.player.takeDamage(dmg)
}

func (e *Enemy) isNearPlayer(range_ int) bool {
	if e.player == nil {
		return false
	}

	horizontal := sdl.Point{e.player.pos.X, e.pos.Y}
	vertical := sdl.Point{e.pos.X, e.player.pos.Y}

	if distance(horizontal, e.pos) > float64(range_) || distance(vertical, e.pos) > float64(range_) {
		return false
	}

	return true
}
