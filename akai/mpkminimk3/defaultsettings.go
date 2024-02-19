package mpkminimk3

// API not perfect yet
func CustomSettings(
	programName string,
	padMidiChannel uint8,
	padBankA PadBank,
	padBankB PadBank,
	keybedSlashControlsMidiChannel uint8,
	knob1 knobConfig,
	knob2 knobConfig,
	knob3 knobConfig,
	knob4 knobConfig,
	knob5 knobConfig,
	knob6 knobConfig,
	knob7 knobConfig,
	knob8 knobConfig,
) Settings {
	return Settings{
		programName: programName,

		padMidiChannel: padMidiChannel,
		padBankA:       padBankA,
		padBankB:       padBankB,

		knob1: knob1,
		knob2: knob2,
		knob3: knob3,
		knob4: knob4,
		knob5: knob5,
		knob6: knob6,
		knob7: knob7,
		knob8: knob8,

		aftertouch:                     aftertouchOff,
		keybedSlashControlsMidiChannel: keybedSlashControlsMidiChannel,
		joyHoriz: joystickAxisConfig{
			mode: joystickModeSingleCc,
			cc1:  69,
		},
		joyVert: joystickAxisConfig{
			mode: joystickModeDualCc,
			cc1:  70,
			cc2:  71,
		},
		octave:                 octave0,
		transpose:              transpose0,
		arpeggiatorStatus:      arpeggiatorDisabled,
		arpeggiatorMode:        arpeggiatorModeOrdr,
		arpeggiatorClock:       arpeggiatorClockInternal,
		arpeggiatorTempoTaps:   arpeggiatorTempoTaps2,
		arpeggiatorLatch:       arpeggiatorLatchOff,
		apreggiatorTempo:       120,
		arpeggiatorSwing50plus: 10, // = 60 in UI
		arpeggiatorOctave:      4,
		arpeggiatorTimeDiv:     arpeggiatorTimeDiv16T,
	}
}

func ExampleSettings() Settings {
	return CustomSettings(
		"Joonas test",
		10,
		AutopopulatePads(16),
		AutopopulatePads(17), // all bank B pad notes shifted by +1 relative to bank A
		1,
		KnobRelative("KNOB 1", 1),
		KnobRelative("KNOB 2", 2),
		KnobRelative("KNOB 3", 3),
		KnobRelative("KNOB 4", 4),
		KnobAbsolute0to127("KNOB 5", 5),
		KnobAbsolute0to127("KNOB 6", 6),
		KnobAbsolute0to127("KNOB 7", 7),
		KnobAbsolute0to127("KNOB 8", 8))
}
