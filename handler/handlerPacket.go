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
}

func (p *PacketHandler) initialContainer() {
	var eth layers.Ethernet
	var arp layers.ARP
	var ipv4 layers.IPv4
	var ipv6 layers.IPv6
	var icmpv4 layers.ICMPv4
	var icmpv6 layers.ICMPv6
	var tcp layers.TCP
	var udp layers.UDP
	var dns layers.DNS

	p.layersContain = gopacket.DecodingLayerContainer(gopacket.DecodingLayerMap{})
	p.layersContain = p.layersContain.Put(&eth)
	p.layersContain = p.layersContain.Put(&arp)
	p.layersContain = p.layersContain.Put(&ipv4)
	p.layersContain = p.layersContain.Put(&ipv6)
	p.layersContain = p.layersContain.Put(&icmpv4)
	p.layersContain = p.layersContain.Put(&icmpv6)
	p.layersContain = p.layersContain.Put(&udp)
	p.layersContain = p.layersContain.Put(&tcp)
	p.layersContain = p.layersContain.Put(&dns)
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
			if layerType == layers.LayerTypeEthernet {
				fmt.Println("this is Ethernet")
			} else if layerType == layers.LayerTypeARP {
				fmt.Println("this is ARP")
			} else if layerType == layers.LayerTypeIPv4 {
				fmt.Println("this is IPV4")
			} else if layerType == layers.LayerTypeIPv6 {
				fmt.Println("this is IPV6")
			} else if layerType == layers.LayerTypeICMPv4 {
				fmt.Println("this is ICMPV4")
			} else if layerType == layers.LayerTypeICMPv6 {
				fmt.Println("this is ICMPV6")
			} else if layerType == layers.LayerTypeTCP {
				fmt.Println("this is TCP")
			} else if layerType == layers.LayerTypeUDP {
				fmt.Println("this is UDP")
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
