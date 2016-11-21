package main

import "math/rand"

const (
	PLAYER_MIN_ATTACK = 1
	PLAYER_MAX_ATTACK = 9999

	PLAYER_MIN_DEFENSE = 1
	PLAYER_MAX_DEFENSE = 9999
)

type Player struct {
	Entity
	hasTreasure bool
}

func newPlayer() Player {
	p := Player{}
	p.Entity = newEntity()
	p.animNormal.load(dpath("player.png"), 0, 1, 100, 0, TILE_SIZE, TILE_SIZE, 0, 0)
	p.animHurt.load(dpath("player.png"), 1, 4, 266, 1, TILE_SIZE, TILE_SIZE, 0, 0)
	p.animDie.load(dpath("player.png"), 2, 4, 266, 1, TILE_SIZE, TILE_SIZE, 0, 0)
	return p
}

func (p *Player) init() {
	p.currentAnim.setTo(&p.animNormal)
	p.isTurn = false
	p.hp, p.maxhp = 50, 50
	p.attack = 10
	p.defense = 10
	if godMode {
		p.attack = PLAYER_MAX_ATTACK
		p.defense = PLAYER_MAX_DEFENSE
	}
	p.hasTreasure = false
	p.animating = false
}

func (p *Player) startTurn() {
	p.isTurn = true
}

func (p *Player) actionMove(x, y int) {
	p.setPos(x, y)
	p.isTurn = false
}

func (p *Player) actionAttack(e *Enemy) {
	if e == nil {
		return
	}

	dmg := (rand.Intn(p.attack) + p.attack/2) - (rand.Intn(e.defense) + e.defense/2)
	if dmg <= 0 {
		dmg = 1
	}

	e.takeDamage(dmg)
}

func (p *Player) bonusAttack(amount int) {
	p.attack += amount
	p.attack = clamp(p.attack, PLAYER_MIN_ATTACK, PLAYER_MAX_ATTACK)
}

func (p *Player) bonusDefense(amount int) {
	p.defense += amount
	p.defense = clamp(p.defense, PLAYER_MIN_DEFENSE, PLAYER_MAX_DEFENSE)
}

func (p *Player) bonusHP(amount int) {
	if amount == 0 {
		return
	}

	realAmount := (amount * p.maxhp) / 100
	p.hp += realAmount
	if p.hp > p.maxhp {
		p.hp = p.maxhp
	}

	// poisons can't kill the player outright
	if p.hp <= 0 {
		p.hp = 1
	}
}

func (p *Player) bonusTreasure() {
	p.hasTreasure = true
}

func (p *Player) takeDamage(dmg int) {
	if godMode {
		return
	}
	p.Entity.takeDamage(dmg)
}
