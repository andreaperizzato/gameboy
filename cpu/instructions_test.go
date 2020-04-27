package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstructions_inc8(t *testing.T) {
	tests := []struct {
		name    string
		initalV uint8
		finalV  uint8
		flags   flags
	}{
		{"result is zero", 0xFF, 0x00, flags{Z: true, H: true}},
		{"half carry", 0xAF, 0xB0, flags{H: true}},
		{"no half carry", 0x0D, 0x0E, flags{}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.regs.B = tC.initalV
			cycles := inc8(regB)(c)

			assert.EqualValues(t, 4, cycles, "cycles")
			assert.Equal(t, tC.finalV, c.regs.B, "value")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_inc16(t *testing.T) {
	c := &CPU{}
	c.regs.D = 0x15
	c.regs.E = 0xFF
	cycles := inc16(regDE)(c)

	assert.EqualValues(t, 8, cycles, "cycles")
	assert.Equal(t, uint8(0x16), c.regs.D, "value")
	assert.Equal(t, uint8(0x00), c.regs.E, "value")
}

func TestInstructions_inc16Ref(t *testing.T) {
	tests := []struct {
		name    string
		initalV uint8
		finalV  uint8
		flags   flags
	}{
		{"result is zero", 0xFF, 0x00, flags{Z: true, H: true}},
		{"half carry", 0xAF, 0xB0, flags{H: true}},
		{"no half carry", 0x0D, 0x0E, flags{}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			mem := make(simpleRAM, 0xFFFF)
			c := &CPU{mem: mem}
			c.regs.H, c.regs.L = 0x11, 0x22
			mem[0x1122] = tC.initalV
			cycles := inc16Ref(regHL)(c)

			assert.EqualValues(t, 12, cycles, "cycles")
			assert.Equal(t, tC.finalV, mem[0x1122], "value")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_dec8(t *testing.T) {
	tests := []struct {
		name    string
		initalV uint8
		finalV  uint8
		flags   flags
	}{
		{"result is zero", 0x01, 0x00, flags{Z: true, N: true}},
		{"half carry", 0xF0, 0xEF, flags{H: true, N: true}},
		{"no half carry", 0xF1, 0xF0, flags{N: true}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.regs.B = tC.initalV
			cycles := dec8(regB)(c)

			assert.EqualValues(t, 4, cycles, "cycles")
			assert.Equal(t, tC.finalV, c.regs.B, "value")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_dec16(t *testing.T) {
	c := &CPU{}
	c.regs.B, c.regs.C = 0x01, 0x00

	cycles := dec16(regBC)(c)

	assert.EqualValues(t, 8, cycles, "cycles")
	// 0x0100 - 0x0001 = 0x00FF
	assert.EqualValues(t, 0x00, c.regs.B, "B")
	assert.EqualValues(t, 0xFF, c.regs.C, "C")
}

func TestInstructions_dec16Ref(t *testing.T) {
	tests := []struct {
		name    string
		initalV uint8
		finalV  uint8
		flags   flags
	}{
		{"result is zero", 0x01, 0x00, flags{Z: true, N: true}},
		{"half carry", 0xF0, 0xEF, flags{H: true, N: true}},
		{"no half carry", 0xF1, 0xF0, flags{N: true}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			mem := make(simpleRAM, 0xFFFF)
			c := &CPU{mem: mem}
			c.regs.H, c.regs.L = 0x11, 0x22
			mem[0x1122] = tC.initalV
			cycles := dec16Ref(regHL)(c)

			assert.EqualValues(t, 12, cycles, "cycles")
			assert.Equal(t, tC.finalV, mem[0x1122], "value")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_sub8(t *testing.T) {
	tests := []struct {
		name  string
		a     uint8
		b     uint8
		res   uint8 // a - b
		flags flags
	}{
		{"result is zero", 0x02, 0x02, 0x00, flags{Z: true, N: true}},
		{"half carry", 0xAA, 0x2F, 0x7B, flags{H: true, N: true}},
		{"no half carry", 0xAA, 0xA1, 0x09, flags{N: true}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.regs.A = tC.a
			c.regs.B = tC.b
			cycles := sub8(regB)(c)

			assert.EqualValues(t, 4, cycles, "cycles")
			assert.Equal(t, tC.res, c.regs.A, "A")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_ld8Const(t *testing.T) {
	c := &CPU{
		regs: registers{
			PC: 0x01,
		},
		mem: simpleRAM{0xAA, 0xBB, 0xCC},
	}
	c.regs.B = 0xFF
	cycles := ld8Const(regB)(c)

	assert.EqualValues(t, 8, cycles, "cycles")
	assert.EqualValues(t, 0x02, c.regs.PC, "PC")
	assert.EqualValues(t, 0xBB, c.regs.B, "B")
}

func TestInstructions_ld88(t *testing.T) {
	c := &CPU{}
	c.regs.A = 0xAA
	c.regs.B = 0xBB
	cycles := ld88(regA, regB)(c)

	assert.EqualValues(t, 4, cycles, "cycles")
	assert.EqualValues(t, 0xBB, c.regs.A, "A")
	assert.EqualValues(t, 0xBB, c.regs.B, "B")
}

func TestInstructions_ld16Const(t *testing.T) {
	c := &CPU{
		regs: registers{
			PC: 0x01,
		},
		mem: simpleRAM{0xAA, 0xBB, 0xCC},
	}
	c.regs.D = 0xFF
	c.regs.E = 0xEE
	cycles := ld16Const(regDE)(c)

	assert.EqualValues(t, 8, cycles, "cycles")
	assert.EqualValues(t, 0x03, c.regs.PC, "PC")
	assert.EqualValues(t, 0xBB, c.regs.E, "E")
	assert.EqualValues(t, 0xCC, c.regs.D, "D")
}

func TestInstructions_jr(t *testing.T) {
	tests := []struct {
		name      string
		flag      func(*CPU) bool
		condition bool
		cycles    uint8
		taken     bool
	}{
		{"unconditional jump", nil, false, 12, true},
		{"conditional jump taken", flagZ, true, 12, true},
		{"conditional jump not taken", flagZ, false, 8, false},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{
				regs: registers{
					PC: 0x02,
				},
				flags: flags{
					Z: true,
				},
				mem: simpleRAM{0x00, 0x00, 0xFE, 0x11, 0x00},
			}
			cycles := jr(tC.flag, tC.condition)(c)

			if tC.taken {
				// The current PC is 0x02 so the next arg will be
				// 0xFE (3th index in memory). Therefore
				// a jump should increment the PC by int8(0xFE)
				// which is -1 to 0x01.
				assert.Equal(t, uint16(0x01), c.regs.PC, "PC")
			} else {
				// The current PC is 0x02 so if we don't take
				// the jump, we expect 0x03 since we read the
				// next arg.
				assert.Equal(t, uint16(0x03), c.regs.PC, "PC")
			}
			assert.Equal(t, tC.cycles, cycles, "cycles")
		})
	}
}

func TestInstructions_ld816Ref(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	mem[0x1122] = 0x99
	c := &CPU{mem: mem}
	c.regs.D, c.regs.E = 0x11, 0x22
	cycles := ld816Ref(regA, regDE, -1)(c)

	assert.EqualValues(t, 8, cycles, "cycles")
	// DE = 0x1122 at which the memory stores 0x99
	// which is the value we expect in A.
	assert.EqualValues(t, 0x99, c.regs.A, "A")
	// DE should've been decremented by one.
	assert.EqualValues(t, 0x11, c.regs.D, "D")
	assert.EqualValues(t, 0x21, c.regs.E, "E")
}

func TestInstructions_ld16Ref8(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	mem[0x1122] = 0x99
	c := &CPU{
		regs: registers{
			A: 0xAA,
		},
		mem: mem,
	}
	c.regs.D, c.regs.E = 0x11, 0x22
	cycles := ld16Ref8(regDE, regA, -1)(c)

	assert.EqualValues(t, 8, cycles, "cycles")
	// A=0xAA and DE=0x1122 so we expect
	// mem[0x1122] to be set to 0xAA (value of A)
	// and DE to be decreased by 1 (we are passing -1 above).
	assert.EqualValues(t, 0xAA, mem[0x1122], "(DE)")
	assert.EqualValues(t, 0x11, c.regs.D, "D")
	assert.EqualValues(t, 0x21, c.regs.E, "E")
}

func TestInstructions_ld8ConstRef8(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := &CPU{mem: mem}
	c.regs.D = 0xAA
	c.regs.PC = 0x0011
	mem[0x0011] = 0x77

	cycles := ld8ConstRef8(regD)(c)

	assert.EqualValues(t, 12, cycles, "cycles")
	// D=0xAA and the argument of the op is n=0x77.
	// We expect the op to write 0xAA at 0xFF00+0x77.
	assert.EqualValues(t, 0xAA, mem[0xFF77], "0xFF77")
}

func TestInstructions_ld16ConstRef8(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := &CPU{mem: mem}
	c.regs.H = 0xAA
	c.regs.PC = 0x0011
	mem[0x0011] = 0x77
	mem[0x0012] = 0x66

	cycles := ld16ConstRef8(regH)(c)

	assert.EqualValues(t, 16, cycles, "cycles")
	// H=0xAA and the arguments of the op are 0x77 and 0x66.
	// We expect the op to write 0xAA at 0x6677.
	assert.EqualValues(t, 0xAA, mem[0x6677], "0x6677")
}

func TestInstructions_ld8Ref8(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := &CPU{mem: mem}
	c.regs.L = 0xCC
	c.regs.A = 0xAA

	cycles := ld8Ref8(regL, regA)(c)

	assert.EqualValues(t, 8, cycles, "cycles")
	// L=0xCC and A=0xAA.
	// We expect the op to write 0xAA at 0xFF00+0xCC.
	assert.EqualValues(t, 0xAA, mem[0xFFCC], "0xFFCC")
}

func TestInstructions_ld88ConstRef(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := &CPU{mem: mem}
	c.regs.PC = 0x0011
	mem[0x0011] = 0x77
	mem[0xFF77] = 0x99

	cycles := ld88ConstRef(regE)(c)

	assert.EqualValues(t, 12, cycles, "cycles")
	// The arg is n=0x77 and mem[0xFF00+0x77] = 0x99.
	// We expect the op to write 0x99 in E.
	assert.EqualValues(t, 0x99, c.regs.E, "E")
}

func TestInstructions_ld16RefConst(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := &CPU{mem: mem}
	c.regs.PC = 0x0011
	mem[0x0011] = 0xAA
	c.regs.D, c.regs.E = 0x11, 0x22

	cycles := ld16RefConst(regDE)(c)

	assert.EqualValues(t, 12, cycles, "cycles")
	// The arg is n=0xAA and DE=0x1122.
	// We expect the op to write 0xAA in memory at 0x1122.
	assert.EqualValues(t, 0xAA, mem[0x1122])
}

func TestInstructions_add16Ref(t *testing.T) {
	tests := []struct {
		name  string
		a     uint8
		b     uint8
		res   uint8 // a + b
		flags flags
	}{
		{"result is zero", 0xFE, 0x02, 0x00, flags{Z: true, C: true, H: true}},
		{"half carry", 0x0C, 0x04, 0x10, flags{H: true}},
		{"no half carry", 0x0C, 0x01, 0x0D, flags{}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			mem := make(simpleRAM, 0xFFFF)
			c := &CPU{mem: mem}
			mem[0xDDEE] = tC.b
			c.regs.D = 0xDD
			c.regs.E = 0xEE
			c.regs.A = tC.a

			cycles := add16Ref(regDE)(c)

			assert.EqualValues(t, 8, cycles, "cycles")
			assert.Equal(t, tC.res, c.regs.A, "A")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_add16(t *testing.T) {
	tests := []struct {
		name  string
		a     uint16
		b     uint16
		res   uint16 // a + b
		flags flags
	}{
		{"carry", 0xFF02, 0xA001, 0x9F03, flags{C: true}},
		{"no carry", 0x0C01, 0x0A02, 0x1603, flags{}},
		{"half carry", 0x0C10, 0x0420, 0x1030, flags{H: true}},
		{"no half carry", 0x0C10, 0x0120, 0x0D30, flags{}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.regs.H, c.regs.L = uint8(tC.a>>8), uint8(tC.a&0xFF)
			c.regs.B, c.regs.C = uint8(tC.b>>8), uint8(tC.b&0xFF)

			cycles := add16(regHL, regBC)(c)

			res := uint16(c.regs.H)<<8 | uint16(c.regs.L)

			assert.EqualValues(t, 8, cycles, "cycles")
			assert.Equal(t, tC.res, res, "HL")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_xor8(t *testing.T) {
	tests := []struct {
		name  string
		a     uint8
		b     uint8
		res   uint8 // a ^ b
		flags flags
	}{
		{"result is zero", 0xAA, 0xAA, 0x00, flags{Z: true}},
		{"result is not zero", 0xAA, 0xCC, 0x66, flags{}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.flags = flags{true, true, true, true}
			c.regs.A = tC.a
			c.regs.C = tC.b

			cycles := xor8(regC)(c)

			assert.EqualValues(t, 4, cycles, "cycles")
			assert.Equal(t, tC.res, c.regs.A, "A")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_cp16Ref(t *testing.T) {
	tests := []struct {
		name  string
		a     uint8
		b     uint8
		flags flags
	}{
		{"result is zero", 0xA1, 0xA1, flags{Z: true, N: true}},
		{"result is not zero", 0xA2, 0xA1, flags{N: true}},
		{"carry", 0xA2, 0xB1, flags{N: true, C: true}},
		{"half carry", 0x2A, 0x1B, flags{H: true, N: true}},
		{"no half carry", 0x2C, 0x1A, flags{N: true}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			mem := make(simpleRAM, 0xFFFF)
			c := &CPU{mem: mem}
			c.flags = flags{true, false, true, true}
			c.regs.A = tC.a
			mem[0xDDEE] = tC.b
			c.regs.D = 0xDD
			c.regs.E = 0xEE

			cycles := cp16Ref(regDE)(c)

			assert.EqualValues(t, 8, cycles, "cycles")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_cpConst(t *testing.T) {
	tests := []struct {
		name  string
		a     uint8
		b     uint8
		flags flags
	}{
		{"result is zero", 0xA1, 0xA1, flags{Z: true, N: true}},
		{"result is not zero", 0xA2, 0xA1, flags{N: true}},
		{"carry", 0xA2, 0xB1, flags{N: true, C: true}},
		{"half carry", 0x2A, 0x1B, flags{H: true, N: true}},
		{"no half carry", 0x2C, 0x1A, flags{N: true}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			mem := make(simpleRAM, 0xFFFF)
			c := &CPU{mem: mem}
			c.flags = flags{true, false, true, true}
			c.regs.A = tC.a
			c.regs.PC = 0xAABB
			mem[0xAABB] = tC.b

			cycles := cpConst()(c)

			assert.EqualValues(t, 8, cycles, "cycles")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_pop16(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := &CPU{mem: mem}
	c.regs.SP = 0x0AA00
	mem[0x0AA00] = 0x11
	mem[0x0AA01] = 0x22             // pop from the stack = increment address
	c.regs.B, c.regs.C = 0xBB, 0xCC // some random value

	cycles := pop16(regBC)(c)

	assert.EqualValues(t, 12, cycles, "cycles")
	assert.Equal(t, uint8(0x22), c.regs.B, "B")
	assert.Equal(t, uint8(0x11), c.regs.C, "C")
}

func TestInstructions_push16(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := &CPU{mem: mem}
	c.regs.SP = 0x0AA02
	c.regs.B, c.regs.C = 0xBB, 0xCC

	cycles := push16(regBC)(c)

	assert.EqualValues(t, 16, cycles, "cycles")
	assert.Equal(t, uint8(0xCC), mem[0xAA00], "C")
	assert.Equal(t, uint8(0xBB), mem[0xAA01], "B")
}

func TestInstructions_ret(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := &CPU{mem: mem}
	c.regs.SP = 0x0AA00
	mem[0x0AA00] = 0x11
	mem[0x0AA01] = 0x22

	cycles := ret()(c)

	assert.EqualValues(t, 16, cycles, "cycles")
	assert.Equal(t, uint16(0x2211), c.regs.PC, "PC")
}

func TestInstructions_call(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := &CPU{mem: mem}
	c.regs.PC = 0x1122
	c.regs.SP = 0xAA02
	mem[0x1122] = 0x44 // these values are the args of the instructions
	mem[0x1123] = 0x55 // and we'll be used to set PC.

	cycles := call()(c)

	assert.EqualValues(t, 24, cycles, "cycles")
	assert.EqualValues(t, 0x5544, c.regs.PC, "PC")
	assert.EqualValues(t, 0xAA00, c.regs.SP, "SP")                  // two bytes have been pushed
	assert.EqualValues(t, 0x24, mem[0xAA00], "Stack - low nibble")  // low nibble of old PC + 2 for args
	assert.EqualValues(t, 0x11, mem[0xAA01], "Stack - high nibble") // high nibble of old PC
}

func TestInstructions_bit8(t *testing.T) {
	tests := []struct {
		name  string
		v     uint8
		b     uint8 // checks bit b in v
		flags flags
	}{
		{"result is zero", 0b11011001, 5, flags{Z: true, H: true}},
		{"result is not zero", 0b11011001, 3, flags{H: true}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.flags = flags{false, true, true, false}
			c.regs.H = tC.v

			cycles := bit8(tC.b, regH)(c)

			assert.EqualValues(t, 8, cycles, "cycles")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_rl8(t *testing.T) {
	tests := []struct {
		name  string
		v     uint8
		carry bool
		res   uint8
		flags flags
	}{
		{"result is zero", 0x00, false, 0x00, flags{Z: true}},
		{"result is not zero", 0b00000001, false, 0b00000010, flags{}},
		{"carry was set", 0b00000001, true, 0b00000011, flags{}},
		{"sets carry", 0b10000001, false, 0b00000010, flags{C: true}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.flags = flags{false, true, true, tC.carry}
			c.regs.H = tC.v

			cycles := rl8(regH)(c)

			assert.EqualValues(t, 8, cycles, "cycles")
			assert.Equal(t, tC.res, c.regs.H, "H")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_rr8(t *testing.T) {
	tests := []struct {
		name  string
		v     uint8
		carry bool
		res   uint8
		flags flags
	}{
		{"result is zero", 0x00, false, 0x00, flags{Z: true}},
		{"result is not zero", 0b00000010, false, 0b00000001, flags{}},
		{"carry was set", 0b10000000, true, 0b11000000, flags{}},
		{"sets carry", 0b10000001, false, 0b01000000, flags{C: true}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.flags = flags{false, true, true, tC.carry}
			c.regs.H = tC.v

			cycles := rr8(regH)(c)

			assert.EqualValues(t, 8, cycles, "cycles")
			assert.Equal(t, tC.res, c.regs.H, "H")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_nop(t *testing.T) {
	c := &CPU{}
	reg := c.regs
	cycles := nop()(c)
	assert.EqualValues(t, 4, cycles, "cycles")
	assert.Equal(t, reg, c.regs, "registers")
}

func TestInstructions_stop(t *testing.T) {
	c := &CPU{mem: make(simpleRAM, 0xFFFF)}
	cycles := stop()(c)
	assert.EqualValues(t, 4, cycles, "cycles")
	assert.EqualValues(t, 0x01, c.regs.PC, "PC") // stop has one ignored arg.
}

func TestInstructions_ld16ConstRefSP(t *testing.T) {
	mem := make(simpleRAM, 0xFFFF)
	c := &CPU{mem: mem}
	c.regs.PC = 0x1122
	c.regs.SP = 0xAA02
	mem[0x1122] = 0x44 // these values are the args of the instructions
	mem[0x1123] = 0x55 // and we'll be used to set PC.
	cycles := ld16ConstRefSP()(c)

	assert.EqualValues(t, 20, cycles, "cycles")
	assert.EqualValues(t, 0x02, mem[0x5544], "SP - low byte")
	assert.EqualValues(t, 0xAA, mem[0x5544+1], "SP - high byte")
}

func TestInstructions_rlc8(t *testing.T) {
	tests := []struct {
		name  string
		v     uint8
		res   uint8
		flags flags
	}{
		{"no carry", 0b00000011, 0b00000110, flags{}},
		{"with carry", 0b11000000, 0b10000001, flags{C: true}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.flags = flags{true, true, true, false}
			c.regs.A = tC.v

			cycles := rlc8(regA)(c)

			assert.EqualValues(t, 4, cycles, "cycles")
			assert.Equal(t, tC.res, c.regs.A, "A")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_rrc8(t *testing.T) {
	tests := []struct {
		name  string
		v     uint8
		res   uint8
		flags flags
	}{
		{"no carry", 0b00001100, 0b00000110, flags{}},
		{"with carry", 0b00000011, 0b10000001, flags{C: true}},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.flags = flags{true, true, true, false}
			c.regs.A = tC.v

			cycles := rrc8(regA)(c)

			assert.EqualValues(t, 4, cycles, "cycles")
			assert.Equal(t, tC.res, c.regs.A, "A")
			assert.Equal(t, tC.flags, c.flags, "flags")
		})
	}
}

func TestInstructions_cpl8(t *testing.T) {
	c := &CPU{}
	c.regs.C = 0b01010001
	cycles := cpl8(regC)(c)

	assert.EqualValues(t, 4, cycles, "cycles")
	assert.EqualValues(t, 0b10101110, c.regs.C, "C")
	assert.True(t, c.flags.H, "Flag H")
	assert.True(t, c.flags.H, "Flag N")
}

func TestInstructions_scf(t *testing.T) {
	c := &CPU{}
	c.flags.H = true
	c.flags.N = true
	cycles := scf()(c)

	assert.EqualValues(t, 4, cycles, "cycles")
	assert.True(t, c.flags.C, "Flag C")
	assert.False(t, c.flags.Z, "Flag Z")
	assert.False(t, c.flags.N, "Flag N")
	assert.False(t, c.flags.H, "Flag H")
}

func TestInstructions_ccf(t *testing.T) {
	tests := []struct {
		name string
		v    bool
	}{
		{"was true", true},
		{"was false", false},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			c := &CPU{}
			c.flags.H = true
			c.flags.N = true
			c.flags.C = tC.v
			cycles := ccf()(c)

			assert.EqualValues(t, 4, cycles, "cycles")
			assert.Equal(t, !tC.v, c.flags.C, "Flag C")
			assert.False(t, c.flags.Z, "Flag Z")
			assert.False(t, c.flags.N, "Flag N")
			assert.False(t, c.flags.H, "Flag H")
		})
	}
}

func TestAllOpcodesDefined(t *testing.T) {
	// https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html

	// Base instruction set.
	invalid := []uint16{0xD3, 0xE3, 0xE4, 0xF4, 0xCB, 0xDB, 0xEB, 0xEC, 0xFC, 0xDD, 0xDE, 0xDF}
	skip := func(opcode uint16) bool {
		for _, c := range invalid {
			if opcode == c {
				return true
			}
		}
		return false
	}
	for opcode := uint16(0); opcode <= 0xFF; opcode++ {
		if skip(opcode) {
			continue
		}
		assert.NotNilf(t, gbcInstructions[opcode].run, "missing opcode 0x%04X", opcode)
	}

	// Extended instruction set.
	for opcode := uint16(0); opcode <= 0xFF; opcode++ {
		if skip(opcode) {
			continue
		}
		exCode := 0xCB00 | opcode
		assert.NotNilf(t, gbcInstructions[exCode].run, "missing opcode 0x%04X", exCode)
	}
}

type simpleRAM []uint8

func (r simpleRAM) Contains(addr uint16) bool {
	return true
}

func (r simpleRAM) Read(addr uint16) uint8 {
	return r[addr]
}

func (r simpleRAM) Write(addr uint16, v uint8) {
	r[addr] = v
}
