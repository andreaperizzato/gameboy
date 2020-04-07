package cpu

import (
	"fmt"

	"github.com/andreaperizzato/gameboy/memory"
)

var ld = map[uint16]Command{
	// LD SP,nn
	0x31: ldInto16Const(0x31, "SP", accessSP()),
	// LD HL,nn
	0x21: ldInto16Const(0x21, "HL", accessHL()),
	// LD DE,nn
	0x11: ldInto16Const(0x11, "DE", accessDE()),
	// LD A,n
	0x3E: ldInto8Const(0x3E, "A", accessA()),
	// LD B,n
	0x06: ldInto8Const(0x3E, "B", accessB()),
	// LD C,n
	0x0E: ldInto8Const(0x0E, "C", accessC()),
	// LD D,n
	0x16: ldInto8Const(0x16, "D", accessD()),
	// LD E,n
	0x1E: ldInto8Const(0x1E, "E", accessE()),
	// LD L,n
	0x2E: ldInto8Const(0x2E, "L", accessL()),
	// LD A,E
	0x7B: ldInto8From8(0x7B, "A", "E", accessA(), accessE()),
	// LD A,B
	0x78: ldInto8From8(0x78, "A", "B", accessA(), accessB()),
	// LD C,A
	0x4F: ldInto8From8(0x4F, "C", "A", accessC(), accessA()),
	// LD D,A
	0x57: ldInto8From8(0x57, "D", "A", accessD(), accessA()),
	// LD H,A
	0x67: ldInto8From8(0x67, "H", "A", accessH(), accessA()),
	// LD A,H
	0x7C: ldInto8From8(0x7C, "A", "H", accessA(), accessH()),
	// LD (0xFF00+C),A
	0xE2: ldInto8RefFrom8(0xE2, "C", "A", accessC(), accessA()),
	// LD (0xFF00+n),A
	0xE0: ldInto8ConstRefFrom8(0xE0, "A", accessA()),
	// LD A,(DE)
	0x1A: ldInto8From16Ref(0x1A, "A", "DE", accessA(), accessDE()),
	// LD (HL),A
	0x77: ldInto16RefFrom8(0x77, "HL", "A", accessHL(), accessA(), 0),
	// LD (HL-),A
	0x32: ldInto16RefFrom8(0x32, "HL-", "A", accessHL(), accessA(), -1),
	// LD (HL+),A
	0x22: ldInto16RefFrom8(0x22, "HL+", "A", accessHL(), accessA(), +1),
	// LD (nn),A
	0xEA: ldInto16ConstRefFrom8(0xEA, "A", accessA()),
	// LD A,(0xFF00+n)
	0xF0: ldInto8From8ConstRef(0xF0, "A", accessA()),
}

func ldInto16Const(opcode uint16, regName string, ac accessor16Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     12,
		ArgsLength: 2,
		Label:      fmt.Sprintf("LD %s,nn", regName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			lower := uint16(args[0])
			upper := uint16(args[1])
			ac(r).Set(upper<<8 | lower)
		},
	}
}

func ldInto8Const(opcode uint16, regName string, ac accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     8,
		ArgsLength: 1,
		Label:      fmt.Sprintf("LD %s,n", regName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			acc := ac(r)
			acc.Set(args[0])
		},
	}
}

func ldInto8From16Ref(opcode uint16, dstRegName, srcRegName string, dst accessor8Creator, src accessor16Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     8,
		ArgsLength: 0,
		Label:      fmt.Sprintf("LD %s,(%s)", dstRegName, srcRegName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			val := m.GetByte(src(r).Get())
			dst(r).Set(val)
		},
	}
}

func ldInto8RefFrom8(opcode uint16, dstRegName, srcRegName string, dst accessor8Creator, src accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     8,
		ArgsLength: 0,
		Label:      fmt.Sprintf("LD (%s),%s", dstRegName, srcRegName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			addr := 0xFF00 + uint16(dst(r).Get())
			m.SetByte(addr, src(r).Get())
		},
	}
}

func ldInto8ConstRefFrom8(opcode uint16, srcRegName string, src accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     12,
		ArgsLength: 1,
		Label:      fmt.Sprintf("LD (n),%s", srcRegName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			addr := 0xFF00 + uint16(args[0])
			m.SetByte(addr, src(r).Get())
		},
	}
}

func ldInto16ConstRefFrom8(opcode uint16, srcRegName string, src accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     12,
		ArgsLength: 2,
		Label:      fmt.Sprintf("LD (nn),%s", srcRegName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			addr := combine(args[1], args[0])
			m.SetByte(addr, src(r).Get())
		},
	}
}

func ldInto8From8(opcode uint16, dstRegName, srcRegName string, dst accessor8Creator, src accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     4,
		ArgsLength: 0,
		Label:      fmt.Sprintf("LD %s,%s", dstRegName, srcRegName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			dst(r).Set(src(r).Get())
		},
	}
}

func ldInto16RefFrom8(opcode uint16, dstRegName, srcRegName string, dst accessor16Creator, src accessor8Creator, addOffset int8) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     8,
		ArgsLength: 0,
		Label:      fmt.Sprintf("LD (%s),%s", dstRegName, srcRegName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			dstVal := dst(r).Get()
			m.SetByte(dstVal, src(r).Get())
			if addOffset == -1 {
				dst(r).Set(dstVal - 1)
			} else if addOffset == 1 {
				dst(r).Set(dstVal + 1)
			}
		},
	}
}

func ldInto8From8ConstRef(opcode uint16, dstRegName string, dst accessor8Creator) Command {
	return Command{
		OpCode:     opcode,
		Cycles:     12,
		ArgsLength: 1,
		Label:      fmt.Sprintf("LD %s,(n)", dstRegName),
		Run: func(r *Registers, m memory.Memory, args []uint8) {
			v := m.GetByte(0xFF00 + uint16(args[0]))
			dst(r).Set(v)
		},
	}
}
