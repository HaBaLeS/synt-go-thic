package engine

type Mixer struct {
	MaxChannels int
	channels    []*Oscillator
	maxAmpl     float64 //num
}

func NewMixer(maxChan int) *Mixer {
	retVal := &Mixer{
		MaxChannels: maxChan,
		channels:    make([]*Oscillator, maxChan),
		maxAmpl:     30000.0 * 0.5,
	}
	for i := 0; i < maxChan; i++ {
		retVal.channels[i] = NewOscillator()
	}
	return retVal
}

func (m *Mixer) Read(p []byte) (n int, err error) {

	for i := 0; i < len(p)/2; i++ {
		val := 0.0
		for _, v := range m.channels {
			val += v.AdvanceOscillatorSine() / 4 //FIXME why is 4 the magic number where we do not clip anymore?
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
