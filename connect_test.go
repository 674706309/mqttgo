package mqtt

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConnectMessageFields(t *testing.T) {
	conn := NewConnect()

	require.Equal(t, 0x4, int(conn.GetProtocolLevel()), "Incorrect version number")

	conn.SetCleanSession(true)
	require.True(t, conn.GetCleanSession(), "Error setting clean session flag.")

	conn.SetCleanSession(false)
	require.False(t, conn.GetCleanSession(), "Error setting clean session flag.")

	conn.SetWillFlag(false)
	require.False(t, conn.GetWillFlag(), "Error setting will flag.")

	conn.SetWillFlag(true)
	require.True(t, conn.GetWillFlag(), "Error setting will flag.")

	conn.SetWillRetain(true)
	require.True(t, conn.GetWillRetain(), "Error setting will retain.")

	conn.SetWillRetain(false)
	require.False(t, conn.GetWillRetain(), "Error setting will retain.")

	conn.SetPasswordFlag(true)
	require.True(t, conn.GetPasswordFlag(), "Error setting password flag.")

	conn.SetPasswordFlag(false)
	require.False(t, conn.GetPasswordFlag(), "Error setting password flag.")

	conn.SetUsernameFlag(true)
	require.True(t, conn.GetUsernameFlag(), "Error setting username flag.")

	conn.SetUsernameFlag(false)
	require.False(t, conn.GetUsernameFlag(), "Error setting username flag.")

	conn.SetWillQos(1)
	require.Equal(t, 1, int(conn.GetWillQos()), "Error setting will QoS.")

	err := conn.SetWillQos(4)
	require.Error(t, err)

	conn.SetClientID([]byte("j0j0jfajf02j0asdjf"))
	require.Equal(t, "j0j0jfajf02j0asdjf", string(conn.GetClientID()), "Error setting client ID.")

	err = conn.SetClientID([]byte("this is good for v3"))
	require.NoError(t, err)

	conn.SetProtocolLevel(0x4)

	err = conn.SetClientID([]byte("this is no good for v4!"))
	require.Error(t, err)

	conn.SetProtocolLevel(0x3)

	conn.SetWillTopic([]byte("willtopic"))
	require.Equal(t, "willtopic", string(conn.GetWillTopic()), "Error setting will topic.")

	require.True(t, conn.GetWillFlag(), "Error setting will flag.")

	conn.SetWillTopic([]byte(""))
	require.Equal(t, "", string(conn.GetWillTopic()), "Error setting will topic.")

	//require.False(t, conn.GetWillFlag(), "Error setting will flag.")

	conn.SetWillMessage([]byte("this is a will message"))
	require.Equal(t, "this is a will message", string(conn.GetWillMessage()), "Error setting will message.")

	require.True(t, conn.GetWillFlag(), "Error setting will flag.")

	conn.SetWillMessage([]byte(""))
	require.Equal(t, "", string(conn.GetWillMessage()), "Error setting will topic.")

	//require.False(t, conn.GetWillFlag(), "Error setting will flag.")

	conn.SetWillTopic([]byte("willtopic"))
	conn.SetWillMessage([]byte("this is a will message"))
	conn.SetWillTopic([]byte(""))
	require.True(t, conn.GetWillFlag(), "Error setting will topic.")

	conn.SetUserName([]byte("myname"))
	require.Equal(t, "myname", string(conn.GetUserName()), "Error setting will message.")

	require.True(t, conn.GetUsernameFlag(), "Error setting will flag.")

	conn.SetUserName([]byte(""))
	require.Equal(t, "", string(conn.GetUserName()), "Error setting will message.")

	require.False(t, conn.GetUsernameFlag(), "Error setting will flag.")

	conn.SetPassword([]byte("myname"))
	require.Equal(t, "myname", string(conn.GetPassword()), "Error setting will message.")

	require.True(t, conn.GetPasswordFlag(), "Error setting will flag.")

	conn.SetPassword([]byte(""))
	require.Equal(t, "", string(conn.GetPassword()), "Error setting will message.")

	require.False(t, conn.GetPasswordFlag(), "Error setting will flag.")
}

func TestConnectMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_CONNECT) << 4,
		60,
		0, // Length MSB (0)
		4, // Length LSB (4)
		'M', 'Q', 'T', 'T',
		4,   // Protocol level 4
		206, // connect flags 11001110, will QoS = 01
		0,   // Keep Alive MSB (0)
		10,  // Keep Alive LSB (10)
		0,   // Client ID MSB (0)
		7,   // Client ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // Will Topic MSB (0)
		4, // Will Topic LSB (4)
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB (0)
		12, // Will Message LSB (12)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB (0)
		7, // Username ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0,  // Password ID MSB (0)
		10, // Password ID LSB (10)
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e', 't',
	}

	conn := NewConnect()
	n, err := conn.decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, 206, int(conn.GetConnectFlags()), "Incorrect flag value.")
	require.Equal(t, 10, int(conn.GetKeepAlive()), "Incorrect KeepAlive value.")
	require.Equal(t, "surgemq", string(conn.GetClientID()), "Incorrect client ID value.")
	require.Equal(t, "will", string(conn.GetWillTopic()), "Incorrect will topic value.")
	require.Equal(t, "send me home", string(conn.GetWillMessage()), "Incorrect will message value.")
	require.Equal(t, "surgemq", string(conn.GetUserName()), "Incorrect username value.")
	require.Equal(t, "verysecret", string(conn.GetPassword()), "Incorrect password value.")
}

func TestConnectMessageDecode2(t *testing.T) {
	// missing last byte 't'
	msgBytes := []byte{
		byte(TYPE_CONNECT) << 4,
		60,
		0, // Length MSB (0)
		4, // Length LSB (4)
		'M', 'Q', 'T', 'T',
		4,   // Protocol level 4
		206, // connect flags 11001110, will QoS = 01
		0,   // Keep Alive MSB (0)
		10,  // Keep Alive LSB (10)
		0,   // Client ID MSB (0)
		7,   // Client ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // Will Topic MSB (0)
		4, // Will Topic LSB (4)
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB (0)
		12, // Will Message LSB (12)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB (0)
		7, // Username ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0,  // Password ID MSB (0)
		10, // Password ID LSB (10)
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e',
	}

	conn := NewConnect()
	_, err := conn.decode(msgBytes)

	require.Error(t, err)
}

func TestConnectMessageDecode3(t *testing.T) {
	// extra bytes
	msgBytes := []byte{
		byte(TYPE_CONNECT) << 4,
		60,
		0, // Length MSB (0)
		4, // Length LSB (4)
		'M', 'Q', 'T', 'T',
		4,   // Protocol level 4
		206, // connect flags 11001110, will QoS = 01
		0,   // Keep Alive MSB (0)
		10,  // Keep Alive LSB (10)
		0,   // Client ID MSB (0)
		7,   // Client ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // Will Topic MSB (0)
		4, // Will Topic LSB (4)
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB (0)
		12, // Will Message LSB (12)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB (0)
		7, // Username ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0,  // Password ID MSB (0)
		10, // Password ID LSB (10)
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e', 't',
		'e', 'x', 't', 'r', 'a',
	}

	conn := NewConnect()
	n, err := conn.decode(msgBytes)

	require.NoError(t, err)
	require.Equal(t, 62, n)
}

func TestConnectMessageDecode4(t *testing.T) {
	// missing client Id, clean session == 0
	msgBytes := []byte{
		byte(TYPE_CONNECT),
		53,
		0, // Length MSB (0)
		4, // Length LSB (4)
		'M', 'Q', 'T', 'T',
		4,   // Protocol level 4
		204, // connect flags 11001110, will QoS = 01
		0,   // Keep Alive MSB (0)
		10,  // Keep Alive LSB (10)
		0,   // Client ID MSB (0)
		0,   // Client ID LSB (0)
		0,   // Will Topic MSB (0)
		4,   // Will Topic LSB (4)
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB (0)
		12, // Will Message LSB (12)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB (0)
		7, // Username ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0,  // Password ID MSB (0)
		10, // Password ID LSB (10)
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e', 't',
	}

	conn := NewConnect()
	_, err := conn.decode(msgBytes)

	require.Error(t, err)
}

func TestConnectMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_CONNECT) << 4,
		60,
		0, // Length MSB (0)
		4, // Length LSB (4)
		'M', 'Q', 'T', 'T',
		4,   // Protocol level 4
		206, // connect flags 11001110, will QoS = 01
		0,   // Keep Alive MSB (0)
		10,  // Keep Alive LSB (10)
		0,   // Client ID MSB (0)
		7,   // Client ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // Will Topic MSB (0)
		4, // Will Topic LSB (4)
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB (0)
		12, // Will Message LSB (12)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB (0)
		7, // Username ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0,  // Password ID MSB (0)
		10, // Password ID LSB (10)
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e', 't',
	}

	conn := NewConnect()
	conn.SetWillFlag(true)
	conn.SetWillQos(1)
	conn.SetProtocolLevel(4)
	conn.SetCleanSession(true)
	conn.SetClientID([]byte("surgemq"))
	conn.SetKeepAlive(10)
	conn.SetWillTopic([]byte("will"))
	conn.SetWillMessage([]byte("send me home"))
	conn.SetUserName([]byte("surgemq"))
	conn.SetPassword([]byte("verysecret"))

	dst := make([]byte, 100)
	n, err := conn.encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, true, conn.GetPasswordFlag(), "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestConnectDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(TYPE_CONNECT) << 4,
		60,
		0, // Length MSB (0)
		4, // Length LSB (4)
		'M', 'Q', 'T', 'T',
		4,   // Protocol level 4
		206, // connect flags 11001110, will QoS = 01
		0,   // Keep Alive MSB (0)
		10,  // Keep Alive LSB (10)
		0,   // Client ID MSB (0)
		7,   // Client ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // Will Topic MSB (0)
		4, // Will Topic LSB (4)
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB (0)
		12, // Will Message LSB (12)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB (0)
		7, // Username ID LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0,  // Password ID MSB (0)
		10, // Password ID LSB (10)
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e', 't',
	}

	conn := NewConnect()
	n, err := conn.decode(msgBytes)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")

	dst := make([]byte, 120)
	n2, err := conn.encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n2, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n2], "Error decoding message.")

	n3, err := conn.decode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n3, "Error decoding message.")
}
