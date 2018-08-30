package mqtt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPubcompMessageFields(t *testing.T) {
	p := NewPubcomp()

	p.Header.SetPacketID(100)

	require.Equal(t, 100, int(p.Header.GetPacketID()))
}

func TestPubcompMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBCOMP << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPubcomp()
	n, err := p.decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_PUBCOMP, p.Header.GetType(), "Error decoding message.")
	require.Equal(t, 7, int(p.Header.GetPacketID()), "Error decoding message.")
}

// test insufficient bytes
func TestPubcompMessageDecode2(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBCOMP << 4),
		2,
		7, // packet ID LSB (7)
	}

	p := NewPubcomp()
	_, err := p.decode(msgBytes)

	require.Error(t, err)
}

func TestPubcompMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBCOMP << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPubcomp()
	p.Header.SetPacketID(7)

	dst := make([]byte, 10)
	n, err := p.encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestPubcompDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBCOMP << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	p := NewPubcomp()
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
