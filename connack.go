package mqtt

import "fmt"

type Connbak struct {
	//固定头
	FixedHeader
	//可变头
	ConnectAcknowledgeFlags bool
	SessionPresent          bool  //当前会话标志
	ConnectReturncode       uint8 //返回码
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

func NewConnack() (msg *Connect) {
	msg.TypeAndFlag = byte(TYPE_FLAG_CONNACK)
	msg.SetRemainingLength(2)
	return
}

//设置当前会话标志
func (c *Connbak) SetSessionPresent(sessionPresent bool) {
	c.SessionPresent = sessionPresent
}

//获取当前会话标志
func (c *Connbak) GetSessionPresent() bool {
	return c.SessionPresent
}
func (c *Connbak) SetReturnCode(code uint8) {
	c.ConnectReturncode = code
}
func (c *Connbak) GetReturnCode() uint8 {
	return c.ConnectReturncode
}
func (c *Connbak) decode(header byte, msg []byte) error {

	if header != TYPE_FLAG_CONNACK {
		return fmt.Errorf("CONNACK固定头信息保留位错误:%v", header)
	}

	c.RemainingLength = 2
	if (msg[0] | 1) != 1 {
		return fmt.Errorf("CONNACK连接确认标志错误:收到%d 应收到0或1", msg[0])
	}
	c.SessionPresent = msg[0] == 1
	c.ConnectReturncode = uint8(msg[1])
	return nil
}
func (c *Connbak) encode() []byte {
	//固定报头
	data := make([]byte, 4, 4)
	data[0] = c.TypeAndFlag
	data[1] = byte(0x2)
	if c.SessionPresent {
		data[2] = 0x1
	}
	data[3] = c.ConnectReturncode
	c.RemainingLength = 2
	return data
}
