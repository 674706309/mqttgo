package mqtt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubscribeMessageFields(t *testing.T) {
	msg := NewSubscribe()

	msg.Header.SetPacketID(100)
	require.Equal(t, 100, int(msg.Header.GetPacketID()), "Error setting packet ID.")

	msg.AddTopic([]byte("/a/b/#/c"), 1)
	require.Equal(t, 1, len(msg.GetTopicFilter()), "Error adding topic.")
	_, found := msg.TopicExists([]byte("a/b"))
	require.False(t, found, "Topic should not exist.")
	_, found = msg.TopicExists([]byte("a/b"))
	msg.RemoveTopic([]byte("/a/b/#/c"))
	require.False(t, found, "Topic should not exist.")
}

func TestSubscribeMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_SUBSCRIBE<<4) | 2,
		36,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // QoS
		0, // topic name MSB (0)
		8, // topic name LSB (8)
		'/', 'a', '/', 'b', '/', '#', '/', 'c',
		1,  // QoS
		0,  // topic name MSB (0)
		10, // topic name LSB (10)
		'/', 'a', '/', 'b', '/', '#', '/', 'c', 'd', 'd',
		2, // QoS
	}

	msg := NewSubscribe()
	n, err := msg.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, TYPE_SUBSCRIBE, msg.Header.GetType(), "Error decoding message.")
	require.Equal(t, 3, len(msg.GetTopicFilter()), "Error decoding topics.")
	_, found := msg.TopicExists([]byte("a/b"))
	require.True(t, found, "Topic 'surgemq' should exist.")
	//require.Equal(t, 0, int(msg.TopicQos([]byte("surgemq"))), "Incorrect topic qos.")
	_, found = msg.TopicExists([]byte("a/b"))
	require.True(t, found, "Topic '/a/b/#/c' should exist.")
	//require.Equal(t, 1, int(msg.TopicQos([]byte("/a/b/#/c"))), "Incorrect topic qos.")
	_, found = msg.TopicExists([]byte("a/b"))
	require.True(t, found, "Topic '/a/b/#/c' should exist.")
	//require.Equal(t, 2, int(msg.TopicQos([]byte("/a/b/#/cdd"))), "Incorrect topic qos.")
}

// test empty topic list
func TestSubscribeMessageDecode2(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_SUBSCRIBE<<4) | 2,
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	msg := NewSubscribe()
	_, err := msg.Decode(msgBytes)

	require.Error(t, err)
}

func TestSubscribeMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_SUBSCRIBE<<4) | 2,
		36,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // QoS
		0, // topic name MSB (0)
		8, // topic name LSB (8)
		'/', 'a', '/', 'b', '/', '#', '/', 'c',
		1,  // QoS
		0,  // topic name MSB (0)
		10, // topic name LSB (10)
		'/', 'a', '/', 'b', '/', '#', '/', 'c', 'd', 'd',
		2, // QoS
	}

	msg := NewSubscribe()
	msg.Header.SetPacketID(7)
	msg.AddTopic([]byte("surgemq"), 0)
	msg.AddTopic([]byte("/a/b/#/c"), 1)
	msg.AddTopic([]byte("/a/b/#/cdd"), 2)

	dst := make([]byte, 100)
	n, err := msg.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestSubscribeDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_SUBSCRIBE<<4) | 2,
		36,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // QoS
		0, // topic name MSB (0)
		8, // topic name LSB (8)
		'/', 'a', '/', 'b', '/', '#', '/', 'c',
		1,  // QoS
		0,  // topic name MSB (0)
		10, // topic name LSB (10)
		'/', 'a', '/', 'b', '/', '#', '/', 'c', 'd', 'd',
		2, // QoS
	}

	msg := NewSubscribe()
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
