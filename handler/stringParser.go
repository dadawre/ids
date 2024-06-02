package handler

import (
	"encoding/hex"
	"fmt"
	"graduationDesign/parser/parserRuler"
	"graduationDesign/structure/ruleStructure"
	"strconv"
)

func (p *PacketHandler) parseMsgStringToString(rule *ruleStructure.Rule) *string {
	result := make([]any, 0, 0)
	var str string
	for _, stringId := range rule.Strings {
		if stringId >= parserRuler.EthSrcMACField && stringId <= parserRuler.EthTypeField {
			str = *p.MacStringParse(stringId)
		} else if stringId >= parserRuler.ARPTypeField && stringId <= parserRuler.ARPDestProtocolAddressField {
			str = *p.ARPStringParse(stringId)
		} else if stringId >= parserRuler.IPv4SrcIPField && stringId <= parserRuler.IPv4VersionField {
			str = *p.IPv4StringParse(stringId)
		} else if stringId >= parserRuler.IPv6DstField && stringId <= parserRuler.IPv6VersionField {
			str = *p.IPv6StringParse(stringId)
		} else if stringId >= parserRuler.TCPSrcPortField && stringId <= parserRuler.TCPUrgPtrField {
			str = *p.TCPStringParse(stringId)
		} else if stringId >= parserRuler.UDPSrcField && stringId <= parserRuler.UDPChecksumField {
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
	if stringId == parserRuler.EthSrcMACField {
		*str = p.eth.SrcMAC.String()
	} else if stringId == parserRuler.EthDstMACField {
		*str = p.eth.DstMAC.String()
	} else if stringId == parserRuler.EthLengthField {
		*str = strconv.Itoa(int(p.eth.Length))
	} else if stringId == parserRuler.EthTypeField {
		*str = p.eth.EthernetType.String()
	} else {
		return nil
	}
	return str
}

func (p *PacketHandler) ARPStringParse(stringID int) *string {
	str := new(string)
	if stringID == parserRuler.ARPTypeField {
		*str = p.arp.AddrType.String()
	} else if stringID == parserRuler.ARPProtocolField {
		*str = p.arp.Protocol.String()
	} else if stringID == parserRuler.ARPHwAddressSizeField {
		*str = strconv.Itoa(int(p.arp.HwAddressSize))
	} else if stringID == parserRuler.ARPProtocolAddressSizeField {
		*str = strconv.Itoa(int(p.arp.ProtAddressSize))
	} else if stringID == parserRuler.ARPOptionField {
		*str = strconv.Itoa(int(p.arp.Operation))
	} else if stringID == parserRuler.ARPSourceHwAddressField {
		*str = hex.EncodeToString(p.arp.SourceHwAddress)
	} else if stringID == parserRuler.ARPSourceProtocolAddressField {
		*str = hex.EncodeToString(p.arp.SourceProtAddress)
	} else if stringID == parserRuler.ARPDestHwAddressField {
		*str = hex.EncodeToString(p.arp.DstHwAddress)
	} else if stringID == parserRuler.ARPDestProtocolAddressField {
		*str = hex.EncodeToString(p.arp.DstProtAddress)
	} else {
		return nil
	}
	return str
}

func (p *PacketHandler) IPv4StringParse(stringID int) *string {
	str := new(string)
	if stringID == parserRuler.IPv4SrcIPField {
		*str = p.ipv4.SrcIP.String()
	} else if stringID == parserRuler.IPv4DstIPField {
		*str = p.ipv4.DstIP.String()
	} else if stringID == parserRuler.IPv4ChecksumField {
		*str = strconv.Itoa(int(p.ipv4.Checksum))
	} else if stringID == parserRuler.IPv4ProtocolField {
		*str = p.ipv4.Protocol.String()
	} else if stringID == parserRuler.IPv4TTLField {
		*str = strconv.Itoa(int(p.ipv4.TTL))
	} else if stringID == parserRuler.IPv4OffsetField {
		*str = strconv.Itoa(int(p.ipv4.FragOffset))
	} else if stringID == parserRuler.IPv4FlagField {
		*str = p.ipv4.Flags.String()
	} else if stringID == parserRuler.IPv4IDField {
		*str = strconv.Itoa(int(p.ipv4.Id))
	} else if stringID == parserRuler.IPv4TotalLengthField {
		*str = strconv.Itoa(int(p.ipv4.Length))
	} else if stringID == parserRuler.IPv4TOSField {
		*str = strconv.Itoa(int(p.ipv4.TOS))
	} else if stringID == parserRuler.IPv4IHLField {
		*str = strconv.Itoa(int(p.ipv4.IHL))
	} else if stringID == parserRuler.IPv4VersionField {
		*str = strconv.Itoa(int(p.ipv4.Version))
	} else {
		return nil
	}
	return str
}

func (p *PacketHandler) IPv6StringParse(stringID int) *string {
	str := new(string)
	if stringID == parserRuler.IPv6DstField {
		*str = p.ipv6.SrcIP.String()
	} else if stringID == parserRuler.IPv6SrcField {
		*str = p.ipv6.DstIP.String()
	} else if stringID == parserRuler.IPv6HopLimitField {
		*str = strconv.Itoa(int(p.ipv6.HopLimit))
	} else if stringID == parserRuler.IPv6NextHeadField {
		*str = p.ipv6.NextHeader.String()
	} else if stringID == parserRuler.IPv6LengthField {
		*str = strconv.Itoa(int(p.ipv6.Length))
	} else if stringID == parserRuler.IPv6FlowLabelField {
		*str = strconv.Itoa(int(p.ipv6.FlowLabel))
	} else if stringID == parserRuler.IPv6TrafficClassField {
		*str = strconv.Itoa(int(p.ipv6.TrafficClass))
	} else if stringID == parserRuler.IPv6VersionField {
		*str = strconv.Itoa(int(p.ipv6.Version))
	} else {
		return nil
	}
	return str
}

func (p *PacketHandler) TCPStringParse(stringID int) *string {
	str := new(string)
	if stringID == parserRuler.TCPSrcPortField {
		*str = strconv.Itoa(int(p.tcp.SrcPort))
	} else if stringID == parserRuler.TCPDstPortField {
		*str = strconv.Itoa(int(p.tcp.DstPort))
	} else if stringID == parserRuler.TCPSequenceField {
		*str = strconv.Itoa(int(p.tcp.Seq))
	} else if stringID == parserRuler.TCPAckSeqField {
		*str = strconv.Itoa(int(p.tcp.Ack))
	} else if stringID == parserRuler.TCPOffsetField {
		*str = strconv.Itoa(int(p.tcp.DataOffset))
	} else if stringID == parserRuler.TCPUrgField {
		if p.tcp.URG {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == parserRuler.TCPAckField {
		if p.tcp.ACK {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == parserRuler.TCPPstField {
		if p.tcp.PSH {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == parserRuler.TCPRstField {
		if p.tcp.RST {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == parserRuler.TCPSynField {
		if p.tcp.SYN {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == parserRuler.TCPFinField {
		if p.tcp.FIN {
			*str = "true"
		} else {
			*str = "false"
		}
	} else if stringID == parserRuler.TCPWindowField {
		*str = strconv.Itoa(int(p.tcp.Window))
	} else if stringID == parserRuler.TCPChecksumField {
		*str = strconv.Itoa(int(p.tcp.Checksum))
	} else if stringID == parserRuler.TCPUrgPtrField {
		*str = strconv.Itoa(int(p.tcp.Urgent))
	} else {
		return nil
	}
	return str
}

func (p *PacketHandler) UDPStringParse(stringID int) *string {
	str := new(string)
	if stringID == parserRuler.UDPSrcField {
		*str = strconv.Itoa(int(p.udp.SrcPort))
	} else if stringID == parserRuler.UDPDstField {
		*str = strconv.Itoa(int(p.udp.DstPort))
	} else if stringID == parserRuler.UDPLengthField {
		*str = strconv.Itoa(int(p.udp.Length))
	} else if stringID == parserRuler.UDPChecksumField {
		*str = strconv.Itoa(int(p.udp.Checksum))
	} else {
		return nil
	}
	return str
}
