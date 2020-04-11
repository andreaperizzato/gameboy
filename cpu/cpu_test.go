package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPU_Tick(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := NewGBC(mem)

	// Test main struction set.
	c.regs.B = 0x00
	c.regs.PC = 0x0000
	mem[0x0000] = 0x04 // INC B
	c.Tick()
	assert.Equal(t, uint8(0x01), c.regs.B, "B after INC B")
	assert.Equal(t, uint16(0x0001), c.regs.PC, "PC after first tick")
	// INC takes 4 cycles
	tick(c, 3)

	// Test extended instruction set.
	mem[0x0001], mem[0x0002] = 0xCB, 0x17 // RL A
	c.regs.A = 0b00000001
	c.Tick()
	assert.Equal(t, uint8(0b00000010), c.regs.A, "A after RL A")
	assert.Equal(t, uint16(0x0003), c.regs.PC, "PC after second tick")
	tick(c, 7)

	// Test unknown instruction.
	mem[0x0003] = 0xFD // this opcode doesn't exist.
	assert.Panics(t, func() { c.Tick() })
}

func tick(c *CPU, times int) {
	for i := 0; i < times; i++ {
		c.Tick()
	}
}
