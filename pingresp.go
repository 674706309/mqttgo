package mqtt

import "fmt"

type Pingresp struct {
	header
}

func NewPingresp() (p *Pingresp) {
	p = &Pingresp{}
	p.header.SetType(TYPE_PINGRESP)
	return
}
func (p Pingresp) String() string {
	return fmt.Sprintf("%s", p.header)
}
func (p *Pingresp) Length() int {
	return p.header.Length()
}
func (p *Pingresp) Encode(dst []byte) (total int, err error) {
	var (
		l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("Pingresp/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	total = 0
	n, err = p.header.encode(dst[total:])
	total += n
	return
}
func (p *Pingresp) Decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.header.decode(src[total:])
	total += n
	return
}
