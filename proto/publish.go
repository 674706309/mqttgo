package mqttgo

import (
	"encoding/binary"
	"fmt"
)

type Publish struct {
	//固定头
	header
	//可变头
	TopicName []byte //主题名
	//有效载荷
	Payload []byte
}

func NewPublish() (c *Publish) {
	c = &Publish{}
	c.header.SetType(TYPE_PUBLISH)
	return
}
func (p Publish) String() string {
	return fmt.Sprintf("%s, Topic=%s, PacketID=%d, QoS=%d, Retained=%t, Dup=%t, Payload=%d",
		p.header, p.GetTopicName(), p.header.GetPacketID(), p.GetQoS(), p.GetRetain(), p.GetDup(), p.GetPayload())
}
func (p *Publish) SetDup(t bool) {
	temp := p.header.GetTypeAndFlag()
	if t {
		p.header.SetTypeAndFlag(temp | 0x8) // 00001000
	} else {
		p.header.SetTypeAndFlag(temp & 247) // 00001000
	}
}
func (p *Publish) GetDup() bool {
	return ((p.header.GetFlag() >> 3) & 0x1) == 1
}
func (p *Publish) SetRetain(t bool) {
	temp := p.header.GetTypeAndFlag()
	if t {
		p.header.SetTypeAndFlag(temp | 0x1) // 00001000
	} else {
		p.header.SetTypeAndFlag(temp & 254) // 00001000
	}
}
func (p *Publish) GetRetain() bool {
	return (p.header.GetFlag() & 0x1) == 1
}
func (p *Publish) SetQoS(t byte) error {
	if t != 0x0 && t != 0x1 && t != 0x2 {
		return fmt.Errorf("publish/SetQoS: Invalid QoS %d", t)
	}
	p.header.SetTypeAndFlag((p.header.GetTypeAndFlag() & 249) | (t << 1)) // 249 = 11111001
	return nil
}
func (p *Publish) GetQoS() byte {
	return (p.header.GetFlag() >> 1) & 0x3
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
	return p.header.Length() + p.GetMessageLength()
}
func (p *Publish) GetMessageLength() int {
	total := 2 + len(p.GetTopicName()) + len(p.GetPayload())
	if p.GetQoS() != 0 {
		total += 2
	}
	return total
}
func (p *Publish) Encode(dst []byte) (total int, err error) {
	var (
		ml, n int
	)
	if len(p.GetTopicName()) == 0 {
		return 0, fmt.Errorf("publish/Encode: Topic name is empty")
	}
	if len(p.GetPayload()) == 0 {
		return 0, fmt.Errorf("publish/Encode: Payload is empty")
	}
	ml = p.GetMessageLength()
	p.SetRemainingLength(uint64(ml))
	total = 0
	n, err = p.encode(dst[total:])
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
		binary.BigEndian.PutUint16(dst[total:], p.header.GetPacketID())
		total += 2
	}
	copy(dst[total:], p.GetPayload())
	total += len(p.GetPayload())
	return total, nil
}
func (p *Publish) Decode(src []byte) (total int, err error) {
	var (
		temp     []byte
		l, hl, n int
	)
	total = 0
	if hl, err = p.decode(src[total:]); err != nil {
		return
	}
	total += hl
	if temp, n, err = ReadBytes(src[total:]); err != nil {
		return
	}
	if err = p.SetTopicName(temp); err != nil {
		return
	}
	total += n
	if p.GetQoS() != 0 {
		p.packetID = binary.BigEndian.Uint16(src[total:])
		total += 2
	}
	l = int(p.GetRemainingLength()) - (total - hl)
	p.Payload = src[total : total+l]
	total += l
	return
}
