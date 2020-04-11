package screen

import (
	"image/color"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const (
	pixelSize   = float64(4)
	pixelWidth  = 160
	pixelHeight = 144
)

type frame [144][161]uint8

type Screen struct {
	frames   []frame
	mux      sync.Mutex
	newFrame frame
	row, col int
	enabled  bool
}

func New() *Screen {
	s := Screen{}
	return &s
}

func (s *Screen) Start() {
	pixelgl.Run(s.run)
}

func (s *Screen) run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Gameboy",
		Bounds: pixel.R(0, 0, pixelSize*pixelWidth, pixelSize*pixelHeight),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		s.flush(win)
		win.Update()
	}
}

// Write outputs a pixel (defined as a color number) to the display.
func (s *Screen) Write(color uint8) {
	s.newFrame[s.row][s.col] = color
	s.col++
}

// HBlank is called whenever all pixels in a scanline have been output.
func (s *Screen) HBlank() {
	s.row++
	s.col = 0
}

// VBlank is called whenever a full frame has been output.
func (s *Screen) VBlank() {
	s.mux.Lock()
	defer s.mux.Unlock()
	c := frame{}
	for i, row := range s.newFrame {
		for j, v := range row {
			c[i][j] = v
		}
	}
	s.frames = append(s.frames, c)

	s.col = 0
	s.row = 0
}

func (s *Screen) Enable(e bool) {
	s.col = 0
	s.row = 0
	s.enabled = e
}

func (s *Screen) IsEnabled() bool {
	return s.enabled
}

func (s *Screen) flush(win *pixelgl.Window) {
	if len(s.frames) == 0 || !s.enabled {
		return
	}

	winHeight := win.Bounds().H()
	imd := imdraw.New(nil)

	frame := s.frames[0]
	for i, row := range frame {
		for j, color := range row {
			imd.Color = getColor(color)
			topLeft := pixel.V(float64(j)*pixelSize, winHeight-float64(i)*pixelSize)
			bottomRight := topLeft.Add(pixel.Vec{X: pixelSize, Y: -pixelSize})
			imd.Push(topLeft)
			imd.Push(bottomRight)
			imd.Rectangle(0)
		}
	}

	imd.Draw(win)
	s.frames = s.frames[1:]
}

var defaultPalette = [4]color.RGBA{
	{0x9B, 0xBC, 0x0F, 0xff}, // White
	{0x8B, 0xAC, 0x0f, 0xff}, // Light gray
	{0x30, 0x62, 0x30, 0xff}, // Dark gray
	{0x0F, 0x38, 0x0F, 0xff}, // Black
}

func getColor(c uint8) color.Color {
	return defaultPalette[c]
}
