package mqtt

import (
	"encoding/binary"
	"fmt"
)

type Header struct {
	TypeAndFlag     uint8
	RemainingLength []uint8

	PacketID uint16 //报文标识符 部分拥有
}

//设置类型
func (h *Header) SetType(t uint8) {
	h.TypeAndFlag = (t << 4) | (h.TypeAndFlag & 0xf)
}

//获取类型
func (h *Header) GetType() uint8 {
	return h.TypeAndFlag >> 4
}

//设置标志
func (h *Header) SetFlag(t uint8) {
	h.TypeAndFlag = (h.TypeAndFlag & 0xf0) | (t & 0xf)
}

//获取标志
func (h *Header) GetFlag() uint8 {
	return h.TypeAndFlag & 0xf
}

//设置类型和标志
func (h *Header) SetTypeAndFlag(t uint8) {
	h.TypeAndFlag = t
}

//获取类型和标志
func (h *Header) GetTypeAndFlag() uint8 {
	return h.TypeAndFlag
}

//设置剩余长度
func (h *Header) SetRemainingLength(x uint64) {
	binary.PutUvarint(h.RemainingLength, x)
}

//获取剩余长度
func (h *Header) GetRemainingLength() []uint8 {
	return h.RemainingLength
}
func (p *Header) SetPacketID(t uint16) {
	p.PacketID = t
}
func (p *Header) GetPacketID() uint16 {
	return p.PacketID
}

//获取头部长度
func (h *Header) getLength() int {
	_, l := binary.Uvarint(h.GetRemainingLength())
	return l + 1
}
func (h Header) String() string {
	return fmt.Sprintf("Type=%q, Flags=%08b, Remaining Length=%d", h.GetType(), h.GetFlag(), h.GetRemainingLength())
}

//头部编码
func (h *Header) encode(dst []byte) (int, error) {
	ml := h.getLength()
	if len(dst) < ml {
		return 0, fmt.Errorf("header/Encode: Insufficient buffer size. Expecting %d, got %d", ml, len(dst))
	}
	total := 0
	l, _ := binary.Uvarint(h.RemainingLength)
	if l > uint64(MaxRemainingLength) || l < 0 {
		return total, fmt.Errorf("header/Encode: Remaining length (%d) out of bound (max %d, min 0)", h.GetRemainingLength(), MaxRemainingLength)
	}
	if !ValidType(h.GetType()) {
		return total, fmt.Errorf("header/Encode: Invalid message type %d", h.GetType())
	}
	dst[total] = h.GetTypeAndFlag()
	total += 1
	n := binary.PutUvarint(dst[total:], l)
	total += n
	return total, nil
}

//头部解码
func (h *Header) decode(src []byte) (int, error) {
	total := 0
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
	ml, m := binary.Uvarint(src[total:])
	total += m
	h.SetRemainingLength(ml)
	if ml > uint64(MaxRemainingLength) || ml < 0 {
		return total, fmt.Errorf("header/Decode: Remaining length (%d) out of bound (max %d, min 0)", ml, MaxRemainingLength)
	}
	if int(ml) > len(src[total:]) {
		return total, fmt.Errorf("header/Decode: Remaining length (%d) is greater than remaining buffer (%d)", ml, len(src[total:]))
	}
	return total, nil
}
