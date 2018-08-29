package mqtt

import (
	"bytes"
	"math"
)

const (
	MaxRemainingLength int32  = math.MaxInt32 // 剩余字段最大长度
	MaxBytes           uint16 = math.MaxInt16
)
const (
	TYPE_RESERVED uint8 = iota
	TYPE_CONNECT
	TYPE_CONNACK
	TYPE_PUBLISH
	TYPE_PUBACK
	TYPE_PUBREC
	TYPE_PUBREL
	TYPE_PUBCOMP
	TYPE_SUBSCRIBE
	TYPE_SUBACK
	TYPE_UNSUBSCRIBE
	TYPE_UNSUBACK
	TYPE_PINGREQ
	TYPE_PINGRESP
	TYPE_DISCONNECT
	TYPE_RESERVED2
)
const (
	TYPE_FLAG_CONNECT     uint8 = TYPE_CONNECT << 4
	TYPE_FLAG_CONNACK     uint8 = TYPE_CONNACK << 4
	TYPE_FLAG_PUBACK      uint8 = TYPE_PUBACK << 4
	TYPE_FLAG_PUBREC      uint8 = TYPE_PUBREC << 4
	TYPE_FLAG_PUBREL      uint8 = TYPE_PUBREL<<4 | 0x2
	TYPE_FLAG_PUBCOMP     uint8 = TYPE_PUBCOMP << 4
	TYPE_FLAG_SUBSCRIBE   uint8 = TYPE_SUBSCRIBE<<4 | 0x2
	TYPE_FLAG_SUBACK      uint8 = TYPE_SUBACK << 4
	TYPE_FLAG_UNSUBSCRIBE uint8 = TYPE_UNSUBSCRIBE<<4 | 0x2
	TYPE_FLAG_UNSUBACK    uint8 = TYPE_UNSUBACK << 4
	TYPE_FLAG_PINGREQ     uint8 = TYPE_PINGREQ << 4
	TYPE_FLAG_PINGRESP    uint8 = TYPE_PINGRESP << 4
	TYPE_FLAG_DISCONNECT  uint8 = TYPE_DISCONNECT << 4
)
const (
	PROTOCOL             = "MQTT"
	PROTOCOL_LEVEL uint8 = 0x4
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
	QosFailure = 0x80
)

func ValidQos(qos byte) bool {
	return qos == QosAtMostOnce || qos == QosAtLeastOnce || qos == QosExactlyOnce
}

var SupportedVersions = map[byte]string{
	0x4: "MQTT",
}

const (
	CONNBAK_RETURN_CODE_OK             uint8 = iota //连接已接受 连接已被服务端接受
	CONNBAK_RETURN_NO_SUPPORT_PROTOCOL              //连接已拒绝，不支持的协议版本 服务端不支持客户端请求的 MQTT 协议级别
	CONNBAK_RETURN_NO_CLIENT_ID                     //连接已拒绝，不合格的客户端标识符 客户端标识符是正确的 UTF-8 编码，但服务端不允许使用
	CONNBAK_RETURN_NO_SERVER                        //连接已拒绝，服务端不可用 网络连接已建立，但 MQTT 服务不可用
	CONNBAK_RETURN_ERROR_UNAME_PWD                  //连接已拒绝，无效的用户名或密码 用户名或密码的数据格式无效
	CONNBAK_RETURN_UNAUTHORIZED                     // 连接已拒绝，未授权 客户端未被授权连接到此服务器
	CONNBAK_RETURN_RESERVED                         // 连接已拒绝，保留码 6-255

)

func ValidTopicName(topic []byte) bool {
	return len(topic) > 0 && bytes.IndexByte(topic, '#') == -1 && bytes.IndexByte(topic, '+') == -1
}
