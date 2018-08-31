package mqttgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPubrecMessageFields(t *testing.T) {
	p := NewPubrec()

	p.SetPacketID(100)

	require.Equal(t, 100, int(p.GetPacketID()))
}

func TestPubrecMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBREC << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPubrec()
	n, err := p.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_PUBREC, p.GetType(), "Error decoding message.")
	require.Equal(t, 7, int(p.GetPacketID()), "Error decoding message.")
}

// test insufficient bytes
func TestPubrecMessageDecode2(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBREC << 4),
		2,
		7, // packet ID LSB (7)
	}

	p := NewPubrec()
	_, err := p.Decode(msgBytes)

	require.Error(t, err)
}

func TestPubrecMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBREC << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPubrec()
	p.SetPacketID(7)

	dst := make([]byte, 10)
	n, err := p.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestPubrecDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBREC << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPubrec()
	n, err := p.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")

	dst := make([]byte, 100)
	n2, err := p.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n2, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n2], "Error decoding message.")

	n3, err := p.Decode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n3, "Error decoding message.")
}
