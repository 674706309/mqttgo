package mqtt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubackMessageFields(t *testing.T) {
	msg := NewSuback()

	msg.header.SetPacketID(100)
	require.Equal(t, 100, int(msg.header.GetPacketID()), "Error setting packet ID.")

	msg.AddReturnCode(1)
	require.Equal(t, 1, len(msg.GetReturnCodes()), "Error adding return code.")

	err := msg.AddReturnCode(0x90)
	require.Error(t, err)
}

func TestSubackMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_SUBACK << 4),
		6,
		0,    // packet ID MSB (0)
		7,    // packet ID LSB (7)
		0,    // return code 1
		1,    // return code 2
		2,    // return code 3
		0x80, // return code 4
	}

	msg := NewSuback()
	n, err := msg.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_SUBACK, msg.header.GetType(), "Error decoding message.")
	require.Equal(t, 4, len(msg.GetReturnCodes()), "Error adding return code.")
}

// test with wrong return code
func TestSubackMessageDecode2(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_SUBACK << 4),
		6,
		0,    // packet ID MSB (0)
		7,    // packet ID LSB (7)
		0,    // return code 1
		1,    // return code 2
		2,    // return code 3
		0x81, // return code 4
	}

	msg := NewSuback()
	_, err := msg.Decode(msgBytes)

	require.Error(t, err)
}

func TestSubackMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_SUBACK << 4),
		6,
		0,    // packet ID MSB (0)
		7,    // packet ID LSB (7)
		0,    // return code 1
		1,    // return code 2
		2,    // return code 3
		0x80, // return code 4
	}

	msg := NewSuback()
	msg.header.SetPacketID(7)
	msg.AddReturnCode(0)
	msg.AddReturnCode(1)
	msg.AddReturnCode(2)
	msg.AddReturnCode(0x80)

	dst := make([]byte, 10)
	n, err := msg.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestSubackDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_SUBACK << 4),
		6,
		0,    // packet ID MSB (0)
		7,    // packet ID LSB (7)
		0,    // return code 1
		1,    // return code 2
		2,    // return code 3
		0x80, // return code 4
	}
	msg := NewSuback()
	n, err := msg.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")

	dst := make([]byte, 100)
	n2, err := msg.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n2, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n2], "Error decoding message.")

	n3, err := msg.Decode(dst[:n2])

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n3, "Error decoding message.")
}
