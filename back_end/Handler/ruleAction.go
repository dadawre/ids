package Handler

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"graduationDesign/back_end/Action"
	"graduationDesign/back_end/DB"
	"graduationDesign/back_end/Initialize/InitializeStructure"
	"graduationDesign/back_end/Initialize/Parameter"
	"graduationDesign/back_end/Structure"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	ReceivePacketCount int32
)

const (
	PacketPassAction = iota + 1
	PacketRecordAction
	PacketAlertAction
	PacketSendAction
	PacketDropAction
)

func incrementPacketCount() {
	count := atomic.AddInt32(&ReceivePacketCount, 1)
	if count%100 == 0 { // 每100个包打印一次计数，减少IO压力
		fmt.Println("当前接收数据包个数:", count)
	}
}

func (p *PacketHandler) GetSrcAndDstAddrString() (string, string) {
	switch p.decoded[1] {
	case layers.LayerTypeIPv4:
		return p.ipv4.SrcIP.String(), p.ipv4.DstIP.String()
	case layers.LayerTypeIPv6:
		return p.ipv6.SrcIP.String(), p.ipv6.DstIP.String()
	case layers.LayerTypeARP:
		// 关键修复：将字节切片转换为 net.IP 再转为字符串
		srcIP := net.IP(p.arp.SourceProtAddress)
		dstIP := net.IP(p.arp.DstProtAddress)
		if len(srcIP) > 0 && len(dstIP) > 0 {
			return srcIP.String(), dstIP.String()
		}
		srcHw := net.HardwareAddr(p.arp.SourceHwAddress).String()
		dstHw := net.HardwareAddr(p.arp.DstHwAddress).String()
		return srcHw, dstHw
	case layers.LayerTypeEthernet:
		src := strings.Replace(p.eth.SrcMAC.String(), ":", "..", -1)
		dst := strings.Replace(p.eth.DstMAC.String(), ":", "..", -1)
		return src, dst
	default:
		return "", ""
	}
}

func (p *PacketHandler) GetSrcAndDstPortAddrString() (string, string) {
	if len(p.decoded) <= 2 {
		return "", ""
	}
	switch p.decoded[2] {
	case layers.LayerTypeTCP:
		return p.tcp.SrcPort.String(), p.tcp.DstPort.String()
	case layers.LayerTypeUDP:
		return p.udp.SrcPort.String(), p.udp.DstPort.String()
	default:
		return "", ""
	}
}

func (p *PacketHandler) GetTcpFlag() string {
	if len(p.decoded) >= 3 && p.decoded[2] == layers.LayerTypeTCP {
		flags := []bool{p.tcp.URG, p.tcp.ACK, p.tcp.PSH, p.tcp.RST, p.tcp.SYN, p.tcp.FIN}
		b := make([]byte, len(flags))
		for i, flag := range flags {
			b[i] = '0'
			if flag {
				b[i] = '1'
			}
		}
		return string(b)
	}
	return ""
}

func (p *PacketHandler) GetTopProtocol() string {
	if len(p.decoded) == 0 {
		return "unknown"
	}
	return p.decoded[len(p.decoded)-1].String()
}

func (p *PacketHandler) ConductionActionLog(rule *Structure.RuleNodeValue) (*DB.ActionDescribeTriggerLog, string) {
	writeContent := p.parseMsgStringToValue(rule)
	dataString := time.Now().Format(time.RFC3339)
	contentString := fmt.Sprintf("时间: %s\n%s", dataString, *writeContent)

	srcIP, dstIP := p.GetSrcAndDstAddrString()
	srcPort, dstPort := p.GetSrcAndDstPortAddrString()

	return &DB.ActionDescribeTriggerLog{
		SrcIPAddr:       srcIP,
		DstIPAddr:       dstIP,
		SrcPortAddr:     srcPort,
		DstPortAddr:     dstPort,
		TopProtocol:     p.GetTopProtocol(),
		RuleID:          rule.ID,
		RuleName:        rule.RuleName,
		DescribeMessage: *writeContent,
		TriggerTime:     dataString,
		ActionType:      rule.Action,
		TcpFlag:         p.GetTcpFlag(),
		PacketLength:    p.packetLen,
	}, contentString
}

func (p *PacketHandler) defaultPass() {
	incrementPacketCount()
}

func (p *PacketHandler) logAndIncrement(log *DB.ActionDescribeTriggerLog) {
	go func(l *DB.ActionDescribeTriggerLog) {
		if err := InitializeStructure.LogSQLManager.InsertRecord(Parameter.SqlRecordTableFullName, l); err != nil {
			fmt.Printf("数据库插入错误: %v\n", err)
		}
		incrementPacketCount()
	}(log)
}

func (p *PacketHandler) PassPacket(rule *Structure.RuleNodeValue) error {
	log, _ := p.ConductionActionLog(rule)
	p.logAndIncrement(log)
	return nil
}

func (p *PacketHandler) RecordPacket(rule *Structure.RuleNodeValue) error {
	log, content := p.ConductionActionLog(rule)

	go func(c string, l *DB.ActionDescribeTriggerLog) {
		folder := filepath.Join(Parameter.LogPath, l.SrcIPAddr)
		if err := CreateFile(folder); err != nil {
			fmt.Printf("创建目录失败: %v\n", err)
			return
		}

		filePath := filepath.Join(folder, l.TopProtocol)
		if err := AppendWriteFile(filePath, "\n"+c); err != nil {
			fmt.Printf("写入文件失败: %v\n", err)
		}
		p.logAndIncrement(l)
		fmt.Println("Recording packet")
	}(content, log)

	return nil
}

func (p *PacketHandler) AlertPacket(rule *Structure.RuleNodeValue) error {
	log, content := p.ConductionActionLog(rule)

	go func(c string, l *DB.ActionDescribeTriggerLog) {
		if err := AppendWriteFile(Parameter.AlertPath, "\n"+c); err != nil {
			fmt.Printf("告警文件写入失败: %v\n", err)
		}
		p.logAndIncrement(l)
		fmt.Println("Alerting packet")
	}(content, log)

	return nil
}

func (p *PacketHandler) SendPacket(rule *Structure.RuleNodeValue) error {
	log, content := p.ConductionActionLog(rule)

	go func(c string, l *DB.ActionDescribeTriggerLog) {
		if err := AppendWriteFile(Parameter.AlertPath, "\n"+c); err != nil {
			fmt.Printf("告警文件写入失败: %v\n", err)
		}

		InitializeStructure.Mail.AddStringEvent(content)

		p.logAndIncrement(l)
		fmt.Println("Sending packet")
	}(content, log)

	return nil
}

func (p *PacketHandler) DropIPv4Packets(rule *Structure.RuleNodeValue) error {
	log, content := p.ConductionActionLog(rule)

	go func(c string, l *DB.ActionDescribeTriggerLog) {
		if err := AppendWriteFile(Parameter.AlertPath, "\n"+c); err != nil {
			fmt.Printf("告警文件写入失败: %v\n", err)
		}

		if rule.Option.BanTime > 0 {
			if err := Action.DropIP(l.SrcIPAddr, strconv.Itoa(rule.Option.BanTime*60)); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := Action.DropIP(l.SrcIPAddr, "900"); err != nil {
				fmt.Println(err)
			}
		}

		InitializeStructure.MainBanIPSet.Insert(Structure.BanIPDescribe{
			BanIP:     l.SrcIPAddr,
			BanRuleID: rule.ID,
			Describe:  c,
			UnbanTime: time.Now().Add(time.Duration(rule.Option.BanTime) * time.Minute),
		})

		p.logAndIncrement(l)

		fmt.Println("Drop packet")
	}(content, log)

	return nil
}
