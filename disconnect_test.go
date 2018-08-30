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
	n, err := d.decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_DISCONNECT, d.Header.GetType(), "Error decoding message.")
}

func TestDisconnectMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_DISCONNECT << 4),
		0,
	}

	d := NewDisconnect()

	dst := make([]byte, 10)
	n, err := d.encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestDisconnectDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_DISCONNECT << 4),
		0,
	}

	d := NewDisconnect()
	n, err := d.decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")

	dst := make([]byte, 100)
	n2, err := d.encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n2, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n2], "Error decoding message.")

	n3, err := d.decode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n3, "Error decoding message.")
}
