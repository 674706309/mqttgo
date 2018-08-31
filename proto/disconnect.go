package mqttgo

import "fmt"

type Disconnect struct {
	header
}

func NewDisconnect() (p *Disconnect) {
	p = &Disconnect{}
	p.header.SetType(TYPE_DISCONNECT)
	return
}
func (p Disconnect) String() string {
	return fmt.Sprintf("%s", p.header)
}
func (p *Disconnect) Length() int {
	return p.header.Length()
}
func (p *Disconnect) Encode(dst []byte) (total int, err error) {
	var (
		l, n int
	)
	l = p.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("disconnect/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	total = 0
	n, err = p.header.encode(dst[total:])
	total += n
	return
}
func (p *Disconnect) Decode(src []byte) (total int, err error) {
	var (
		n int
	)
	total = 0
	n, err = p.header.decode(src[total:])
	total += n
	return
}
