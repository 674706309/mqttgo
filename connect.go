package mqtt

import (
	"encoding/binary"
	"fmt"
	"regexp"
)

var ClientIDRegexp *regexp.Regexp

type Connect struct {
	//固定头
	Header FixedHeader
	//可变头
	ProtocolName  []byte //协议
	ProtocolLevel byte   //协议级别
	// 7: UserNameFlag 密码标志
	// 6: PasswordFlag 用户名标志
	// 5: WillRetain 遗嘱保留标志
	// 4-3: WillQoS 遗嘱服务质量等级
	// 2: WillFlag 遗嘱标志
	// 1: CleanSession 清理会话
	// 0: Reserved 服务端必须验证CONNECT控制报文的保留标志位（第0位）是否为0，如果不为0必须断开客户端连接
	ConnectFlags byte   //连接标志 上述具体内容
	KeepAlive    uint16 //保持时间单位秒
	//有效载荷
	ClientID    []byte //客户端标识符
	WillTopic   []byte //遗嘱主题
	WillMessage []byte //遗嘱消息
	UserName    []byte //用户名
	Password    []byte //密码
}

func init() {
	ClientIDRegexp = regexp.MustCompile("^[0-9a-zA-Z _]*$")
}
func NewConnect() (c *Connect) {
	c.Header.SetType(TYPE_CONNECT)
	c.Header.SetFlag(TYPE_FLAG_CONNECT)
	//msg.SetProtocolName(bytes.Join([][]byte{{0x00, 0x04}, []byte(PROTOCOL)}, []byte("")))
	//msg.SetProtocolLevel(PROTOCOL_LEVEL)
	return
}
func (c *Connect) SetProtocolName(t []byte) {
	c.ProtocolName = t
}
func (c *Connect) GetProtocolName() []byte {
	return c.ProtocolName
}
func (c *Connect) SetProtocolLevel(t byte) {
	c.ProtocolLevel = t
}
func (c *Connect) GetProtocolLevel() byte {
	return c.ProtocolLevel
}
func (c *Connect) SetConnectFlags(t byte) {
	c.ConnectFlags = t
}
func (c *Connect) GetConnectFlags() byte {
	return c.ConnectFlags
}
func (c *Connect) SetCleanSession(v bool) {
	if v {
		c.ConnectFlags |= 0x2 // 00000010
	} else {
		c.ConnectFlags &= 253 // 11111101
	}
}
func (c *Connect) GetCleanSession() bool {
	return ((c.ConnectFlags >> 1) & 0x1) == 1
}
func (c *Connect) SetWillFlag(v bool) {
	if v {
		c.ConnectFlags |= 0x4 // 00000100
	} else {
		c.ConnectFlags &= 251 // 11111011
	}
}
func (c *Connect) GetWillFlag() bool {
	return ((c.ConnectFlags >> 2) & 0x1) == 1
}
func (c *Connect) SetWillQos(qos byte) error {
	if qos != QosAtMostOnce && qos != QosAtLeastOnce && qos != QosExactlyOnce {
		return fmt.Errorf("connect/SetWillQos: Invalid QoS level %d", qos)
	}
	c.ConnectFlags = (c.ConnectFlags & 231) | (qos << 3) // 231 = 11100111
	return nil
}
func (c *Connect) GetWillQos() byte {
	return (c.ConnectFlags >> 3) & 0x3
}
func (c *Connect) SetWillRetain(v bool) {
	if v {
		c.ConnectFlags |= 32 // 00100000
	} else {
		c.ConnectFlags &= 223 // 11011111
	}
}
func (c *Connect) GetWillRetain() bool {
	return ((c.ConnectFlags >> 5) & 0x1) == 1
}
func (c *Connect) SetUsernameFlag(v bool) {
	if v {
		c.ConnectFlags |= 128 // 10000000
	} else {
		c.ConnectFlags &= 127 // 01111111
	}
}
func (c *Connect) GetUsernameFlag() bool {
	return ((c.ConnectFlags >> 7) & 0x1) == 1
}
func (c *Connect) SetPasswordFlag(v bool) {
	if v {
		c.ConnectFlags |= 64 // 01000000
	} else {
		c.ConnectFlags &= 191 // 10111111
	}
}
func (c *Connect) GetPasswordFlag() bool {
	return ((c.ConnectFlags >> 6) & 0x1) == 1
}
func (c *Connect) SetKeepAlive(t uint16) {
	c.KeepAlive = t
}
func (c *Connect) GetKeepAlive() uint16 {
	return c.KeepAlive
}
func (c *Connect) SetClientID(t []byte) {
	c.ClientID = t
}
func (c *Connect) GetClientID() []byte {
	return c.ClientID
}
func (c *Connect) ValidClientID(t []byte) bool {

	return ClientIDRegexp.Match(t)
}
func (c *Connect) SetWillTopic(t []byte) {
	c.WillTopic = t
}
func (c *Connect) GetWillTopic() []byte {
	return c.WillTopic
}
func (c *Connect) SetWillMessage(t []byte) {
	c.WillMessage = t
}
func (c *Connect) GetWillMessage() []byte {
	return c.WillMessage
}
func (c *Connect) SetUserName(t []byte) {
	c.UserName = t
}
func (c *Connect) GetUserName() []byte {
	return c.UserName
}
func (c *Connect) SetPassword(t []byte) {
	c.Password = t
}
func (c *Connect) GetPassword() []byte {
	return c.Password
}
func (c *Connect) GetLength() int {
	total := 0
	// n bytes protocol name
	// 1 byte protocol version
	// 1 byte connect flags
	// 2 bytes keep alive timer
	total += len(c.GetProtocolName()) + 1 + 1 + 1 + 2
	total += 2 + len(c.GetClientID())
	if c.GetWillFlag() {
		total += 2 + len(c.GetWillTopic()) + 2 + len(c.GetWillMessage())
	}
	if c.GetUsernameFlag() && len(c.GetUserName()) > 0 {
		total += 2 + len(c.GetUserName())
	}
	if c.GetPasswordFlag() && len(c.GetPassword()) > 0 {
		total += 2 + len(c.GetPassword())
	}
	return total
}
func (c *Connect) String() string {
	return fmt.Sprintf("%s, ProtocolName=%q, ProtocolLevel=%d, ConnectFlags=%08b, KeepAlive=%d, ClientID=%q, WillTopic=%q, WillMessage=%q, UserName=%q, Password=%q",
		c.Header,
		c.GetProtocolName(),
		c.GetProtocolLevel(),
		c.GetConnectFlags(),
		c.GetKeepAlive(),
		c.GetClientID(),
		c.GetWillTopic(),
		c.GetWillMessage(),
		c.GetUserName(),
		c.GetPassword(),
	)
}
func (c *Connect) encode(dst []byte) (int, error) {
	if t := c.Header.GetType(); t != TYPE_CONNECT {
		return 0, fmt.Errorf("connect/Encode: Invalid message type. Expecting %d, got %d", TYPE_CONNECT, t)
	}
	if _, ok := SupportedVersions[c.GetProtocolLevel()]; !ok {
		return 0, fmt.Errorf("ErrInvalidProtocolVersion")
	}
	ml := c.GetLength()
	c.Header.SetRemainingLength(uint64(ml))
	total := 0
	n, err := c.Header.encode(dst[total:])
	total += n
	if err != nil {
		return total, err
	}
	n, err = c.encodeMessage(dst[total:])
	total += n
	if err != nil {
		return total, err
	}
	return total, nil
}

func (c *Connect) encodeMessage(dst []byte) (int, error) {
	total := 0
	n, err := WriteBytes(dst[total:], c.GetProtocolName())
	total += n
	if err != nil {
		return total, err
	}
	dst[total] = c.GetProtocolLevel()
	total += 1
	dst[total] = c.GetConnectFlags()
	total += 1
	binary.BigEndian.PutUint16(dst[total:], c.GetKeepAlive())
	total += 2
	n, err = WriteBytes(dst[total:], c.GetClientID())
	total += n
	if err != nil {
		return total, err
	}
	if c.GetWillFlag() {
		n, err = WriteBytes(dst[total:], c.GetWillTopic())
		total += n
		if err != nil {
			return total, err
		}

		n, err = WriteBytes(dst[total:], c.GetWillMessage())
		total += n
		if err != nil {
			return total, err
		}
	}
	if c.GetUsernameFlag() && len(c.GetUserName()) > 0 {
		n, err = WriteBytes(dst[total:], c.GetUserName())
		total += n
		if err != nil {
			return total, err
		}
	}
	if c.GetPasswordFlag() && len(c.GetPassword()) > 0 {
		n, err = WriteBytes(dst[total:], c.GetPassword())
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}
func (c *Connect) decode(src []byte) (int, error) {
	total := 0
	n, err := c.Header.decode(src[total:])
	if err != nil {
		return total + n, err
	}
	total += n
	if n, err = c.decodeMessage(src[total:]); err != nil {
		return total + n, err
	}
	total += n
	return total, nil
}

func (c *Connect) decodeMessage(src []byte) (int, error) {
	var (
		err      error
		temp     []byte
		n, total int
	)
	n, total = 0, 0
	n, err = ReadBytes(src[total:], temp)
	total += n
	if err != nil {
		return total, err
	}
	if _, ok := SupportedVersions[src[total]]; !ok {
		return total, fmt.Errorf("ErrInvalidProtocolVersion")
	} else {
		c.SetProtocolLevel(src[total])
		total++
	}
	c.SetConnectFlags(src[total])
	total++
	if c.GetConnectFlags()&0x1 != 0 {
		return total, fmt.Errorf("connect/decodeMessage: Connect Flags reserved bit 0 is not 0")
	}
	if c.GetWillQos() > QosExactlyOnce {
		return total, fmt.Errorf("connect/decodeMessage: Invalid QoS level (%d) for %d message", c.GetWillQos(), c.Header.GetType())
	}
	if !c.GetWillFlag() && (c.GetWillRetain() || c.GetWillQos() != QosAtMostOnce) {
		return total, fmt.Errorf("connect/decodeMessage: Protocol violation: If the Will Flag (%t) is set to 0 the Will QoS (%d) and Will Retain (%t) fields MUST be set to zero", c.GetWillFlag(), c.GetWillQos(), c.GetWillRetain())
	}
	if c.GetUsernameFlag() && !c.GetPasswordFlag() {
		return total, fmt.Errorf("connect/decodeMessage: Username flag is set but Password flag is not set")
	}
	if len(src[total:]) < 2 {
		return 0, fmt.Errorf("connect/decodeMessage: Insufficient buffer size. Expecting %d, got %d", 2, len(src[total:]))
	}
	c.SetKeepAlive(binary.BigEndian.Uint16(src[total:]))
	total += 2
	n, err = ReadBytes(src[total:], temp)
	c.SetClientID(temp)
	total += n
	if err != nil {
		return total, err
	}
	if len(c.GetClientID()) == 0 && !c.GetCleanSession() {
		return total, fmt.Errorf("ErrIdentifierRejected")
	}
	if len(c.GetClientID()) > 0 && !c.ValidClientID(c.GetClientID()) {
		return total, fmt.Errorf("ErrIdentifierRejected")
	}
	if c.GetWillFlag() {
		n, err = ReadBytes(src[total:], temp)
		total += n
		if err != nil {
			return total, err
		}
		n, err = ReadBytes(src[total:], temp)
		total += n
		if err != nil {
			return total, err
		}
	}
	if c.GetUsernameFlag() && len(src[total:]) > 0 {
		n, err = ReadBytes(src[total:], temp)
		c.SetUserName(temp)
		total += n
		if err != nil {
			return total, err
		}
	}
	if c.GetPasswordFlag() && len(src[total:]) > 0 {
		n, err = ReadBytes(src[total:], temp)
		c.SetPassword(temp)
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}
