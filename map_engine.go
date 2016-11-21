package main

import (
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
)

const (
	MAP_W = 14
	MAP_H = 8
)

const (
	CUR_UP = iota
	CUR_DOWN
	CUR_LEFT
	CUR_RIGHT
)

const (
	GAME_PLAY = iota
	GAME_WIN
	GAME_LOSE
)

const (
	TILE_FLOOR = iota
	TILE_STAIRS_UP
	TILE_STAIRS_DOWN
	TILE_STAIRS_BLOCKED
	TILE_WALL_T
	TILE_WALL_B
	TILE_WALL_TL
	TILE_WALL_TR
	TILE_WALL_L
	TILE_WALL_R
	TILE_WALL_BL
	TILE_WALL_BR
	TILE_ROCK
	TILE_CRATE
)

const (
	CONTEXT_WALKABLE = iota
	CONTEXT_ENEMY
)

type MapEngine struct {
	gameState int
	player    *Player

	enemies     [][]*Enemy
	activeEnemy int

	powerups [][]*Powerup

	cursor    Image
	cursorPos sdl.Point

	currentLevel int
	tileset      Image
	levels       [][MAP_H][MAP_W]int
	contextTiles [MAP_H][MAP_W]int

	msg       Text
	firstTurn bool

	hud Image

	dungeonDepth int

	sfxCursor   Sound
	sfxPowerup  Sound
	sfxTreasure Sound
	sfxWin      Sound
}

func newMapEngine() *MapEngine {
	p := &MapEngine{}
	p.tileset.load(renderEngine, dpath("tileset.png"))
	p.cursor.load(renderEngine, dpath("cursor.png"))

	p.hud.load(renderEngine, dpath("hud.png"))

	soundEngine.loadSound(&p.sfxCursor, dpath("cursor.wav"))
	soundEngine.loadSound(&p.sfxPowerup, dpath("powerup.wav"))
	soundEngine.loadSound(&p.sfxTreasure, dpath("treasure.wav"))
	soundEngine.loadSound(&p.sfxWin, dpath("victory.wav"))

	return p
}

func (p *MapEngine) init(player *Player) {
	p.clear()

	p.player = player
	p.player.init()

	p.activeEnemy = -1
	p.dungeonDepth = dungeonDepth
	p.currentLevel = 0

	p.msg.setText("URLD = cursor. Enter = action.")
	p.msg.setPos(2, 0)
	p.firstTurn = true

	p.gameState = GAME_PLAY

	p.nextLevel()
	p.playerStartTurn()
}

func (p *MapEngine) clear() {
	p.enemies = nil
	p.powerups = nil
	p.levels = nil
}

func (p *MapEngine) prevLevel() {
	if p.currentLevel-1 >= 0 {
		p.currentLevel--
		for i := 0; i < MAP_H; i++ {
			for j := 0; j < MAP_W; j++ {
				if p.levels[p.currentLevel][i][j] == TILE_STAIRS_DOWN {
					p.player.setPos(j, i)
					p.cursorPos = p.player.pos
					if p.player.hasTreasure {
						p.spawnEnemies()
						p.spawnPowerups()
						p.levels[p.currentLevel][i][j] = TILE_STAIRS_BLOCKED
					}
					p.setContextTiles()
					return
				}
			}
		}
	}
}

func (p *MapEngine) nextLevel() {
	if p.currentLevel+1 < len(p.levels) {
		// load existing level
		p.currentLevel++
		foundStairs := false
		for i := 0; i < MAP_H; i++ {
			for j := 0; j < MAP_W; j++ {
				if p.levels[p.currentLevel][i][j] == TILE_STAIRS_UP {
					p.player.setPos(j, i)
					foundStairs = true
					break
				}
			}
			if foundStairs {
				break
			}
		}
	} else {
		// generate new level
		playerPos := sdl.Point{1, 1}
		stairsPos := sdl.Point{int32(MAP_W - 2), int32(rand.Intn(MAP_H-2) + 1)}

		if p.player != nil {
			playerPos = sdl.Point{1, int32(rand.Intn(MAP_H-2) + 1)}
			p.player.setPos(int(playerPos.X), int(playerPos.Y))
		}

		var room [MAP_H][MAP_W]int
		for i := 0; i < MAP_H; i++ {
			for j := 0; j < MAP_W; j++ {
				room[i][j] = TILE_FLOOR

				// borders
				switch {
				case i == 0 && j == 0:
					room[i][j] = TILE_WALL_TL
				case i == MAP_H-1 && j == 0:
					room[i][j] = TILE_WALL_BL
				case i == 0 && j == MAP_W-1:
					room[i][j] = TILE_WALL_TR
				case i == MAP_H-1 && j == MAP_W-1:
					room[i][j] = TILE_WALL_BR
				case i == 0:
					room[i][j] = TILE_WALL_T
				case i == MAP_H-1:
					room[i][j] = TILE_WALL_B
				case j == 0:
					room[i][j] = TILE_WALL_L
				case j == MAP_W-1:
					room[i][j] = TILE_WALL_R
				default:
					if rand.Int()%4 == 0 {
						room[i][j] = rand.Int()%2 + TILE_ROCK
					}
				}

				// up stairs
				if i == int(playerPos.Y) && j == int(playerPos.X) {
					room[i][j] = TILE_STAIRS_UP
				}

				if p.dungeonDepth > len(p.levels)+1 {
					// down stairs
					if i == int(stairsPos.Y) && j == int(stairsPos.X) {
						room[i][j] = TILE_STAIRS_DOWN
					}
				}
			}
		}

		p1 := playerPos
		p2 := stairsPos
		for p1 != p2 {
			if p1.X < MAP_W-2 {
				if rand.Intn(2) == 0 {
					p1.X++
				} else {
					if p1.Y > 1 && rand.Int()%2 == 0 {
						p1.Y--
					} else if p1.Y < MAP_H-2 && rand.Int()%2 == 0 {
						p1.Y++
					}
				}
			} else {
				if p1.Y > 1 && rand.Int()%2 == 0 {
					p1.Y--
				} else if p1.Y < MAP_H-2 && rand.Int()%2 == 0 {
					p1.Y++
				}
			}

			if room[p1.Y][p1.X] != TILE_STAIRS_UP && room[p1.Y][p1.X] != TILE_STAIRS_DOWN {
				room[p1.Y][p1.X] = TILE_FLOOR
			}
		}

		p.levels = append(p.levels, room)
		p.currentLevel = len(p.levels) - 1
		p.spawnEnemies()
		p.spawnPowerups()
		p.spawnTreasure(stairsPos)
	}

	p.cursorPos = p.player.pos
	p.setContextTiles()
}

func (p *MapEngine) playerStartTurn() {
	if p.player.isAnimating() {
		return
	}

	p.player.startTurn()
	p.setContextTiles()
	p.cursorPos = p.player.pos
	if !p.firstTurn {
		p.msg.setText("Player turn")
	}
}

func (p *MapEngine) playerAction() bool {
	if p.firstTurn {
		p.firstTurn = false
	}

	if p.player.pos == p.cursorPos {
		return false
	}

	dist := distance(p.player.pos, p.cursorPos)
	if dist <= 1 && p.contextTiles[p.cursorPos.Y][p.cursorPos.X] == CONTEXT_WALKABLE {
		p.player.actionMove(int(p.cursorPos.X), int(p.cursorPos.Y))
		p.player.sfxMove.play()
		p.checkPowerup()
		if p.levels[p.currentLevel][p.cursorPos.Y][p.cursorPos.X] == TILE_STAIRS_UP {
			if p.player.hasTreasure && p.currentLevel == 0 {
				p.gameState = GAME_WIN
				p.msg.setText("You win! R to play again.")
				p.sfxWin.play()
			} else {
				p.prevLevel()
			}
			return false
		} else if p.levels[p.currentLevel][p.cursorPos.Y][p.cursorPos.X] == TILE_STAIRS_DOWN {
			p.nextLevel()
			return false
		}

		return true
	} else if dist <= 1 && p.contextTiles[p.cursorPos.Y][p.cursorPos.X] == CONTEXT_ENEMY {
		e := p.getEnemy(int(p.cursorPos.X), int(p.cursorPos.Y))
		if e != nil {
			p.player.actionAttack(e)
		}
		return true
	}

	return false
}

func (p *MapEngine) spawnEnemies() {
	if len(p.enemies) < len(p.levels) {
		p.enemies = append(p.enemies, make([][]*Enemy, len(p.levels)-len(p.enemies))...)
	}

	spawnCount := p.currentLevel + 1
	if p.player.hasTreasure {
		spawnCount = p.dungeonDepth + (p.dungeonDepth - spawnCount)
	}
	spawnCount = rand.Int()%spawnCount + 3

	failCount := 10
	for spawnCount > 0 {
		var spawnPos sdl.Point
		spawnPos.Y = int32(rand.Intn(MAP_H-2) + 1)

		if p.player.hasTreasure {
			spawnPos.X = int32(rand.Intn(MAP_W/2) + 1)
		} else {
			spawnPos.X = int32(rand.Intn(MAP_W/2) + MAP_W/2 - 1)
		}

		if p.isWalkable(int(spawnPos.X), int(spawnPos.Y)) {
			e := newEnemy(ENEMY_SLIME)
			e.setPos(int(spawnPos.X), int(spawnPos.Y))
			e.setPlayer(p.player)
			p.enemies[p.currentLevel] = append(p.enemies[p.currentLevel], e)
		} else {
			failCount--
		}

		if failCount == 0 {
			break
		}

		spawnCount--
	}
}

func (p *MapEngine) spawnPowerups() {
	if len(p.powerups) < len(p.levels) {
		p.powerups = append(p.powerups, make([][]*Powerup, len(p.levels)-len(p.powerups))...)
	}

	spawnCount := p.currentLevel + 5
	if p.player.hasTreasure {
		spawnCount = p.dungeonDepth + (p.dungeonDepth - spawnCount)
	}
	spawnCount = rand.Int()%spawnCount + 3
	if spawnCount > 3 {
		spawnCount = 3
	}

	failCount := 10
	for spawnCount > 0 {
		spawnPos := sdl.Point{
			int32(rand.Intn(MAP_W-3) + 1),
			int32(rand.Intn(MAP_H-2) + 1),
		}

		typ := rand.Intn(3)

		if p.isWalkable(int(spawnPos.X), int(spawnPos.Y)) && !p.isPowerup(int(spawnPos.X), int(spawnPos.Y)) {
			up := newPowerup(typ)
			up.setPos(int(spawnPos.X), int(spawnPos.Y))
			p.powerups[p.currentLevel] = append(p.powerups[p.currentLevel], up)
		} else {
			failCount--
		}

		if failCount == 0 {
			break
		}

		spawnCount--
	}
}

func (p *MapEngine) spawnTreasure(pos sdl.Point) {
	if p.dungeonDepth > len(p.levels) {
		return
	}

	if len(p.powerups) < len(p.levels) {
		p.powerups = append(p.powerups, make([][]*Powerup, len(p.levels)-len(p.powerups))...)
	}

	spawnPos := pos
	spawnCount := 1
	typ := POWERUP_TREASURE

	for spawnCount > 0 {
		if p.isWalkable(int(spawnPos.X), int(spawnPos.Y)) && !p.isPowerup(int(spawnPos.X), int(spawnPos.Y)) {
			up := newPowerup(typ)
			up.setPos(int(spawnPos.X), int(spawnPos.Y))
			p.powerups[p.currentLevel] = append(p.powerups[p.currentLevel], up)
			spawnCount--
		}
	}
}

func (p *MapEngine) setContextTiles() {
	for i := 0; i < MAP_H; i++ {
		for j := 0; j < MAP_W; j++ {
			p.contextTiles[i][j] = -1
			dist := distance(p.player.pos, sdl.Point{int32(j), int32(i)})

			if int(dist) <= 1 && p.isEnemy(j, i) {
				p.contextTiles[i][j] = CONTEXT_ENEMY
			} else if dist <= 1 && p.isWalkable(j, i) {
				p.contextTiles[i][j] = CONTEXT_WALKABLE
			}
		}
	}
}

func (p *MapEngine) isWalkable(x, y int) bool {
	if x < 1 || x > MAP_W-1 {
		return false
	}
	if y < 1 || y > MAP_H-1 {
		return false
	}

	tile := p.levels[p.currentLevel][y][x]

	if tile != TILE_FLOOR && tile != TILE_STAIRS_UP && tile != TILE_STAIRS_DOWN && tile != TILE_STAIRS_BLOCKED {
		return false
	}

	if p.isEnemy(x, y) {
		return false
	}

	if p.player != nil && p.player.pos == (sdl.Point{int32(x), int32(y)}) {
		return false
	}

	return true
}

func (p *MapEngine) isEnemy(x, y int) bool {
	if x < 1 || x > MAP_W-1 {
		return false
	}
	if y < 1 || y > MAP_H-1 {
		return false
	}

	e := p.getEnemy(x, y)
	if e != nil {
		return true
	}

	return false
}

func (p *MapEngine) isPowerup(x, y int) bool {
	if x < 1 || x > MAP_W-1 {
		return false
	}
	if y < 1 || y > MAP_H-1 {
		return false
	}

	for _, up := range p.powerups[p.currentLevel] {
		if up.pos == (sdl.Point{int32(x), int32(y)}) {
			return true
		}
	}
	return false
}

func (p *MapEngine) moveCursor(direction int) {
	switch direction {
	case CUR_UP:
		p.cursorPos.Y--
	case CUR_DOWN:
		p.cursorPos.Y++
	case CUR_LEFT:
		p.cursorPos.X--
	case CUR_RIGHT:
		p.cursorPos.X++
	}

	if p.cursorPos.X < 0 {
		p.cursorPos.X = 0
	} else if p.cursorPos.X >= MAP_W {
		p.cursorPos.X = MAP_W - 1
	}

	if p.cursorPos.Y < 0 {
		p.cursorPos.Y = 0
	} else if p.cursorPos.Y > MAP_H-1 {
		p.cursorPos.Y = MAP_H - 1
	}

	p.sfxCursor.play()
}

func (p *MapEngine) getEnemy(x, y int) *Enemy {
	for _, e := range p.enemies[p.currentLevel] {
		if e.pos == (sdl.Point{int32(x), int32(y)}) {
			return e
		}
	}
	return nil
}

func (p *MapEngine) render() {
	for i := 0; i < MAP_H; i++ {
		for j := 0; j < MAP_W; j++ {
			tile := p.levels[p.currentLevel][i][j]
			p.tileset.setClip(tile*TILE_SIZE, 0, TILE_SIZE, TILE_SIZE)
			p.tileset.setPos(j*TILE_SIZE, i*TILE_SIZE)
			p.tileset.render()
		}
	}

	// context tiles
	if p.player != nil && p.player.isTurn {
		for i := 0; i < MAP_H; i++ {
			for j := 0; j < MAP_W; j++ {
				tile := p.contextTiles[i][j]
				p.tileset.setClip(tile*TILE_SIZE, 16, TILE_SIZE, TILE_SIZE)
				p.tileset.setPos(j*TILE_SIZE, i*TILE_SIZE)
				p.tileset.render()
			}
		}
	}

	// render cursor
	p.cursor.setPos(int(p.cursorPos.X*TILE_SIZE), int(p.cursorPos.Y*TILE_SIZE))
	p.cursor.render()

	if p.player != nil {
		p.player.render()
	}

	// render powerups
	for _, up := range p.powerups[p.currentLevel] {
		up.logic()
		up.render()
	}

	// render enemies
	for _, e := range p.enemies[p.currentLevel] {
		e.render()
	}

	// render hud/status text
	p.hud.render()
	p.msg.render()
}

func (p *MapEngine) enemyStartTurn() {
	if len(p.enemies[p.currentLevel]) > 0 {
		p.player.isTurn = false
		p.msg.setText("Enemy turn")
		p.activeEnemy = 0
		p.enemies[p.currentLevel][p.activeEnemy].startTurn()
	}
}

func (p *MapEngine) enemyAction() bool {
	if p.activeEnemy == -1 {
		return true
	}

	for i := 0; i < len(p.enemies[p.currentLevel]); i++ {
		e := p.enemies[p.currentLevel][i]
		if e.hp == 0 {
			if !e.isAnimating() {
				p.removeEnemy(e)
				if len(p.enemies[p.currentLevel]) == 0 {
					p.activeEnemy = -1
				}
			}
			return false
		}
	}

	e := p.enemies[p.currentLevel][p.activeEnemy]
	if e != nil && e.isTurn {
		p.cursorPos = e.pos
	}

	if e != nil && e.isTurn && !e.isAnimating() {
		if e.actionTicks > 0 {
			e.actionTicks--
		} else {
			// move torwards player
			if !e.isNearPlayer(1) {
				dest := e.pos
				if e.isNearPlayer(4) {
					if e.pos.X > p.player.pos.X {
						dest.X--
					} else if e.pos.X < p.player.pos.X {
						dest.X++
					}

					if e.pos.Y > p.player.pos.Y {
						dest.Y--
					} else if e.pos.Y < p.player.pos.Y {
						dest.Y++
					}
				} else {
					dx, dy, failCount := 0, 0, 3
					for dx == 0 && dy == 0 && failCount > 0 {
						dx = rand.Intn(3) - 1
						dy = rand.Intn(3) - 1
						failCount--
					}
					dest.X += int32(dx)
					dest.Y += int32(dy)
				}
				if p.isWalkable(int(dest.X), int(dest.Y)) {
					e.setPos(int(dest.X), int(dest.Y))
					e.sfxMove.play()
				}
			} else {
				// attack
				e.actionAttack()
			}

			e.isTurn = false
		}

		if p.player.hp == 0 {
			p.gameState = GAME_LOSE
			p.msg.setText("You lose. R to retry.")
			return false
		}
	}

	if e != nil && !e.isTurn {
		if p.activeEnemy+1 < len(p.enemies[p.currentLevel]) {
			p.activeEnemy++
			p.enemies[p.currentLevel][p.activeEnemy].startTurn()
			return false
		} else {
			p.activeEnemy = -1
			return true
		}
	}

	return false
}

func (p *MapEngine) checkPowerup() {
	for i, up := range p.powerups[p.currentLevel] {
		if up.pos == p.player.pos {
			switch up.typ {
			case POWERUP_ATTACK:
				p.player.bonusAttack(up.amount)
			case POWERUP_DEFENSE:
				p.player.bonusDefense(up.amount)
			case POWERUP_POTION:
				p.player.bonusHP(up.amount)
			case POWERUP_TREASURE:
				p.player.bonusTreasure()
			}

			if up.typ == POWERUP_TREASURE {
				p.sfxTreasure.play()
			} else {
				p.sfxPowerup.play()
			}

			p.powerups[p.currentLevel] = append(p.powerups[p.currentLevel][:i], p.powerups[p.currentLevel][i+1:]...)
			return
		}
	}
}

func (p *MapEngine) removeEnemy(e *Enemy) {
	if e == nil {
		return
	}

	for i, enemy := range p.enemies[p.currentLevel] {
		if e == enemy {
			p.enemies[p.currentLevel] = append(p.enemies[p.currentLevel][:i], p.enemies[p.currentLevel][i+1:]...)
			break
		}
	}

	p.setContextTiles()
}

func (p *MapEngine) getCurrentLevel() int {
	return p.currentLevel + 1
}
