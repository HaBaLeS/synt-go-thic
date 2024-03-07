package mpk

import (
	"math"
)

type KnobMode int

const (
	KNOB_MODE_RELATIVE KnobMode = iota
	KNOB_MODE_8BIT
)

type MidiKnob struct {
	id          string
	displayName string

	appmode   KnobMode
	sysExMode int

	minval int
	maxval int

	low, high uint8

	ccId       uint8
	currentVal float64
}

func (k *MidiKnob) Range(min, max uint8) {
	//TODO allow setting only before settings are applied via Sysex
	k.low = min
	k.high = max
}

func (k *MidiKnob) Percent() float64 {
	rng := k.high - k.low
	perc := 100 / float64(rng) * k.currentVal

	return perc
}

func (k *MidiKnob) Value() int {
	return int(k.currentVal)
}

func NewRelativeKnob(id, name string, ccId uint8) *MidiKnob {
	retVal := &MidiKnob{
		id:          id,
		displayName: name,
		appmode:     KNOB_MODE_RELATIVE,
		sysExMode:   1,
		minval:      0,
		maxval:      math.MaxInt, //LOW Value not max of the know-> UIk
		ccId:        ccId,
		currentVal:  0,
		low:         0,
		high:        0,
	}
	return retVal
}

func New8BitKnob(id, name string, ccId uint8) *MidiKnob {
	return &MidiKnob{
		id:          id,
		displayName: name,
		ccId:        ccId,
		appmode:     KNOB_MODE_8BIT,
		sysExMode:   0,
		minval:      0,
		maxval:      127,
		currentVal:  0, //undefined .. is there a way to read the know value without touching it?
		low:         0,
		high:        127,
	}
}
