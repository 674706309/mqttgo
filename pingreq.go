package mqtt

import (
	"fmt"
)

type Pingreq struct {
	Header header
}

func NewPingreq() (p *Pingreq) {
	p = &Pingreq{}
	p.Header.SetType(TYPE_PINGREQ)
	return
}
func (p Pingreq) String() string {
	return fmt.Sprintf("%s", p.Header)
}
func (p *Pingreq) Length() int {
	return p.Header.Length()
}
func (p *Pingreq) encode(dst []byte) (total int, err error) {
	var (
		l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("pingreq/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	total = 0
	n, err = p.Header.encode(dst[total:])
	total += n
	return
}
func (p *Pingreq) decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.Header.decode(src[total:])
	total += n
	return
}
