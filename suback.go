package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Suback struct {
	header

	ReturnCode []byte
}

func NewSuback() (s *Suback) {
	s = &Suback{}
	s.header.SetType(TYPE_SUBACK)
	return
}
func (s Suback) String() string {
	return fmt.Sprintf("%s, Packet ID=%d, Return Codes=%v", s.header, s.header.GetPacketID(), s.GetReturnCodes())
}
func (s *Suback) GetReturnCodes() []byte {
	return s.ReturnCode
}
func (s *Suback) AddReturnCodes(ret []byte) error {
	for _, c := range ret {
		if c != QosAtMostOnce && c != QosAtLeastOnce && c != QosExactlyOnce && c != QosFailure {
			return fmt.Errorf("suback/AddReturnCode: Invalid return code %d. Must be 0, 1, 2, 0x80", c)
		}
		s.ReturnCode = append(s.ReturnCode, c)
	}
	return nil
}
func (s *Suback) AddReturnCode(ret byte) error {
	return s.AddReturnCodes([]byte{ret})
}

func (s *Suback) Length() int {
	return s.header.Length() + s.GetRemainingLength()
}
func (s *Suback) GetRemainingLength() int {
	return 2 + len(s.ReturnCode)
}
func (s *Suback) Encode(dst []byte) (total int, err error) {
	var (
		ml, l, n int
	)
	l = s.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("Suback/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	ml = s.GetRemainingLength()
	s.header.SetRemainingLength(uint64(ml))
	total = 0
	n, err = s.header.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	binary.BigEndian.PutUint16(dst[total:], s.header.GetPacketID())
	total += 2
	copy(dst[total:], s.ReturnCode)
	total += len(s.ReturnCode)
	return
}
func (s *Suback) Decode(src []byte) (total int, err error) {
	var (
		hl int
	)

	total = 0
	hl, err = s.header.decode(src[total:])
	total += hl
	if err != nil {
		return
	}
	s.header.SetPacketID(binary.BigEndian.Uint16(src[total:]))
	total += 2
	l := int(s.header.GetRemainingLength()) - (total - hl)
	s.ReturnCode = s.ReturnCode[:0:0]
	err = s.AddReturnCodes(src[total : total+l])
	total += len(s.ReturnCode)
	return
}
