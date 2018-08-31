package mqttgo

import (
	"encoding/binary"
	"fmt"
)

type Pubrec struct {
	//固定头
	header
}

func NewPubrec() (p *Pubrec) {
	p = &Pubrec{}
	p.header.SetType(TYPE_PUBREC)
	p.header.SetRemainingLength(2)
	return
}
func (p Pubrec) String() string {
	return fmt.Sprintf("%s, PacketID=%d", p.header, p.header.GetPacketID())
}
func (p *Pubrec) Length() int {
	return p.header.Length() + int(p.header.GetRemainingLength())
}
func (p *Pubrec) Encode(dst []byte) (total int, err error) {
	var (
		l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("Pubrec/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	total = 0
	n, err = p.header.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	binary.BigEndian.PutUint16(dst[total:], p.header.GetPacketID())
	total += 2
	return
}
func (p *Pubrec) Decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.header.decode(src[total:])
	total += n
	if err != nil {
		return
	}
	p.header.SetPacketID(binary.BigEndian.Uint16(src[total:]))
	total += 2
	return
}
