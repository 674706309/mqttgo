package mqtt

import "fmt"

type Pingresp struct {
	Header Header
}

func NewPingresp() (p *Pingresp) {
	p.Header.SetType(TYPE_PINGRESP)
	p.Header.SetFlag(TYPE_FLAG_PINGRESP)
	return
}
func (p Pingresp) String() string {
	return fmt.Sprintf("%s", p.Header)
}
func (p *Pingresp) Length() int {
	return p.Header.Length()
}
func (p *Pingresp) encode(dst []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.Header.encode(dst[total:])
	total += n
	return
}
func (p *Pingresp) decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.Header.decode(src[total:])
	total += n
	return
}
