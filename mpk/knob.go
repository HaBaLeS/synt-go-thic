package mpk

import (
	"math"
)

type KnobMode int

const (
	KNOB_MODE_RELATIVE KnobMode = iota
	KNOB_MODE_LIMITED
	KNOB_MODE_8BIT
)

type Knob struct {
	id          string
	displayName string

	appmode   KnobMode
	sysExMode int

	minval int
	maxval int

	low, high uint8

	ccId       uint8
	currentVal int
}

func NewRelativeKnob(id, name string, startVal int, ccId uint8) *Knob {
	retVal := &Knob{
		id:          id,
		displayName: name,
		appmode:     KNOB_MODE_RELATIVE,
		sysExMode:   1,
		minval:      0,
		maxval:      math.MaxInt, //LOW Value not max of the know-> UIk
		ccId:        ccId,
		currentVal:  startVal,
		low:         0,
		high:        0,
	}
	return retVal
}

func New8BitKnob(id, name string, ccId uint8) *Knob {
	return &Knob{
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

func NewLimitedKnob(id, name string, ccId uint8) *Knob {
	return &Knob{
		id:          id,
		displayName: name,
		ccId:        ccId,
		appmode:     KNOB_MODE_LIMITED,
		sysExMode:   0,
		minval:      0,
		maxval:      127,
		currentVal:  0, //undefined .. is there a way to read the know value without touching it?
		low:         0,
		high:        127,
	}
}
