package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Pubrel struct {
	//固定头
	header
}

func NewPubrel() (p *Pubrel) {
	p = &Pubrel{}
	p.header.SetType(TYPE_PUBREL)
	p.header.SetRemainingLength(2)
	return
}
func (p Pubrel) String() string {
	return fmt.Sprintf("%s, PacketID=%d", p.header, p.header.GetPacketID())
}
func (p *Pubrel) Length() int {
	return p.header.Length() + int(p.header.GetRemainingLength())
}
func (p *Pubrel) Encode(dst []byte) (total int, err error) {
	var (
		l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("Suback/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	total = 0
	n, err = p.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	binary.BigEndian.PutUint16(dst[total:], p.header.GetPacketID())
	total += 2
	return
}
func (p *Pubrel) Decode(src []byte) (total int, err error) {
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
