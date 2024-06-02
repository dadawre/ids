package ruleMatch

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"graduationDesign/parser/parserRuler"
	"graduationDesign/structure/ruleStructure"
	"net"
)

// ExecuteLogic evaluates the boolean expression given bools and operators
// in the order of logical AND first, then logical OR.
func ExecuteLogic(bools []bool, operators []int8) (bool, error) {
	if len(bools) == 0 {
		return false, errors.New("the bool array is empty")
	}
	//if len(operators) != len(bools)-1 {
	//	return false, errors.New("the number of operators must be one less than the number of bools")
	//}

	// Step 1: Process all && operators
	andProcessed := []bool{bools[0]}
	var andOps []int8
	for i := 0; i < len(operators); i++ {
		if operators[i] == 1 {
			andProcessed[len(andProcessed)-1] = andProcessed[len(andProcessed)-1] && bools[i+1]
		} else {
			andProcessed = append(andProcessed, bools[i+1])
			andOps = append(andOps, operators[i])
		}
	}

	// Step 2: Process all || operators
	result := andProcessed[0]
	for i := 0; i < len(andOps); i++ {
		if andOps[i] == 2 {
			result = result || andProcessed[i+1]
		} else {
			return false, fmt.Errorf("invalid operator after && processing: %d", andOps[i])
		}
	}

	return result, nil
}

// 比较两个字节数组是否相等
func bytesEqual[T []byte | net.HardwareAddr](a, b T) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// MapLayerTypeToMatchFunc LayerType映射到对应协议的匹配函数
var MapLayerTypeToMatchFunc = map[gopacket.LayerType]func(protocol *ruleStructure.Protocol, payloadPtr interface{}) bool{
	layers.LayerTypeEthernet: MACPayloadMatch,
	layers.LayerTypeARP:      ARPPayLoadMatch,
	layers.LayerTypeIPv4:     IPv4PayLoadMatch,
	layers.LayerTypeIPv6:     IPv6PayLoadMatch,
	layers.LayerTypeTCP:      TCPPayLoadMatch,
	layers.LayerTypeUDP:      UDPPayLoadMatch,
}

// MACPayloadMatch 对Mac协议进行匹配
func MACPayloadMatch(protocol *ruleStructure.Protocol, payloadPtr interface{}) bool {
	mac := payloadPtr.(*layers.Ethernet)

	for _, field := range protocol.Fields {
		if field.NameId == parserRuler.EthSrcMACField {
			if !bytesEqual(mac.SrcMAC, *field.ValuePtr.(*net.HardwareAddr)) {
				return false
			}
		} else if field.NameId == parserRuler.EthDstMACField {
			if !bytesEqual(mac.DstMAC, *field.ValuePtr.(*net.HardwareAddr)) {
				return false
			}
		} else if field.NameId == parserRuler.EthLengthField {
			if mac.Length != field.ValuePtr.(uint16) {
				return false
			}
		} else if field.NameId == parserRuler.EthTypeField {
			if mac.EthernetType != field.ValuePtr.(layers.EthernetType) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func ARPPayLoadMatch(protocol *ruleStructure.Protocol, payloadPtr interface{}) bool {
	arp := payloadPtr.(*layers.ARP)
	for _, field := range protocol.Fields {
		if field.NameId == parserRuler.ARPTypeField {
			if field.ValuePtr.(layers.LinkType) != arp.AddrType {
				return false
			}
		} else if field.NameId == parserRuler.ARPProtocolField {
			if field.ValuePtr.(layers.EthernetType) != arp.Protocol {
				return false
			}
		} else if field.NameId == parserRuler.ARPHwAddressSizeField {
			if field.ValuePtr.(uint8) != arp.HwAddressSize {
				return false
			}
		} else if field.NameId == parserRuler.ARPProtocolAddressSizeField {
			if field.ValuePtr.(uint8) != arp.ProtAddressSize {
				return false
			}
		} else if field.NameId == parserRuler.ARPOptionField {
			if field.ValuePtr.(uint16) != arp.Operation {
				return false
			}
		} else if field.NameId == parserRuler.ARPSourceHwAddressField {
			if !(bytesEqual(field.ValuePtr.([]byte), arp.SourceHwAddress)) {
				return false
			}
		} else if field.NameId == parserRuler.ARPDestHwAddressField {
			if !(bytesEqual(field.ValuePtr.([]byte), arp.DstHwAddress)) {
				return false
			}
		} else if field.NameId == parserRuler.ARPSourceProtocolAddressField {
			if !(bytesEqual(field.ValuePtr.([]byte), arp.SourceProtAddress)) {
				return false
			}
		} else if field.NameId == parserRuler.ARPDestProtocolAddressField {
			if !(bytesEqual(field.ValuePtr.([]byte), arp.DstProtAddress)) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func IPv6PayLoadMatch(protocol *ruleStructure.Protocol, payloadPtr interface{}) bool {
	ipv6 := payloadPtr.(*layers.IPv6)
	for _, field := range protocol.Fields {
		if field.NameId == parserRuler.IPv6DstField {
			if !(field.ValuePtr.(*net.IPNet).Contains(ipv6.DstIP)) {
				return false
			}
		} else if field.NameId == parserRuler.IPv6SrcField {
			if !(field.ValuePtr.(*net.IPNet).Contains(ipv6.SrcIP)) {
				return false
			}
		} else if field.NameId == parserRuler.IPv6VersionField {
			if field.ValuePtr.(uint8) != ipv6.Version {
				return false
			}
		} else if field.NameId == parserRuler.IPv6TrafficClassField {
			if field.ValuePtr.(uint8) != ipv6.TrafficClass {
				return false
			}
		} else if field.NameId == parserRuler.IPv6HopLimitField {
			if field.ValuePtr.(uint8) != ipv6.HopLimit {
				return false
			}
		} else if field.NameId == parserRuler.IPv6FlowLabelField {
			if field.ValuePtr.(uint32) != ipv6.FlowLabel {
				return false
			}
		} else if field.NameId == parserRuler.IPv6LengthField {
			if field.ValuePtr.(uint16) != ipv6.Length {
				return false
			}
		} else if field.NameId == parserRuler.IPv6NextHeadField {
			if field.ValuePtr.(layers.IPProtocol) != ipv6.NextHeader {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func IPv4PayLoadMatch(protocol *ruleStructure.Protocol, payloadPtr interface{}) bool {
	ipv4 := payloadPtr.(*layers.IPv4)
	for _, field := range protocol.Fields {
		if field.NameId == parserRuler.IPv4SrcIPField {
			if !(field.ValuePtr.(*net.IPNet).Contains(ipv4.SrcIP)) {
				return false
			}
		} else if field.NameId == parserRuler.IPv4DstIPField {
			if !(field.ValuePtr.(*net.IPNet).Contains(ipv4.DstIP)) {
				return false
			}
		} else if field.NameId == parserRuler.IPv4VersionField {
			if field.ValuePtr.(uint8) != ipv4.Version {
				return false
			}
		} else if field.NameId == parserRuler.IPv4IHLField {
			if field.ValuePtr.(uint8) != ipv4.IHL {
				return false
			}
		} else if field.NameId == parserRuler.IPv4TOSField {
			if field.ValuePtr.(uint8) != ipv4.TOS {
				return false
			}
		} else if field.NameId == parserRuler.IPv4TTLField {
			if field.ValuePtr.(uint8) != ipv4.TTL {
				return false
			}
		} else if field.NameId == parserRuler.IPv4FlagField {
			if field.ValuePtr.(layers.IPv4Flag) != ipv4.Flags {
				return false
			}
		} else if field.NameId == parserRuler.IPv4ProtocolField {
			if field.ValuePtr.(layers.IPProtocol) != ipv4.Protocol {
				return false
			}
		} else if field.NameId == parserRuler.IPv4TotalLengthField {
			if field.ValuePtr.(uint16) != ipv4.Length {
				return false
			}
		} else if field.NameId == parserRuler.IPv4IDField {
			if field.ValuePtr.(uint16) != ipv4.Id {
				return false
			}
		} else if field.NameId == parserRuler.IPv4OffsetField {
			if field.ValuePtr.(uint16) != ipv4.FragOffset {
				return false
			}
		} else if field.NameId == parserRuler.IPv4ChecksumField {
			if field.ValuePtr.(uint16) != ipv4.Checksum {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func TCPPayLoadMatch(protocol *ruleStructure.Protocol, payloadPtr interface{}) bool {
	tcp := payloadPtr.(*layers.TCP)
	for _, field := range protocol.Fields {
		if field.NameId == parserRuler.TCPSrcPortField {
			if field.ValuePtr.(layers.TCPPort) != tcp.SrcPort {
				return false
			}
		} else if field.NameId == parserRuler.TCPDstPortField {
			if field.ValuePtr.(layers.TCPPort) != tcp.DstPort {
				return false
			}
		} else if field.NameId == parserRuler.TCPAckSeqField {
			if field.ValuePtr.(uint32) != tcp.Ack {
				return false
			}
		} else if field.NameId == parserRuler.TCPSequenceField {
			if field.ValuePtr.(uint32) != tcp.Seq {
				return false
			}
		} else if field.NameId == parserRuler.TCPAckField {
			if field.ValuePtr.(bool) != tcp.ACK {
				return false
			}
		} else if field.NameId == parserRuler.TCPSynField {
			if field.ValuePtr.(bool) != tcp.SYN {
				return false
			}
		} else if field.NameId == parserRuler.TCPUrgField {
			if field.ValuePtr.(bool) != tcp.URG {
				return false
			}
		} else if field.NameId == parserRuler.TCPPstField {
			if field.ValuePtr.(bool) != tcp.PSH {
				return false
			}
		} else if field.NameId == parserRuler.TCPFinField {
			if field.ValuePtr.(bool) != tcp.FIN {
				return false
			}
		} else if field.NameId == parserRuler.TCPRstField {
			if field.ValuePtr.(bool) != tcp.RST {
				return false
			}
		} else if field.NameId == parserRuler.TCPWindowField {
			if field.ValuePtr.(uint16) != tcp.Window {
				return false
			}
		} else if field.NameId == parserRuler.TCPChecksumField {
			if field.ValuePtr.(uint16) != tcp.Checksum {
				return false
			}
		} else if field.NameId == parserRuler.TCPUrgPtrField {
			if field.ValuePtr.(uint16) != tcp.Urgent {
				return false
			}
		} else if field.NameId == parserRuler.TCPOffsetField {
			if field.ValuePtr.(uint8) != tcp.DataOffset {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func UDPPayLoadMatch(protocol *ruleStructure.Protocol, payloadPtr interface{}) bool {
	udp := payloadPtr.(*layers.UDP)
	for _, field := range protocol.Fields {
		if field.NameId == parserRuler.UDPSrcField {
			if field.ValuePtr.(layers.UDPPort) != udp.SrcPort {
				return false
			}
		} else if field.NameId == parserRuler.UDPDstField {
			if field.ValuePtr.(layers.UDPPort) != udp.DstPort {
				return false
			}
		} else if field.NameId == parserRuler.UDPChecksumField {
			if field.ValuePtr.(uint16) != udp.Checksum {
				return false
			}
		} else if field.NameId == parserRuler.UDPLengthField {
			if field.ValuePtr.(uint16) != udp.Length {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
