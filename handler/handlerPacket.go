package handler

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"os"
)

type PacketHandler struct {
	handlerName   string
	isStart       bool
	packetChan    chan *gopacket.Packet
	protocolMap   map[gopacket.LayerType]interface{}
	layersContain gopacket.DecodingLayerContainer
	eth           *layers.Ethernet
	arp           *layers.ARP
	ipv4          *layers.IPv4
	ipv6          *layers.IPv6
	icmpv4        *layers.ICMPv4
	icmpv6        *layers.ICMPv6
	tcp           *layers.TCP
	udp           *layers.UDP
	dns           *layers.DNS
}

func (p *PacketHandler) initialContainer() {

	p.eth = &layers.Ethernet{}
	p.arp = &layers.ARP{}
	p.ipv4 = &layers.IPv4{}
	p.ipv6 = &layers.IPv6{}
	p.icmpv4 = &layers.ICMPv4{}
	p.icmpv6 = &layers.ICMPv6{}
	p.tcp = &layers.TCP{}
	p.udp = &layers.UDP{}
	p.dns = &layers.DNS{}

	p.layersContain = gopacket.DecodingLayerContainer(gopacket.DecodingLayerMap{})
	p.layersContain = p.layersContain.Put(p.eth)
	p.layersContain = p.layersContain.Put(p.arp)
	p.layersContain = p.layersContain.Put(p.ipv4)
	p.layersContain = p.layersContain.Put(p.ipv6)
	p.layersContain = p.layersContain.Put(p.icmpv4)
	p.layersContain = p.layersContain.Put(p.icmpv6)
	p.layersContain = p.layersContain.Put(p.udp)
	p.layersContain = p.layersContain.Put(p.tcp)
	p.layersContain = p.layersContain.Put(p.dns)
}

func (p *PacketHandler) HandlePacket() {

	decoder := p.layersContain.LayersDecoder(layers.LayerTypeEthernet, gopacket.NilDecodeFeedback)
	decoded := make([]gopacket.LayerType, 0, 20)

	for packet := range p.packetChan {
		lt, err := decoder((*packet).Data(), &decoded)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not decode layers: %v\n", err)
		}

		if lt != gopacket.LayerTypeZero {
			fmt.Fprintf(os.Stderr, "unknown layer type: %v\n", lt)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not decode layers: %v\n", err)
		}
		fmt.Println("handlerName:" + p.handlerName)
		//fmt.Println(decoded)
		for _, layerType := range decoded {
			fmt.Println("------------------------------------------")
			if layerType == layers.LayerTypeEthernet {
				fmt.Println("this is Ethernet")
				fmt.Println(p.eth.EthernetType)
				fmt.Println(p.eth.Length)
				fmt.Println(p.eth.SrcMAC)
				fmt.Println(p.eth.DstMAC)
				fmt.Println(p.eth.Payload)  // Payload是除去头部的所有字节载荷
				fmt.Println(p.eth.Contents) // Contents是头部的所有字节载荷
			} else if layerType == layers.LayerTypeARP {
				fmt.Println("this is ARP")
			} else if layerType == layers.LayerTypeIPv4 {
				fmt.Println("this is IPV4")
				fmt.Println(p.ipv4.SrcIP)
				fmt.Println(p.ipv4.DstIP)
				fmt.Println(p.ipv4.Checksum)
				fmt.Println(p.ipv4.Protocol)
				fmt.Println(p.ipv4.TTL)
				fmt.Println(p.ipv4.FragOffset)
				fmt.Println(p.ipv4.Flags)
				fmt.Println(p.ipv4.Id)
				fmt.Println(p.ipv4.Length)
				fmt.Println(p.ipv4.TOS)
				fmt.Println(p.ipv4.IHL)
				fmt.Println(p.ipv4.Version)
				//fmt.Println("Src IP :", p.ipv4.SrcIP, "  Dest IP :", p.ipv4.DstIP)
				fmt.Printf("IPv4 Source IP: %s, IPv4 Destination IP: %s\n", p.ipv4.SrcIP, p.ipv4.DstIP)
			} else if layerType == layers.LayerTypeIPv6 {
				fmt.Println("this is IPV6")
			} else if layerType == layers.LayerTypeICMPv4 {
				fmt.Println("this is ICMPV4")
			} else if layerType == layers.LayerTypeICMPv6 {
				fmt.Println("this is ICMPV6")
			} else if layerType == layers.LayerTypeTCP {
				fmt.Println("this is TCP")
				fmt.Println(p.tcp.SrcPort)
				fmt.Println(p.tcp.DstPort)
				fmt.Println(p.tcp.Seq)
				fmt.Println(p.tcp.Ack)
				fmt.Println(p.tcp.DataOffset)
				fmt.Println(p.tcp.URG)
				fmt.Println(p.tcp.ACK)
				fmt.Println(p.tcp.PSH)
				fmt.Println(p.tcp.RST)
				fmt.Println(p.tcp.SYN)
				fmt.Println(p.tcp.FIN)
				fmt.Println(p.tcp.Window)
				fmt.Println(p.tcp.Checksum)
				fmt.Println(p.tcp.Urgent)
				fmt.Println(p.tcp.Options)
			} else if layerType == layers.LayerTypeUDP {
				fmt.Println("this is UDP")
				fmt.Println(p.udp.SrcPort)
				fmt.Println(p.udp.DstPort)
				fmt.Println(p.udp.Length)
				fmt.Println(p.udp.Checksum)
			} else if layerType == layers.LayerTypeDNS {
				fmt.Println("this is DNS")
			} else {
				fmt.Println("unknow Protocol")
			}
		}
		fmt.Println("----------------------------------------------")
	}
}

func (p *PacketHandler) StartHandle() bool {
	if p.isStart == false {
		go p.HandlePacket()
		p.isStart = true
		return true
	}
	return false
}

func NewHandlerPacket(PacketChan chan *gopacket.Packet, handlerName string) *PacketHandler {
	newHandlerPacket := &PacketHandler{
		packetChan:  PacketChan,
		handlerName: handlerName,
		isStart:     false,
	}
	newHandlerPacket.initialContainer()
	return newHandlerPacket
}
