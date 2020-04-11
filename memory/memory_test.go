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
