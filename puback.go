package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Puback struct {
	//固定头
	Header Header
}

func NewPuback() (p *Puback) {
	p.Header.SetType(TYPE_PUBLISH)
	p.Header.SetFlag(TYPE_FLAG_PUBACK)
	p.Header.SetRemainingLength(2)
	return
}
func (p Puback) String() string {
	return fmt.Sprintf("%s, PacketID=%d", p.Header, p.Header.GetPacketID())
}
func (p *Puback) GetLength() int {
	ml, _ := binary.Uvarint(p.Header.GetRemainingLength())
	return p.Header.getLength() + int(ml)
}
func (p *Puback) encode(dst []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.Header.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	binary.BigEndian.PutUint16(dst[total:], p.Header.GetPacketID())
	total += 2
	return
}
func (p *Puback) decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.Header.decode(src[total:])
	total += n
	if err != nil {
		return
	}
	p.Header.SetPacketID(binary.BigEndian.Uint16(src[total:]))
	total += 2
	return
}
