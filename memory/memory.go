package memory

type Memory interface {
	GetByte(addr uint16) uint8
	SetByte(addr uint16, v uint8)
}
