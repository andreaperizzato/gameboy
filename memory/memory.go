package memory

// AddressSpace is a region in memory.
type AddressSpace interface {
	Contains(addr uint16) bool
	Read(addr uint16) uint8
	Write(addr uint16, v uint8)
}

// Register is a special byte in memory.
type Register struct {
	Address uint16
	mask    uint8
	shift   uint8
	mem     AddressSpace
}

// NewRegister returns a register mapping the byte at the address in the address space.
func NewRegister(mem AddressSpace, addr uint16) Register {
	return NewRegisterWithMask(mem, addr, 0xFF)
}

// NewRegisterWithMask returns a Register with a mask which is applied
// to the byte to only consider some bits. An example when you might need to apply
// a mask is when you want to access the volume of register NR12 which is
// stored in the 4 most significant bits:
// volume := NewRegisterWithMask(m, 0xFF12, 0xF0)
// Calling volume.Get() would return the value of the high nibble and
// volume.Set(v) would only set those bits.
func NewRegisterWithMask(mem AddressSpace, addr uint16, mask uint8) Register {
	// the shift is equal to the first non-zero bit from the least significant
	// and it's the number of bits we need to shift to get the masked value.
	shift := uint8(8)
	for i := 0; i < 8; i++ {
		if (mask>>i)&1 == 1 {
			shift = uint8(i)
			break
		}
	}
	return Register{
		Address: addr,
		mem:     mem,
		mask:    mask,
		shift:   shift,
	}
}

// Set sets the value of the register.
func (r Register) Set(v uint8) {
	curr := r.mem.Read(r.Address)
	new := (v << r.shift) | curr&(0xFF-r.mask)
	r.mem.Write(r.Address, new)
}

// Get gets the value of the register.
func (r Register) Get() uint8 {
	v := r.mem.Read(r.Address)
	return (v & r.mask) >> r.shift
}

// GetBit gets the value of the b-th bit of the register.
// The mask is not applied.
//
// Deprecated: prefer to use a register with mask.
func (r Register) GetBit(b uint8) bool {
	v := r.mem.Read(r.Address)
	return (v>>b)&0x01 == 0x01
}

// Register16 is a 16-bit register made up of two 8-bit registers.
type Register16 struct {
	low  Register
	high Register
}

// NewRegister16 creates a new Register16.
func NewRegister16(mem AddressSpace, lowAddr, highAddr uint16) Register16 {
	return NewRegister16WithMask(mem, lowAddr, highAddr, 0xFF)
}

// NewRegister16WithMask creates a new Register16 applying a mask to the high byte register.
func NewRegister16WithMask(mem AddressSpace, lowAddr, highAddr uint16, highMask uint8) Register16 {
	return Register16{
		low:  NewRegister(mem, lowAddr),
		high: NewRegisterWithMask(mem, highAddr, highMask),
	}
}

// Set sets the value of the register.
func (r Register16) Set(v uint16) {
	r.low.Set(uint8(v))
	r.high.Set(uint8((v & 0xFF00) >> 8))
}

// Get gets the value of the register.
func (r Register16) Get() uint16 {
	l := r.low.Get()
	h := r.high.Get()
	return uint16(h)<<8 | uint16(l)
}
