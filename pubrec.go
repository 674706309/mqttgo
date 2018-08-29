package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Pubrec struct {
	//固定头
	Header Header
}

func NewPubrec() (p *Pubrec) {
	p.Header.SetType(TYPE_PUBREC)
	p.Header.SetFlag(TYPE_FLAG_PUBREC)
	p.Header.SetRemainingLength(2)
	return
}
func (p Pubrec) String() string {
	return fmt.Sprintf("%s, PacketID=%d", p.Header, p.Header.GetPacketID())
}
func (p *Pubrec) Length() int {
	return p.Header.Length() + int(p.Header.GetRemainingLength())
}
func (p *Pubrec) encode(dst []byte) (total int, err error) {
	var (
		l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("Pubrec/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
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
func (p *Pubrec) decode(src []byte) (total int, err error) {
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
