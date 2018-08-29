package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Pubrel struct {
	//固定头
	Header Header
}

func NewPubrel() (p *Pubrel) {
	p.Header.SetType(TYPE_PUBREL)
	p.Header.SetFlag(TYPE_FLAG_PUBREL)
	p.Header.SetRemainingLength(2)
	return
}
func (p Pubrel) String() string {
	return fmt.Sprintf("%s, PacketID=%d", p.Header, p.Header.GetPacketID())
}
func (p *Pubrel) Length() int {
	return p.Header.Length() + int(p.Header.GetRemainingLength())
}
func (p *Pubrel) encode(dst []byte) (total int, err error) {
	var (
		l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("Suback/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
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
func (p *Pubrel) decode(src []byte) (total int, err error) {
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
