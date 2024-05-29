package handler

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"graduationDesign/ruleAction"
	"graduationDesign/ruleMatch"
	"graduationDesign/structure/ruleStructure"
)

type PacketHandler struct {
	handlerName string
	isStart     bool
	packetChan  chan *gopacket.Packet
	//protocolMap   map[gopacket.LayerType]interface{}
	layersContain                   gopacket.DecodingLayerContainer
	eth                             *layers.Ethernet
	arp                             *layers.ARP
	ipv4                            *layers.IPv4
	ipv6                            *layers.IPv6
	icmpv4                          *layers.ICMPv4
	icmpv6                          *layers.ICMPv6
	tcp                             *layers.TCP
	udp                             *layers.UDP
	dns                             *layers.DNS
	decoder                         gopacket.DecodingLayerFunc
	decoded                         []gopacket.LayerType
	ruleLinkedList                  *ruleStructure.RuleLinkedList
	MapLayerTypeToProtocolStructure map[gopacket.LayerType]interface{}
	RuleBools                       []bool
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

	p.decoder = p.layersContain.LayersDecoder(layers.LayerTypeEthernet, gopacket.NilDecodeFeedback)
	p.decoded = make([]gopacket.LayerType, 0, 20)

	p.MapLayerTypeToProtocolStructure = map[gopacket.LayerType]interface{}{
		layers.LayerTypeEthernet: p.eth,
		layers.LayerTypeARP:      p.arp,
		layers.LayerTypeIPv4:     p.ipv4,
		layers.LayerTypeIPv6:     p.ipv6,
		layers.LayerTypeTCP:      p.tcp,
		layers.LayerTypeUDP:      p.udp,
	}
	p.RuleBools = make([]bool, 20, 20)
}

func (p *PacketHandler) ClearProtocolStructure() {
	for _, layerType := range p.decoded {
		if layerType == layers.LayerTypeEthernet {
			p.eth = &layers.Ethernet{}
		} else if layerType == layers.LayerTypeARP {
			p.arp = &layers.ARP{}
		} else if layerType == layers.LayerTypeIPv4 {
			p.ipv4 = &layers.IPv4{}
		} else if layerType == layers.LayerTypeIPv6 {
			p.ipv6 = &layers.IPv6{}
		} else if layerType == layers.LayerTypeTCP {
			p.tcp = &layers.TCP{}
		} else if layerType == layers.LayerTypeUDP {
			p.udp = &layers.UDP{}
		}
	}
}

func (p *PacketHandler) PayloadMatch() {
	for node := p.ruleLinkedList.GetHead(); node != nil; node = node.Next {
		for index, protocol := range node.RuleValue.Protocols {
			p.RuleBools[index] = ruleMatch.MapLayerTypeToMatchFunc[protocol.NameId](protocol, p.MapLayerTypeToProtocolStructure[protocol.NameId])
			logic, _ := ruleMatch.ExecuteLogic(p.RuleBools, node.RuleValue.Operator)
			if logic {
				p.ClearProtocolStructure()
				ruleAction.PacketActionMap[node.RuleValue.Action](node.RuleValue)
			}
		}
	}
	ruleAction.PassPacket(nil)
	p.ClearProtocolStructure()
}

func (p *PacketHandler) HandlePacket() {

	//decoder := p.layersContain.LayersDecoder(layers.LayerTypeEthernet, gopacket.NilDecodeFeedback)
	//decoded := make([]gopacket.LayerType, 0, 20)

	for packet := range p.packetChan {
		p.decoder((*packet).Data(), &p.decoded)
		p.PayloadMatch()
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

func NewHandlerPacket(PacketChan chan *gopacket.Packet, handlerName string, ruleList *ruleStructure.RuleLinkedList) *PacketHandler {
	newHandlerPacket := &PacketHandler{
		packetChan:     PacketChan,
		handlerName:    handlerName,
		isStart:        false,
		ruleLinkedList: ruleList,
	}
	newHandlerPacket.initialContainer()
	return newHandlerPacket
}
