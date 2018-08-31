package mqtt

import (
	"github.com/zentures/message"
	"testing"
)

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func BenchmarkConnack_Decode(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_CONNACK << 4),
		2,
		0, // session not present
		0, // connection accepted
	}
	for i := 0; i < b.N; i++ {
		c := NewConnack()
		c.Decode(msgBytes)
		//n, err := c.Decode(msgBytes)
		//require.NoError(b, err, "Error decoding message.")
		//require.Equal(b, len(msgBytes), n, "Error decoding message.")
	}
}
func BenchmarkConnack_Decode1(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_CONNACK << 4),
		2,
		0, // session not present
		0, // connection accepted
	}
	for i := 0; i < b.N; i++ {
		c := message.NewConnackMessage()
		c.Decode(msgBytes)
		//n, err := c.Decode(msgBytes)
		//require.NoError(b, err, "Error decoding message.")
		//require.Equal(b, len(msgBytes), n, "Error decoding message.")
	}
}

func BenchmarkConnect_Decode(b *testing.B) {
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
	//flag.Parse()
	//if *cpuprofile != "" {
	//	f, err := os.Create(*cpuprofile)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	pprof.StartCPUProfile(f)
	//	defer pprof.StopCPUProfile()
	//}
	for i := 0; i < b.N; i++ {
		conn := NewConnect()
		conn.Decode(msgBytes)
		//n, err := conn.Decode(msgBytes)

		//require.NoError(b, err, "Error decoding message.")
		//require.Equal(b, len(msgBytes), n, "Error decoding message.")
	}
}
func BenchmarkConnect_Decode1(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		conn := message.NewConnectMessage()
		conn.Decode(msgBytes)
		//n, err := conn.Decode(msgBytes)

		//require.NoError(b, err, "Error decoding message.")
		//require.Equal(b, len(msgBytes), n, "Error decoding message.")
	}
}

func BenchmarkDisconnect_Decode(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_DISCONNECT << 4),
		0,
	}
	for i := 0; i < b.N; i++ {
		d := NewDisconnect()
		d.Decode(msgBytes)
	}
}
func BenchmarkDisconnect_Decode1(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_DISCONNECT << 4),
		0,
	}
	for i := 0; i < b.N; i++ {
		d := message.NewDisconnectMessage()
		d.Decode(msgBytes)
	}
}

func BenchmarkPingreq_Decode(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PINGREQ << 4),
		0,
	}
	for i := 0; i < b.N; i++ {
		p := NewPingreq()
		p.Decode(msgBytes)
	}
}
func BenchmarkPingreq_Decode1(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PINGREQ << 4),
		0,
	}
	for i := 0; i < b.N; i++ {
		p := message.NewPingreqMessage()
		p.Decode(msgBytes)
	}
}

func BenchmarkPingresp_Decode(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PINGRESP << 4),
		0,
	}
	for i := 0; i < b.N; i++ {
		p := NewPingresp()
		p.Decode(msgBytes)
	}

}
func BenchmarkPingresp_Decode1(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PINGRESP << 4),
		0,
	}
	for i := 0; i < b.N; i++ {
		p := message.NewPingrespMessage()
		p.Decode(msgBytes)
	}
}

func BenchmarkPuback_Decode(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}
	for i := 0; i < b.N; i++ {
		p := NewPuback()
		p.Decode(msgBytes)
	}
}
func BenchmarkPuback_Decode1(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}
	for i := 0; i < b.N; i++ {
		p := message.NewPubackMessage()
		p.Decode(msgBytes)
	}
}

func BenchmarkPubcomp_Decode(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PUBCOMP << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}
	for i := 0; i < b.N; i++ {
		p := NewPubcomp()
		p.Decode(msgBytes)
	}
}
func BenchmarkPubcomp_Decode1(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PUBCOMP << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}
	for i := 0; i < b.N; i++ {
		p := message.NewPubcompMessage()
		p.Decode(msgBytes)
	}
}

func BenchmarkPublish_Decode(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		p := NewPublish()
		p.Decode(msgBytes)
	}
}
func BenchmarkPublish_Decode1(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		p := message.NewPublishMessage()
		p.Decode(msgBytes)
	}
}

func BenchmarkPubrec_Decode(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PUBREC << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}
	for i := 0; i < b.N; i++ {
		p := NewPubrec()
		p.Decode(msgBytes)
	}
}
func BenchmarkPubrec_Decode1(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PUBREC << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}
	for i := 0; i < b.N; i++ {
		p := message.NewPubrecMessage()
		p.Decode(msgBytes)
	}
}

func BenchmarkPubrel_Decode(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PUBREL<<4) | 2,
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}
	for i := 0; i < b.N; i++ {
		p := NewPubrel()
		p.Decode(msgBytes)
	}
}
func BenchmarkPubrel_Decode1(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_PUBREL<<4) | 2,
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}
	for i := 0; i < b.N; i++ {
		p := message.NewPubrelMessage()
		p.Decode(msgBytes)
	}
}

func BenchmarkSuback_Decode(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		msg := NewSuback()
		msg.Decode(msgBytes)
	}
}
func BenchmarkSuback_Decode1(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		msg := message.NewSubackMessage()
		msg.Decode(msgBytes)
	}
}

func BenchmarkSubscribe_Decode(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		msg := NewSubscribe()
		msg.Decode(msgBytes)
	}
}
func BenchmarkSubscribe_Decode1(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		msg := message.NewSubscribeMessage()
		msg.Decode(msgBytes)
	}
}

func BenchmarkUnSuback_Decode(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_UNSUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}
	for i := 0; i < b.N; i++ {
		msg := NewUnSuback()
		msg.Decode(msgBytes)
	}
}
func BenchmarkUnSuback_Decode1(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_UNSUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}
	for i := 0; i < b.N; i++ {
		msg := message.NewUnsubackMessage()
		msg.Decode(msgBytes)
	}
}

func BenchmarkUnSubscribe_Decode(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_UNSUBSCRIBE<<4) | 2,
		33,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // topic name MSB (0)
		8, // topic name LSB (8)
		'/', 'a', '/', 'b', '/', '#', '/', 'c',
		0,  // topic name MSB (0)
		10, // topic name LSB (10)
		'/', 'a', '/', 'b', '/', '#', '/', 'c', 'd', 'd',
	}
	for i := 0; i < b.N; i++ {
		msg := NewUnSubscribe()
		msg.Decode(msgBytes)
	}
}
func BenchmarkUnSubscribe_Decode1(b *testing.B) {
	msgBytes := []byte{
		byte(TYPE_UNSUBSCRIBE<<4) | 2,
		33,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // topic name MSB (0)
		8, // topic name LSB (8)
		'/', 'a', '/', 'b', '/', '#', '/', 'c',
		0,  // topic name MSB (0)
		10, // topic name LSB (10)
		'/', 'a', '/', 'b', '/', '#', '/', 'c', 'd', 'd',
	}
	for i := 0; i < b.N; i++ {
		msg := message.NewUnsubscribeMessage()
		msg.Decode(msgBytes)
	}
}
