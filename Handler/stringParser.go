package Handler

import (
	"encoding/hex"
	"fmt"
	"graduationDesign/Parser/ParserRuler"
	"graduationDesign/structure/ruleStructure"
	"strconv"
)

// 将字符串的字段值解析成对应类型的值
func (p *PacketHandler) parseMsgStringToValue(rule *ruleStructure.Rule) *string {
	result := make([]any, 0, 0)
	var str string
	for _, stringId := range rule.Strings {
		if stringId >= ParserRuler.EthSrcMACField && stringId <= ParserRuler.EthTypeField {
			str = *p.MacStringParse(stringId)
		} else if stringId >= ParserRuler.ARPTypeField && stringId <= ParserRuler.ARPDestProtocolAddressField {
			str = *p.ARPStringParse(stringId)
		} else if stringId >= ParserRuler.IPv4SrcIPField && stringId <= ParserRuler.IPv4VersionField {
			str = *p.IPv4StringParse(stringId)
		} else if stringId >= ParserRuler.IPv6DstField && stringId <= ParserRuler.IPv6VersionField {
			str = *p.IPv6StringParse(stringId)
		} else if stringId >= ParserRuler.TCPSrcPortField && stringId <= ParserRuler.TCPUrgPtrField {
			str = *p.TCPStringParse(stringId)
		} else if stringId >= ParserRuler.UDPSrcField && stringId <= ParserRuler.UDPChecksumField {
			str = *p.UDPStringParse(stringId)
		} else {
			return nil
		}
		result = append(result, str)
	}

	resultStr := fmt.Sprintf(rule.Msg+"\n", result...)
	//fmt.Println(rule.Msg)
	//fmt.Println(result)
	//fmt.Println(resultStr)
	return &resultStr
}

func (p *PacketHandler) MacStringParse(stringId int) *string {
	str := new(string)
	if stringId == ParserRuler.EthSrcMACField {
		*str = p.eth.SrcMAC.String()
	} else if stringId == ParserRuler.EthDstMACField {
		*str = p.eth.DstMAC.String()
	} else if stringId == ParserRuler.EthLengthField {
		*str = strconv.Itoa(int(p.eth.Length))
	} else if stringId == ParserRuler.EthTypeField {
		*str = p.eth.EthernetType.String()
	} else {
		return nil
	}
	return str
}

func (p *PacketHandler) ARPStringParse(stringID int) *string {
	str := new(string)
	if stringID == ParserRuler.ARPTypeField {
		*str = p.arp.AddrType.String()
	} else if stringID == ParserRuler.ARPProtocolField {
		*str = p.arp.Protocol.String()
	} else if stringID == ParserRuler.ARPHwAddressSizeField {
		*str = strconv.Itoa(int(p.arp.HwAddressSize))
	} else if stringID == ParserRuler.ARPProtocolAddressSizeField {
		*str = strconv.Itoa(int(p.arp.ProtAddressSize))
	} else if stringID == ParserRuler.ARPOptionField {
		*str = strconv.Itoa(int(p.arp.Operation))
	} else if stringID == ParserRuler.ARPSourceHwAddressField {
		*str = hex.EncodeToString(p.arp.SourceHwAddress)
	} else if stringID == ParserRuler.ARPSourceProtocolAddressField {
		*str = hex.EncodeToString(p.arp.SourceProtAddress)
	} else if stringID == ParserRuler.ARPDestHwAddressField {
		*str = hex.EncodeToString(p.arp.DstHwAddress)
	} else if stringID == ParserRuler.ARPDestProtocolAddressField {
		*str = hex.EncodeToString(p.arp.DstProtAddress)
	} else {
		return nil
	}
	return str
}

func (p *PacketHandler) IPv4StringParse(stringID int) *string {
	str := new(string)
	if stringID == ParserRuler.IPv4SrcIPField {
		*str = p.ipv4.SrcIP.String()
	} else if stringID == ParserRuler.IPv4DstIPField {
		*str = p.ipv4.DstIP.String()
	} else if stringID == ParserRuler.IPv4ChecksumField {
		*str = strconv.Itoa(int(p.ipv4.Checksum))
	} else if stringID == ParserRuler.IPv4ProtocolField {
		*str = p.ipv4.Protocol.String()
	} else if stringID == ParserRuler.IPv4TTLField {
		*str = strconv.Itoa(int(p.ipv4.TTL))
	} else if stringID == ParserRuler.IPv4OffsetField {
		*str = strconv.Itoa(int(p.ipv4.FragOffset))
	} else if stringID == ParserRuler.IPv4FlagField {
		*str = p.ipv4.Flags.String()
	} else if stringID == ParserRuler.IPv4IDField {
		*str = strconv.Itoa(int(p.ipv4.Id))
	} else if stringID == ParserRuler.IPv4TotalLengthField {
		*str = strconv.Itoa(int(p.ipv4.Length))
	} else if stringID == ParserRuler.IPv4TOSField {
		*str = strconv.Itoa(int(p.ipv4.TOS))
	} else if stringID == ParserRuler.IPv4IHLField {
		*str = strconv.Itoa(int(p.ipv4.IHL))
	} else if stringID == ParserRuler.IPv4VersionField {
		*str = strconv.Itoa(int(p.ipv4.Version))
	} else {
		return nil
	}
	return str
}

func (p *PacketHandler) IPv6StringParse(stringID int) *string {
	str := new(string)
	if stringID == ParserRuler.IPv6DstField {
		*str = p.ipv6.SrcIP.String()
	} else if stringID == ParserRuler.IPv6SrcField {
		*str = p.ipv6.DstIP.String()
	} else if stringID == ParserRuler.IPv6HopLimitField {
		*str = strconv.Itoa(int(p.ipv6.HopLimit))
	} else if stringID == ParserRuler.IPv6NextHeadField {
		*str = p.ipv6.NextHeader.String()
	} else if stringID == ParserRuler.IPv6LengthField {
		*str = strconv.Itoa(int(p.ipv6.Length))
	} else if stringID == ParserRuler.IPv6FlowLabelField {
		*str = strconv.Itoa(int(p.ipv6.FlowLabel))
	} else if stringID == ParserRuler.IPv6TrafficClassField {
		*str = strconv.Itoa(int(p.ipv6.TrafficClass))
	} else if stringID == ParserRuler.IPv6VersionField {
		*str = strconv.Itoa(int(p.ipv6.Version))
	} else {
		return nil
	}
	return str
}

func (p *PacketHandler) TCPStringParse(stringID int) *string {
	str := new(string)
	if stringID == ParserRuler.TCPSrcPortField {
		*str = strconv.Itoa(int(p.tcp.SrcPort))
	} else if stringID == ParserRuler.TCPDstPortField {
		*str = strconv.Itoa(int(p.tcp.DstPort))
	} else if stringID == ParserRuler.TCPSequenceField {
		*str = strconv.Itoa(int(p.tcp.Seq))
	} else if stringID == ParserRuler.TCPAckSeqField {
		*str = strconv.Itoa(int(p.tcp.Ack))
	} else if stringID == ParserRuler.TCPOffsetField {
		*str = strconv.Itoa(int(p.tcp.DataOffset))
	} else if stringID == ParserRuler.TCPUrgField {
		if p.tcp.URG {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == ParserRuler.TCPAckField {
		if p.tcp.ACK {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == ParserRuler.TCPPstField {
		if p.tcp.PSH {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == ParserRuler.TCPRstField {
		if p.tcp.RST {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == ParserRuler.TCPSynField {
		if p.tcp.SYN {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == ParserRuler.TCPFinField {
		if p.tcp.FIN {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == ParserRuler.TCPWindowField {
		*str = strconv.Itoa(int(p.tcp.Window))
	} else if stringID == ParserRuler.TCPChecksumField {
		*str = strconv.Itoa(int(p.tcp.Checksum))
	} else if stringID == ParserRuler.TCPUrgPtrField {
		*str = strconv.Itoa(int(p.tcp.Urgent))
	} else {
		return nil
	}
	return str
}

func (p *PacketHandler) UDPStringParse(stringID int) *string {
	str := new(string)
	if stringID == ParserRuler.UDPSrcField {
		*str = strconv.Itoa(int(p.udp.SrcPort))
	} else if stringID == ParserRuler.UDPDstField {
		*str = strconv.Itoa(int(p.udp.DstPort))
	} else if stringID == ParserRuler.UDPLengthField {
		*str = strconv.Itoa(int(p.udp.Length))
	} else if stringID == ParserRuler.UDPChecksumField {
		*str = strconv.Itoa(int(p.udp.Checksum))
	} else {
		return nil
	}
	return str
}
