package apu

import (
	"time"

	"github.com/andreaperizzato/gameboy/memory"
	"github.com/faiface/beep"
)

// Pulse implements the pulse voices available on the Gameboy.
type Pulse struct {
	SampleRate beep.SampleRate

	Frequency memory.Register16
	Envelope  memory.Register
	Volume    memory.Register
	Duty      memory.Register
	Restart   memory.Register

	waveT time.Duration
	envT  time.Duration
	vol   int
}

func (p *Pulse) envelopePeriod() time.Duration {
	sw := p.Envelope.Get()
	return time.Duration(sw) * time.Second / time.Duration(64)
}

func (p *Pulse) wavePeriod() time.Duration {
	freq := 131072 / (2048 - int(p.Frequency.Get()))
	return time.Second / time.Duration(freq)
}

func (p *Pulse) highDuration() time.Duration {
	waveP := p.wavePeriod()
	duty := float64(0.125)
	if p.Duty.Get() > 0 {
		duty = 0.25 * float64(p.Duty.Get())
	}
	return time.Duration(float64(waveP) * duty)
}

// Stream generates audio samples.
func (p *Pulse) Stream(samples [][2]float64) (n int, ok bool) {
	envelopePeriod := p.envelopePeriod()
	period := p.wavePeriod()
	highDuration := p.highDuration()
	samplingTime := p.SampleRate.D(1)
	for i := range samples {
		v := float64(p.vol) / 15
		if p.waveT > highDuration {
			v = 0
		}
		samples[i][0] = v
		samples[i][1] = v

		// Since the wave is period, there is no need to
		// have waveT growing indefinitely and we can subtract the period
		// as soon as it gets bigger.
		// In fact, f(t) = f(t+T) where T is the period and f() is the waveform.
		p.waveT += samplingTime
		if p.waveT > period {
			p.waveT -= period
		}

		// Same as above for the evelope period.
		// Moreover, once every period, we need to decrease the volume by one,
		// if not already zero.
		p.envT += samplingTime
		if p.envT > envelopePeriod {
			p.envT -= envelopePeriod
			if p.vol > 0 {
				p.vol--
			}
		}

		// When the restart flag is set, we need to start the sound again.
		if p.Restart.Get() == 1 {
			p.envT = 0
			p.waveT = 0
			p.vol = int(p.Volume.Get())
			p.Restart.Set(0)
		}
	}
	return len(samples), true
}

// Err returns a streaming error which cannot occur.
func (p *Pulse) Err() error {
	return nil
}
