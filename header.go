package mqtt

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	maxRemainingLength int32 = math.MaxInt32 // bytes, or 256 MB
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

//获取所有固定头标志
func (fh *FixedHeader) GetFlag() uint8 {
	return fh.TypeAndFlag & 0xf
}
func (fh *FixedHeader) GetTypeAndFlag() uint8 {
	return fh.TypeAndFlag
}
func (fh *FixedHeader) SetRemainingLength(x uint64) {
	binary.PutUvarint(fh.RemainingLength, x)
}
func (fh *FixedHeader) GetRemainingLength() []uint8 {
	return fh.RemainingLength
}
func (fh *FixedHeader) getLength() int {
	_, l := binary.Uvarint(fh.GetRemainingLength())
	return l + 1
}
func (fh FixedHeader) String() string {
	return fmt.Sprintf("Type=%q, Flags=%08b, Remaining Length=%d", fh.GetType(), fh.GetFlag(), fh.GetRemainingLength())
}
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

	if fh.GetType() > TYPE_RESERVED && fh.GetType() < TYPE_RESERVED2 {
		return total, fmt.Errorf("header/Encode: Invalid message type %d", fh.GetType())
	}

	dst[total] = fh.GetTypeAndFlag()
	total += 1
	n := binary.PutUvarint(dst[total:], l)
	total += n

	return total, nil
}
func (fh *FixedHeader) decode(src []byte) (int, error) {

}
