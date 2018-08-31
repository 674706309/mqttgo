package mqtt

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublishMessageHeaderFields(t *testing.T) {
	p := NewPublish()
	p.header.SetFlag(11)
	//p.mtypeflags[0] |= 11

	require.True(t, p.GetDup(), "Incorrect DUP flag.")
	require.True(t, p.GetRetain(), "Incorrect RETAIN flag.")
	require.Equal(t, 1, int(p.GetQoS()), "Incorrect QoS.")

	p.SetDup(false)

	require.False(t, p.GetDup(), "Incorrect DUP flag.")

	p.SetRetain(false)

	require.False(t, p.GetRetain(), "Incorrect RETAIN flag.")

	err := p.SetQoS(2)

	require.NoError(t, err, "Error setting QoS.")
	require.Equal(t, 2, int(p.GetQoS()), "Incorrect QoS.")

	err = p.SetQoS(3)

	require.Error(t, err)

	err = p.SetQoS(0)

	require.NoError(t, err, "Error setting QoS.")
	require.Equal(t, 0, int(p.GetQoS()), "Incorrect QoS.")

	p.SetDup(true)

	require.True(t, p.GetDup(), "Incorrect DUP flag.")

	p.SetRetain(true)

	require.True(t, p.GetRetain(), "Incorrect RETAIN flag.")
}

func TestPublishMessageFields(t *testing.T) {
	p := NewPublish()

	p.SetTopicName([]byte("coolstuff"))

	require.Equal(t, "coolstuff", string(p.GetTopicName()), "Error setting message topic.")

	err := p.SetTopicName([]byte("coolstuff/#"))

	require.Error(t, err)

	p.header.SetPacketID(100)

	require.Equal(t, 100, int(p.header.GetPacketID()), "Error setting acket ID.")

	p.SetPayload([]byte("this is a payload to be sent"))

	require.Equal(t, []byte("this is a payload to be sent"), p.GetPayload(), "Error setting payload.")
}

func TestPublishMessageDecode1(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBLISH<<4) | 2,
		23,
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
	}

	p := NewPublish()
	n, err := p.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, 7, int(p.header.GetPacketID()), "Error decoding message.")
	require.Equal(t, "surgemq", string(p.GetTopicName()), "Error deocding topic name.")
	require.Equal(t, []byte{'s', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e'}, p.GetPayload(), "Error deocding payload.")
}

// test insufficient bytes
func TestPublishMessageDecode2(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBLISH<<4) | 2,
		26,
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
	}

	p := NewPublish()
	_, err := p.Decode(msgBytes)

	require.Error(t, err)
}

// test qos = 0 and no client id
func TestPublishMessageDecode3(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBLISH << 4),
		21,
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
	}

	p := NewPublish()
	_, err := p.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
}

func TestPublishMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBLISH<<4) | 2,
		23,
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
	}

	p := NewPublish()
	p.SetTopicName([]byte("surgemq"))
	p.SetQoS(1)
	p.header.SetPacketID(7)
	p.SetPayload([]byte{'s', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e'})

	dst := make([]byte, 100)
	n, err := p.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test empty topic name
func TestPublishMessageEncode2(t *testing.T) {
	p := NewPublish()
	p.SetTopicName([]byte(""))
	p.header.SetPacketID(7)
	p.SetPayload([]byte{'s', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e'})

	dst := make([]byte, 100)
	_, err := p.Encode(dst)
	require.Error(t, err)
}

// test encoding qos = 0 and no packet id
func TestPublishMessageEncode3(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBLISH << 4),
		21,
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
	}

	p := NewPublish()
	p.SetTopicName([]byte("surgemq"))
	p.SetQoS(0)
	p.SetPayload([]byte{'s', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e'})

	dst := make([]byte, 100)
	n, err := p.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test large message
func TestPublishMessageEncode4(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBLISH << 4),
		137,
		8,
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
	}

	payload := make([]byte, 1024)
	msgBytes = append(msgBytes, payload...)

	p := NewPublish()
	p.SetTopicName([]byte("surgemq"))
	p.SetQoS(0)
	p.SetPayload(payload)

	dst := make([]byte, 1100)
	n, err := p.Encode(dst)
	dd := make([]byte, 2)
	binary.PutUvarint(dd, uint64(p.header.GetRemainingLength()))

	require.Equal(t, len(msgBytes), p.Length())

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test from github issue #2, @mrdg
func TestPublishDecodeEncodeEquiv2(t *testing.T) {
	msgBytes := []byte{50, 18, 0, 9, 103, 114, 101, 101, 116, 105, 110, 103, 115, 0, 1, 72, 101, 108, 108, 111}

	p := NewPublish()
	n, err := p.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")

	dst := make([]byte, 100)
	n2, err := p.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n2, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// Decode, Encode, and Decode again
func TestPublishDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_PUBLISH<<4) | 2,
		23,
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
	}

	p := NewPublish()

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
