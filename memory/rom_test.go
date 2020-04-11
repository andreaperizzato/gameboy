package memory_test

import (
	"testing"

	"github.com/andreaperizzato/gameboy/memory"
	"github.com/stretchr/testify/assert"
)

func TestROM_ReadOnly(t *testing.T) {
	rom := memory.NewGBCBootROM()
	v := rom.Read(0x07)
	rom.Write(0x07, v+1)
	assert.Equal(t, v, rom.Read(0x07))
}
