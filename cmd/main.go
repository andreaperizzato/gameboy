package main

import (
	"io/ioutil"

	"github.com/andreaperizzato/gameboy/cpu"
)

// https://gbdev.gg8.se/wiki/articles/Gameboy_Bootstrap_ROM
// https://blog.ryanlevick.com/DMG-01/public/book/cpu/reading_and_writing_memory.html?highlight=call#the-stack
// https://rednex.github.io/rgbds/gbz80.7.html#CALL_n16
// https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html

func main() {
	startup, err := ioutil.ReadFile("../DMG_ROM.bin")
	if err != nil {
		panic(err)
	}

	cpu.Run(startup)
}
