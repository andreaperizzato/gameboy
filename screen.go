package gameboy

// import (
// 	"crypto/rand"
// 	"image/color"

// 	"github.com/faiface/pixel"
// 	"github.com/faiface/pixel/imdraw"
// 	"github.com/faiface/pixel/pixelgl"
// )

// const (
// 	pixelSize = float64(4)
// 	pixelCount = 128
// )

// func run() {
// 	winSize := pixelCount * pixelSize
// 	cfg := pixelgl.WindowConfig{
// 		Title:  "Pixel Rocks!",
// 		Bounds: pixel.R(0, 0, winSize, winSize),
// 	}
// 	win, err := pixelgl.NewWindow(cfg)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for !win.Closed() {
// 		drawRandomCrap(win)
// 		win.Update()
// 	}
// }

// func drawRandomCrap(win *pixelgl.Window) {
// 	winSize := win.Bounds().W()
// 	imd := imdraw.New(nil)
// 	for r := 0; r < pixelCount; r++ {
// 		for c := 0; c < pixelCount; c++ {
// 			b := make([]byte, 1)
// 			rand.Read(b)
// 			imd.Color = color.RGBA{R: b[0], G: 12, B: 123, A: 255}
// 			topLeft := pixel.V(winSize - float64(r + 1) * pixelSize, float64(c) * pixelSize)
// 			bottomRight := topLeft.Add(pixel.Vec{X: pixelSize, Y: pixelSize})
// 			imd.Push(topLeft)
// 			imd.Push(bottomRight)
// 			imd.Rectangle(0)
// 		}
// 	}
// 	imd.Draw(win)
// }

// func main() {
// 	pixelgl.Run(run)
// }
