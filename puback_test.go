package mqtt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPubackMessageFields(t *testing.T) {
	p := NewPuback()

	p.Header.SetPacketID(100)

	require.Equal(t, 100, int(p.Header.GetPacketID()))
}

func TestPubackMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPuback()
	n, err := p.decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_PUBACK, p.Header.GetType(), "Error decoding message.")
	require.Equal(t, 7, int(p.Header.GetPacketID()), "Error decoding message.")
}

// test insufficient bytes
func TestPubackMessageDecode2(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBACK << 4),
		2,
		7, // packet ID LSB (7)
	}

	p := NewPuback()
	_, err := p.decode(msgBytes)

	require.Error(t, err)
}

func TestPubackMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPuback()
	p.Header.SetPacketID(7)

	dst := make([]byte, 10)
	n, err := p.encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestPubackDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPuback()
	n, err := p.decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")

	dst := make([]byte, 100)
	n2, err := p.encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n2, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n2], "Error decoding message.")

	n3, err := p.decode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n3, "Error decoding message.")
}
