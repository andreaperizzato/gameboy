package ppu

import (
	"github.com/andreaperizzato/gameboy/memory"
)

// More info about the fetcher can be found here:
// https://www.youtube.com/watch?v=HyzD8pNlpwI&t=49m17s
// https://blog.tigris.fr/2019/09/15/writing-an-emulator-the-first-pixel/

// The fetcher does 4 things in sequence:
// 1. Reads the background tile id from memory
// 2. Reads the first part of the tile data (first byte)
// 3. Reads the second part of the tile data (second byte)
// 4. Constructs 8 new pixels and puts them in the fifo queue

// fetcherState is the state of the fetcher.
type fetcherState uint8

const (
	// loads the id of the current tile from memory.
	readTileID fetcherState = iota
	// load first byte of the tile data.
	readTileData0
	// load second byte of the tile data.
	readTileData1
	// build pixels and push them to the queue.
	pushToFIFO
)

// Fetcher fetches tile data from memory and generates new pixels.
type Fetcher struct {
	Q         FIFO
	mem       memory.AddressSpace
	ticks     int
	mapAddr   uint16
	state     fetcherState
	tileLine  uint8
	tileIndex uint8
	tileID    uint8
	tileData  []uint8

	// BGP - BG Palette Data
	// https://gbdev.io/pandocs/#ff47-bgp-bg-palette-data-r-w-non-cgb-mode-only
	bgp memory.Register
}

// NewFetcher creates a new fetcher.
func NewFetcher(m memory.AddressSpace) *Fetcher {
	return &Fetcher{
		mem:      m,
		Q:        NewFIFOQueue(16),
		tileData: make([]uint8, 8),
		bgp:      memory.NewRegister(m, 0xFF47),
	}
}

// Tick runs an interation of the fetcher.
func (f *Fetcher) Tick() {
	// The fetcher runs at half the speed of the PPU.
	f.ticks++
	if f.ticks < 2 {
		return
	}
	f.ticks = 0

	switch f.state {
	case readTileID:
		f.tileID = f.mem.Read(f.mapAddr + uint16(f.tileIndex))
		f.state = readTileData0

	case readTileData0:
		f.readTileData(0, readTileData1)

	case readTileData1:
		f.readTileData(1, pushToFIFO)

	case pushToFIFO:
		if f.Q.Size() <= 8 {
			// We stored pixel bits from lest significant (rightmost)
			// to most (leftmost) in the data array, we must push them
			// in reverse order.
			for i := 7; i >= 0; i-- {
				f.Q.Push(f.tileData[i])
			}
			// Advance to the next tile in the map's row.
			f.tileIndex++
			// We're done, back from the beginning.
			f.state = readTileID
		}
	}
}

// Start fetching a line of pixels starting from the given tile address in the
// background map. Here, tileLine indicates which row of pixels to pick from
// each tile we read.
func (f *Fetcher) Start(mapAddr uint16, tileLine uint8) {
	f.tileIndex = 0
	f.mapAddr = mapAddr
	f.tileLine = tileLine
	f.state = readTileID

	f.Q.Clear()
}

func (f *Fetcher) readTileData(bitPlane uint8, nextState fetcherState) {
	// A tile's graphical data takes 16 bytes (2B per row of 8px).
	// Tile data starts at address 0x8000 so we first compute an offset
	// to find out where the data for the tile we want starts.
	offset := 0x8000 + uint16(f.tileID)*16
	// Then, from that starting offset, we compute the final address
	// to read by finding out which of the 8px (ie 2B) rows of the tile we want.
	addr := offset + uint16(f.tileLine)*2
	// Finally, read the first or second byte of graphical data depending
	// on what state we're in.
	data := f.mem.Read(addr + uint16(bitPlane))
	for bitPos := uint(0); bitPos <= 7; bitPos++ {
		if bitPlane == 0 {
			f.tileData[bitPos] = (data >> bitPos) & 1
		} else {
			f.tileData[bitPos] |= ((data >> bitPos) & 1) << 1
			f.tileData[bitPos] = f.mapColor(f.tileData[bitPos])
		}
	}
	f.state = nextState
}

func (f *Fetcher) mapColor(col uint8) uint8 {
	bgp := f.bgp.Get()
	// BGP is 0bAABBCCDD
	// where the nth pair of bits is the shade for the nth color.
	// E.g. BGP=0b10110001 and col=0x03 then the shade is 0b10:
	// (0b10110001 >> 6) & 0x00000011 = 0b00000010 & 0x00000011 = 0b10
	return bgp >> (col * 2) & 0x03
}
