package mpk

import (
	"fmt"
	"github.com/HaBaLeS/synt-go-thic/akai/midiutil"
)

type Settings struct {
	programName string

	padMidiChannel uint8
	padBankA       PadBank
	padBankB       PadBank

	knob1 *Knob
	knob2 *Knob
	knob3 *Knob
	knob4 *Knob
	knob5 *Knob
	knob6 *Knob
	knob7 *Knob
	knob8 *Knob

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

type PadBank struct {
	pad1 Pad
	pad2 Pad
	pad3 Pad
	pad4 Pad
	pad5 Pad
	pad6 Pad
	pad7 Pad
	pad8 Pad
}

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

func (s *Settings) SysExStore(prog Program) ([]byte, error) {
	programNameBin := [16]byte{}
	copy(programNameBin[:], s.programName)

	if len(s.programName) > len(programNameBin) {
		return nil, fmt.Errorf(
			"programName parameter max length %d; got %d",
			len(programNameBin),
			len(s.programName))
	}

	// since "read program 1 msg" looks like:
	//     0x7f, 0x49, 0x66, 0x00, 0x01, 0x01
	//
	// - I suspect 0x7f, 0x49 is Akai's preamble (possibly for this model)
	// - 0x64 for writing program, 0x66 for reading program?
	conf := []byte{0x7f, 0x49, 0x64, 0x01, 0x76, byte(prog)}
	conf = append(conf, programNameBin[:]...)

	conf = append(conf, []byte{
		s.padMidiChannel - 1,
		byte(s.aftertouch),
		s.keybedSlashControlsMidiChannel - 1,
		byte(s.octave),
		byte(s.arpeggiatorStatus),
		byte(s.arpeggiatorMode),
		byte(s.arpeggiatorTimeDiv),
		byte(s.arpeggiatorClock),
		byte(s.arpeggiatorLatch),
		s.arpeggiatorSwing50plus,
		byte(s.arpeggiatorTempoTaps),
		0x00, // TODO: what is this?
		s.apreggiatorTempo,
		s.arpeggiatorOctave - 1,
		byte(s.joyHoriz.mode),
		s.joyHoriz.cc1,
		s.joyHoriz.cc2,
		byte(s.joyVert.mode),
		s.joyVert.cc1,
		s.joyVert.cc2,
		s.padBankA.pad1.note,
		s.padBankA.pad1.cc,
		s.padBankA.pad1.pc,
		s.padBankA.pad2.note,
		s.padBankA.pad2.cc,
		s.padBankA.pad2.pc,
		s.padBankA.pad3.note,
		s.padBankA.pad3.cc,
		s.padBankA.pad3.pc,
		s.padBankA.pad4.note,
		s.padBankA.pad4.cc,
		s.padBankA.pad4.pc,
		s.padBankA.pad5.note,
		s.padBankA.pad5.cc,
		s.padBankA.pad5.pc,
		s.padBankA.pad6.note,
		s.padBankA.pad6.cc,
		s.padBankA.pad6.pc,
		s.padBankA.pad7.note,
		s.padBankA.pad7.cc,
		s.padBankA.pad7.pc,
		s.padBankA.pad8.note,
		s.padBankA.pad8.cc,
		s.padBankA.pad8.pc,
		s.padBankB.pad1.note,
		s.padBankB.pad1.cc,
		s.padBankB.pad1.pc,
		s.padBankB.pad2.note,
		s.padBankB.pad2.cc,
		s.padBankB.pad2.pc,
		s.padBankB.pad3.note,
		s.padBankB.pad3.cc,
		s.padBankB.pad3.pc,
		s.padBankB.pad4.note,
		s.padBankB.pad4.cc,
		s.padBankB.pad4.pc,
		s.padBankB.pad5.note,
		s.padBankB.pad5.cc,
		s.padBankB.pad5.pc,
		s.padBankB.pad6.note,
		s.padBankB.pad6.cc,
		s.padBankB.pad6.pc,
		s.padBankB.pad7.note,
		s.padBankB.pad7.cc,
		s.padBankB.pad7.pc,
		s.padBankB.pad8.note,
		s.padBankB.pad8.cc,
		s.padBankB.pad8.pc,
	}...)

	for _, knob := range []*Knob{
		s.knob1, s.knob2, s.knob3, s.knob4,
		s.knob5, s.knob6, s.knob7, s.knob8,
	} {
		nameBin := [16]byte{}
		copy(nameBin[:], knob.displayName)

		if len(knob.displayName) > len(nameBin) {
			return nil, fmt.Errorf(
				"knob name max length %d; got %d",
				len(nameBin),
				len(knob.displayName))
		}

		// mode=abs|relative, CC, low, high, knob name (16 bytes)
		// total 20 bytes for each knob
		knobBytes := append([]byte{byte(knob.sysExMode), knob.ccId, knob.low, knob.high}, nameBin[:]...)

		conf = append(conf, knobBytes...)

	}

	// lol seems like an afterthought
	conf = append(conf, byte(s.transpose))

	return midiutil.ManufacturerAkai.SysEx(conf), nil
}

// /---------------- UNDERSTAND AND GET RIF OF
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

type joystickMode uint8

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

func AutopopulatePads(start uint8) PadBank {
	mkPad := func(n uint8) Pad {
		return Pad{n, n, n}
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
