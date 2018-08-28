package mqtt

type Publish struct {
	//固定头
	FixedHeader
	//可变头
	TopicName        []byte //主题名
	PacketIdentifier uint16 //报文标识符
	//有效载荷
	Payload []byte
}
