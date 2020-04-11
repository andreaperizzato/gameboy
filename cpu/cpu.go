package cpu

import (
	"fmt"
	"strings"

	"github.com/andreaperizzato/gameboy/memory"
)

type registers struct {
	A  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	H  uint8
	L  uint8
	SP uint16
	PC uint16
}

type flags struct {
	Z bool
	N bool
	H bool
	C bool
}

type runnable func(c *CPU) uint8

type instruction struct {
	name string
	// run executes the instruction and returns the cycles to wait.
	run runnable
}

// CPU emulates a CPU.
type CPU struct {
	mem   memory.AddressSpace
	regs  registers
	flags flags
	instr map[uint16]instruction
	wait  uint8
	b     strings.Builder
}

// NewGBC creats a new CPU with the GBC instruction set.
func NewGBC(mmu memory.AddressSpace) *CPU {
	return &CPU{
		mem:   mmu,
		instr: gbcInstructions,
	}
}

// Tick executes one CPU step.
func (c *CPU) Tick() {
	if c.wait > 0 {
		c.wait--
		return
	}

	initialPC := c.regs.PC
	opcode := uint16(nextArg(c))
	if opcode == 0xCB {
		opcode = 0xCB00 | uint16(nextArg(c))
	}
	cmd, found := c.instr[opcode]
	if !found {
		err := fmt.Sprintf("opcode 0x%04X not found at 0x%04X", opcode, initialPC)
		panic(err)
	}
	c.wait = cmd.run(c)
	c.wait--
}
