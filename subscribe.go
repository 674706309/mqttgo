package mqtt

type Subscribe struct {
	//固定头
	FixedHeader
	//可变头
	packetId uint16 //报文标识符
	TopicFilter
}
