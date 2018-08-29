package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Pubcomp struct {
	//固定头
	Header Header
}

func NewPubcomp() (p *Pubcomp) {
	p.Header.SetType(TYPE_PUBCOMP)
	p.Header.SetFlag(TYPE_FLAG_PUBCOMP)
	p.Header.SetRemainingLength(2)
	return
}
func (p Pubcomp) String() string {
	return fmt.Sprintf("%s, PacketID=%d", p.Header, p.Header.GetPacketID())
}
func (p *Pubcomp) GetLength() int {
	ml, _ := binary.Uvarint(p.Header.GetRemainingLength())
	return p.Header.Length() + int(ml)
}
func (p *Pubcomp) encode(dst []byte) (total int, err error) {
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
func (p *Pubcomp) decode(src []byte) (total int, err error) {
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
