package mqtt

import "math"

const (
	maxRemainingLength int32 = math.MaxInt32 // 剩余字段最大长度
)

func ValidType(t uint8) bool {
	return t > TYPE_RESERVED && t < TYPE_RESERVED2
}
func DefaultFlags(t uint8) uint8 {
	switch t {
	case TYPE_RESERVED:
		return 0
	case TYPE_CONNECT:
		return 0
	case TYPE_CONNACK:
		return 0
	case TYPE_PUBLISH:
		return 0
	case TYPE_PUBACK:
		return 0
	case TYPE_PUBREC:
		return 0
	case TYPE_PUBREL:
		return 2
	case TYPE_PUBCOMP:
		return 0
	case TYPE_SUBSCRIBE:
		return 2
	case TYPE_SUBACK:
		return 0
	case TYPE_UNSUBSCRIBE:
		return 2
	case TYPE_UNSUBACK:
		return 0
	case TYPE_PINGREQ:
		return 0
	case TYPE_PINGRESP:
		return 0
	case TYPE_DISCONNECT:
		return 0
	case TYPE_RESERVED2:
		return 0
	}

	return 0
}

const (
	QosAtMostOnce byte = iota
	QosAtLeastOnce
	QosExactlyOnce
)

func ValidQos(qos byte) bool {
	return qos == QosAtMostOnce || qos == QosAtLeastOnce || qos == QosExactlyOnce
}
