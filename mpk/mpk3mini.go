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

	/*setupError := akai.AkaiSendSettingsWithDriver()
	if setupError != nil {
		panic(setupError)
	}*/

	retVal := &MPK3Mini{

		k1: NewRelativeKnob("k1", "Knpf 1", 1000, 70), // Only ASCII
		k2: NewRelativeKnob("k2", "OhneUmlauf 2", 1000, 71),
		k3: NewRelativeKnob("k3", "i3 panel width", 1000, 72),
		k4: NewRelativeKnob("k4", "KNOB 4", 1000, 73),

		k5: New8BitKnob("k5", "max", 74),
		k6: New8BitKnob("k6", "min", 75),
		k7: New8BitKnob("k7", "range", 76),
		k8: New8BitKnob("k8", "normal", 77),
	}

	retVal.k5.high = 64
	retVal.k6.low = 32
	retVal.k7.low = 32
	retVal.k7.high = 100

	retVal.initSysExProgramm()

	//akai.AkaiSendSettingsWithDriver()

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
	k1        *Knob
	k2        *Knob
	k3        *Knob
	k4        *Knob
	k5        *Knob
	k6        *Knob
	k7        *Knob
	k8        *Knob
	keyEvents []*MidiKey
	stopFunc  func()

	in drivers.In
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

		knob1: m.k1,
		knob2: m.k2,
		knob3: m.k3,
		knob4: m.k4,
		knob5: m.k5,
		knob6: m.k6,
		knob7: m.k7,
		knob8: m.k8,

		padBankA: AutopopulatePads(36),
		padBankB: AutopopulatePads(36 + 8),

		//aftertouch:                     aftertouchOff,
		aftertouch: aftertouchOff,
		//aftertouch:                     aftertouchOff,
		keybedSlashControlsMidiChannel: 1, //SEEMS Important
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

	data, err := setting.SysExStore(Program1)
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
