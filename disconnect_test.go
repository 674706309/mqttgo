package mqtt

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDisconnectMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_DISCONNECT << 4),
		0,
	}

	d := NewDisconnect()
	n, err := d.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_DISCONNECT, d.header.GetType(), "Error decoding message.")
}

func TestDisconnectMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_DISCONNECT << 4),
		0,
	}

	d := NewDisconnect()

	dst := make([]byte, 10)
	n, err := d.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// Decode, Encode, and Decode again
func TestDisconnectDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_DISCONNECT << 4),
		0,
	}

	d := NewDisconnect()
	n, err := d.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")

	dst := make([]byte, 100)
	n2, err := d.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n2, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n2], "Error decoding message.")

	n3, err := d.Decode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n3, "Error decoding message.")
}
