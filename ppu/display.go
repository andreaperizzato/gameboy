package ppu

// Display output pixels on a screen.
type Display interface {
	// Write outputs a pixel (defined as a color number) to the display.
	Write(color uint8)
	// HBlank is called whenever all pixels in a scanline have been output.
	HBlank()
	// VBlank is called whenever a full frame has been output.
	VBlank()
	// Enable enables/disables the display.
	Enable(bool)
	// IsEnabled returns true when the display is enabled.
	IsEnabled() bool
}
