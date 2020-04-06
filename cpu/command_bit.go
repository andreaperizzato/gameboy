package cpu

import (
	"fmt"

	"github.com/andreaperizzato/gameboy/memory"
)

var bit = map[uint16]Command{
	// BIT 7,H
	0xCB7C: bit8(0xCB7C, "H", 7, accessH()),
	// RL C
	0xCB11: rl8(0xCB11, "C", accessC()),
	// RL A
	0xCB17: rl8(0xCB17, "A", accessA()),
	// RLA
	0x17: rl8(0x17, "A", accessA()),
}

func bit8(opcode uint16, regName string, pos uint8, src accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     4,
		ArgsLength: 0,
		Label:      fmt.Sprintf("BIT %d,%s", pos, regName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			v := src(r).Get()
			set := (v>>pos)&0x01 == 0x01
			r.SetFlagZ(set)
			r.SetFlagN(false)
			r.SetFlagH(true)
		},
	}
}

func rl8(opcode uint16, srcRegName string, src accessor8Creator) Command {
	cycles := uint8(4)
	divider := ""
	if opcode>>8 == 0xCB { // extended set
		cycles = 8
		divider = " "
	}
	return Command{
		OpCode:     opcode,
		Cycles:     cycles,
		ArgsLength: 0,
		Label:      fmt.Sprintf("RL%s%s", divider, srcRegName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			v := src(r).Get()
			res := v << 1
			if r.FlagC() {
				res |= 0x01
			}
			r.SetFlagZ(res == 0)
			r.SetFlagN(false)
			r.SetFlagH(false)
			r.SetFlagC(v&(1<<7) != 0)
		},
	}
}
