package Handler

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"graduationDesign/back_end/Action"
	Structure2 "graduationDesign/back_end/Initialize/InitializeStructure"
	"graduationDesign/back_end/Initialize/Parameter"
	"graduationDesign/back_end/Structure"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

var (
	ReceivePacketCount int32
)

// Action标识符
const (
	PacketPassAction = iota + 1
	PacketRecordAction
	PacketAlertAction
	PacketSendAction
	PacketDropAction
)

// 增加线程安全的递增函数
func incrementPacketCount() {
	atomic.AddInt32(&ReceivePacketCount, 1)
	fmt.Println("当前接收数据包个数:", atomic.LoadInt32(&ReceivePacketCount))
}

func (p *PacketHandler) GetSrcAndDstAddrString() (string, string) {
	if p.decoded[1] == layers.LayerTypeIPv4 {
		return p.ipv4.SrcIP.String(), p.ipv4.DstIP.String()
	} else if p.decoded[1] == layers.LayerTypeIPv6 {
		return p.ipv6.SrcIP.String(), p.ipv4.DstIP.String()
	} else if p.decoded[0] == layers.LayerTypeEthernet {
		return strings.Replace(p.eth.SrcMAC.String(), ":", "..", -1), strings.Replace(p.eth.DstMAC.String(), ":", "..", -1)
	}
	return "", ""
}

// 记录Action的日志
func recordActionLog(srcAddr, dstAddr string, action int8, message string) Structure.ActionDescribeTriggerLog {
	return Structure.ActionDescribeTriggerLog{
		SrcAddr:         srcAddr,
		DstAddr:         dstAddr,
		ActionType:      action,
		DescribeMessage: message,
	}
}
func (p *PacketHandler) defaultPass() {
	incrementPacketCount()
}

func (p *PacketHandler) PassPacket(rule *Structure.RuleNodeValue) error {
	// 写入日志内容
	WriteContentPtr := p.parseMsgStringToValue(rule)
	ContentString := fmt.Sprintf("时间: %s\n", time.Now().Format(time.RFC3339)) + *WriteContentPtr

	srcAddr, dstAddr := p.GetSrcAndDstAddrString()
	ActionDescribeLog := recordActionLog(srcAddr, dstAddr, rule.Action, ContentString)

	Structure2.ActionDescribeTriggerLogArray = append(Structure2.ActionDescribeTriggerLogArray, ActionDescribeLog)
	incrementPacketCount()
	return nil
}

func (p *PacketHandler) RecordPacket(rule *Structure.RuleNodeValue) error {
	// 写入日志内容
	WriteContentPtr := p.parseMsgStringToValue(rule)
	ContentString := fmt.Sprintf("时间: %s\n", time.Now().Format(time.RFC3339)) + *WriteContentPtr

	srcAddr, dstAddr := p.GetSrcAndDstAddrString()
	ActionDescribeLog := recordActionLog(srcAddr, dstAddr, rule.Action, ContentString)

	//顶端协议
	PacketTopProtocol := p.decoded[len(p.decoded)-1].String()

	PacketFileName := ActionDescribeLog.SrcAddr
	folderString := filepath.Join(Parameter.LogPath, PacketFileName)

	err := CreateFile(folderString)
	if err != nil {
		fmt.Println("Error creating file:", err.Error())
	}

	go func() {
		err := AppendWriteFile(filepath.Join(folderString, PacketTopProtocol), "\n"+ContentString)
		if err != nil {
			fmt.Println("Error appending to file:", err.Error())
		}

		Structure2.ActionDescribeTriggerLogArray = append(Structure2.ActionDescribeTriggerLogArray, ActionDescribeLog)
		incrementPacketCount()
		fmt.Println("Recording packet")
	}()

	return nil
}

func (p *PacketHandler) AlertPacket(rule *Structure.RuleNodeValue) error {
	WriteContentPtr := p.parseMsgStringToValue(rule)
	ContentString := fmt.Sprintf("时间: %s\n", time.Now().Format(time.RFC3339)) + *WriteContentPtr
	srcAddr, dstAddr := p.GetSrcAndDstAddrString()
	ActionDescribeLog := recordActionLog(srcAddr, dstAddr, rule.Action, ContentString)

	go func() {
		err := AppendWriteFile(Parameter.AlertPath, "\n"+ContentString)
		if err != nil {
			fmt.Println("Error appending to alert file:", err.Error())
		}

		Structure2.ActionDescribeTriggerLogArray = append(Structure2.ActionDescribeTriggerLogArray, ActionDescribeLog)
		incrementPacketCount()
		fmt.Println("Alerting packet")
	}()
	return nil
}

func (p *PacketHandler) SendPacket(rule *Structure.RuleNodeValue) error {
	WriteContentPtr := p.parseMsgStringToValue(rule)
	srcAddr, dstAddr := p.GetSrcAndDstAddrString()
	ContentString := fmt.Sprintf("时间: %s\n", time.Now().Format(time.RFC3339)) + *WriteContentPtr
	ActionDescribeLog := recordActionLog(srcAddr, dstAddr, rule.Action, ContentString)

	go func() {
		err := AppendWriteFile(Parameter.AlertPath, "\n"+ContentString)
		if err != nil {
			fmt.Println("Error appending to alert file:", err.Error())
		}

		Structure2.ActionDescribeTriggerLogArray = append(Structure2.ActionDescribeTriggerLogArray, ActionDescribeLog)
		incrementPacketCount()
		fmt.Println("Alerting packet")
	}()

	go func() {
		Structure2.Mail.SetSubject(fmt.Sprintf("ids id:%d, ids name:%s", Parameter.IdsID, Parameter.IdsName))
		Structure2.Mail.SetBody(ContentString)
		err := Structure2.Mail.SendMail()
		if err != nil {
			fmt.Println("Error sending mail:", err.Error())
		}
		fmt.Println("Sending packet")
	}()

	return nil
}

// DropIPv4Packets 封IP
func (p *PacketHandler) DropIPv4Packets(rule *Structure.RuleNodeValue) error {
	WriteContentPtr := p.parseMsgStringToValue(rule)
	srcAddr, dstAddr := p.GetSrcAndDstAddrString()
	ContentString := fmt.Sprintf("时间: %s\n", time.Now().Format(time.RFC3339)) + *WriteContentPtr
	ActionDescribeLog := recordActionLog(srcAddr, dstAddr, rule.Action, ContentString)

	go func() {
		err := AppendWriteFile(Parameter.AlertPath, "\n"+ContentString)
		if err != nil {
			fmt.Println("Error appending to alert file:", err.Error())
		}

		// Append to the log array
		Structure2.ActionDescribeTriggerLogArray = append(Structure2.ActionDescribeTriggerLogArray, ActionDescribeLog)
		incrementPacketCount()
		fmt.Println("Alerting packet")
	}()

	err := Action.DropIP(p.ipv4.SrcIP.String())
	if err != nil {
		return err
	}

	Structure2.MainBanIPSet.Insert(Structure.BanIPDescribe{
		BanIP:     p.ipv4.SrcIP.String(),
		BanRuleID: rule.ID,
		Describe:  ContentString,
		UnbanTime: time.Now().Add(time.Duration(rule.Option.BanTime) * time.Minute),
	})
	return nil
}
