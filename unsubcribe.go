package mqtt

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type UnSubscribe struct {
	//固定头
	Header Header

	TopicFilter [][]byte
}

func NewUnSubscribe() (u *UnSubscribe) {
	u.Header.SetType(TYPE_UNSUBSCRIBE)
	u.Header.SetFlag(TYPE_FLAG_UNSUBSCRIBE)
	return
}
func (u UnSubscribe) String() (str string) {
	str = fmt.Sprintf("%s, Packet ID=%d", u.Header, u.Header.GetPacketID())
	for i, t := range u.GetTopicFilter() {
		str += fmt.Sprintf(", Topic[%d]=%q", i, string(t))
	}
	return
}
func (u *UnSubscribe) AddTopic(topic []byte) {
	var (
		found bool
	)
	_, found = u.TopicExists(topic)
	if found {
		return
	}
	u.TopicFilter = append(u.TopicFilter, topic)
}
func (u *UnSubscribe) RemoveTopic(topic []byte) (err error) {
	var (
		i     int
		found bool
	)
	i, found = u.TopicExists(topic)
	if found {
		u.TopicFilter = append(u.TopicFilter[:i], u.TopicFilter[i+1:]...)
		return nil
	}
	return fmt.Errorf("Topic don't Exists")
}
func (u *UnSubscribe) TopicExists(topic []byte) (int, bool) {
	for i, t := range u.GetTopicFilter() {
		if bytes.Equal(t, topic) {
			return i, true
		}
	}
	return -1, false
}
func (u *UnSubscribe) GetTopicFilter() [][]byte {
	return u.TopicFilter
}
func (u *UnSubscribe) Length() int {
	return u.Header.Length() + u.GetRemainingLength()
}
func (u *UnSubscribe) GetRemainingLength() (total int) {
	total = 2
	for _, t := range u.GetTopicFilter() {
		total += 2 + len(t) + 1
	}
	return
}
func (u *UnSubscribe) encode(dst []byte) (total int, err error) {
	var (
		hl, ml, n int
		t         []byte
	)
	hl = u.Header.Length()
	ml = u.GetRemainingLength()
	if len(dst) < hl+ml {
		return 0, fmt.Errorf("subscribe/Encode: Insufficient buffer size. Expecting %d, got %d", hl+ml, len(dst))
	}
	u.Header.SetRemainingLength(uint64(ml))
	total = 0
	n, err = u.Header.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	binary.BigEndian.PutUint16(dst[total:], u.Header.GetPacketID())
	total += n
	for _, t = range u.GetTopicFilter() {
		n, err = WriteBytes(dst[total:], t)
		total += n
		if err != nil {
			return
		}
	}
	return
}
func (u *UnSubscribe) decode(src []byte) (total int, err error) {
	var (
		hl, ml, n int
		temp      []byte
	)

	total = 0
	hl, err = u.Header.decode(src[total:])
	total += hl
	if err != nil {
		return
	}
	u.Header.SetPacketID(binary.BigEndian.Uint16(src[total:]))
	total += 2
	ml = int(u.Header.GetRemainingLength()) - (total - hl)
	for ml > 0 {
		n, err = ReadBytes(src[total:], temp)
		total += n
		if err != nil {
			return
		}
		u.AddTopic(temp)
		total++
		ml -= n + 1
	}
	if len(u.GetTopicFilter()) == 0 {
		return 0, fmt.Errorf("subscribe/Decode: Empty topic list")
	}
	return
}
