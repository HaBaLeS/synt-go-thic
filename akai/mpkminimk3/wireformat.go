package mpkminimk3

import (
	"fmt"
	"github.com/HaBaLeS/synt-go-thic/akai/midiutil"
)

// WARN: NOT tested if read from RAM works
func requestConfigMsg(prog Program) []byte {
	/*
		7f 49 66 00 01 01  <- read program 1
		7f 49 66 00 01 08  <- read program 8
	*/
	return midiutil.ManufacturerAkai.SysEx([]byte{0x7f, 0x49, 0x66, 0x00, 0x01, byte(prog)})
}

func (a Settings) SysExStore(prog Program) ([]byte, error) {
	programNameBin := [16]byte{}
	copy(programNameBin[:], a.programName)

	if len(a.programName) > len(programNameBin) {
		return nil, fmt.Errorf(
			"programName parameter max length %d; got %d",
			len(programNameBin),
			len(a.programName))
	}

	// since "read program 1 msg" looks like:
	//     0x7f, 0x49, 0x66, 0x00, 0x01, 0x01
	//
	// - I suspect 0x7f, 0x49 is Akai's preamble (possibly for this model)
	// - 0x64 for writing program, 0x66 for reading program?
	conf := []byte{0x7f, 0x49, 0x64, 0x01, 0x76, byte(prog)}
	conf = append(conf, programNameBin[:]...)

	conf = append(conf, []byte{
		a.padMidiChannel - 1,
		byte(a.aftertouch),
		a.keybedSlashControlsMidiChannel - 1,
		byte(a.octave),
		byte(a.arpeggiatorStatus),
		byte(a.arpeggiatorMode),
		byte(a.arpeggiatorTimeDiv),
		byte(a.arpeggiatorClock),
		byte(a.arpeggiatorLatch),
		a.arpeggiatorSwing50plus,
		byte(a.arpeggiatorTempoTaps),
		0x00, // TODO: what is this?
		a.apreggiatorTempo,
		a.arpeggiatorOctave - 1,
		byte(a.joyHoriz.mode),
		a.joyHoriz.cc1,
		a.joyHoriz.cc2,
		byte(a.joyVert.mode),
		a.joyVert.cc1,
		a.joyVert.cc2,
		a.padBankA.pad1.note,
		a.padBankA.pad1.cc,
		a.padBankA.pad1.pc,
		a.padBankA.pad2.note,
		a.padBankA.pad2.cc,
		a.padBankA.pad2.pc,
		a.padBankA.pad3.note,
		a.padBankA.pad3.cc,
		a.padBankA.pad3.pc,
		a.padBankA.pad4.note,
		a.padBankA.pad4.cc,
		a.padBankA.pad4.pc,
		a.padBankA.pad5.note,
		a.padBankA.pad5.cc,
		a.padBankA.pad5.pc,
		a.padBankA.pad6.note,
		a.padBankA.pad6.cc,
		a.padBankA.pad6.pc,
		a.padBankA.pad7.note,
		a.padBankA.pad7.cc,
		a.padBankA.pad7.pc,
		a.padBankA.pad8.note,
		a.padBankA.pad8.cc,
		a.padBankA.pad8.pc,
		a.padBankB.pad1.note,
		a.padBankB.pad1.cc,
		a.padBankB.pad1.pc,
		a.padBankB.pad2.note,
		a.padBankB.pad2.cc,
		a.padBankB.pad2.pc,
		a.padBankB.pad3.note,
		a.padBankB.pad3.cc,
		a.padBankB.pad3.pc,
		a.padBankB.pad4.note,
		a.padBankB.pad4.cc,
		a.padBankB.pad4.pc,
		a.padBankB.pad5.note,
		a.padBankB.pad5.cc,
		a.padBankB.pad5.pc,
		a.padBankB.pad6.note,
		a.padBankB.pad6.cc,
		a.padBankB.pad6.pc,
		a.padBankB.pad7.note,
		a.padBankB.pad7.cc,
		a.padBankB.pad7.pc,
		a.padBankB.pad8.note,
		a.padBankB.pad8.cc,
		a.padBankB.pad8.pc,
	}...)

	for _, knob := range []knobConfig{
		a.knob1, a.knob2, a.knob3, a.knob4,
		a.knob5, a.knob6, a.knob7, a.knob8,
	} {
		nameBin := [16]byte{}
		copy(nameBin[:], knob.name)

		if len(knob.name) > len(nameBin) {
			return nil, fmt.Errorf(
				"knob name max length %d; got %d",
				len(nameBin),
				len(knob.name))
		}

		// mode=abs|relative, CC, low, high, knob name (16 bytes)
		// total 20 bytes for each knob
		knobBytes := append([]byte{byte(knob.mode), knob.cc, knob.low, knob.high}, nameBin[:]...)

		conf = append(conf, knobBytes...)

	}

	// lol seems like an afterthought
	conf = append(conf, byte(a.transpose))

	return midiutil.ManufacturerAkai.SysEx(conf), nil
}
