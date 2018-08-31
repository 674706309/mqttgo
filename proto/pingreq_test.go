package mqttgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPingreqMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PINGREQ << 4),
		0,
	}

	p := NewPingreq()
	n, err := p.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_PINGREQ, p.header.GetType(), "Error decoding message.")
}

func TestPingreqMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PINGREQ << 4),
		0,
	}

	p := NewPingreq()

	dst := make([]byte, 10)
	n, err := p.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// Decode, Encode, and Decode again
func TestPingreqDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PINGREQ << 4),
		0,
	}

	p := NewPingreq()
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
