package action

import "fmt"

var ReceivePacketCount = 0

// PacketAction 定义数据包操作的枚举类型
type PacketAction int

const (
	PacketPassAction PacketAction = iota + 1
	PacketRecordAction
	PacketAlertAction
	PacketSendAction
	PacketDropAction
)

func passPacket(...interface{}) {
	ReceivePacketCount++
}

// 定义每个枚举值对应的函数
func recordPacket(...interface{}) {
	fmt.Println("Recording packet")
}

func alertPacket(...interface{}) {
	fmt.Println("Alerting packet")
}

func sendPacket(...interface{}) {
	fmt.Println("Sending packet")
}

func dropPacket(...interface{}) {
	fmt.Println("Dropping packet")
}

// 创建一个map，将枚举值映射到对应的函数
var packetActionMap = map[PacketAction]func(...interface{}){
	PacketPassAction:   passPacket,
	PacketRecordAction: recordPacket,
	PacketAlertAction:  alertPacket,
	PacketSendAction:   sendPacket,
	PacketDropAction:   dropPacket,
}
