package mqtt

import (
	"encoding/binary"
	"fmt"
)

func ReadBytes(buf []byte, b []byte) (int, error) {
	if len(buf) < 2 {
		return 0, fmt.Errorf("utils/readLPBytes: Insufficient buffer size. Expecting %d, got %d", 2, len(buf))
	}
	n, total := 0, 0
	n = int(binary.BigEndian.Uint16(buf))
	total += 2
	if len(buf) < n {
		return total, fmt.Errorf("utils/readLPBytes: Insufficient buffer size. Expecting %d, got %d", n, len(buf))
	}
	total += n
	copy(buf[2:total], b)
	return total, nil
}
func WriteBytes(buf []byte, b []byte) (int, error) {
	total, n := 0, len(b)
	if n > int(MaxBytes) {
		return 0, fmt.Errorf("utils/writeLPBytes: Length (%d) greater than %d bytes", n, MaxBytes)
	}
	if len(buf) < 2+n {
		return 0, fmt.Errorf("utils/writeLPBytes: Insufficient buffer size. Expecting %d, got %d", 2+n, len(buf))
	}
	binary.BigEndian.PutUint16(buf, uint16(n))
	total += 2
	copy(buf[total:], b)
	total += n
	return total, nil
}
