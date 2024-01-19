package engine

import (
	"math"
	"math/rand"
)

type WaveForm int

const (
	SINE WaveForm = iota
	SQUARE
	SAW
	TRIANGLE
	NOISE
)
const (
	TWO_PI float64 = 2 * math.Pi
)

type Oscillator struct {
	phase      float64
	Frequency  float64
	samplerate float64
	vel        int
	velGain    int
	harmonics  int
}

func NewOscillator() *Oscillator {
	return &Oscillator{
		phase:      0.0,
		Frequency:  0.0,
		samplerate: 48000.0,
		harmonics:  1.0,
	}
}

func (o *Oscillator) AdvanceOscillatorSine() float64 {
	o.phase += TWO_PI * o.Frequency / o.samplerate
	//keep have between 0 and 2*PI to not rund in to float issues with large numbers
	for o.phase > TWO_PI {
		//See https://blog.demofox.org/2012/05/19/diy-synthesizer-chapter-2-common-wave-forms/ why this is a for instead of if
		o.phase -= TWO_PI
	}
	for o.phase < 0 {
		o.phase += TWO_PI
	}

	vg := math.Min(float64(o.vel+o.velGain), 127.0)
	return math.Sin(o.phase) / (127.0 / vg)
}

func (o *Oscillator) AdvanceOscilatorSquare() float64 {
	o.phase += o.Frequency / o.samplerate

	for o.phase > 1.0 {
		o.phase -= 1.0
	}

	for o.phase < 0.0 {
		o.phase += 1.0
	}

	val := 0.0
	if o.phase <= 0.5 {
		val = -1.0
	} else {
		val = 1.0
	}
	vg := math.Min(float64(o.vel+o.velGain), 127.0)
	return val / (127.0 / vg)
}

func (o *Oscillator) AdvanceOscilatorSaw() float64 {
	o.phase += o.Frequency / o.samplerate

	for o.phase > 1.0 {
		o.phase -= 1.0
	}

	for o.phase < 0.0 {
		o.phase += 1.0
	}

	vg := math.Min(float64(o.vel+o.velGain), 127.0)
	return ((o.phase * 2.0) - 1.0) / (127.0 / vg)
}

func (o *Oscillator) AdvanceOscilatorSawBandLimited() float64 {
	o.phase += o.Frequency / o.samplerate

	for o.phase > 1.0 {
		o.phase -= 1.0
	}

	for o.phase < 0.0 {
		o.phase += 1.0
	}

	var vg float64
	for i := 1; i <= o.harmonics; i++ {
		vg += math.Sin(o.phase*float64(i)) / float64(i)
	}

	//adjust the volume
	vg = vg * 2.0 / math.Pi

	return vg
	//vg := math.Min(float64(o.vel+o.velGain), 127.0)
	//return ((o.phase * 2.0) - 1.0) / (127.0 / vg)
}

func (o *Oscillator) AdvanceOscilatorTriangle() float64 {
	o.phase += o.Frequency / o.samplerate

	for o.phase > 1.0 {
		o.phase -= 1.0
	}

	for o.phase < 0.0 {
		o.phase += 1.0
	}

	val := 0.0
	if o.phase <= 0.5 {
		val = o.phase * 2
	} else {
		val = (1.0 - o.phase) * 2
	}
	vg := math.Min(float64(o.vel+o.velGain), 127.0)
	return ((val * 2.0) - 1.0) / (127.0 / vg)
}

func (o *Oscillator) AdvanceOscilatorNoise(lastVal float64) float64 {
	lastSeed := int(o.phase)
	o.phase += o.Frequency / o.samplerate
	seed := int(o.phase)

	for o.phase > 2.0 {
		o.phase -= 1.0
	}

	if seed != lastSeed {
		val := rand.Float64()
		val = (val * 2.0) - 1.0

		//uncomment the below to make it slightly more intense
		/*
			if(fValue < 0)
				fValue = -1.0f;
			else
				fValue = 1.0f;
		*/
		vg := math.Min(float64(o.vel+o.velGain), 127.0)
		return val / (127.0 / vg)
	} else {
		return lastVal
	}
}

func (o *Oscillator) Velocity(vel int) {
	o.vel = vel
}

func (o *Oscillator) VelGain(gain int) {
	o.velGain = gain
}

func (o *Oscillator) Harmonics(h int) {
	o.harmonics = h
}
