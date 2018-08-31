package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Puback struct {
	//固定头
	header
}

func NewPuback() (p *Puback) {
	p = &Puback{}
	p.header.SetType(TYPE_PUBACK)
	return
}
func (p Puback) String() string {
	return fmt.Sprintf("%s, PacketID=%d", p.header, p.GetPacketID())
}
func (p *Puback) Length() int {
	return p.header.Length() + p.GetMessageLength()
}
func (p *Puback) GetMessageLength() int {
	return 2
}
func (p *Puback) Encode(dst []byte) (total int, err error) {
	var (
		ml, l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("pingreq/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	total = 0
	ml = p.GetMessageLength()
	p.SetRemainingLength(uint64(ml))
	n, err = p.header.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	binary.BigEndian.PutUint16(dst[total:], p.header.packetID)
	total += 2
	return
}
func (p *Puback) Decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.header.decode(src[total:])
	total += n
	if err != nil {
		return
	}
	p.header.packetID = binary.BigEndian.Uint16(src[total:])
	total += 2
	return
}
