package ppu

import (
	"github.com/andreaperizzato/gameboy/memory"
)

// ppuState is a state the PPU can be in.
type ppuState uint8

// All Possible PPU States.
const (
	oamSearch ppuState = iota
	pixelTransfer
	hBlank
	vBlank
)

const (
	lcdcDisplayEnabledBit = 7
)

// PPU is the Gameboy Picture Processing Unit.
type PPU struct {
	// state is the current state of the PPU.
	state ppuState
	// ticks is the clock ticks counter for the current line.
	ticks uint
	// x is the number of pixels already output in the current scanline.
	x uint8

	Screen  Display
	Fetcher *Fetcher
	mem     memory.AddressSpace

	// LY Y-Coordinate
	// https://gbdev.io/pandocs/#ff44-ly-lcdc-y-coordinate-r
	ly memory.Register

	// SCY - Scroll Y
	// https://gbdev.io/pandocs/#ff42-scy-scroll-y-r-w-ff43-scx-scroll-x-r-w
	scy memory.Register

	// LDCD - LCD Control Register
	lcdc memory.Register
}

// New creates anew PPU.
func New(m memory.AddressSpace, screen Display) *PPU {
	return &PPU{
		Fetcher: NewFetcher(m),
		Screen:  screen,
		state:   oamSearch,
		mem:     m,
		ly:      memory.NewRegister(m, 0xFF44),
		scy:     memory.NewRegister(m, 0xFF42),
		lcdc:    memory.NewRegister(m, 0xFF40),
	}
}

// Tick advances the PPU state by one step.
func (p *PPU) Tick() {
	if !p.lcdc.GetBit(lcdcDisplayEnabledBit) {
		return
	}

	p.ticks++
	switch p.state {
	case oamSearch:
		// collect sprite data
		// Here we need to scan the OAM (obj attribute memory)
		// from 0xFE00 to 0xFE9F to mix sprites with the current line.
		// This always takes 40 ticks.
		if p.ticks == 40 {
			p.x = 0
			y := p.scy.Get() + p.ly.Get()
			tileLine := y % 8
			tileMapRowAddr := 0x9800 + uint16(y/8)*32
			p.Fetcher.Start(tileMapRowAddr, tileLine)
			p.state = pixelTransfer
		}

	case pixelTransfer:
		// Fetch pixel data into the FIFO queue.
		p.Fetcher.Tick()
		if p.Fetcher.Q.Size() < 8 {
			return
		}
		// Put a pixel from the FIFO on the screen if we have any.
		if pxColor, ok := p.Fetcher.Q.Pop(); ok {
			p.Screen.Write(pxColor)
			p.x++
		}
		if p.x == 160 {
			p.Screen.HBlank()
			p.state = hBlank
		}

	case hBlank:
		// A full scanline takes 456 ticks to complete. At the end
		// of a scanline, the PPU goes black to the initial OAM Search state.
		// When we reach line 144, we switch to VBlank state.
		if p.ticks == 456 {
			p.ticks = 0
			p.ly.Set(p.ly.Get() + 1)
			if p.ly.Get() == 144 {
				p.Screen.VBlank()
				p.state = vBlank
			} else {
				p.state = oamSearch
			}
		}

	case vBlank:
		// According to https://gbdev.io/pandocs/#lcdc-7-lcd-display-enable
		// switching the display on and off can only be done when in VBlank.
		wasOn := p.Screen.IsEnabled()
		isOn := p.lcdc.GetBit(lcdcDisplayEnabledBit)
		if wasOn && !isOn {
			// Turn off.
			p.Screen.Enable(false)
			p.Fetcher.Q.Clear()
			p.x = 0
			p.ly.Set(0)
			return
		}
		if !wasOn && isOn {
			// Turn on.
			p.Fetcher.Q.Clear()
			p.x = 0
			p.ly.Set(0)
			p.Screen.Enable(true)
			p.state = oamSearch
			return
		}

		if p.ticks == 456 {
			p.ticks = 0
			p.ly.Set(p.ly.Get() + 1)
			if p.ly.Get() == 144 {
				p.ly.Set(0)
				p.state = oamSearch
			}
		}
	}
}
