package memory

import "fmt"

// RAM is a Random Access Memory.
type RAM struct {
	bytes  []uint8
	offset uint16
}

// NewRAM creates a new RAM with the given size
// and address space starting at the given offset.
func NewRAM(size uint16, offset uint16) *RAM {
	return &RAM{
		bytes:  make([]uint8, size),
		offset: offset,
	}
}

// Contains returns true when the address is part of the address space.
func (r *RAM) Contains(addr uint16) bool {
	return addr >= r.offset && addr < r.offset+uint16(len(r.bytes))
}

func (r *RAM) validateAddress(addr uint16) {
	if !r.Contains(addr) {
		msg := fmt.Sprintf("accessing invalid memory address 0x%04x in range (0x%04x, 0x%04x)", addr, r.offset, r.offset+uint16(len(r.bytes)))
		panic(msg)
	}
}

// Read returns the byte at the given address.
// Panics when the address is not in the address space.
func (r *RAM) Read(addr uint16) uint8 {
	r.validateAddress(addr)
	return r.bytes[addr-r.offset]
}

// Write writes a value at the given address.
// Panics when the address is not in the address space.
func (r *RAM) Write(addr uint16, v uint8) {
	r.validateAddress(addr)
	r.bytes[addr-r.offset] = v
}
