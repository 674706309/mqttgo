package mqtt

import "fmt"

type connack struct {
	//固定头
	header
	//可变头
	acknowledgeFlags uint8
	returnCode       uint8 //返回码
}

func NewConnack() (c *connack) {
	c = &connack{}
	c.SetType(TYPE_CONNACK)
	c.header.SetRemainingLength(2)
	return
}
func (c connack) String() string {
	return fmt.Sprintf("%s, SessionPresent=%t, ReturnCode=%q\n", c.header, c.GetSessionPresent(), c.GetReturnCode())
}
func (c *connack) SetSessionPresent(v bool) {
	if v {
		c.acknowledgeFlags |= 0x1 // 00000001
	} else {
		c.acknowledgeFlags &= 254 // 11111110
	}
}

func (c *connack) GetSessionPresent() bool {
	return (c.acknowledgeFlags & 0x1) == 1
}
func (c *connack) SetReturnCode(code uint8) {
	c.returnCode = code
}
func (c *connack) GetReturnCode() uint8 {
	return c.returnCode
}
func (c *connack) Length() int {
	return c.header.Length() + c.GetRemainingLength()
}
func (c *connack) GetRemainingLength() int {
	return 2
}
func (c *connack) Encode(dst []byte) (total int, err error) {
	var (
		l, ml int
	)
	l = c.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("connack/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	ml = c.GetRemainingLength()
	c.SetRemainingLength(uint64(ml))
	total = 0
	n, err := c.encode(dst[total:])
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
func (c *connack) Decode(src []byte) (total int, err error) {
	var (
		n int
		b byte
	)
	total = 0
	n, err = c.decode(src)
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
