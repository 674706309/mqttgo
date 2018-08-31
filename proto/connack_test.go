package mqttgo

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConnackMessageFields(t *testing.T) {
	c := NewConnack()

	c.SetSessionPresent(true)
	require.True(t, c.GetSessionPresent(), "Error setting session present flag.")

	c.SetSessionPresent(false)
	require.False(t, c.GetSessionPresent(), "Error setting session present flag.")

	c.SetReturnCode(CONNBAK_ACCEPT)
	require.Equal(t, CONNBAK_ACCEPT, c.GetReturnCode(), "Error setting return code.")
}

func TestConnackMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_CONNACK << 4),
		2,
		0, // session not present
		0, // connection accepted
	}

	c := NewConnack()

	n, err := c.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.False(t, c.GetSessionPresent(), "Error decoding session present flag.")
	require.Equal(t, CONNBAK_ACCEPT, c.GetReturnCode(), "Error decoding return code.")
}

// testing wrong message length
func TestConnackMessageDecode2(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_CONNACK << 4),
		3,
		0, // session not present
		0, // connection accepted
	}

	c := NewConnack()

	_, err := c.Decode(msgBytes)
	require.Error(t, err, "Error decoding message.")
}

// testing wrong message size
func TestConnackMessageDecode3(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_CONNACK << 4),
		2,
		0, // session not present
	}

	c := NewConnack()

	_, err := c.Decode(msgBytes)
	require.Error(t, err, "Error decoding message.")
}

// testing wrong reserve bits
func TestConnackMessageDecode4(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_CONNACK << 4),
		2,
		64, // <- wrong size
		0,  // connection accepted
	}

	c := NewConnack()

	_, err := c.Decode(msgBytes)
	require.Error(t, err, "Error decoding message.")
}

// testing invalid return code
func TestConnackMessageDecode5(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_CONNACK << 4),
		2,
		0,
		6, // <- wrong code
	}

	c := NewConnack()

	_, err := c.Decode(msgBytes)
	require.Error(t, err, "Error decoding message.")
}

func TestConnackMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_CONNACK << 4),
		2,
		1, // session present
		0, // connection accepted
	}

	c := NewConnack()
	c.SetReturnCode(CONNBAK_ACCEPT)
	c.SetSessionPresent(true)

	dst := make([]byte, 10)
	n, err := c.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error encoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error encoding connack message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestConnackDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_CONNACK << 4),
		2,
		0, // session not present
		0, // connection accepted
	}

	c := NewConnack()
	n, err := c.Decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")

	dst := make([]byte, 100)
	n2, err := c.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n2, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n2], "Error decoding message.")

	n3, err := c.Decode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n3, "Error decoding message.")
}
