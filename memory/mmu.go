package memory

const bootstrapCompletedAddr = uint16(0xFF50)

// MMU manages access to the memory
type MMU struct {
	spaces []AddressSpace
}

// NewMMU creates a new MMU.
func NewMMU(boot *ROM, spaces ...AddressSpace) *MMU {
	return &MMU{
		spaces: append([]AddressSpace{boot}, spaces...),
	}
}

func (c *MMU) spaceForAddr(addr uint16) AddressSpace {
	for _, s := range c.spaces {
		if s.Contains(addr) {
			return s
		}
	}
	return nil
}

// Contains returns true when the address is part of the address space.
func (c *MMU) Contains(addr uint16) bool {
	return c.spaceForAddr(addr) != nil
}

// Read returns the byte at the given address.
// Panics when the address is not in the address space.
func (c *MMU) Read(addr uint16) uint8 {
	if s := c.spaceForAddr(addr); s != nil {
		return s.Read(addr)
	}
	return 0xFF
}

// Write writes a value at the given address.
// Panics when the address is not in the address space.
func (c *MMU) Write(addr uint16, v uint8) {
	if addr == bootstrapCompletedAddr {
		c.disableBootRom()
	}
	if s := c.spaceForAddr(addr); s != nil {
		s.Write(addr, v)
	}
}

func (c *MMU) disableBootRom() {
	// the boot rom is the first address space.
	c.spaces = c.spaces[1:]
}
