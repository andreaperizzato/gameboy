package cpu

type Registers struct {
	A  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	F  uint8
	H  uint8
	L  uint8
	SP uint16
	PC uint16
}

// combine combines two 8-bit numbers (x, y)
// into a 16-bit one (xy)
func combine(msb, lsb uint8) uint16 {
	return uint16(lsb) | (uint16(msb) << 8)
}

// setPair sets splits v into xy.
func setPair(x, y *uint8, v uint16) {
	*x = uint8((v & 0xFF00) >> 8)
	*y = uint8(v & 0xFF)
}

func (r *Registers) BC() uint16 {
	return combine(r.B, r.C)
}

func (r *Registers) SetBC(v uint16) {
	setPair(&r.B, &r.C, v)
}

func (r *Registers) AF() uint16 {
	return combine(r.A, r.F)
}

func (r *Registers) SetAF(v uint16) {
	setPair(&r.A, &r.F, v)
}

func (r *Registers) DE() uint16 {
	return combine(r.D, r.E)
}

func (r *Registers) SetDE(v uint16) {
	setPair(&r.D, &r.E, v)
}

func (r *Registers) HL() uint16 {
	return combine(r.H, r.L)
}

func (r *Registers) SetHL(v uint16) {
	setPair(&r.H, &r.L, v)
}

func bitValue(r uint8, pos uint8) bool {
	return (r>>pos)&0x01 == 1
}

func setBit(r *uint8, pos uint8, v bool) {
	if v {
		*r |= 1 << pos
	} else {
		*r &^= (1 << pos) // this is *r = *r & ^(1 << pos) where ^ is `not`
	}
	return
}

func (r *Registers) FlagZ() bool {
	return bitValue(r.F, 7)
}

func (r *Registers) SetFlagZ(v bool) {
	setBit(&r.F, 7, v)
}

func (r *Registers) FlagN() bool {
	return bitValue(r.F, 6)
}

func (r *Registers) SetFlagN(v bool) {
	setBit(&r.F, 6, v)
}

func (r *Registers) FlagH() bool {
	return bitValue(r.F, 5)
}

func (r *Registers) SetFlagH(v bool) {
	setBit(&r.F, 5, v)
}

func (r *Registers) FlagC() bool {
	return bitValue(r.F, 4)
}

func (r *Registers) SetFlagC(v bool) {
	setBit(&r.F, 4, v)
}

type reg8Accessor interface {
	Get() uint8
	Set(uint8)
}

type accessor8 struct {
	reg *uint8
}

func (a *accessor8) Get() uint8 {
	return *a.reg
}

func (a *accessor8) Set(v uint8) {
	*a.reg = v
}

type accessor8Creator func(r *Registers) reg8Accessor

type reg16Accessor interface {
	Get() uint16
	Set(uint16)
}

type accessor16 struct {
	reg *uint16
}

func (a *accessor16) Get() uint16 {
	return *a.reg
}

func (a *accessor16) Set(v uint16) {
	*a.reg = v
}

type accessor16Composite struct {
	set func(v uint16)
	get func() uint16
}

func (a *accessor16Composite) Get() uint16 {
	return a.get()
}

func (a *accessor16Composite) Set(v uint16) {
	a.set(v)
}

type accessor16Creator func(r *Registers) reg16Accessor

func accessA() accessor8Creator {
	return func(r *Registers) reg8Accessor {
		return &accessor8{reg: &r.A}
	}
}

func accessB() accessor8Creator {
	return func(r *Registers) reg8Accessor {
		return &accessor8{reg: &r.B}
	}
}

func accessC() accessor8Creator {
	return func(r *Registers) reg8Accessor {
		return &accessor8{reg: &r.C}
	}
}

func accessD() accessor8Creator {
	return func(r *Registers) reg8Accessor {
		return &accessor8{reg: &r.D}
	}
}

func accessE() accessor8Creator {
	return func(r *Registers) reg8Accessor {
		return &accessor8{reg: &r.E}
	}
}

func accessH() accessor8Creator {
	return func(r *Registers) reg8Accessor {
		return &accessor8{reg: &r.H}
	}
}

func accessL() accessor8Creator {
	return func(r *Registers) reg8Accessor {
		return &accessor8{reg: &r.L}
	}
}

func accessHL() accessor16Creator {
	return func(r *Registers) reg16Accessor {
		return &accessor16Composite{
			get: r.HL,
			set: r.SetHL,
		}
	}
}

func accessDE() accessor16Creator {
	return func(r *Registers) reg16Accessor {
		return &accessor16Composite{
			get: r.DE,
			set: r.SetDE,
		}
	}
}

func accessBC() accessor16Creator {
	return func(r *Registers) reg16Accessor {
		return &accessor16Composite{
			get: r.BC,
			set: r.SetBC,
		}
	}
}

func accessSP() accessor16Creator {
	return func(r *Registers) reg16Accessor {
		return &accessor16{reg: &r.SP}
	}
}

type flagAccessor interface {
	Get() bool
}

type boolFlagAccessor struct {
	get func() bool
}

func (a *boolFlagAccessor) Get() bool {
	return a.get()
}

type flagAccessorCreator func(r *Registers) flagAccessor

func accessFlagZ() flagAccessorCreator {
	return func(r *Registers) flagAccessor {
		return &boolFlagAccessor{
			get: r.FlagZ,
		}
	}
}

func accessFlagC() flagAccessorCreator {
	return func(r *Registers) flagAccessor {
		return &boolFlagAccessor{
			get: r.FlagC,
		}
	}
}
