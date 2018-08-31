package mqttgo

import (
	"encoding/binary"
	"fmt"
)

type UnSuback struct {
	header
}

func NewUnSuback() (u *UnSuback) {
	u = &UnSuback{}
	u.header.SetType(TYPE_UNSUBACK)
	return
}
func (u UnSuback) String() string {
	return fmt.Sprintf("%s, PacketID=%d", u.header, u.header.GetPacketID())
}
func (u *UnSuback) Length() int {
	return u.header.Length() + u.GetRemainingLength()
}
func (u *UnSuback) GetRemainingLength() int {
	return 2
}
func (u *UnSuback) Encode(dst []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	u.header.SetRemainingLength(uint64(u.GetRemainingLength()))
	n, err = u.header.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	binary.BigEndian.PutUint16(dst[total:], u.header.GetPacketID())
	total += n
	return
}
func (u *UnSuback) Decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = u.header.decode(src[total:])
	total += n
	if err != nil {
		return
	}
	u.header.SetPacketID(binary.BigEndian.Uint16(src[total:]))
	total += 2
	return
}
