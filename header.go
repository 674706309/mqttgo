package mqtt

import (
	"encoding/binary"
	"fmt"
)

type FixedHeader struct {
	TypeAndFlag     uint8
	RemainingLength []uint8
}

//设置类型
func (fh *FixedHeader) SetType(t uint8) {
	fh.TypeAndFlag = (t << 4) | (fh.TypeAndFlag & 0xf)
}

//获取类型
func (fh *FixedHeader) GetType() uint8 {
	return fh.TypeAndFlag >> 4
}

//设置标志
func (fh *FixedHeader) SetFlag(t uint8) {
	fh.TypeAndFlag = (fh.TypeAndFlag & 0xf0) | (t & 0xf)
}

//获取标志
func (fh *FixedHeader) GetFlag() uint8 {
	return fh.TypeAndFlag & 0xf
}

//设置类型和标志
func (fh *FixedHeader) SetTypeAndFlag(t uint8) {
	fh.TypeAndFlag = t
}

//获取类型和标志
func (fh *FixedHeader) GetTypeAndFlag() uint8 {
	return fh.TypeAndFlag
}

//设置剩余长度
func (fh *FixedHeader) SetRemainingLength(x uint64) {
	binary.PutUvarint(fh.RemainingLength, x)
}

//获取剩余长度
func (fh *FixedHeader) GetRemainingLength() []uint8 {
	return fh.RemainingLength
}

//获取头部长度
func (fh *FixedHeader) getLength() int {
	_, l := binary.Uvarint(fh.GetRemainingLength())
	return l + 1
}
func (fh FixedHeader) String() string {
	return fmt.Sprintf("Type=%q, Flags=%08b, Remaining Length=%d", fh.GetType(), fh.GetFlag(), fh.GetRemainingLength())
}

//头部编码
func (fh *FixedHeader) encode(dst []byte) (int, error) {
	ml := fh.getLength()
	if len(dst) < ml {
		return 0, fmt.Errorf("header/Encode: Insufficient buffer size. Expecting %d, got %d", ml, len(dst))
	}
	total := 0
	l, _ := binary.Uvarint(fh.RemainingLength)
	if l > uint64(maxRemainingLength) || l < 0 {
		return total, fmt.Errorf("header/Encode: Remaining length (%d) out of bound (max %d, min 0)", fh.GetRemainingLength(), maxRemainingLength)
	}
	if !ValidType(fh.GetType()) {
		return total, fmt.Errorf("header/Encode: Invalid message type %d", fh.GetType())
	}
	dst[total] = fh.GetTypeAndFlag()
	total += 1
	n := binary.PutUvarint(dst[total:], l)
	total += n
	return total, nil
}

//头部解码
func (fh *FixedHeader) decode(src []byte) (int, error) {
	total := 0
	fh.SetTypeAndFlag(src[total])
	if !ValidType(fh.GetType()) {
		return total, fmt.Errorf("header/Decode: Invalid message type %d", fh.GetType())
	}
	if fh.GetType() != TYPE_PUBLISH && fh.GetFlag() != DefaultFlags(fh.GetType()) {
		return total, fmt.Errorf("header/Decode: Invalid message (%d) flags. Expecting %d, got %d", fh.GetType(), DefaultFlags(fh.GetType()), fh.GetFlag())
	}
	if fh.GetType() == TYPE_PUBLISH && !ValidQos((fh.GetFlag()>>1)&0x3) {
		return total, fmt.Errorf("header/Decode: Invalid QoS (%d) for PUBLISH message", (fh.GetFlag()>>1)&0x3)
	}
	total++
	ml, m := binary.Uvarint(src[total:])
	total += m
	fh.SetRemainingLength(ml)
	if ml > uint64(maxRemainingLength) || ml < 0 {
		return total, fmt.Errorf("header/Decode: Remaining length (%d) out of bound (max %d, min 0)", ml, maxRemainingLength)
	}
	if int(ml) > len(src[total:]) {
		return total, fmt.Errorf("header/Decode: Remaining length (%d) is greater than remaining buffer (%d)", ml, len(src[total:]))
	}
	return total, nil
}
