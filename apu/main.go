package main

import (
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func SinWave(sr beep.SampleRate) beep.Streamer {
	ts := float64(sr.D(1)) / float64(time.Second)
	freq := 600
	ns := int(sr) / freq
	buf := make([]float64, ns)
	for i := range buf {
		buf[i] = math.Sin(2 * math.Pi * float64(freq) * ts * float64(i))
	}
	x := 0
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			v := buf[x]
			samples[i][0] = v
			samples[i][1] = v
			x++
			if x == len(buf) {
				x = 0
			}
		}
		return len(samples), true
	})
}

func main() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(SinWave(sr))
	select {}
}
