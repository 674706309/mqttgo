package mqtt

type Puback struct {
	//固定头
	FixedHeader
	//可变头
	PacketIdentifier uint16 //报文标识符
}
