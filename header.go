package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Header struct {
	TypeAndFlag     byte
	RemainingLength []byte

	PacketID uint16 //报文标识符 部分拥有
}

func (h Header) String() string {
	return fmt.Sprintf("Type=%08d, Flags=%08d, RemainingLength=%d", h.GetType(), h.GetFlag(), h.GetRemainingLength())
}

//设置类型
func (h *Header) SetType(t byte) (err error) {
	if !ValidType(t) {
		return fmt.Errorf("header/SetType: Invalid control packet type %d", t)
	}
	h.TypeAndFlag = (t << 4) | (h.TypeAndFlag & 0xf)
	h.SetFlag(DefaultFlags(h.GetType()))
	return
}

//获取类型
func (h *Header) GetType() byte {
	return h.TypeAndFlag >> 4
}

//设置标志
func (h *Header) SetFlag(t byte) (err error) {
	if !(DefaultFlags(h.GetType()) == t) {
		return fmt.Errorf("Flag error")
	}
	h.TypeAndFlag = (t & 0xf) | (h.TypeAndFlag & 0xf0)
	return
}

//获取标志
func (h *Header) GetFlag() byte {
	return h.TypeAndFlag & 0xf
}

//设置类型和标志
func (h *Header) SetTypeAndFlag(t byte) {
	h.TypeAndFlag = t
}

//获取类型和标志
func (h *Header) GetTypeAndFlag() byte {
	return h.TypeAndFlag
}

//设置剩余长度
func (h *Header) SetRemainingLength(t uint64) (err error) {
	if t > uint64(MaxRemainingLength) || t < 0 {
		return fmt.Errorf("header/SetLength: Remaining length (%d) out of bound (max %d, min 0)", t, MaxRemainingLength)
	}
	h.RemainingLength = make([]byte, 8)
	binary.PutUvarint(h.RemainingLength, t)
	return
}

//获取剩余长度
func (h *Header) GetRemainingLength() uint64 {
	l, _ := binary.Uvarint(h.RemainingLength)
	return l
}

//设置报文标识符
func (p *Header) SetPacketID(t uint16) {
	p.PacketID = t
}

//获取报文标识符
func (p *Header) GetPacketID() uint16 {
	return p.PacketID
}

//获取头部长度
func (h *Header) Length() int {
	return 1 + int(h.GetRemainingLength())
}

//头部编码
func (h *Header) encode(dst []byte) (total int, err error) {
	var (
		l uint64
		n int
	)
	total = 0
	l = h.GetRemainingLength()
	//if l > uint64(MaxRemainingLength) || l < 0 {
	//	return total, fmt.Errorf("header/Encode: Remaining length (%d) out of bound (max %d, min 0)", h.GetRemainingLength(), MaxRemainingLength)
	//}
	//if !ValidType(h.GetType()) {
	//	return total, fmt.Errorf("header/Encode: Invalid message type %d", h.GetType())
	//}
	dst[total] = h.GetTypeAndFlag()
	total += 1
	n = binary.PutUvarint(dst[total:], l)
	total += n
	return
}

//头部解码
func (h *Header) decode(src []byte) (total int, err error) {
	var (
		n  int
		ml uint64
	)
	total = 0
	h.SetTypeAndFlag(src[total])
	if !ValidType(h.GetType()) {
		return total, fmt.Errorf("header/Decode: Invalid message type %d", h.GetType())
	}
	if h.GetType() != TYPE_PUBLISH && h.GetFlag() != DefaultFlags(h.GetType()) {
		return total, fmt.Errorf("header/Decode: Invalid message (%d) flags. Expecting %d, got %d", h.GetType(), DefaultFlags(h.GetType()), h.GetFlag())
	}
	if h.GetType() == TYPE_PUBLISH && !ValidQos((h.GetFlag()>>1)&0x3) {
		return total, fmt.Errorf("header/Decode: Invalid QoS (%d) for PUBLISH message", (h.GetFlag()>>1)&0x3)
	}
	total++
	ml, n = binary.Uvarint(src[total:])
	total += n
	h.SetRemainingLength(ml)
	if ml > uint64(MaxRemainingLength) || ml < 0 {
		return total, fmt.Errorf("header/Decode: Remaining length (%d) out of bound (max %d, min 0)", ml, MaxRemainingLength)
	}
	if int(ml) > len(src[total:]) {
		return total, fmt.Errorf("header/Decode: Remaining length (%d) is greater than remaining buffer (%d)", ml, len(src[total:]))
	}
	return
}
