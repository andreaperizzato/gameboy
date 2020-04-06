package cpu

import (
	"fmt"

	"github.com/andreaperizzato/gameboy/memory"
)

var jump = map[uint16]Command{
	// PUSH BC
	0xC5: push16(0xC5, "BC", accessBC()),
	// POP BC
	0xC1: pop16(0xC1, "BC", accessBC()),
	// JR NZ,nn
	0x20: jrConditional(0x20, "NZ", accessFlagZ(), false),
	// JR Z,nn
	0x28: jrConditional(0x28, "Z", accessFlagZ(), true),
	// JR NC,nn
	0x30: jrConditional(0x30, "NC", accessFlagC(), false),
	// JR C,nn
	0x38: jrConditional(0x38, "C", accessFlagC(), true),
	// JR nn
	0x18: Command{
		OpCode:     0x18,
		Cycles:     8,
		ArgsLength: 1,
		Label:      "JR n",
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			delta := int8(args[0]) // args[0] is a signed byte.
			if delta < 0 {
				r.PC -= uint16(-delta)
			} else {
				r.PC += uint16(delta)
			}
		},
	},
	// CALL nn
	0xCD: Command{
		OpCode:     0xCD,
		Cycles:     8,
		ArgsLength: 2,
		Label:      "CALL nn",
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			r.SP--
			m.SetByte(r.SP, uint8(r.PC))
			r.SP--
			m.SetByte(r.SP, uint8(r.PC>>8))
			r.PC = combine(args[1], args[0])
		},
	},
	// RET
	0xC9: Command{
		OpCode:     0xC9,
		Cycles:     16,
		ArgsLength: 0,
		Label:      "RET",
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			msb := m.GetByte(r.SP)
			r.SP++
			lsb := m.GetByte(r.SP)
			r.SP++
			v := combine(msb, lsb)
			r.PC = v
		},
	},
}

func jrConditional(opcode uint16, conditionName string, flag flagAccessorCreator, jumpWhen bool) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     8,
		ArgsLength: 1,
		Label:      fmt.Sprintf("JR %s,n", conditionName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			f := flag(r)
			if f.Get() == jumpWhen {
				delta := int8(args[0]) // args[0] is a signed byte.
				if delta < 0 {
					r.PC -= uint16(-delta)
				} else {
					r.PC += uint16(delta)
				}
			}
		},
	}
}

func push16(opcode uint16, srcRegName string, src accessor16Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     16,
		ArgsLength: 0,
		Label:      fmt.Sprintf("PUSH %s", srcRegName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			v := src(r).Get()
			r.SP--
			m.SetByte(r.SP, uint8(v))
			r.SP--
			m.SetByte(r.SP, uint8(v>>8))
		},
	}
}

func pop16(opcode uint16, dstRegName string, dst accessor16Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     12,
		ArgsLength: 0,
		Label:      fmt.Sprintf("POP %s", dstRegName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			msb := m.GetByte(r.SP)
			r.SP++
			lsb := m.GetByte(r.SP)
			r.SP++
			v := combine(msb, lsb)
			dst(r).Set(v)
		},
	}
}
