package memory_test

import (
	"testing"

	"github.com/andreaperizzato/gameboy/memory"
	"github.com/stretchr/testify/assert"
)

func TestRAM_Contains(t *testing.T) {
	tests := []struct {
		name string
		addr uint16
		exp  bool
	}{
		{"before range", 0x0F, false},
		{"first address", 0x10, true},
		{"within range", 0x11, true},
		{"last address", 0x12, true},
		{"after address", 0x13, false},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			ram := memory.NewRAM(3, 0x10)
			assert.Equal(t, tC.exp, ram.Contains(tC.addr))
		})
	}
}

func TestRAM_ReadWrite(t *testing.T) {
	ram := memory.NewRAM(3, 0x10)
	// Valid address.
	ram.Write(0x11, 0xAA)
	assert.Equal(t, uint8(0xAA), ram.Read(0x11))

	// Invalid address.
	assert.Panics(t, func() { ram.Write(0x33, 0xFF) })
	assert.Panics(t, func() { ram.Read(0x33) })
}
