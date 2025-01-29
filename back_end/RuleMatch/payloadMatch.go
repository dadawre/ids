package RuleMatch

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"graduationDesign/back_end/OSController"
	"graduationDesign/back_end/Parser/ParserRuler"
	"graduationDesign/back_end/Structure"
	"net"
)

// MapLayerTypeToMatchFunc LayerType映射到对应协议的匹配函数
var MapLayerTypeToMatchFunc = map[gopacket.LayerType]func(protocol *Structure.Protocol, payloadPtr interface{}) bool{
	layers.LayerTypeEthernet: MACPayloadMatch,
	layers.LayerTypeARP:      ARPPayLoadMatch,
	layers.LayerTypeIPv4:     IPv4PayLoadMatch,
	layers.LayerTypeIPv6:     IPv6PayLoadMatch,
	layers.LayerTypeTCP:      TCPPayLoadMatch,
	layers.LayerTypeUDP:      UDPPayLoadMatch,
	10000:                    OSPayLoadMatch,
}

// MACPayloadMatch 对Mac协议进行匹配
func MACPayloadMatch(protocol *Structure.Protocol, payloadPtr interface{}) bool {
	mac := payloadPtr.(*layers.Ethernet)

	for _, field := range protocol.Fields {
		if field.NameId == ParserRuler.EthSrcMACField {
			if !bytesEqual(mac.SrcMAC, *field.ValuePtr.(*net.HardwareAddr)) {
				return false
			}
		} else if field.NameId == ParserRuler.EthDstMACField {
			if !bytesEqual(mac.DstMAC, *field.ValuePtr.(*net.HardwareAddr)) {
				return false
			}
		} else if field.NameId == ParserRuler.EthLengthField {
			if mac.Length != field.ValuePtr.(uint16) {
				return false
			}
		} else if field.NameId == ParserRuler.EthTypeField {
			if mac.EthernetType != field.ValuePtr.(layers.EthernetType) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func ARPPayLoadMatch(protocol *Structure.Protocol, payloadPtr interface{}) bool {
	arp := payloadPtr.(*layers.ARP)
	for _, field := range protocol.Fields {
		if field.NameId == ParserRuler.ARPTypeField {
			if field.ValuePtr.(layers.LinkType) != arp.AddrType {
				return false
			}
		} else if field.NameId == ParserRuler.ARPProtocolField {
			if field.ValuePtr.(layers.EthernetType) != arp.Protocol {
				return false
			}
		} else if field.NameId == ParserRuler.ARPHwAddressSizeField {
			if field.ValuePtr.(uint8) != arp.HwAddressSize {
				return false
			}
		} else if field.NameId == ParserRuler.ARPProtocolAddressSizeField {
			if field.ValuePtr.(uint8) != arp.ProtAddressSize {
				return false
			}
		} else if field.NameId == ParserRuler.ARPOptionField {
			if field.ValuePtr.(uint16) != arp.Operation {
				return false
			}
		} else if field.NameId == ParserRuler.ARPSourceHwAddressField {
			if !(bytesEqual(field.ValuePtr.([]byte), arp.SourceHwAddress)) {
				return false
			}
		} else if field.NameId == ParserRuler.ARPDestHwAddressField {
			if !(bytesEqual(field.ValuePtr.([]byte), arp.DstHwAddress)) {
				return false
			}
		} else if field.NameId == ParserRuler.ARPSourceProtocolAddressField {
			if !(bytesEqual(field.ValuePtr.([]byte), arp.SourceProtAddress)) {
				return false
			}
		} else if field.NameId == ParserRuler.ARPDestProtocolAddressField {
			if !(bytesEqual(field.ValuePtr.([]byte), arp.DstProtAddress)) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func IPv6PayLoadMatch(protocol *Structure.Protocol, payloadPtr interface{}) bool {
	ipv6 := payloadPtr.(*layers.IPv6)
	for _, field := range protocol.Fields {
		if field.NameId == ParserRuler.IPv6DstField {
			if !(field.ValuePtr.(*net.IPNet).Contains(ipv6.DstIP)) {
				return false
			}
		} else if field.NameId == ParserRuler.IPv6SrcField {
			if !(field.ValuePtr.(*net.IPNet).Contains(ipv6.SrcIP)) {
				return false
			}
		} else if field.NameId == ParserRuler.IPv6VersionField {
			if field.ValuePtr.(uint8) != ipv6.Version {
				return false
			}
		} else if field.NameId == ParserRuler.IPv6TrafficClassField {
			if field.ValuePtr.(uint8) != ipv6.TrafficClass {
				return false
			}
		} else if field.NameId == ParserRuler.IPv6HopLimitField {
			if field.ValuePtr.(uint8) != ipv6.HopLimit {
				return false
			}
		} else if field.NameId == ParserRuler.IPv6FlowLabelField {
			if field.ValuePtr.(uint32) != ipv6.FlowLabel {
				return false
			}
		} else if field.NameId == ParserRuler.IPv6LengthField {
			if field.ValuePtr.(uint16) != ipv6.Length {
				return false
			}
		} else if field.NameId == ParserRuler.IPv6NextHeadField {
			if field.ValuePtr.(layers.IPProtocol) != ipv6.NextHeader {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func IPv4PayLoadMatch(protocol *Structure.Protocol, payloadPtr interface{}) bool {
	ipv4 := payloadPtr.(*layers.IPv4)
	for _, field := range protocol.Fields {
		if field.NameId == ParserRuler.IPv4SrcIPField {
			if !(field.ValuePtr.(*net.IPNet).Contains(ipv4.SrcIP)) {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4DstIPField {
			if !(field.ValuePtr.(*net.IPNet).Contains(ipv4.DstIP)) {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4VersionField {
			if field.ValuePtr.(uint8) != ipv4.Version {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4IHLField {
			if field.ValuePtr.(uint8) != ipv4.IHL {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4TOSField {
			if field.ValuePtr.(uint8) != ipv4.TOS {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4TTLField {
			if field.ValuePtr.(uint8) != ipv4.TTL {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4FlagField {
			if field.ValuePtr.(layers.IPv4Flag) != ipv4.Flags {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4ProtocolField {
			if field.ValuePtr.(layers.IPProtocol) != ipv4.Protocol {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4TotalLengthField {
			if field.ValuePtr.(uint16) != ipv4.Length {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4IDField {
			if field.ValuePtr.(uint16) != ipv4.Id {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4OffsetField {
			if field.ValuePtr.(uint16) != ipv4.FragOffset {
				return false
			}
		} else if field.NameId == ParserRuler.IPv4ChecksumField {
			if field.ValuePtr.(uint16) != ipv4.Checksum {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func TCPPayLoadMatch(protocol *Structure.Protocol, payloadPtr interface{}) bool {
	tcp := payloadPtr.(*layers.TCP)
	for _, field := range protocol.Fields {
		if field.NameId == ParserRuler.TCPSrcPortField {
			if field.ValuePtr.(layers.TCPPort) != tcp.SrcPort {
				return false
			}
		} else if field.NameId == ParserRuler.TCPDstPortField {
			if field.ValuePtr.(layers.TCPPort) != tcp.DstPort {
				return false
			}
		} else if field.NameId == ParserRuler.TCPAckSeqField {
			if field.ValuePtr.(uint32) != tcp.Ack {
				return false
			}
		} else if field.NameId == ParserRuler.TCPSequenceField {
			if field.ValuePtr.(uint32) != tcp.Seq {
				return false
			}
		} else if field.NameId == ParserRuler.TCPAckField {
			if field.ValuePtr.(bool) != tcp.ACK {
				return false
			}
		} else if field.NameId == ParserRuler.TCPSynField {
			if field.ValuePtr.(bool) != tcp.SYN {
				return false
			}
		} else if field.NameId == ParserRuler.TCPUrgField {
			if field.ValuePtr.(bool) != tcp.URG {
				return false
			}
		} else if field.NameId == ParserRuler.TCPPstField {
			if field.ValuePtr.(bool) != tcp.PSH {
				return false
			}
		} else if field.NameId == ParserRuler.TCPFinField {
			if field.ValuePtr.(bool) != tcp.FIN {
				return false
			}
		} else if field.NameId == ParserRuler.TCPRstField {
			if field.ValuePtr.(bool) != tcp.RST {
				return false
			}
		} else if field.NameId == ParserRuler.TCPWindowField {
			if field.ValuePtr.(uint16) != tcp.Window {
				return false
			}
		} else if field.NameId == ParserRuler.TCPChecksumField {
			if field.ValuePtr.(uint16) != tcp.Checksum {
				return false
			}
		} else if field.NameId == ParserRuler.TCPUrgPtrField {
			if field.ValuePtr.(uint16) != tcp.Urgent {
				return false
			}
		} else if field.NameId == ParserRuler.TCPOffsetField {
			if field.ValuePtr.(uint8) != tcp.DataOffset {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func UDPPayLoadMatch(protocol *Structure.Protocol, payloadPtr interface{}) bool {
	udp := payloadPtr.(*layers.UDP)
	for _, field := range protocol.Fields {
		if field.NameId == ParserRuler.UDPSrcField {
			if field.ValuePtr.(layers.UDPPort) != udp.SrcPort {
				return false
			}
		} else if field.NameId == ParserRuler.UDPDstField {
			if field.ValuePtr.(layers.UDPPort) != udp.DstPort {
				return false
			}
		} else if field.NameId == ParserRuler.UDPChecksumField {
			if field.ValuePtr.(uint16) != udp.Checksum {
				return false
			}
		} else if field.NameId == ParserRuler.UDPLengthField {
			if field.ValuePtr.(uint16) != udp.Length {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func OSPayLoadMatch(protocol *Structure.Protocol, payloadPtr interface{}) bool {
	os := payloadPtr.(*OSController.OSState)

	for _, field := range protocol.Fields {
		if field.NameId == ParserRuler.OSCpuUsageRateField {
			if field.ValuePtr.(float64) > os.TotalCPUPercent {
				return false
			}
		} else if field.NameId == ParserRuler.OSGpuUsageRateField {
			if field.ValuePtr.(float64) > os.TotalGPUPercent {
				return false
			}
		} else if field.NameId == ParserRuler.OSMemoryUsageRateField {
			if field.ValuePtr.(float64) > os.TotalMemPercent {
				return false
			}
		}
	}
	return true
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
