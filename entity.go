package main

import "github.com/qeedquan/go-media/sdl"

type Entity struct {
	pos    sdl.Point
	isTurn bool

	hp      int
	maxhp   int
	attack  int
	defense int

	animNormal  Animation
	animHurt    Animation
	animDie     Animation
	currentAnim Animation

	hpBar     Image
	animating bool

	sfxHurt Sound
	sfxDie  Sound
	sfxMove Sound
}

func newEntity() Entity {
	e := Entity{}
	e.hpBar.load(renderEngine, dpath("hp_bar.png"))

	soundEngine.loadSound(&e.sfxHurt, dpath("hurt.wav"))
	soundEngine.loadSound(&e.sfxDie, dpath("die.wav"))
	soundEngine.loadSound(&e.sfxMove, dpath("move.wav"))
	return e
}

func (e *Entity) setPos(x, y int) {
	e.pos = sdl.Point{int32(x), int32(y)}
}

func (e *Entity) takeDamage(dmg int) {
	if dmg <= 0 {
		return
	}

	e.hp -= dmg
	if e.hp <= 0 {
		e.hp = 0
		e.currentAnim.setTo(&e.animDie)
		e.sfxDie.play()
		e.animating = true
	} else {
		e.currentAnim.setTo(&e.animHurt)
		e.sfxHurt.play()
		e.animating = true
	}
}

func (e *Entity) isAnimating() bool {
	return e.animating
}

func (e *Entity) render() {
	e.currentAnim.logic()

	if e.animating && e.currentAnim.isFinished() {
		e.currentAnim.setTo(&e.animNormal)
		e.animating = false
	}

	if e.hp > 0 || e.animating {
		e.currentAnim.setPos(int(e.pos.X)*TILE_SIZE, int(e.pos.Y)*TILE_SIZE)
		e.currentAnim.render()

		e.hpBar.setPos(int(e.pos.X)*TILE_SIZE, int(e.pos.Y)*TILE_SIZE)

		barLength := int(float64(e.hp) / float64(e.maxhp) * float64(e.hpBar.getWidth()))
		if e.hp > 0 && barLength < 1 {
			barLength = 1
		}

		e.hpBar.setClip(0, 0, barLength, e.hpBar.getHeight())
		e.hpBar.render()
	}
}
