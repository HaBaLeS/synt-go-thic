package engine

import (
	"math"
)

type WaveForm int

const (
	SINE WaveForm = iota
	RECT
)
const (
	TWO_PI float64 = 2 * math.Pi
)

type Oscillator struct {
	WaveForm   WaveForm
	phase      float64
	Frequency  float64
	samplerate float64
}

func NewOscillator() *Oscillator {
	return &Oscillator{
		WaveForm:   SINE,
		phase:      0.0,
		Frequency:  0.0,
		samplerate: 48000.0,
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
	return math.Sin(o.phase)
}
