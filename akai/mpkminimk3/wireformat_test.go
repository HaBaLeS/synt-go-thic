package mpkminimk3

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"github.com/function61/gokit/testing/assert"
)

func TestWireFormat(t *testing.T) {
	confBytes, err := ExampleSettings().SysExStore(Program8)
	assert.Ok(t, err)

	assert.Assert(t, len(confBytes) == 254)

	assert.EqualString(
		t,
		sha1Hex(confBytes),
		"e03da67c6e4a33feac4083246e8dad009ab5a35d")
}

func sha1Hex(input []byte) string {
	sha1Sum := sha1.Sum(input)
	return hex.EncodeToString(sha1Sum[:])
}
