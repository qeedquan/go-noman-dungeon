package main

import "github.com/qeedquan/go-media/sdl"

type Animation struct {
	image       Image
	frames      []int
	filename    string
	index       int
	frameIndex  int
	loops       int
	timesPlayed int
	frameSize   sdl.Point
	frameOffset sdl.Point
}

func (a *Animation) setTo(other *Animation) {
	if a.filename != other.filename {
		other.image.ref()
	}
	*a = *other
}

func (a *Animation) load(filename string, index, frameCount, duration, loops, frameW, frameH, offsetX, offsetY int) {
	a.filename = filename
	a.image.load(renderEngine, filename)
	a.index = index
	a.loops = loops
	a.frameSize = sdl.Point{int32(frameW), int32(frameH)}
	a.frameOffset = sdl.Point{int32(offsetX), int32(offsetY)}

	a.frames = a.frames[:0]
	duration = int(float64(duration)*MAX_FRAMES_PER_SEC/1000 + 0.5)
	if frameCount > 0 && duration%frameCount == 0 {
		divided := duration / frameCount
		for i := 0; i < frameCount; i++ {
			for j := 0; j < divided; j++ {
				a.frames = append(a.frames, i)
			}
		}
	} else {
		x0, y0 := 0, 0
		x1, y1 := duration-1, frameCount-1

		dx := x1 - x0
		dy := y1 - y0
		D := 2*dy - dx

		a.frames = append(a.frames, y0)

		x, y := x0+1, y0

		for ; x <= x1; x++ {
			if D > 0 {
				y++
				D += 2*dy - 2*dx
			} else {
				D += 2 * dy
			}
			a.frames = append(a.frames, y)
		}
	}

	if len(a.frames) > 0 {
		a.image.setClip(a.frames[a.frameIndex]*int(a.frameSize.X), a.index*int(a.frameSize.Y), int(a.frameSize.X), int(a.frameSize.Y))
	}
}

func (a *Animation) logic() {
	if len(a.frames) > 0 {
		a.image.setClip(a.frames[a.frameIndex]*int(a.frameSize.X), a.index*int(a.frameSize.Y), int(a.frameSize.X), int(a.frameSize.Y))
	}

	if a.loops == 0 || a.timesPlayed < a.loops {
		if a.frameIndex < len(a.frames)-1 {
			a.frameIndex++
		} else {
			if a.loops == 0 || a.timesPlayed+1 < a.loops {
				a.frameIndex = 0
			}
			a.timesPlayed++
		}
	}
}

func (a *Animation) render() {
	a.image.render()
}

func (a *Animation) setPos(x, y int) {
	a.image.setPos(x-int(a.frameOffset.X), y-int(a.frameOffset.Y))
}

func (a *Animation) isLastFrame() bool {
	return a.frameIndex == len(a.frames)-1
}

func (a *Animation) isFinished() bool {
	return a.loops != 0 && a.isLastFrame() && a.timesPlayed == a.loops
}
