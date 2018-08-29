package mqtt

import "fmt"

type Connack struct {
	//固定头
	Header Header
	//可变头
	AcknowledgeFlags uint8
	ReturnCode       uint8 //返回码
}

func NewConnack() (c *Connack) {
	c.Header.SetType(TYPE_CONNACK)
	c.Header.SetFlag(TYPE_FLAG_CONNACK)
	c.Header.SetRemainingLength(2)
	return
}
func (c Connack) String() string {
	return fmt.Sprintf("%s, SessionPresent=%t, ReturnCode=%q\n", c.Header, c.GetSessionPresent(), c.GetReturnCode())
}
func (c *Connack) SetSessionPresent(v bool) {
	if v {
		c.AcknowledgeFlags |= 0x1 // 00000001
	} else {
		c.AcknowledgeFlags &= 254 // 11111110
	}
}

func (c *Connack) GetSessionPresent() bool {
	return (c.AcknowledgeFlags & 0x1) == 1
}
func (c *Connack) SetReturnCode(code uint8) {
	c.ReturnCode = code
}
func (c *Connack) GetReturnCode() uint8 {
	return c.ReturnCode
}
func (c *Connack) Length() int {
	return c.Header.Length() + c.GetRemainingLength()
}
func (c *Connack) GetRemainingLength() int {
	return 2
}
func (c *Connack) encode(dst []byte) (total int, err error) {
	var (
		l, ml int
	)
	l = c.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("connack/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	ml = c.GetRemainingLength()
	c.Header.SetRemainingLength(uint64(ml))
	total = 0
	n, err := c.Header.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	if c.GetSessionPresent() {
		dst[total] = 1
	}
	total++
	if c.GetReturnCode() > 5 {
		return total, fmt.Errorf("connack/Encode: Invalid CONNACK return code (%d)", c.GetReturnCode())
	}
	dst[total] = c.GetReturnCode()
	total++
	return
}
func (c *Connack) decode(src []byte) (total int, err error) {
	var (
		n int
		b byte
	)
	total = 0
	n, err = c.Header.decode(src)
	total += n
	if err != nil {
		return
	}
	b = src[total]
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
	return
}
