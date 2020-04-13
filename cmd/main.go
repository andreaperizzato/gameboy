package main

import (
	"github.com/andreaperizzato/gameboy/cpu"
	"github.com/andreaperizzato/gameboy/memory"
	"github.com/andreaperizzato/gameboy/ppu"
	"github.com/andreaperizzato/gameboy/screen"
)

// https://gbdev.gg8.se/wiki/articles/Gameboy_Bootstrap_ROM
// https://blog.ryanlevick.com/DMG-01/public/book/cpu/reading_and_writing_memory.html?highlight=call#the-stack
// https://rednex.github.io/rgbds/gbz80.7.html#CALL_n16
// https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html
// https://blog.tigris.fr/2019/09/15/writing-an-emulator-the-first-pixel/
// https://gbdev.io/pandocs

var logo = []uint8{
	0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
	0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
	0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
	0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
	0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
	0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
}

func main() {
	ram := memory.NewRAM(0xFFFF, 0)
	// Just write the logo so it shows up.
	for i, b := range logo {
		ram.Write(0x0104+uint16(i), b)
	}

	mmu := memory.NewMMU(memory.NewGBCBootROM(), ram)
	cpux := cpu.NewGBC(mmu)
	scrx := screen.New()
	ppux := ppu.New(mmu, scrx)

	go func() {
		for {
			cpux.Tick()
			ppux.Tick()
			ppux.Tick()
			ppux.Tick()
			ppux.Tick()
		}
	}()
	scrx.Start()
}
