package mqtt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnsubackMessageFields(t *testing.T) {
	msg := NewUnSuback()

	msg.Header.SetPacketID(100)

	require.Equal(t, 100, int(msg.Header.GetPacketID()))
}

func TestUnsubackMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_UNSUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	msg := NewUnSuback()
	n, err := msg.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_UNSUBACK, msg.Header.GetType(), "Error decoding message.")
	require.Equal(t, 7, int(msg.Header.GetPacketID()), "Error decoding message.")
}

// test insufficient bytes
func TestUnsubackMessageDecode2(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_UNSUBACK << 4),
		2,
		7, // packet ID LSB (7)
	}

	msg := NewUnSuback()
	_, err := msg.Decode(msgBytes)

	require.Error(t, err)
}

func TestUnsubackMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_UNSUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	msg := NewUnSuback()
	msg.Header.SetPacketID(7)

	dst := make([]byte, 10)
	n, err := msg.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestUnsubackDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_UNSUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	msg := NewUnSuback()
	n, err := msg.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")

	dst := make([]byte, 100)
	n2, err := msg.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n2, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n2], "Error decoding message.")

	n3, err := msg.Decode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n3, "Error decoding message.")
}
