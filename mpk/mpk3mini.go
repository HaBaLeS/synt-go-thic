package mpk

import (
	"fmt"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"log"
	"os"
	//_ "gitlab.com/gomidi/midi/v2/drivers/midicat"
	//_ "gitlab.com/gomidi/midi/v2/drivers/midicatdrv"
	//_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func NewMPK3Mini() (*MPK3Mini, error) {

	retVal := &MPK3Mini{
		knobs:       make(map[string]*MidiKnob),
		activeKnobs: make(map[uint8]*MidiKnob),
	}

	retVal.defineRelativeKnob("k1", "Knpf 1", 70) // Only ASCII
	retVal.defineRelativeKnob("k2", "OhneUmlauf 2", 71)
	retVal.defineRelativeKnob("k3", "i3 panel width", 72)
	retVal.defineRelativeKnob("k4", "KNOB 4", 73)

	retVal.defineAbsolutKnob("k5", "maxx", 74)
	retVal.defineAbsolutKnob("k6", "minn", 75)
	retVal.defineAbsolutKnob("k7", "rrange", 76)
	retVal.defineAbsolutKnob("k8", "normal", 77)

	retVal.MidiKnob("k5").Range(0, 64)
	retVal.MidiKnob("k6").Range(32, 127)
	retVal.MidiKnob("k7").Range(32, 50)

	retVal.initSysExProgramm()

	inPorts := midi.GetInPorts()
	fmt.Printf("Found Midi Device: %v", inPorts)
	in, err := midi.FindInPort("MPK mini 3")
	if err != nil {
		return nil, fmt.Errorf("can't find Midi Device")
	}
	retVal.in = in

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel, val uint8
		switch {
		case msg.GetSysEx(&bt):
			fmt.Printf("got sysex: % X\n", bt)
		case msg.GetNoteStart(&ch, &key, &vel):
			fmt.Printf("Press %v -> Key %d in Oct%d -> Power: %d\n", key, midi.Note(key).Base(), midi.Note(key).Octave(), int(vel))
			note := midi.Note(key)
			mk := &MidiKey{
				Idx:       int(key),
				Number:    int(note.Base()),
				Octave:    int(note.Octave()),
				Velocity:  int(vel),
				Name:      note.String(),
				EventType: EventPress,
			}
			retVal.keyEvents = append(retVal.keyEvents, mk)
		case msg.GetNoteEnd(&ch, &key):
			fmt.Printf("Release %v -> Key %d in Oct %d\n", key, midi.Note(key).Base(), midi.Note(key).Octave())
			note := midi.Note(key)
			mk := &MidiKey{
				Idx:       int(key),
				Number:    int(note.Base()),
				Octave:    int(note.Octave()),
				Velocity:  int(vel),
				Name:      note.String(),
				EventType: EventRelease,
			}
			retVal.keyEvents = append(retVal.keyEvents, mk)

		case msg.GetControlChange(&ch, &key, &val):
			fmt.Printf("controll msg: %d val: %d, channel %d\n", key, val, ch)
			knob := retVal.KnobByCC(key)
			if knob.appmode == KNOB_MODE_8BIT {
				knob.currentVal = int(val)
			} else if knob.appmode == KNOB_MODE_RELATIVE {
				if val >= 64 {
					knob.currentVal++
				} else {
					knob.currentVal--
				}
			}
			fmt.Printf("Knob %s value= %d", knob.displayName, knob.currentVal)
			//knob := retVal.knobMap[int(key)]
			//if knob != nil {
			//	knob.val = int(val)
			//}

		default:
			fmt.Printf("Other: %s\n", msg.String())
		}
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}

	retVal.keyEvents = make([]*MidiKey, 0)

	retVal.stopFunc = stop
	return retVal, nil
}

type MPK3Mini struct {
	activeKnobs map[uint8]*MidiKnob //current program knobs
	knobs       map[string]*MidiKnob
	keyEvents   []*MidiKey
	stopFunc    func()

	in drivers.In
}

func (m *MPK3Mini) defineAbsolutKnob(id, name string, ccid uint8) {
	m.knobs[id] = New8BitKnob(id, name, ccid)
	m.activeKnobs[ccid] = m.knobs[id]
}

func (m *MPK3Mini) defineRelativeKnob(id, name string, ccid uint8) {
	m.knobs[id] = NewRelativeKnob(id, name, ccid)
	m.activeKnobs[ccid] = m.knobs[id]
}

// MidiKnob fetches knobs by id this can be all defined knobs even in not active programs
func (m *MPK3Mini) MidiKnob(id string) *MidiKnob {
	return m.knobs[id]
}

// KnobByCC fetches the knobs by the incoming midi message
func (m *MPK3Mini) KnobByCC(ccid uint8) *MidiKnob {
	return m.activeKnobs[ccid]
}

func (m *MPK3Mini) initSysExProgramm() {

	outPorts := midi.GetOutPorts()
	fmt.Printf("Found OUT Midi Device: %v", outPorts)
	midiOut := outPorts[1]
	err := midiOut.Open()
	defer midiOut.Close()
	if err != nil {
		panic(err)
	}

	setting := &Settings{

		programName:    "gtfo",
		padMidiChannel: 10,

		knob1: m.MidiKnob("k1"),
		knob2: m.MidiKnob("k2"),
		knob3: m.MidiKnob("k3"),
		knob4: m.MidiKnob("k4"),
		knob5: m.MidiKnob("k5"),
		knob6: m.MidiKnob("k6"),
		knob7: m.MidiKnob("k7"),
		knob8: m.MidiKnob("k8"),

		padBankA: AutopopulatePads(36),
		padBankB: AutopopulatePads(36 + 8),

		//aftertouch:                     aftertouchOff,
		aftertouch:                     aftertouchOff, //aftertouch:                     aftertouchOff,
		keybedSlashControlsMidiChannel: 1,             //SEEMS Important
		joyHoriz: joystickAxisConfig{
			mode: joystickModeSingleCc, //up down same cc
			cc1:  69,
		},
		joyVert: joystickAxisConfig{
			mode: joystickModeDualCc, //up down 2 different CC
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

	data, err := setting.SysExStore(ProgramRAM)
	if err != nil {
		panic(err)
	}

	f, _ := os.Create("/tmp/my.hex")
	f.Write(data)
	f.Close()

	log.Printf("Sending SysEx\n")
	err = midiOut.Send(data)
	if err != nil {
		panic(err)
	}
	log.Printf("Sending Done\n")

}

type MidiEvent int

const (
	EventRelease MidiEvent = iota
	EventPress
)

type MidiKey struct {
	Name      string
	Octave    int
	Number    int
	Velocity  int
	EventType MidiEvent
	Idx       int
}

func (m *MPK3Mini) Stop() {
	m.stopFunc()
	midi.CloseDriver()
}

func (m *MPK3Mini) KnobVal(s string) int {
	/*for _, v := range m.knobMap {
		if v.name == s {
			return v.val
		}
	}*/
	return -1
}

func (m *MPK3Mini) MidiKeys() []*MidiKey {
	return m.keyEvents
}

func (m *MPK3Mini) ClearEvents() {
	m.keyEvents = m.keyEvents[0:0]
}
