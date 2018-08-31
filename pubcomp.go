package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Pubcomp struct {
	//固定头
	header
}

func NewPubcomp() (p *Pubcomp) {
	p = &Pubcomp{}
	p.header.SetType(TYPE_PUBCOMP)
	return
}
func (p Pubcomp) String() string {
	return fmt.Sprintf("%s, PacketID=%d", p.header, p.header.GetPacketID())
}
func (p *Pubcomp) Length() int {
	return p.header.Length() + p.GetMessageLength()
}
func (p *Pubcomp) GetMessageLength() int {
	return 2
}
func (p *Pubcomp) Encode(dst []byte) (total int, err error) {
	var (
		ml, l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("Pubcomp/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	total = 0
	ml = p.GetMessageLength()
	p.SetRemainingLength(uint64(ml))
	n, err = p.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	binary.BigEndian.PutUint16(dst[total:], p.GetPacketID())
	total += 2
	return
}
func (p *Pubcomp) Decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.decode(src[total:])
	total += n
	if err != nil {
		return
	}
	p.SetPacketID(binary.BigEndian.Uint16(src[total:]))
	total += 2
	return
}
