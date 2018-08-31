package mqttgo

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPingrespMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PINGRESP << 4),
		0,
	}

	p := NewPingresp()
	n, err := p.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_PINGRESP, p.header.GetType(), "Error decoding message.")
}

func TestPingrespMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PINGRESP << 4),
		0,
	}

	p := NewPingresp()

	dst := make([]byte, 10)
	n, err := p.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// Decode, Encode, and Decode again
func TestPingrespDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PINGRESP << 4),
		0,
	}

	p := NewPingresp()
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
