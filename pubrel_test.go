package mqtt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPubrelMessageFields(t *testing.T) {
	p := NewPubrel()

	p.SetPacketID(100)

	require.Equal(t, 100, int(p.GetPacketID()))
}

func TestPubrelMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBREL<<4) | 2,
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPubrel()
	n, err := p.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_PUBREL, p.GetType(), "Error decoding message.")
	require.Equal(t, 7, int(p.GetPacketID()), "Error decoding message.")
}

// test insufficient bytes
func TestPubrelMessageDecode2(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBREL<<4) | 2,
		2,
		7, // packet ID LSB (7)
	}

	p := NewPubrel()
	_, err := p.Decode(msgBytes)

	require.Error(t, err)
}

func TestPubrelMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBREL<<4) | 2,
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPubrel()
	p.SetPacketID(7)

	dst := make([]byte, 10)
	n, err := p.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestPubrelDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBREL<<4) | 2,
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPubrel()
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
