package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Publish struct {
	//固定头
	Header Header
	//可变头
	TopicName []byte //主题名
	//有效载荷
	Payload []byte
}

func NewPublish() (c *Publish) {
	c.Header.SetType(TYPE_PUBLISH)
	return
}
func (p Publish) String() string {
	return fmt.Sprintf("%s, Topic=%q, Packet ID=%d, QoS=%d, Retained=%t, Dup=%t, Payload=%t",
		p.Header, p.GetTopicName(), p.Header.GetPacketID(), p.GetQoS(), p.GetRetain(), p.GetDup(), p.GetPayload())
}
func (p *Publish) SetDup(t bool) {
	temp := p.Header.GetTypeAndFlag()
	if t {
		p.Header.SetTypeAndFlag(temp | 0x8) // 00001000
	} else {
		p.Header.SetTypeAndFlag(temp & 247) // 00001000
	}
}
func (p *Publish) GetDup() bool {
	return ((p.Header.GetFlag() >> 3) & 0x1) == 1
}
func (p *Publish) SetRetain(t bool) {
	temp := p.Header.GetTypeAndFlag()
	if t {
		p.Header.SetTypeAndFlag(temp | 0x1) // 00001000
	} else {
		p.Header.SetTypeAndFlag(temp & 254) // 00001000
	}
}
func (p *Publish) GetRetain() bool {
	return (p.Header.GetFlag() & 0x1) == 1
}
func (p *Publish) SetQoS(t byte) error {
	if t != 0x0 && t != 0x1 && t != 0x2 {
		return fmt.Errorf("publish/SetQoS: Invalid QoS %d", t)
	}
	p.Header.SetTypeAndFlag((p.Header.GetTypeAndFlag() & 249) | (t << 1)) // 249 = 11111001
	return nil
}
func (p *Publish) GetQoS() byte {
	return (p.Header.GetFlag() >> 1) & 0x3
}
func (p *Publish) SetTopicName(t []byte) error {
	if !ValidTopicName(t) {
		return fmt.Errorf("publish/SetTopic: Invalid topic name (%s). Must not be empty or contain wildcard characters", string(t))
	}
	p.TopicName = t
	return nil
}
func (p *Publish) GetTopicName() []byte {
	return p.TopicName
}
func (p *Publish) SetPayload(t []byte) {
	p.Payload = t
}
func (p *Publish) GetPayload() []byte {
	return p.Payload
}
func (p *Publish) Length() int {
	total := 2 + len(p.GetTopicName()) + len(p.GetPayload())
	if p.GetQoS() != 0 {
		total += 2
	}
	return total
}
func (p *Publish) encode(dst []byte) (total int, err error) {
	var (
		ml, n int
	)
	if len(p.GetTopicName()) == 0 {
		return 0, fmt.Errorf("publish/Encode: Topic name is empty")
	}
	if len(p.GetPayload()) == 0 {
		return 0, fmt.Errorf("publish/Encode: Payload is empty")
	}
	ml = p.Length()
	p.Header.SetRemainingLength(uint64(ml))
	total = 0
	n, err = p.Header.encode(dst[total:])
	total += n
	if err != nil {
		return
	}
	n, err = WriteBytes(dst[total:], p.GetTopicName())
	total += n
	if err != nil {
		return
	}
	if p.GetQoS() != 0 {
		binary.BigEndian.PutUint16(dst[total:], p.Header.GetPacketID())
		total += 2
	}
	copy(dst[total:], p.GetPayload())
	total += len(p.GetPayload())
	return total, nil
}
func (p *Publish) decode(src []byte) (total int, err error) {
	var (
		temp     []byte
		l, hl, n int
	)
	total = 0
	hl, err = p.Header.decode(src[total:])
	total += hl
	if err != nil {
		return
	}
	n, err = ReadBytes(src[total:], temp)
	p.SetTopicName(temp)
	total += n
	if err != nil {
		return
	}
	if !ValidTopicName(p.GetTopicName()) {
		return total, fmt.Errorf("publish/Decode: Invalid topic name (%s). Must not be empty or contain wildcard characters", string(p.GetTopicName()))
	}
	if p.GetQoS() != 0 {
		p.Header.SetPacketID(binary.BigEndian.Uint16(src[total:]))
		total += 2
	}
	l = int(p.Header.GetRemainingLength()) - (total - hl)
	p.SetPayload(src[total : total+l])
	total += l
	return
}
