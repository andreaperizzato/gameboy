package screen

import (
	"image/color"
	"log"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	pixelScale   = 4
	screenWidth  = 160
	screenHeight = 144
)

// Screen is an implementation of a Gameboy Classic display
// using OpenGL.
type Screen struct {
	enabled bool

	window   *pixelgl.Window
	picture  *pixel.PictureData
	row, col int
}

// New returns a new screen. You must call Start() to show it.
func New() *Screen {
	s := Screen{
		picture: pixel.MakePictureData(pixel.R(0, 0, screenWidth, screenHeight)),
	}
	return &s
}

// Start presents a window and shows it.
// Calling start is a blocking operation and retuns when the window is closed.
func (s *Screen) Start() {
	pixelgl.Run(s.run)
}

func (s *Screen) run() {
	// Create a new window.
	cfg := pixelgl.WindowConfig{
		Title:  "Gameboy",
		Bounds: pixel.R(0, 0, screenWidth*pixelScale, screenHeight*pixelScale),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}
	s.window = win

	// Set the window transform so that the content
	// fills the window and it is centered.
	xScale := s.window.Bounds().W() / screenWidth
	yScale := s.window.Bounds().H() / screenHeight
	scale := math.Min(yScale, xScale)
	s.window.SetMatrix(
		pixel.IM.
			Scaled(pixel.ZV, scale).
			Moved(win.Bounds().Size().Scaled(0.5)),
	)

	for !win.Closed() {
		spr := pixel.NewSprite(pixel.Picture(s.picture), pixel.R(0, 0, screenWidth, screenHeight))
		spr.Draw(s.window, pixel.IM)
		s.window.Update()
	}
}

// Write outputs a pixel (defined as a color number) to the display.
func (s *Screen) Write(color uint8) {
	if !s.enabled {
		return
	}
	s.picture.Pix[(screenHeight-1-s.row)*screenWidth+s.col] = getColor(color)
	s.col++
}

// HBlank is called whenever all pixels in a scanline have been output.
func (s *Screen) HBlank() {
	if !s.enabled {
		return
	}
	s.row++
	s.col = 0
}

// VBlank is called whenever a full frame has been output.
func (s *Screen) VBlank() {
	if !s.enabled {
		return
	}
	s.col, s.row = 0, 0
}

// Enable enables the screen.
func (s *Screen) Enable(e bool) {
	s.col, s.row = 0, 0
	s.enabled = e
}

// IsEnabled returns true when the screen is enabled.
func (s *Screen) IsEnabled() bool {
	return s.enabled
}

// This is the default pallet for the gameboy classic.
// Values come from https://en.wikipedia.org/wiki/Game_Boy.
var defaultPalette = [4]color.RGBA{
	{0x9B, 0xBC, 0x0F, 0xff}, // White
	{0x8B, 0xAC, 0x0f, 0xff}, // Light gray
	{0x30, 0x62, 0x30, 0xff}, // Dark gray
	{0x0F, 0x38, 0x0F, 0xff}, // Black
}

func getColor(c uint8) color.RGBA {
	return defaultPalette[c]
}
