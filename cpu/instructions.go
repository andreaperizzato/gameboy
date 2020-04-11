package cpu

func inc8(reg reg8) runnable {
	return func(c *CPU) uint8 {
		// Flags: Z 0 H -
		set, get := reg(c)
		set(add(c, get(), 1, true))
		return 4
	}
}

func inc16(reg reg16) runnable {
	return func(c *CPU) uint8 {
		// Flags: Z 0 H -
		set, get := reg(c)
		set(get() + 1)
		return 8
	}
}

func dec8(reg reg8) runnable {
	return func(c *CPU) uint8 {
		// Flags: Z 1 H -
		set, get := reg(c)
		set(sub(c, get(), 1, true))
		return 4
	}
}

func ld8Const(reg reg8) runnable {
	return func(c *CPU) uint8 {
		set, _ := reg(c)
		set(nextArg(c))
		return 8
	}
}

func ld16Const(reg reg16) runnable {
	return func(c *CPU) uint8 {
		set, _ := reg(c)
		low := nextArg(c)
		high := nextArg(c)
		v := uint16(low) | (uint16(high) << 8)
		set(v)
		return 8
	}
}

func ld88(dst, src reg8) runnable {
	return func(c *CPU) uint8 {
		_, get := src(c)
		set, _ := dst(c)
		set(get())
		return 4
	}
}

// jr implements conditional and normal jumps.
// set flag = nil to ignore the condition.
func jr(flag func(c *CPU) bool, condition bool) runnable {
	return func(c *CPU) uint8 {
		offset := int8(nextArg(c)) // args[0] is a signed byte.
		if flag == nil || flag(c) == condition {
			// we need int16 as we have a possible subtraction
			// between two 16-bit numbers.
			c.regs.PC = uint16(int16(c.regs.PC) + int16(offset))
			return 12
		}
		// According to https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html
		// 'JR n' always taks 12 cycles, but JR c,n takes 8 when not jumping.
		return 8
	}
}

// ld816Ref implements instructions such as LD A,(HL)
// loading in A the byte in memory at the address pointed by HL.
func ld816Ref(reg reg8, ptr reg16) runnable {
	return func(c *CPU) uint8 {
		_, getPtr := ptr(c)
		v := c.mem.Read(getPtr())
		setReg, _ := reg(c)
		setReg(v)
		return 8
	}
}

// ld16Ref8 implements instructions as 'LD (HL),A' and 'LD (HL+),A'
// where it writes the value in A at the location in memory pointed by HL
// and then increments HL by the given offset.
func ld16Ref8(ptr reg16, src reg8, offset int16) runnable {
	return func(c *CPU) uint8 {
		_, getSrc := src(c)
		setPtr, getPtr := ptr(c)
		c.mem.Write(getPtr(), getSrc())
		newPtr := uint16(int16(getPtr()) + offset)
		setPtr(newPtr)
		return 8
	}
}

// ld8ConstRef8 implements instructions like `LD (0x00FF+n),A`
// sometimes also called `LD (n),A` or `LDH (n),A`.
func ld8ConstRef8(src reg8) runnable {
	return func(c *CPU) uint8 {
		_, get := src(c)
		offset := nextArg(c)
		c.mem.Write(0xFF00+uint16(offset), get())
		return 12
	}
}

// ld8Ref8 implements instructions like `LD (0x00FF+C),A`.
func ld8Ref8(dst, src reg8) runnable {
	return func(c *CPU) uint8 {
		_, getSrc := src(c)
		_, getDst := dst(c)
		c.mem.Write(0xFF00+uint16(getDst()), getSrc())
		return 8
	}
}

// ld16ConstRef8 implements instructions as 'LD (nn),A'.
func ld16ConstRef8(reg reg8) runnable {
	return func(c *CPU) uint8 {
		_, get := reg(c)
		low := nextArg(c)
		high := nextArg(c)
		v := uint16(low) | (uint16(high) << 8)
		c.mem.Write(v, get())
		return 16
	}
}

// ld88ConstRef implements instructions such as LD A,(n)
// sometimes noted as LD A,(0xFF00+n) or LDH A,(n)
func ld88ConstRef(reg reg8) runnable {
	return func(c *CPU) uint8 {
		offset := nextArg(c)
		v := c.mem.Read(0xFF00 + uint16(offset))
		set, _ := reg(c)
		set(v)
		return 12
	}
}

// add16Ref implements instructions like 'ADD (HL)',
// adding to A the value pointed by HL.
func add16Ref(reg reg16) runnable {
	return func(c *CPU) uint8 {
		// Flags: Z 0 H C
		_, get := reg(c)
		c.regs.A = add(c, c.regs.A, c.mem.Read(get()), false)
		return 8
	}
}

func xor8(reg reg8) runnable {
	return func(c *CPU) uint8 {
		// Flags: Z 0 0 0
		_, get := reg(c)
		c.regs.A ^= get()
		c.flags.Z = c.regs.A == 0
		c.flags.N, c.flags.C, c.flags.H = false, false, false
		return 4
	}
}

func sub8(reg reg8) runnable {
	return func(c *CPU) uint8 {
		// Flags: Z 1 H C
		_, get := reg(c)
		c.regs.A = sub(c, c.regs.A, get(), true)
		return 4
	}
}

func cp16Ref(reg reg16) runnable {
	return func(c *CPU) uint8 {
		// Flags: Z 1 H C
		_, get := reg(c)
		v := c.mem.Read(get())
		_ = sub(c, c.regs.A, v, false)
		return 8
	}
}

func cpConst() runnable {
	return func(c *CPU) uint8 {
		// Flags: Z 1 H C
		v := nextArg(c)
		_ = sub(c, c.regs.A, v, false)
		return 8
	}
}

func pop16(reg reg16) runnable {
	return func(c *CPU) uint8 {
		set, _ := reg(c)
		set(pop(c))
		return 12
	}
}

func push16(reg reg16) runnable {
	return func(c *CPU) uint8 {
		_, get := reg(c)
		push(c, get())
		return 16
	}
}

func ret() runnable {
	return func(c *CPU) uint8 {
		c.regs.PC = pop(c)
		return 16
	}
}

func call() runnable {
	return func(c *CPU) uint8 {
		low := nextArg(c)
		high := nextArg(c)
		push(c, c.regs.PC)
		c.regs.PC = uint16(low) | uint16(high)<<8
		return 24
	}
}

// bit8 implements instuctions like 'BIT 7,A'
func bit8(pos uint8, reg reg8) runnable {
	return func(c *CPU) uint8 {
		// Flags: Z 0 1 -
		_, get := reg(c)
		c.flags.N = false
		c.flags.H = true
		c.flags.Z = (1<<pos)&get() == 0
		return 8
	}
}

// rl8 implements instuctions like 'RL A'.
// Rotates A through the carry flag.
func rl8(reg reg8) runnable {
	return func(c *CPU) uint8 {
		// Flags: Z 0 0 C
		c.flags.N = false
		c.flags.H = false

		set, get := reg(c)
		v := get()
		res := (v << 1) & 0xFF
		if c.flags.C {
			res |= 1
		}
		// If the 7th bit was set, carry will happen.
		c.flags.C = v&(1<<7) > 0
		c.flags.Z = res == 0
		set(res)
		return 8
	}
}

// add computes a+b and sets the flags accordingly.
// when ignoreCarry is set, flag C is left as it is.
func add(c *CPU, a, b uint8, ignoreCarry bool) uint8 {
	c.flags.N = false
	// https://robdor.com/2016/08/10/gameboy-emulator-half-carry-flag/
	c.flags.H = a&0x0F+b&0x0F == 0x10
	if !ignoreCarry {
		// if the sum is greater than 0xFF, overflow will occur.
		c.flags.C = uint16(a)+uint16(b) > 0xFF
	}
	res := a + b
	c.flags.Z = res == 0
	return res
}

// sub computes a-b and sets the flags accordingly.
// when ignoreCarry is set, flag C is left as it is.
func sub(c *CPU, a, b uint8, ignoreCarry bool) uint8 {
	// Flags: Z 1 H C/-
	c.flags.N = true
	// if low-b is bigger then low-a then the lower nibble will wrap.
	c.flags.H = b&0x0F > a&0x0F
	if !ignoreCarry {
		// if b > a, then the result will wrap.
		c.flags.C = b > a
	}
	res := a - b
	c.flags.Z = res == 0
	return res
}

func nextArg(c *CPU) uint8 {
	v := c.mem.Read(c.regs.PC)
	c.regs.PC++
	return v
}

func pop(c *CPU) uint16 {
	low := c.mem.Read(c.regs.SP)
	high := c.mem.Read(c.regs.SP + 1)
	c.regs.SP += 2
	return uint16(low) | uint16(high)<<8
}

func push(c *CPU, v uint16) {
	c.regs.SP -= 2
	c.mem.Write(c.regs.SP, uint8(v&0xFF)) // low nibble
	c.mem.Write(c.regs.SP+1, uint8(v>>8)) // high nibble
}
