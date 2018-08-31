package mqtt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Sirupsen/logrus"
)

type Subscribe struct {
	//固定头
	Header header

	TopicFilter  [][]byte
	RequestedQoS []byte
}

func NewSubscribe() (s *Subscribe) {
	s = &Subscribe{}
	s.Header.SetType(TYPE_SUBSCRIBE)
	return
}
func (s Subscribe) String() (str string) {
	str = fmt.Sprintf("%s, Packet ID=%d", s.Header, s.Header.GetPacketID())
	qos := s.GetQos()
	for i, t := range s.GetTopicFilter() {
		str += fmt.Sprintf(", Topic[%d]=%q/%d", i, string(t), qos[i])
	}
	return
}
func (s *Subscribe) AddTopic(topic []byte, qos byte) error {
	if !ValidQos(qos) {
		return fmt.Errorf("Invalid QoS %d", qos)
	}
	var (
		i     int
		found bool
	)
	i, found = s.TopicExists(topic)
	if found {
		s.RequestedQoS[i] = qos
		return nil
	}
	s.TopicFilter = append(s.TopicFilter, topic)
	s.RequestedQoS = append(s.RequestedQoS, qos)
	return nil
}
func (s *Subscribe) RemoveTopic(topic []byte) (err error) {
	var (
		i     int
		found bool
	)
	i, found = s.TopicExists(topic)
	if found {
		s.TopicFilter = append(s.TopicFilter[:i], s.TopicFilter[i+1:]...)
		s.RequestedQoS = append(s.RequestedQoS[:i], s.RequestedQoS[i+1:]...)
		return nil
	}
	return fmt.Errorf("Topic don't Exists")
}
func (s *Subscribe) TopicExists(topic []byte) (int, bool) {
	for i, t := range s.GetTopicFilter() {
		if bytes.Equal(t, topic) {
			return i, true
		}
	}
	return -1, false
}
func (s *Subscribe) GetTopicFilter() [][]byte {
	return s.TopicFilter
}
func (s *Subscribe) AddQos(t byte) {
	s.RequestedQoS = append(s.RequestedQoS, t)
}
func (s *Subscribe) GetQos() []byte {
	return s.RequestedQoS
}
func (s *Subscribe) Length() int {
	return s.Header.Length() + s.GetRemainingLength()
}
func (s *Subscribe) GetRemainingLength() (total int) {
	total = 2
	for _, t := range s.GetTopicFilter() {
		total += 2 + len(t) + 1
	}
	return
}
func (s *Subscribe) Encode(dst []byte) (total int, err error) {
	var (
		ml, l, n int
		qos      []byte
	)
	l = s.Length()
	if len(dst) < l {
		return 0, fmt.Errorf("Subscribe/Encode: Insufficient buffer size. Expecting %d, got %d", l, len(dst))
	}
	ml = s.GetRemainingLength()
	s.Header.SetRemainingLength(uint64(ml))
	total = 0
	n, err = s.Header.encode(dst[total:])
	total += n
	if err != nil {
		return total, err
	}
	binary.BigEndian.PutUint16(dst[total:], s.Header.GetPacketID())
	total += n
	qos = s.GetQos()
	for i, t := range s.GetTopicFilter() {
		n, err = WriteBytes(dst[total:], t)
		total += n
		if err != nil {
			return
		}
		dst[total] = qos[i]
		total++
	}
	return
}
func (s *Subscribe) Decode(src []byte) (total int, err error) {
	var (
		hl, ml, n int
		temp      []byte
	)
	total = 0
	hl, err = s.Header.decode(src[total:])
	total += hl
	if err != nil {
		return
	}
	s.Header.SetPacketID(binary.BigEndian.Uint16(src[total:]))
	total += 2
	ml = int(s.Header.GetRemainingLength()) - (total - hl)
	for ml > 0 {
		temp, n, err = ReadBytes(src[total:])
		total += n
		if err != nil {
			return
		}
		s.AddTopic(temp, src[total])
		//s.AddQos(src[total])
		logrus.Infoln(s.GetQos())
		total++
		ml -= n + 1
	}
	if len(s.GetTopicFilter()) == 0 {
		return 0, fmt.Errorf("subscribe/Decode: Empty topic list")
	}
	return total, nil
}
