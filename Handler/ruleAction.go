package Handler

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"graduationDesign/Initialize"
	"graduationDesign/structure/ruleStructure"
	"path/filepath"
	"strings"
)

var ReceivePacketCount = 0

var RuleMailSetting *MailSetting

// Action标识符
const (
	PacketPassAction = iota + 1
	PacketRecordAction
	PacketAlertAction
	PacketSendAction
	PacketDropAction
)

func (p *PacketHandler) PassPacket(rule *ruleStructure.Rule) {
	ReceivePacketCount++
	fmt.Println(ReceivePacketCount)
}

func (p *PacketHandler) RecordPacket(rule *ruleStructure.Rule) {
	PacketTopProtocol := p.decoded[len(p.decoded)-1].String() //顶端协议
	var PacketFileName string

	if p.decoded[1] == layers.LayerTypeIPv4 {
		PacketFileName = p.ipv4.SrcIP.String()
	} else if p.decoded[1] == layers.LayerTypeIPv6 {
		PacketFileName = p.ipv6.SrcIP.String()
	} else if p.decoded[0] == layers.LayerTypeEthernet {
		PacketFileName = strings.Replace(p.eth.SrcMAC.String(), ":", "..", -1)
	}

	folderString := filepath.Join(LogPath, PacketFileName)

	err := CreateFile(folderString)
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Println(*p.parseMsgStringToString(rule))
	WriteContentPtr := p.parseMsgStringToString(rule)
	go func() {
		err := AppendWriteFile(filepath.Join(folderString, PacketTopProtocol), *WriteContentPtr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	ReceivePacketCount++
	fmt.Println("Recording packet")
}

func (p *PacketHandler) AlertPacket(rule *ruleStructure.Rule) {
	WriteContentPtr := p.parseMsgStringToString(rule)
	go func() {
		err := AppendWriteFile(Initialize.AlertPath, *WriteContentPtr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	ReceivePacketCount++
	fmt.Println("Alerting packet")
}

func (p *PacketHandler) SendPacket(rule *ruleStructure.Rule) {
	p.AlertPacket(rule)
	WriteContentPtr := p.parseMsgStringToString(rule)
	RuleMailSetting.SetBody(*WriteContentPtr)
	err := RuleMailSetting.SendMail()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Sending packet")
}

func (p *PacketHandler) DropPacket(rule *ruleStructure.Rule) {
	fmt.Println("Dropping packet")
}
