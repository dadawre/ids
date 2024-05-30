package handler

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"graduationDesign/structure/ruleStructure"
	"path/filepath"
)

var ReceivePacketCount = 0

var RuleMailSetting *MailSetting

const (
	PacketPassAction = iota + 1
	PacketRecordAction
	PacketAlertAction
	PacketSendAction
	PacketDropAction
)

func (p *PacketHandler) PassPacket(rule *ruleStructure.Rule) {
	ReceivePacketCount++
	fmt.Println("Passing packet  ", ReceivePacketCount)
}

// RecordPacket 定义每个枚举值对应的函数
func (p *PacketHandler) RecordPacket(rule *ruleStructure.Rule) {
	PacketTopProtocol := p.decoded[len(p.decoded)-1].String()
	var PacketID string

	if p.decoded[1] == layers.LayerTypeIPv4 {
		PacketID = p.ipv4.SrcIP.String()
	} else if p.decoded[1] == layers.LayerTypeIPv6 {
		PacketID = p.ipv6.SrcIP.String()
	} else if p.decoded[0] == layers.LayerTypeEthernet {
		PacketID = p.eth.SrcMAC.String()
	}

	folderString := filepath.Join(LogPath, PacketID)

	err := CreateFile(folderString)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = AppendWriteFile(filepath.Join(folderString, PacketTopProtocol), rule.Msg+"\n")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Recording packet")
}

func (p *PacketHandler) AlertPacket(rule *ruleStructure.Rule) {
	p.RecordPacket(rule)
	RuleMailSetting.SetBody(rule.Msg)
	err := RuleMailSetting.SendMail()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Alerting packet")
}

func (p *PacketHandler) SendPacket(rule *ruleStructure.Rule) {
	fmt.Println("Sending packet")
}

func (p *PacketHandler) DropPacket(rule *ruleStructure.Rule) {
	fmt.Println("Dropping packet")
}
