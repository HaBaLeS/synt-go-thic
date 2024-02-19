package mpk

type Pad struct {
	note uint8 // has max length, asserted by wire format
	cc   uint8
	pc   uint8
}
