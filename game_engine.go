package main

import "fmt"

type GameEngine struct {
	done      bool
	mapEngine *MapEngine
	player    Player

	statusText Text

	treasureStatus Image
}

func newGameEngine() *GameEngine {
	g := &GameEngine{}
	g.player = newPlayer()
	g.mapEngine = newMapEngine()
	g.mapEngine.init(&g.player)

	g.statusText.setPos(2, VIEW_H-14)

	g.treasureStatus.load(renderEngine, dpath("treasure_status.png"))
	g.treasureStatus.setPos(VIEW_W-TILE_SIZE, VIEW_H-TILE_SIZE)

	return g
}

func (p *GameEngine) logic() {
	if inputEngine.pressing[EXIT] {
		p.done = true
		return
	}

	if inputEngine.pressing[FULLSCREEN_TOGGLE] && !inputEngine.lock[FULLSCREEN_TOGGLE] {
		inputEngine.lock[FULLSCREEN_TOGGLE] = true
		renderEngine.toggleFullscreen()
	}

	text := fmt.Sprintf("Att: %v | Def: %v | Lvl: %v", p.player.attack, p.player.defense, p.mapEngine.getCurrentLevel())
	p.statusText.setText(text)

	if p.mapEngine.gameState == GAME_PLAY {
		if p.player.isTurn {
			if inputEngine.pressing[UP] && !inputEngine.lock[UP] {
				inputEngine.lock[UP] = true
				p.mapEngine.moveCursor(CUR_UP)
			}

			if inputEngine.pressing[DOWN] && !inputEngine.lock[DOWN] {
				inputEngine.lock[DOWN] = true
				p.mapEngine.moveCursor(CUR_DOWN)
			}

			if inputEngine.pressing[LEFT] && !inputEngine.lock[LEFT] {
				inputEngine.lock[LEFT] = true
				p.mapEngine.moveCursor(CUR_LEFT)
			}

			if inputEngine.pressing[RIGHT] && !inputEngine.lock[RIGHT] {
				inputEngine.lock[RIGHT] = true
				p.mapEngine.moveCursor(CUR_RIGHT)
			}

			if inputEngine.pressing[ACTION] && !inputEngine.lock[ACTION] {
				inputEngine.lock[ACTION] = true
				if p.mapEngine.playerAction() {
					p.mapEngine.enemyStartTurn()
				}
			}
		} else {
			if p.mapEngine.enemyAction() {
				p.mapEngine.playerStartTurn()
			}
		}
	} else {
		// win and lose have the same game logic
		if inputEngine.pressing[RESTART] && !inputEngine.lock[RESTART] {
			inputEngine.lock[RESTART] = true
			p.mapEngine.init(&p.player)
		}
	}
}

func (p *GameEngine) render() {
	p.mapEngine.render()
	p.statusText.render()

	if p.player.hasTreasure {
		p.treasureStatus.setClip(0, TILE_SIZE, TILE_SIZE, TILE_SIZE)
	} else {
		p.treasureStatus.setClip(0, 0, TILE_SIZE, TILE_SIZE)
	}

	p.treasureStatus.render()
}
