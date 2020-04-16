package apu

import (
	"time"

	"github.com/andreaperizzato/gameboy/memory"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

// APU implements the gameboy audio processing unit.
type APU struct {
	pulse1 *Pulse
}

// NewAPU creates a new APU.
func NewAPU(mem memory.AddressSpace) *APU {
	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/30))
	if err != nil {
		panic(err)
	}
	pulse1 := &Pulse{
		SampleRate: sr,
		Duty:       memory.NewRegisterWithMask(mem, 0xFF11, 0xC0),
		Frequency:  memory.NewRegister16WithMask(mem, 0xFF13, 0xFF14, 0x07),
		Envelope:   memory.NewRegisterWithMask(mem, 0xFF12, 0x07),
		Volume:     memory.NewRegisterWithMask(mem, 0xFF12, 0xF0),
		Restart:    memory.NewRegisterBit(mem, 0xFF14, 7),
	}
	return &APU{
		pulse1: pulse1,
	}
}

// Start starts playing sounds.
func (apu *APU) Start() {
	speaker.Play(apu.pulse1)
}
