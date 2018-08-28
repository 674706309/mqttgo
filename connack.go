package mqtt

import "fmt"

type Connack struct {
	//固定头
	Header FixedHeader
	//可变头
	ConnectAcknowledgeFlags uint8
	ConnectReturnCode       uint8 //返回码
}

func NewConnack() (c *Connack) {
	c.Header.SetType(TYPE_FLAG_CONNECT)
	c.Header.SetFlag(TYPE_FLAG_CONNACK)
	c.Header.SetRemainingLength(2)
	return
}
func (c Connack) String() string {
	return fmt.Sprintf("%s, SessionPresent=%t, ReturnCode=%q\n", c.Header, c.GetSessionPresent(), c.GetReturnCode())
}
func (c *Connack) SetSessionPresent(v bool) {
	if v {
		c.ConnectAcknowledgeFlags |= 0x1 // 00000001
	} else {
		c.ConnectAcknowledgeFlags &= 254 // 11111110
	}
}

func (c *Connack) GetSessionPresent() bool {
	return (c.ConnectAcknowledgeFlags & 0x1) == 1
}
func (c *Connack) SetReturnCode(code uint8) {
	c.ConnectReturnCode = code
}
func (c *Connack) GetReturnCode() uint8 {
	return c.ConnectReturnCode
}
func (c *Connack) GetLength() int {
	return 2
}
func (c *Connack) encode(dst []byte) (int, error) {
	ml := c.GetLength()
	c.Header.SetRemainingLength(uint64(ml))
	total := 0
	n, err := c.Header.encode(dst[total:])
	total += n
	if err != nil {
		return 0, err
	}
	if c.GetSessionPresent() {
		dst[total] = 1
	}
	total++
	if c.ConnectReturnCode > 5 {
		return total, fmt.Errorf("connack/Encode: Invalid CONNACK return code (%d)", c.ConnectReturnCode)
	}
	dst[total] = c.ConnectReturnCode
	total++
	return total, nil
}
func (c *Connack) decode(src []byte) (int, error) {
	total := 0
	n, err := c.Header.decode(src)
	total += n
	if err != nil {
		return total, err
	}
	b := src[total]
	if b&254 != 0 {
		return 0, fmt.Errorf("connack/Decode: Bits 7-1 in Connack Acknowledge Flags byte (1) are not 0")
	}
	c.SetSessionPresent(b&0x1 == 1)
	total++
	b = src[total]
	if b > 5 {
		return 0, fmt.Errorf("connack/Decode: Invalid CONNACK return code (%d)", b)
	}
	c.SetReturnCode(b)
	total++
	return total, nil
}
