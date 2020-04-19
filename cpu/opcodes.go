package cpu

var gbcInstructions = map[uint16]instruction{
	// 0x
	0x00: {"NOP", nop()},
	0x01: {"LD BC,nn", ld16Const(regBC)},
	0x02: {"LD (BC),A", ld16Ref8(regBC, regA, 0)},
	0x03: {"INC BC", inc16(regBC)},
	0x04: {"INC B", inc8(regB)},
	0x05: {"DEC B", dec8(regB)},
	0x06: {"LD B,n", ld8Const(regB)},
	0x07: {"RLCA", rlc8(regA)},
	0x08: {"LD (nn),SP", ld16ConstRefSP()},
	0x09: {"ADD HL,BC", add16(regHL, regBC)},
	0x0A: {"LD A,(BC)", ld816Ref(regA, regBC, 0)},
	0x0B: {"DEC BC", dec16(regBC)},
	0x0C: {"INC C", inc8(regC)},
	0x0D: {"DEC C", dec8(regC)},
	0x0E: {"LD C,n", ld8Const(regC)},
	0x0F: {"RRCA", rrc8(regA)},
	// 1x
	0x10: {"STOP", nop()}, // Stop updates the CPU state, nothing to do here.
	0x11: {"LD DE,nn", ld16Const(regDE)},
	0x12: {"LD (DE),A", ld16Ref8(regDE, regA, 0)},
	0x13: {"INC DE", inc16(regDE)},
	0x14: {"INC D", inc8(regD)},
	0x15: {"DEC D", dec8(regD)},
	0x16: {"LD D,n", ld8Const(regD)},
	0x17: {"RLA", rl8(regA)},
	0x18: {"JR n", jr(nil, true)},
	0x19: {"ADD HL,DE", add16(regHL, regDE)},
	0x1A: {"LD A,(DE)", ld816Ref(regA, regDE, 0)},
	0x1B: {"DEC DE", dec16(regDE)},
	0x1C: {"INC E", inc8(regE)},
	0x1D: {"DEC E", dec8(regE)},
	0x1E: {"LD E,n", ld8Const(regE)},
	0x1F: {"RRA", rr8(regA)},
	// 2x
	0x20: {"JR NZ,n", jr(flagZ, false)},
	0x21: {"LD HL,nn", ld16Const(regHL)},
	0x22: {"LD (HL+),A", ld16Ref8(regHL, regA, 1)},
	0x23: {"INC HL", inc16(regHL)},
	0x24: {"INC H", inc8(regH)},
	0x25: {"DEC H", dec8(regH)},
	0x26: {"LD H,n", ld8Const(regH)},
	0x27: {"DAA", daa()},
	0x28: {"JR Z,n", jr(flagZ, true)},
	0x29: {"ADD HL,HL", add16(regHL, regHL)},
	0x2A: {"LD A,(HL+)", ld816Ref(regA, regHL, 1)},
	0x2B: {"DEC HL", dec16(regHL)},
	0x2C: {"INC L", inc8(regL)},
	0x2D: {"DEC L", dec8(regL)},
	0x2E: {"LD L,n", ld8Const(regL)},
	0x2F: {"CPL", cpl8(regA)},
	0x30: {"JR NC,n", jr(flagC, false)},
	0x31: {"LD SP,nn", ld16Const(regSP)},
	0x32: {"LD (HL-),A", ld16Ref8(regHL, regA, -1)},
	0x38: {"JR C,n", jr(flagC, true)},
	0x3D: {"DEC A", dec8(regA)},
	0x3E: {"LD A,n", ld8Const(regA)},
	0x4F: {"LD C,A", ld88(regC, regA)},
	0x57: {"LD D,A", ld88(regD, regA)},
	0x67: {"LD H,A", ld88(regH, regA)},
	0x77: {"LD (HL),A", ld16Ref8(regHL, regA, 0)},
	0x78: {"LD A,B", ld88(regA, regB)},
	0x7B: {"LD A,E", ld88(regA, regE)},
	0x7C: {"LD A,H", ld88(regA, regH)},
	0x7D: {"LD A,L", ld88(regA, regL)},
	0x86: {"ADD (HL)", add16Ref(regHL)},
	0x90: {"SUB B", sub8(regB)},
	0xAF: {"XOR A", xor8(regA)},
	0xBE: {"CP (HL)", cp16Ref(regHL)},
	0xC1: {"POP BC", pop16(regBC)},
	0xC5: {"PUSH BC", push16(regBC)},
	0xC9: {"RET", ret()},
	0xCD: {"CALL nn", call()},
	0xE0: {"LD (n),A", ld8ConstRef8(regA)},
	0xE2: {"LD (C),A", ld8Ref8(regC, regA)},
	0xEA: {"LD (nn),A", ld16ConstRef8(regA)},
	0xF0: {"LD A,(n)", ld88ConstRef(regA)},
	0xFE: {"CP n", cpConst()},

	// Extended set.
	0xCB11: {"RL C", rl8(regC)},
	0xCB17: {"RL A", rl8(regA)},
	0xCB7C: {"BIT 7,H", bit8(7, regH)},
}
