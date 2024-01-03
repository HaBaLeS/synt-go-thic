package engine

import (
	"github.com/HaBaLeS/synt-go-thic/mpk"
	"math"
	"time"
)

type Mixer struct {
	MaxChannels  int
	channels     []*Oscillator
	maxAmpl      float64 //num
	midi         *mpk.MPK3Mini
	WaveForm     WaveForm
	WaveFormName string
}

func NewMixer(maxChan int, midi *mpk.MPK3Mini) *Mixer {
	retVal := &Mixer{
		MaxChannels: maxChan,
		channels:    make([]*Oscillator, maxChan),
		maxAmpl:     float64(math.MaxInt16) / 2.0,
		midi:        midi,
	}
	for i := 0; i < maxChan; i++ {
		retVal.channels[i] = NewOscillator()
	}

	go retVal.AsynMidiMonitor()

	return retVal
}

func (m *Mixer) Read(p []byte) (n int, err error) {
	noiseLastVal := 0.0
	for i := 0; i < len(p)/2; i++ {
		val := 0.0
		for _, v := range m.channels {
			switch m.WaveForm {
			case SINE:
				val += v.AdvanceOscillatorSine() / 4 //FIXME why is 4 the magic number where we do not clip anymore?
			case SQUARE:
				val += v.AdvanceOscilatorSquare() / 4
			case SAW:
				val += v.AdvanceOscilatorSaw() / 4
			case TRIANGLE:
				val += v.AdvanceOscilatorTriangle() / 4
			case NOISE:
				noiseLastVal = v.AdvanceOscilatorNoise(noiseLastVal)
				val += noiseLastVal / 4
			}
		}
		ampl := int(val * m.maxAmpl)
		p[2*i] = byte(ampl)
		p[2*i+1] = byte(ampl >> 8)
	}
	return len(p), nil

}

func (m *Mixer) Channel(i int) *Oscillator {
	return m.channels[i]
}

func (m *Mixer) AsynMidiMonitor() {
	for true {
		time.Sleep(100 * time.Millisecond)
		velGain := m.midi.KnobVal("K1")
		for _, v := range m.channels {
			v.VelGain(velGain)
		}

		//WaveForm Selector
		m.WaveForm = WaveForm(m.midi.KnobVal("K2") / 25)
		switch m.WaveForm {
		case SQUARE:
			m.WaveFormName = "Square Wave Oscillator"
		case SINE:
			m.WaveFormName = "Sine Wave Oscillator"
		case SAW:
			m.WaveFormName = "Saw Wave Oscillator"
		case TRIANGLE:
			m.WaveFormName = "Triangle Wave Oscillator"
		case NOISE:
			m.WaveFormName = "Noise Oscillator"
		}

	}
}
