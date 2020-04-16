package memory_test

import (
	"testing"

	"github.com/andreaperizzato/gameboy/memory"
	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {
	ram := memory.NewRAM(3, 0)
	r := memory.NewRegister(ram, 0x01)

	r.Set(0xAA)
	assert.Equal(t, uint8(0xAA), ram.Read(0x01))
	assert.Equal(t, uint8(0xAA), r.Get())

	r.Set(0b00100000)
	assert.True(t, r.GetBit(5))
}

func Test_RegisterWithMask(t *testing.T) {
	ram := memory.NewRAM(3, 0)
	r := memory.NewRegisterWithMask(ram, 0x01, 0b11100000)
	ram.Write(0x01, 0b00000111)

	r.Set(0b101)
	assert.Equal(t, uint8(0b101), r.Get())
	assert.Equal(t, uint8(0b10100111), ram.Read(0x01))
}

func Test_Register16(t *testing.T) {
	ram := memory.NewRAM(3, 0)
	r := memory.NewRegister16(ram, 0x01, 0x02)

	r.Set(0xAABB)
	assert.Equal(t, uint8(0xBB), ram.Read(0x01))
	assert.Equal(t, uint8(0xAA), ram.Read(0x02))
	assert.Equal(t, uint16(0xAABB), r.Get())
}

func Test_Register16WithMask(t *testing.T) {
	ram := memory.NewRAM(3, 0)
	// using 0b00000111 as mask, effectively makes it a 11-bit register (8+3).
	r := memory.NewRegister16WithMask(ram, 0x01, 0x02, 0b00000111)
	ram.Write(0x02, 0b11000000)

	r.Set(0b10100111111) // 0x53F
	assert.Equal(t, uint8(0x3F), ram.Read(0x01))
	assert.Equal(t, uint8(0xC5), ram.Read(0x02)) // we have C5 because we don't change the other bits
	assert.Equal(t, uint16(0x053F), r.Get())
}
