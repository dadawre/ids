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

//var RuleMailSetting *MailSetting

// Action标识符
const (
	PacketPassAction = iota + 1
	PacketRecordAction
	PacketAlertAction
	PacketSendAction
	PacketDropAction
)

func (p *PacketHandler) PassPacket(rule *ruleStructure.Rule) error {
	ReceivePacketCount++
	fmt.Println(ReceivePacketCount)
	return nil
}

func (p *PacketHandler) RecordPacket(rule *ruleStructure.Rule) error {
	PacketTopProtocol := p.decoded[len(p.decoded)-1].String() //顶端协议
	var PacketFileName string

	if p.decoded[1] == layers.LayerTypeIPv4 {
		PacketFileName = p.ipv4.SrcIP.String()
	} else if p.decoded[1] == layers.LayerTypeIPv6 {
		PacketFileName = p.ipv6.SrcIP.String()
	} else if p.decoded[0] == layers.LayerTypeEthernet {
		PacketFileName = strings.Replace(p.eth.SrcMAC.String(), ":", "..", -1)
	}

	folderString := filepath.Join(Initialize.LogPath, PacketFileName)

	err := CreateFile(folderString)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	//fmt.Println(*p.parseMsgStringToValue(rule))
	WriteContentPtr := p.parseMsgStringToValue(rule)
	go func() {
		err := AppendWriteFile(filepath.Join(folderString, PacketTopProtocol), *WriteContentPtr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	ReceivePacketCount++
	fmt.Println("Recording packet")
	return nil
}

func (p *PacketHandler) AlertPacket(rule *ruleStructure.Rule) error {
	WriteContentPtr := p.parseMsgStringToValue(rule)
	go func() {
		err := AppendWriteFile("", *WriteContentPtr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	ReceivePacketCount++
	fmt.Println("Alerting packet")
	return nil
}

func (p *PacketHandler) SendPacket(rule *ruleStructure.Rule) error {
	p.AlertPacket(rule)
	WriteContentPtr := p.parseMsgStringToValue(rule)
	Initialize.Mail.SetBody(*WriteContentPtr)
	err := Initialize.Mail.SendMail()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Sending packet")
	return nil
}

func (p *PacketHandler) DropPacket(rule *ruleStructure.Rule) error {
	fmt.Println("Dropping packet")
	return nil
}
