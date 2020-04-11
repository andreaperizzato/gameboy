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
	mem     AddressSpace
}

// NewRegister returns a register mapping the byte at the address in the address space.
func NewRegister(mem AddressSpace, addr uint16) Register {
	return Register{
		Address: addr,
		mem:     mem,
	}
}

// Set sets the value of the register.
func (r Register) Set(v uint8) {
	r.mem.Write(r.Address, v)
}

// Get gets the value of the register.
func (r Register) Get() uint8 {
	return r.mem.Read(r.Address)
}

// GetBit gets the value of the b-th bit of the register.
func (r Register) GetBit(b uint8) bool {
	return (r.Get()>>b)&0x01 == 0x01
}
