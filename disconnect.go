package mqtt

import "fmt"

type Disconnect struct {
	Header Header
}

func NewDisconnect() (p *Disconnect) {
	p.Header.SetType(TYPE_DISCONNECT)
	p.Header.SetFlag(TYPE_FLAG_DISCONNECT)
	return
}
func (p Disconnect) String() string {
	return fmt.Sprintf("%s", p.Header)
}
func (p *Disconnect) Length() int {
	return p.Header.Length()
}
func (p *Disconnect) encode(dst []byte) (total int, err error) {
	var (
		l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("disconnect/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	total = 0
	n, err = p.Header.encode(dst[total:])
	total += n
	return
}
func (p *Disconnect) decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.Header.decode(src[total:])
	total += n
	return
}
