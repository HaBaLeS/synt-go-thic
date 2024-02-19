// AKAI MPK mini mk3
package mpkminimk3

type joystickMode uint8

const (
	joystickModePitchBend joystickMode = 0
	joystickModeSingleCc  joystickMode = 1
	joystickModeDualCc    joystickMode = 2
)

type aftertouch uint8

const (
	aftertouchOff        aftertouch = 0
	aftertouchChannel    aftertouch = 1
	aftertouchPolyphonic aftertouch = 2
)

type Program uint8

const (
	ProgramRAM Program = 0
	Program1   Program = 1
	Program2   Program = 2
	Program3   Program = 3
	Program4   Program = 4
	Program5   Program = 5
	Program6   Program = 6
	Program7   Program = 7
	Program8   Program = 8
)

type arpeggiatorStatus byte

const (
	arpeggiatorDisabled arpeggiatorStatus = 0x0
	arpeggiatorEnabled  arpeggiatorStatus = 0b01111111 // TODO: why seven on-bits?
)

type arpeggiatorMode uint8

const (
	arpeggiatorModeUp   arpeggiatorMode = 0
	arpeggiatorModeDown arpeggiatorMode = 1
	arpeggiatorModeExcl arpeggiatorMode = 2
	arpeggiatorModeIncl arpeggiatorMode = 3
	arpeggiatorModeOrdr arpeggiatorMode = 4
	arpeggiatorModeRand arpeggiatorMode = 5
)

type arpeggiatorClock uint8

const (
	arpeggiatorClockInternal arpeggiatorClock = 0
	arpeggiatorClockExternal arpeggiatorClock = 1
)

type arpeggiatorTempoTaps uint8 // only three options available in UI

const (
	arpeggiatorTempoTaps2 arpeggiatorTempoTaps = 2
	arpeggiatorTempoTaps3 arpeggiatorTempoTaps = 3
	arpeggiatorTempoTaps4 arpeggiatorTempoTaps = 4
)

type arpeggiatorLatch uint8

const (
	arpeggiatorLatchOff arpeggiatorLatch = 0
	arpeggiatorLatchOn  arpeggiatorLatch = 1
)

type arpeggiatorTimeDiv uint8

const (
	arpeggiatorTimeDiv4   arpeggiatorTimeDiv = 0
	arpeggiatorTimeDiv4T  arpeggiatorTimeDiv = 1
	arpeggiatorTimeDiv8   arpeggiatorTimeDiv = 2
	arpeggiatorTimeDiv8T  arpeggiatorTimeDiv = 3
	arpeggiatorTimeDiv16  arpeggiatorTimeDiv = 4
	arpeggiatorTimeDiv16T arpeggiatorTimeDiv = 5
	arpeggiatorTimeDiv32  arpeggiatorTimeDiv = 6
	arpeggiatorTimeDiv32T arpeggiatorTimeDiv = 7
)

type octave uint8

const (
	octaveMinus4 octave = 0
	octaveMinus3 octave = 1
	octaveMinus2 octave = 2
	octaveMinus1 octave = 3
	octave0      octave = 4
	octave1      octave = 5
	octave2      octave = 6
	octave3      octave = 7
	octave4      octave = 8
)

type transpose uint8

const (
	transposeMinus12 transpose = 0
	transposeMinus11 transpose = 1
	transposeMinus10 transpose = 2
	transposeMinus9  transpose = 3
	transposeMinus8  transpose = 4
	transposeMinus7  transpose = 5
	transposeMinus6  transpose = 6
	transposeMinus5  transpose = 7
	transposeMinus4  transpose = 8
	transposeMinus3  transpose = 9
	transposeMinus2  transpose = 10
	transposeMinus1  transpose = 11
	transpose0       transpose = 12
	transpose1       transpose = 13
	transpose2       transpose = 14
	transpose3       transpose = 15
	transpose4       transpose = 16
	transpose5       transpose = 17
	transpose6       transpose = 18
	transpose7       transpose = 19
	transpose8       transpose = 20
	transpose9       transpose = 21
	transpose10      transpose = 22
	transpose11      transpose = 23
	transpose12      transpose = 24
)

type joystickAxisConfig struct {
	mode joystickMode
	cc1  uint8
	cc2  uint8
}

type padConfig struct {
	note uint8 // has max length, asserted by wire format
	cc   uint8
	pc   uint8
}

type PadBank struct {
	pad1 padConfig
	pad2 padConfig
	pad3 padConfig
	pad4 padConfig
	pad5 padConfig
	pad6 padConfig
	pad7 padConfig
	pad8 padConfig
}

type knobMode uint8

const (
	knobModeAbsolute knobMode = 0
	knobModeRelative knobMode = 1
)

type knobConfig struct {
	name string   // has max length, asserted by wire format
	mode knobMode // when relative, each left turn sends low value, and left turn high value (TODO: confirm. seems low/high are ignored?)
	cc   uint8
	low  uint8
	high uint8
}

type Settings struct {
	programName string

	padMidiChannel uint8
	padBankA       PadBank
	padBankB       PadBank

	knob1 knobConfig
	knob2 knobConfig
	knob3 knobConfig
	knob4 knobConfig
	knob5 knobConfig
	knob6 knobConfig
	knob7 knobConfig
	knob8 knobConfig

	aftertouch                     aftertouch
	keybedSlashControlsMidiChannel uint8
	joyHoriz                       joystickAxisConfig
	joyVert                        joystickAxisConfig

	octave    octave
	transpose transpose

	arpeggiatorStatus      arpeggiatorStatus
	arpeggiatorMode        arpeggiatorMode
	arpeggiatorOctave      uint8
	arpeggiatorSwing50plus uint8 // valid values [0,25]. in UI 50 = 0. 75 = 25. "swing has to do with how far a sequence deviates from the metronomic grid"
	arpeggiatorClock       arpeggiatorClock
	arpeggiatorLatch       arpeggiatorLatch
	apreggiatorTempo       uint8
	arpeggiatorTempoTaps   arpeggiatorTempoTaps
	arpeggiatorTimeDiv     arpeggiatorTimeDiv
}

// makes absolute knob with values ranging between (0, 127)
func KnobAbsolute0to127(name string, cc uint8) knobConfig {
	return mkKnob(name, knobModeAbsolute, cc, 0, 127)
}

// sends decrease = (1, 2, more?), increase = (127, 126, less?) depending on decrease/increase
// but b/c of this weirdness I recommend the application uses RelativeKnobWasIncrease() to detect direction
func KnobRelative(name string, cc uint8) knobConfig {
	return knobConfig{name, knobModeRelative, cc, 0, 0}
}

func RelativeKnobWasIncrease(val uint8) bool {
	// for some odd reason decreases are communicated with high value, increases with low value..
	// also sometimes the low value is 1,2 and sometimes the high value is 126.127..
	return val < 64
}

func mkKnob(name string, mode knobMode, cc uint8, low uint8, high uint8) knobConfig {
	return knobConfig{name, mode, cc, low, high}
}

func AutopopulatePads(start uint8) PadBank {
	mkPad := func(n uint8) padConfig {
		return padConfig{n, n, n}
	}

	return PadBank{
		pad1: mkPad(start + 0),
		pad2: mkPad(start + 1),
		pad3: mkPad(start + 2),
		pad4: mkPad(start + 3),
		pad5: mkPad(start + 4),
		pad6: mkPad(start + 5),
		pad7: mkPad(start + 6),
		pad8: mkPad(start + 7),
	}
}
