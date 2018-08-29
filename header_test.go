package mqtt

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMessageHeaderFields(t *testing.T) {
	header := &Header{}

	header.SetRemainingLength(33)
	require.Equal(t, uint64(33), header.GetRemainingLength())

	err := header.SetRemainingLength(268435456)
	require.Error(t, err)

	err = header.SetType(TYPE_RESERVED)
	require.Error(t, err)

	err = header.SetType(TYPE_PUBREL)
	require.NoError(t, err)
	require.Equal(t, TYPE_PUBREL, header.GetType())

	err = header.SetFlag(TYPE_FLAG_PUBREL)
	require.NoError(t, err)
	require.Equal(t, 2, int(header.GetFlag()))
}

// Not enough bytes
func TestMessageHeaderDecode(t *testing.T) {
	buf := []byte{0x6f, 193, 2}
	header := &Header{}

	_, err := header.decode(buf)
	require.Error(t, err)
}

// Remaining length too big
func TestMessageHeaderDecode2(t *testing.T) {
	buf := []byte{0x62, 0xff, 0xff, 0xff, 0xff}
	header := &Header{}

	_, err := header.decode(buf)
	require.Error(t, err)
}

func TestMessageHeaderDecode3(t *testing.T) {
	buf := []byte{0x62, 0xff}
	header := &Header{}

	_, err := header.decode(buf)
	require.Error(t, err)
}

func TestMessageHeaderDecode4(t *testing.T) {
	buf := []byte{0x62, 0xff, 0xff, 0xff, 0x7f}
	header := &Header{
		TypeAndFlag: byte(6<<4 | 2),
		//mtype:      6,
		//flags:      2,
	}

	n, err := header.decode(buf)

	require.Error(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, MaxRemainingLength, int32(header.GetRemainingLength()))
}

func TestMessageHeaderDecode5(t *testing.T) {
	buf := []byte{0x62, 0xff, 0x7f}
	header := &Header{
		TypeAndFlag: byte(6<<4 | 2),
		//mtype:      6,
		//flags:      2,
	}

	n, err := header.decode(buf)
	require.Error(t, err)
	require.Equal(t, 3, n)
}

func TestMessageHeaderEncode1(t *testing.T) {
	header := &Header{}
	headerBytes := []byte{0x62, 193, 2}

	err := header.SetType(TYPE_PUBREL)
	require.NoError(t, err)

	err = header.SetRemainingLength(321)
	require.NoError(t, err)

	buf := make([]byte, 3)
	n, err := header.encode(buf)
	require.NoError(t, err)
	require.Equal(t, 3, n)
	require.Equal(t, headerBytes, buf)
}

func TestMessageHeaderEncode2(t *testing.T) {
	header := &Header{}

	err := header.SetType(TYPE_PUBREL)
	require.NoError(t, err)

	header.SetRemainingLength(268435456)

	buf := make([]byte, 5)
	_, err = header.encode(buf)

	require.Error(t, err)
}

func TestMessageHeaderEncode3(t *testing.T) {
	header := &Header{}
	headerBytes := []byte{0x62, 0xff, 0xff, 0xff, 0x7f}

	err := header.SetType(TYPE_PUBREL)

	require.NoError(t, err)

	err = header.SetRemainingLength(uint64(MaxRemainingLength))
	require.NoError(t, err)

	buf := make([]byte, 5)
	n, err := header.encode(buf)

	require.NoError(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, headerBytes, buf)
}

func TestMessageHeaderEncode4(t *testing.T) {
	header := &Header{
		TypeAndFlag: byte(TYPE_RESERVED2) << 4,
		//mtype:      6,
		//flags:      2,
	}

	buf := make([]byte, 5)
	_, err := header.encode(buf)
	require.Error(t, err)
}

/*
// This test is to ensure that an empty message is at least 2 bytes long
func TestMessageHeaderEncode5(t *testing.T) {
	msg := NewPingreqMessage()

	dst, n, err := msg.encode()
	if err != nil {
		t.Errorf("Error encoding PINGREQ message: %v", err)
	} else if n != 2 {
		t.Errorf("Incorrect result. Expecting length of 2 bytes, got %d.", dst.(*bytes.Buffer).Len())
	}
}
*/