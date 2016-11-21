package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/qeedquan/go-media/sdl"
)

const (
	MAX_FRAMES_PER_SEC = 60
	TILE_SIZE          = 16
)

var (
	dataDir      string
	dungeonDepth int
	godMode      bool
)

var (
	renderEngine *RenderEngine
	inputEngine  *InputEngine
	soundEngine  *SoundEngine
	gameEngine   *GameEngine
)

func main() {
	runtime.LockOSThread()
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())
	parseFlags()
	err := sdl.Init(sdl.INIT_EVERYTHING &^ sdl.INIT_AUDIO)
	if err != nil {
		log.Fatal("sdl: ", err)
	}

	renderEngine = newRenderEngine()
	inputEngine = newInputEngine()
	soundEngine = newSoundEngine()
	gameEngine = newGameEngine()
	mainLoop()
}

func parseFlags() {
	flag.StringVar(&dataDir, "data", "data", "data directory")
	flag.IntVar(&dungeonDepth, "depth", 5, "dungeon depth (2-9999)")
	flag.BoolVar(&godMode, "god", false, "god mode")
	flag.Parse()
	if dungeonDepth <= 1 || dungeonDepth > 9999 {
		log.Fatal("flag: invalid dungeon depth!")
	}
}

func mainLoop() {
	done := false
	delay := uint32(math.Floor(1000.0/MAX_FRAMES_PER_SEC + 0.5))
	logicTicks := sdl.GetTicks()

	for !done {
		loops := 0
		nowTicks := sdl.GetTicks()
		prevTicks := nowTicks

		for nowTicks > logicTicks && loops < MAX_FRAMES_PER_SEC {
			inputEngine.logic()
			gameEngine.logic()

			done = gameEngine.done || inputEngine.done
			logicTicks += delay
			loops++
		}

		renderEngine.clear()
		gameEngine.render()
		renderEngine.commitFrame()

		nowTicks = sdl.GetTicks()
		if nowTicks-prevTicks < delay {
			sdl.Delay(delay - (nowTicks - prevTicks))
		}
	}
}
