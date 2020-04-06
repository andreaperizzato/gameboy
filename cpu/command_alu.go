package cpu

import (
	"fmt"

	"github.com/andreaperizzato/gameboy/memory"
)

var alu = map[uint16]Command{
	// XOR A
	0xAF: xor8(0xAF, "A", accessA()),
	// INC B
	0x04: inc8(0x04, "B", accessB()),
	// INC C
	0x0C: inc8(0x0C, "C", accessC()),
	// INC H
	0x24: inc8(0x24, "H", accessH()),
	// INC DE
	0x13: inc16(0x13, "DE", accessDE()),
	// INC HL
	0x23: inc16(0x23, "HL", accessHL()),
	// CP n
	0xFE: cpConst(0xFE),
	// DEC A
	0x3D: dec8(0x3D, "A", accessA()),
	// DEC B
	0x05: dec8(0x05, "B", accessB()),
	// DEC C
	0x0D: dec8(0x0D, "C", accessC()),
	// DEC D
	0x15: dec8(0x15, "D", accessD()),
	// DEC E
	0x1D: dec8(0x1D, "E", accessE()),
	// SUB B
	0x90: sub8(0x90, "B", accessB()),
}

type regGetter8 func(r *Registers) uint8

func xor8(opcode uint16, regName string, a8c accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     4,
		ArgsLength: 0,
		Label:      fmt.Sprintf("XOR %s", regName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			regAcc := a8c(r)
			r.A ^= regAcc.Get()
			r.SetFlagN(false)
			r.SetFlagH(false)
			r.SetFlagC(false)
			r.SetFlagZ(r.A == 0x00)
		},
	}
}

func inc8(opcode uint16, regName string, a8c accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     4,
		ArgsLength: 0,
		Label:      fmt.Sprintf("INC %s", regName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			regAcc := a8c(r)
			res := regAcc.Get() + 1
			regAcc.Set(res)
			r.SetFlagN(false)
			r.SetFlagZ(res == 0x00)
			r.SetFlagH(res&0x0F == 0x0F)
		},
	}
}

func inc16(opcode uint16, regName string, dst accessor16Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     8,
		ArgsLength: 0,
		Label:      fmt.Sprintf("INC %s", regName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			regAcc := dst(r)
			regAcc.Set(regAcc.Get() + 1)
		},
	}
}

func cpConst(opcode uint16) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     8,
		ArgsLength: 1,
		Label:      "CP n",
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			n := args[0]
			r.SetFlagZ(r.A == n)
			r.SetFlagN(true)
			r.SetFlagH((0x0F & r.A) < (0x0F & n))
			r.SetFlagC(r.A < n)
		},
	}
}

func dec8(opcode uint16, regName string, reg accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     4,
		ArgsLength: 0,
		Label:      fmt.Sprintf("INC %s", regName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			regAcc := reg(r)
			val := regAcc.Get()
			res := val - 1
			regAcc.Set(res)
			r.SetFlagZ(res == 0x00)
			r.SetFlagN(true)
			r.SetFlagH(val&0x0F == 0x00)
		},
	}
}

func sub8(opcode uint16, regName string, src accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     4,
		ArgsLength: 0,
		Label:      fmt.Sprintf("SUB %s", regName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			v := src(r).Get()
			a := r.A
			r.A = a - v
			r.SetFlagZ(r.A == 0x00)
			r.SetFlagN(true)
			r.SetFlagC(v > a)
			r.SetFlagH((0x0F & v) > (0x0F & a))
		},
	}
}
