package mqtt

import "bytes"

type Connect struct {
	//固定头
	FixedHeader
	//可变头
	ProtocolName  []byte //协议
	ProtocolLevel byte   //协议级别
	ConnectFlags  byte   //连接标志
	KeepAlive     uint16 //保持时间单位秒
	CleanSession  bool   //清理会话
	WillFlag      bool   //遗嘱标志
	QoSWillQoS    uint8  //遗嘱服务质量等级
	WillRetain    bool   //遗嘱保留标志
	UserNameFlag  bool   //用户名标志
	PasswordFlag  bool   //密码标志
	//有效载荷
	ClientIdentifier []byte //客户端标识符
	WillTopic        []byte //遗嘱主题
	WillMessage      []byte //遗嘱消息
	UserName         []byte //用户名
	Password         []byte //密码
}

func NewConnect() (msg *Connect) {
	msg.TypeAndFlag = byte(TYPE_FLAG_CONNECT)
	msg.ProtocolName = bytes.Join([][]byte{{0x00, 0x04}, []byte(PROTOCOL)}, []byte(""))
	msg.ProtocolLevel = PROTOCOL_LEVEL
	return
}
func (c *Connect) SetProtocolName(flag []byte) {
	c.ProtocolName = flag
}
func (c *Connect) GetProtocolName() []byte {
	return c.ProtocolName
}
func (c *Connect) SetProtocolLevel(flag byte) {
	c.ProtocolLevel = flag
}
func (c *Connect) GetProtocolLevel() byte {
	return c.ProtocolLevel
}
func (c *Connect) SetConnectFlags(flag byte) {
	c.ConnectFlags = flag
}
func (c *Connect) GetConnectFlags() byte {
	return c.ConnectFlags
}
func (c *Connect) SetCleanSession(flag bool) {
	c.CleanSession = flag
}
func (c *Connect) GetCleanSession() bool {
	return c.CleanSession
}
func (c *Connect) SetWillFlag(flag bool) {
	c.WillFlag = flag
}
func (c *Connect) GetWillFlag() bool {
	return c.WillFlag
}
func (c *Connect) SetQoSWillQoS(flag byte) {
	c.QoSWillQoS = flag
}
func (c *Connect) GetQoSWillQoS() byte {
	return c.QoSWillQoS
}
func (c *Connect) SetWillRetain(flag bool) {
	c.WillRetain = flag
}
func (c *Connect) GetWillRetain() bool {
	return c.WillRetain
}
func (c *Connect) SetUserNameFlag(flag bool) {
	c.UserNameFlag = flag
}
func (c *Connect) GetUserNameFlag() bool {
	return c.UserNameFlag
}
func (c *Connect) SetPasswordFlag(flag bool) {
	c.PasswordFlag = flag
}
func (c *Connect) GetPasswordFlag() bool {
	return c.PasswordFlag
}
func (c *Connect) SetKeepAlive(flag uint16) {
	c.KeepAlive = flag
}
func (c *Connect) GetKeepAlive() uint16 {
	return c.KeepAlive
}
func (c *Connect) SetClientIdentifier(flag []byte) {
	c.ClientIdentifier = flag
}
func (c *Connect) GetClientIdentifier() []byte {
	return c.ClientIdentifier
}
func (c *Connect) SetWillTopic(flag []byte) {
	c.WillTopic = flag
}
func (c *Connect) GetWillTopic() []byte {
	return c.WillTopic
}
func (c *Connect) SetWillMessage(flag []byte) {
	c.WillMessage = flag
}
func (c *Connect) GetWillMessage() []byte {
	return c.WillMessage
}
func (c *Connect) SetUserName(flag []byte) {
	c.UserName = flag
}
func (c *Connect) GetUserName() []byte {
	return c.UserName
}
func (c *Connect) SetPassword(flag []byte) {
	c.Password = flag
}
func (c *Connect) GetPassword() []byte {
	return c.Password
}
