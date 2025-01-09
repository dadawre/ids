package Mode

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"time"
)

type DecodedProtocol struct {
	ethLayer    layers.Ethernet
	ipv4Layer   layers.IPv4
	ipv6Layer   layers.IPv6
	icmpv4Layer layers.ICMPv4
	icmpv6Layer layers.ICMPv6
	tcpLayer    layers.TCP
	udpLayer    layers.UDP
	payload     gopacket.Payload
}

func PacketSniffMode(packetSource *gopacket.PacketSource) {
	defer inactiveHandle.CleanUp()
	defer handle.Close()
	// 创建 DecodingLayerParser 和用于存储解析结果的 slice
	var (
		ethLayer    = layers.Ethernet{}
		ipv4Layer   = layers.IPv4{}
		ipv6Layer   = layers.IPv6{}
		icmpv4Layer = layers.ICMPv4{}
		icmpv6Layer = layers.ICMPv6{}
		tcpLayer    = layers.TCP{}
		udpLayer    = layers.UDP{}
		payload     = gopacket.Payload{}
		protocols   = DecodedProtocol{
			ethLayer:    ethLayer,
			ipv4Layer:   ipv4Layer,
			ipv6Layer:   ipv6Layer,
			icmpv4Layer: icmpv4Layer,
			icmpv6Layer: icmpv6Layer,
			tcpLayer:    tcpLayer,
			udpLayer:    udpLayer,
			payload:     payload,
		}
		decoded []gopacket.LayerType
		parser  = gopacket.NewDecodingLayerParser(
			layers.LayerTypeEthernet,
			&protocols.ethLayer,
			&protocols.ipv4Layer,
			&protocols.ipv6Layer,
			&protocols.icmpv4Layer,
			&protocols.icmpv6Layer,
			&protocols.tcpLayer,
			&protocols.udpLayer,
			&protocols.payload,
		)
	)

	// 捕获数据包
	for packet := range packetSource.Packets() {
		decoded = decoded[:0] // 重置解析结果
		err := parser.DecodeLayers(packet.Data(), &decoded)
		if err != nil {

			fmt.Println("解析数据包错误:", err)
			continue
		}
		printDecodedPacket(decoded, &protocols)
	}
}

// printDecodedPacket 打印解析后的协议层信息
func printDecodedPacket(decoded []gopacket.LayerType, protocols *DecodedProtocol) {
	fmt.Printf("\n时间: %s\n", time.Now().Format(time.RFC3339))

	// 遍历已解码的层
	for _, layerType := range decoded {
		switch layerType {
		case layers.LayerTypeEthernet:
			fmt.Println("以太网层:")
			fmt.Printf("  源MAC: %s, 目标MAC: %s\n", protocols.ethLayer.SrcMAC, protocols.ethLayer.DstMAC)
		case layers.LayerTypeIPv4:
			fmt.Println("IPv4层:")
			fmt.Printf("  源IP: %s, 目标IP: %s\n", protocols.ipv4Layer.SrcIP, protocols.ipv4Layer.DstIP)
			fmt.Printf("  协议: %s\n", protocols.ipv4Layer.Protocol)
		case layers.LayerTypeIPv6:
			fmt.Println("IPv6层:")
			fmt.Printf("  源IP: %s, 目标IP: %s\n", protocols.ipv6Layer.SrcIP, protocols.ipv6Layer.DstIP)
			fmt.Printf("  协议: %s\n", protocols.ipv6Layer.NextHeader)
		case layers.LayerTypeTCP:
			fmt.Println("TCP层:")
			fmt.Printf("  源端口: %d, 目标端口: %d\n", protocols.tcpLayer.SrcPort, protocols.tcpLayer.DstPort)
			fmt.Printf("  标志: SYN=%t, ACK=%t, FIN=%t, PSH=%t, RST=%t, URG=%t\n", protocols.tcpLayer.SYN,
				protocols.tcpLayer.ACK, protocols.tcpLayer.FIN, protocols.tcpLayer.PSH, protocols.tcpLayer.RST, protocols.tcpLayer.URG)
		case layers.LayerTypeUDP:
			fmt.Println("UDP层:")
			fmt.Printf("  源端口: %d, 目标端口: %d\n", protocols.udpLayer.SrcPort, protocols.udpLayer.DstPort)
		case gopacket.LayerTypePayload:
			fmt.Println("应用层数据:")
			fmt.Printf("  %s\n", string(protocols.payload.Payload()))
		default:
			fmt.Printf("未知层: %s\n", layerType)
		}
	}
}
