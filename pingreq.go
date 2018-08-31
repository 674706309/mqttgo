package mqtt

import (
	"fmt"
)

type Pingreq struct {
	header
}

func NewPingreq() (p *Pingreq) {
	p = &Pingreq{}
	p.header.SetType(TYPE_PINGREQ)
	return
}
func (p Pingreq) String() string {
	return fmt.Sprintf("%s", p.header)
}
func (p *Pingreq) Length() int {
	return p.header.Length()
}
func (p *Pingreq) Encode(dst []byte) (total int, err error) {
	var (
		l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("pingreq/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	total = 0
	n, err = p.header.encode(dst[total:])
	total += n
	return
}
func (p *Pingreq) Decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.header.decode(src[total:])
	total += n
	return
}
