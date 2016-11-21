package main

import "github.com/qeedquan/go-media/sdl/sdlmixer"

type Sound struct {
	soundEngine *SoundEngine
	filename    string
	chunk       *sdlmixer.Chunk
	music       *sdlmixer.Music
}

func (s *Sound) init(soundEngine *SoundEngine) {
	s.soundEngine = soundEngine
}

func (s *Sound) play() {
	if s.chunk != nil {
		s.chunk.PlayChannel(-1, 0)
	} else if s.soundEngine != nil && s.music != nil {
		s.soundEngine.playMusic(s)
	}
}
