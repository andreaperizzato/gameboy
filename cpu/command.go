package cpu

import (
	"fmt"
	"log"
	"time"

	"github.com/andreaperizzato/gameboy/memory"
)

type Operation func(r *Registers, m memory.Memory, args []uint8)

type Command struct {
	OpCode     uint16
	Cycles     uint8
	ArgsLength uint8
	Label      string
	Run        Operation
}

var commands map[uint16]Command

func addCommands(dest, cs map[uint16]Command) {
	for k, v := range cs {
		if _, found := dest[k]; found {
			panic(fmt.Sprintf("duplicate opcode 0x%04X", k))
		}
		dest[k] = v
	}
}

func init() {
	commands = map[uint16]Command{}
	addCommands(commands, alu)
	addCommands(commands, bit)
	addCommands(commands, jump)
	addCommands(commands, ld)
}

type ram []uint8

func (r ram) GetByte(addr uint16) uint8 {
	return r[addr]
}

func (r ram) SetByte(addr uint16, v uint8) {
	r[addr] = v
}

func Run(rom []uint8) {
	reg := Registers{}
	mem := make(ram, 0xFFFF)

	copy(mem, rom)

	t := time.NewTimer(time.Second * 3)
	reg.PC = 0x0000
	reg.SP = 0x0
	for true {
		opcodePC := reg.PC
		opcode := uint16(rom[reg.PC])
		reg.PC++
		if opcode == 0xCB {
			opcode = 0xCB00 | uint16(rom[reg.PC])
			reg.PC++
		}
		cmd, found := commands[opcode]
		if !found {
			log.Fatalf("opcode 0x%04X not found at 0x%04X", opcode, reg.PC)
		}
		args := make([]uint8, cmd.ArgsLength)
		argsStr := make([]string, cmd.ArgsLength)
		for j := uint16(0); j < uint16(cmd.ArgsLength); j++ {
			args[j] = rom[reg.PC]
			argsStr[j] = fmt.Sprintf("0x%02X", args[j])
			reg.PC++
		}
		_ = opcodePC
		// fmt.Printf("0x%04X - %s - %v\n", opcodePC, cmd.Label, argsStr)

		cmd.Run(&reg, mem, args)

		// This is simulating the screen ($0064: wait for screen frame)
		select {
		case <-t.C:
			mem.SetByte(0xFF44, 0x90)
		default:
			break
		}

	}
}
