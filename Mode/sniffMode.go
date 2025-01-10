package Mode

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net/http"
	"time"
)

type DecodedProtocol struct {
	ethLayer    layers.Ethernet
	arpLayer    layers.ARP
	ipv4Layer   layers.IPv4
	ipv6Layer   layers.IPv6
	icmpv4Layer layers.ICMPv4
	icmpv6Layer layers.ICMPv6
	tcpLayer    layers.TCP
	udpLayer    layers.UDP
	dnsLayer    layers.DNS
	payload     gopacket.Payload
}

func PacketSniffMode(packetSource *gopacket.PacketSource) {
	defer inactiveHandle.CleanUp()
	defer handle.Close()
	// 创建 DecodingLayerParser 和用于存储解析结果的 slice
	var (
		ethLayer    = layers.Ethernet{}
		arpLayer    = layers.ARP{}
		ipv4Layer   = layers.IPv4{}
		ipv6Layer   = layers.IPv6{}
		icmpv4Layer = layers.ICMPv4{}
		icmpv6Layer = layers.ICMPv6{}
		tcpLayer    = layers.TCP{}
		udpLayer    = layers.UDP{}
		payload     = gopacket.Payload{}
		protocols   = DecodedProtocol{
			ethLayer:    ethLayer,
			arpLayer:    arpLayer,
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
			&protocols.arpLayer,
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
		case layers.LayerTypeARP:
			fmt.Println("ARP层:")
			fmt.Printf("  地址类型: %d\n  硬件源地址: %s, 硬件目标地址: %s\n",
				protocols.arpLayer.AddrType, protocols.arpLayer.SourceHwAddress, protocols.arpLayer.DstHwAddress)
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
			if isHTTPPayload(protocols.payload.Payload()) {
				fmt.Println("  检测到HTTP数据流:")
				parseHTTP(protocols.payload.Payload())
			}
		case layers.LayerTypeUDP:
			fmt.Println("UDP层:")
			fmt.Printf("  源端口: %d, 目标端口: %d\n", protocols.udpLayer.SrcPort, protocols.udpLayer.DstPort)
		case layers.LayerTypeDNS:
			fmt.Println("DNS层:")
			for _, query := range protocols.dnsLayer.Questions {
				fmt.Printf("  查询: %s 类型: %d\n", string(query.Name), query.Type)
			}
			for _, answer := range protocols.dnsLayer.Answers {
				fmt.Printf("  回答: %s 类型: %d 数据: %s\n", string(answer.Name), answer.Type, answer.IP)
			}
		case gopacket.LayerTypePayload:
			fmt.Println("应用层数据:")
			fmt.Printf("  %s\n", string(protocols.payload.Payload()))
		default:
			fmt.Printf("未知层: %s\n", layerType)
		}
	}
}

func isHTTPPayload(payload []byte) bool {
	// 将载荷转换为字符串，便于解析
	data := string(payload)

	// 检查是否匹配 HTTP 请求的起始格式
	httpMethods := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}
	for _, method := range httpMethods {
		if len(data) > len(method) && data[:len(method)] == method {
			return true
		}
	}

	// 检查是否匹配 HTTP 响应的起始格式
	if len(data) > 5 && data[:5] == "HTTP/" {
		return true
	}

	return false
}

// parseHTTP 解析 HTTP 数据并按行输出
func parseHTTP(data []byte) {
	fmt.Println("HTTP数据:")
	// 将字节流转换为 *bytes.Reader
	reader := bytes.NewReader(data)

	// 使用 bufio.NewReader 包装 *bytes.Reader
	bufReader := bufio.NewReader(reader)

	req, err := http.ReadRequest(bufReader)
	if err != nil {
		fmt.Printf("  无法解析为HTTP请求: %s\n", err)
		return
	}

	// 打印HTTP请求头
	fmt.Printf("  方法: %s\n", req.Method)
	fmt.Printf("  URL: %s\n", req.URL.String())
	fmt.Printf("  协议: %s\n", req.Proto)
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("  Header: %s: %s\n", key, value)
		}
	}

	// 打印HTTP正文（如果有）
	if req.Body != nil {
		body := new(bytes.Buffer)
		_, err := body.ReadFrom(req.Body)
		if err == nil {
			fmt.Printf("  正文: %s\n", body.String())
		}
	}
}
