package ruleAction

import (
	"fmt"
	"graduationDesign/structure/ruleStructure"
)

var ReceivePacketCount = 0

const (
	PacketPassAction = iota + 1
	PacketRecordAction
	PacketAlertAction
	PacketSendAction
	PacketDropAction
)

func PassPacket(rule *ruleStructure.Rule) {
	ReceivePacketCount++
	fmt.Println("Passing packet  ", ReceivePacketCount)
}

// 定义每个枚举值对应的函数
func recordPacket(rule *ruleStructure.Rule) {
	fmt.Println("Recording packet")
}

func alertPacket(rule *ruleStructure.Rule) {
	fmt.Println("Alerting packet")
}

func sendPacket(rule *ruleStructure.Rule) {
	fmt.Println("Sending packet")
}

func dropPacket(rule *ruleStructure.Rule) {
	fmt.Println("Dropping packet")
}

// 创建一个map，将枚举值映射到对应的函数
var PacketActionMap = map[int8]func(rule *ruleStructure.Rule){
	PacketPassAction:   PassPacket,
	PacketRecordAction: recordPacket,
	PacketAlertAction:  alertPacket,
	PacketSendAction:   sendPacket,
	PacketDropAction:   dropPacket,
}
