// MIDI utilities
package midiutil

// https://www.midi.org/specifications-old/item/manufacturer-id-numbers
var (
	ManufacturerAkai = Manufacturer{[]byte{0x47}}
)

type Manufacturer struct {
	manufacturer []byte // 1-3 bytes
}

// http://midi.teragonaudio.com/tech/midispec/sysex.htm
func (s Manufacturer) SysEx(payload []byte) []byte {
	// SysEx messages look like 0xF0 <manufacturer ID> <payload> 0xF7

	msg := append([]byte{0xF0}, s.manufacturer...)
	msg = append(msg, payload...)
	msg = append(msg, 0xF7)
	return msg
}
