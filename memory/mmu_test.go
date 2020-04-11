package memory_test

import (
	"testing"

	"github.com/andreaperizzato/gameboy/memory"
	"github.com/stretchr/testify/assert"
)

func TestMMU_Contains(t *testing.T) {
	tests := []struct {
		name string
		addr uint16
		exp  bool
	}{
		{"within rom", 0x01, true},
		{"in ram 1", 0x04, true},
		{"between ram 1 and 2", 0x30, false},
		{"in ram 2", 0x51, true},
		{"outside", 0x11, false},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			boot := memory.NewROM([]uint8{0xAA, 0xBB, 0xCC}, 0)
			ram1 := memory.NewRAM(5, 0x03)
			ram2 := memory.NewRAM(5, 0x50)
			m := memory.NewMMU(boot, ram1, ram2)
			assert.Equal(t, tC.exp, m.Contains(tC.addr))
		})
	}
}

func TestMMU_ReadWrite(t *testing.T) {
	boot := memory.NewROM([]uint8{}, 0)
	ram := memory.NewRAM(3, 0)

	m := memory.NewMMU(boot, ram)

	// Read at valid address.
	m.Write(0x01, 0xAA)
	assert.Equal(t, uint8(0xAA), m.Read(0x01))
	// Read at invalid address returns 0xFF.
	assert.Equal(t, uint8(0xFF), m.Read(0xABCD))

	// Write at invalid address has no effect.
	m.Write(0xABCD, 0xAA)
}

func TestMMU_DisablingBootROM(t *testing.T) {
	boot := memory.NewROM([]uint8{0xAA, 0xBB, 0xCC}, 0)
	ram1 := memory.NewRAM(3, 0)
	ram1.Write(0x00, 0x11)
	ram1.Write(0x01, 0x22)
	ram1.Write(0x02, 0x33)

	ram2 := memory.NewRAM(5, 0x03)
	ram2.Write(0x03, 0x0D)

	m := memory.NewMMU(boot, ram1, ram2)
	// With boot enabled
	expBytes := []uint8{0xAA, 0xBB, 0xCC, 0x0D}
	for i := uint16(0); i < 4; i++ {
		assert.Equal(t, expBytes[i], m.Read(i))
	}
	// With boot disabled
	m.Write(0xFF50, 0x01)
	expBytes = []uint8{0x11, 0x22, 0x33, 0x0D}
	for i := uint16(0); i < 4; i++ {
		assert.Equal(t, expBytes[i], m.Read(i))
	}
}
