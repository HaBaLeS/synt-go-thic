package mpk

import (
	"fmt"
	"gitlab.com/gomidi/midi/v2"
)

func NewMPK3Mini() *MPK3Mini {
	inPorts := midi.GetInPorts()
	fmt.Printf("Found Midi Device: %v", inPorts)
	in, err := midi.FindInPort("MPK mini 3")
	if err != nil {
		panic(fmt.Errorf("can't find Midi Device"))

	}

	retVal := &MPK3Mini{}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel, val uint8
		switch {
		case msg.GetSysEx(&bt):
			fmt.Printf("got sysex: % X\n", bt)
		case msg.GetNoteStart(&ch, &key, &vel):
			//fmt.Printf("starting note %s on channel %v with velocity %v\n", midi.Note(key), ch, vel)
			fmt.Printf("Key %d in Oct %d\n", midi.Note(key).Base(), midi.Note(key).Octave())
			note := midi.Note(key)
			mk := &MidiKey{
				Number:   int(note.Base()),
				Octave:   int(note.Octave()),
				Velocity: int(vel),
				Name:     note.String(),
			}
			retVal.keyMap[mk.Name] = mk
		case msg.GetNoteEnd(&ch, &key):
			note := midi.Note(key)
			delete(retVal.keyMap, note.String())
			//fmt.Printf("ending note %s on channel %v\n", midi.Note(key), ch)

		case msg.GetControlChange(&ch, &key, &val):
			//fmt.Printf("controll msg: %d val: %d, channel %d\n", key, val, ch)
			knob := retVal.knobMap[int(key)]
			if knob != nil {
				knob.val = int(val)
			}

		default:
			fmt.Printf("Other: %s\n", msg.String())
		}
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil
	}

	retVal.knobMap = make(map[int]*Knob)
	retVal.knobMap[70] = &Knob{"K1", 70, 0}
	retVal.knobMap[71] = &Knob{"K2", 71, 0}
	retVal.knobMap[72] = &Knob{"K3", 72, 0}
	retVal.knobMap[73] = &Knob{"K4", 73, 0}
	retVal.knobMap[74] = &Knob{"K5", 74, 0}
	retVal.knobMap[75] = &Knob{"K6", 75, 0}
	retVal.knobMap[76] = &Knob{"K7", 76, 0}
	retVal.knobMap[77] = &Knob{"K8", 77, 0}

	retVal.keyMap = make(map[string]*MidiKey)

	retVal.stopFunc = stop
	return retVal
}

type MPK3Mini struct {
	knobMap  map[int]*Knob
	keyMap   map[string]*MidiKey
	stopFunc func()
}

type Knob struct {
	name string
	ccId int
	val  int
}

type MidiKey struct {
	Name     string
	Octave   int
	Number   int
	Velocity int
}

type Pad struct {
}

func (m *MPK3Mini) Stop() {
	m.stopFunc()
	midi.CloseDriver()
}

func (m *MPK3Mini) KnobVal(s string) int {
	for _, v := range m.knobMap {
		if v.name == s {
			return v.val
		}
	}
	return -1
}

func (m *MPK3Mini) MidiKeys() []*MidiKey {
	//fixme check maps package https://stackoverflow.com/questions/21362950/getting-a-slice-of-keys-from-a-map
	keys := make([]*MidiKey, len(m.keyMap))
	i := 0
	for _, k := range m.keyMap {
		keys[i] = k
		i++
	}
	return keys
}
