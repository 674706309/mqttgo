package mqtt

import (
	"encoding/binary"
	"fmt"
)

type UnSuback struct {
	Header header
}

func NewUnSuback() (u *UnSuback) {
	u = &UnSuback{}
	u.Header.SetType(TYPE_SUBACK)
	return
}
func (u UnSuback) String() string {
	return fmt.Sprintf("%s, PacketID=%d", u.Header, u.Header.GetPacketID())
}
func (u *UnSuback) Length() int {
	return u.Header.Length() + u.GetRemainingLength()
}
func (u *UnSuback) GetRemainingLength() int {
	return 2
}
func (u *UnSuback) Encode(dst []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = u.Header.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	binary.BigEndian.PutUint16(dst[total:], u.Header.GetPacketID())
	total += n
	return
}
func (u *UnSuback) Decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = u.Header.decode(src[total:])
	total += n
	if err != nil {
		return
	}
	u.Header.SetPacketID(binary.BigEndian.Uint16(src[total:]))
	total += 2
	return
}
