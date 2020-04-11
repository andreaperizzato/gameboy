package cpu

type setter8 func(v uint8)
type getter8 func() uint8
type reg8 func(c *CPU) (setter8, getter8)

func access8(r *uint8) (setter8, getter8) {
	set := func(v uint8) { *r = v }
	get := func() uint8 { return *r }
	return set, get
}

type setter16 func(v uint16)
type getter16 func() uint16
type reg16 func(c *CPU) (setter16, getter16)

func access16Pair(high *uint8, low *uint8) (setter16, getter16) {
	set := func(v uint16) {
		*low = uint8(v & 0xFF)
		*high = uint8((v >> 8) & 0xFF)
	}
	get := func() uint16 {
		return uint16(*low) | uint16(*high)<<8
	}
	return set, get
}

func access16(r *uint16) (setter16, getter16) {
	set := func(v uint16) { *r = v }
	get := func() uint16 { return *r }
	return set, get
}

func regA(c *CPU) (setter8, getter8) {
	return access8(&c.regs.A)
}

func regB(c *CPU) (setter8, getter8) {
	return access8(&c.regs.B)
}

func regC(c *CPU) (setter8, getter8) {
	return access8(&c.regs.C)
}

func regD(c *CPU) (setter8, getter8) {
	return access8(&c.regs.D)
}

func regH(c *CPU) (setter8, getter8) {
	return access8(&c.regs.H)
}

func regL(c *CPU) (setter8, getter8) {
	return access8(&c.regs.L)
}

func regE(c *CPU) (setter8, getter8) {
	return access8(&c.regs.E)
}

func regDE(c *CPU) (setter16, getter16) {
	return access16Pair(&c.regs.D, &c.regs.E)
}

func regHL(c *CPU) (setter16, getter16) {
	return access16Pair(&c.regs.H, &c.regs.L)
}

func regBC(c *CPU) (setter16, getter16) {
	return access16Pair(&c.regs.B, &c.regs.C)
}

func regSP(c *CPU) (setter16, getter16) {
	return access16(&c.regs.SP)
}

func flagZ(c *CPU) bool {
	return c.flags.Z
}

func flagC(c *CPU) bool {
	return c.flags.C
}
