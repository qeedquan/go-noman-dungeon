package main

import (
	"log"

	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

type soundCache struct {
	filename string
	chunk    *sdlmixer.Chunk
	music    *sdlmixer.Music
	ref      int
}

type SoundEngine struct {
	cache         []soundCache
	musicFilename string
}

func newSoundEngine() *SoundEngine {
	err := sdlmixer.OpenAudio(44100, sdlmixer.DEFAULT_FORMAT, 2, 1024)
	if err != nil {
		log.Print("sdl: ", err)
	}

	sdlmixer.AllocateChannels(128)
	return &SoundEngine{}
}

func (s *SoundEngine) loadSound(sound *Sound, filename string) {
	if sound == nil {
		return
	}

	chunk := s.cacheLookup(filename)
	if chunk != nil {
		sound.filename = filename
		sound.chunk = chunk
		sound.init(s)
		return
	}

	chunk, err := sdlmixer.LoadWAV(filename)
	if err == nil {
		sound.filename = filename
		sound.chunk = chunk
		s.cacheStore(sound, false)
	}
	sound.init(s)

	if err != nil {
		log.Println("sound: ", err)
	}
}

func (s *SoundEngine) loadMusic(sound *Sound, filename string) {
	if sound == nil {
		return
	}

	music := s.cacheLookupMusic(filename)
	if music != nil {
		sound.filename = filename
		sound.music = music
		sound.init(s)
		return
	}

	music, err := sdlmixer.LoadMUS(filename)
	if err == nil {
		sound.filename = filename
		sound.music = music
		s.cacheStore(sound, true)
	}
	sound.init(s)

	if err != nil {
		log.Println("sound: ", err)
	}

}

func (s *SoundEngine) cacheStore(sound *Sound, isMusic bool) {
	if sound.filename != "" {
		var snd soundCache

		snd.filename = sound.filename
		snd.ref = 1

		if !isMusic && sound.chunk != nil {
			snd.chunk = sound.chunk
			snd.music = nil
			s.cache = append(s.cache, snd)
		} else if isMusic && sound.music != nil {
			snd.music = sound.music
			snd.chunk = nil
			s.cache = append(s.cache, snd)
		}
	}
}

func (s *SoundEngine) cacheLookup(filename string) *sdlmixer.Chunk {
	for i := range s.cache {
		c := &s.cache[i]
		if c.music == nil && filename == c.filename {
			c.ref++
			return c.chunk
		}
	}

	return nil
}

func (s *SoundEngine) cacheLookupMusic(filename string) *sdlmixer.Music {
	for i := range s.cache {
		c := &s.cache[i]
		if c.music != nil && filename == c.filename {
			c.ref++
			return c.music
		}
	}

	return nil
}

func (s *SoundEngine) playMusic(sound *Sound) {
	if sound != nil && sound.music != nil && sound.filename != s.musicFilename {
		sdlmixer.HaltMusic()
		sound.music.Play(-1)
	}
}
