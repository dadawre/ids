package handler

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"graduationDesign/ruleMatch"
	"graduationDesign/structure/ruleStructure"
)

var LogPath string

// PacketHandler 进行数据包处理的函数，功能包括提取数据包各层协议并与规则进行匹配
type PacketHandler struct {
	handlerName                     string
	isStart                         bool
	packetChan                      chan *gopacket.Packet
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
	PacketActionMap                 map[int8]func(rule *ruleStructure.Rule)
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
	p.RuleBools = make([]bool, 10, 10)

	// PacketActionMap 创建一个map，将枚举值映射到对应的函数
	p.PacketActionMap = map[int8]func(rule *ruleStructure.Rule){
		PacketPassAction:   p.PassPacket,
		PacketRecordAction: p.RecordPacket,
		PacketAlertAction:  p.AlertPacket,
		PacketSendAction:   p.SendPacket,
		PacketDropAction:   p.DropPacket,
	}
}

// PayloadMatch 进行匹配操作
func (p *PacketHandler) PayloadMatch() {
	for node := p.ruleLinkedList.GetHead(); node != nil; node = node.Next {
		for index, protocol := range node.RuleValue.Protocols {
			for indexProtocol := range p.decoded {
				if protocol.NameId == p.decoded[indexProtocol] {
					p.RuleBools[index] = ruleMatch.MapLayerTypeToMatchFunc[protocol.NameId](protocol, p.MapLayerTypeToProtocolStructure[protocol.NameId])
				}
			}
		}
		logic, _ := ruleMatch.ExecuteLogic(p.RuleBools, node.RuleValue.Operator)

		if logic {
			fmt.Println(p.RuleBools, "    ", node.RuleValue.Operator)
			p.PacketActionMap[node.RuleValue.Action](node.RuleValue)
			return
		}
	}
	p.PassPacket(nil)
}

// HandlePacket 一个协程函数负责接收数据包并进行解包
func (p *PacketHandler) HandlePacket() {

	//decoder := p.layersContain.LayersDecoder(layers.LayerTypeEthernet, gopacket.NilDecodeFeedback)
	//decoded := make([]gopacket.LayerType, 0, 20)

	for packet := range p.packetChan {
		p.decoder((*packet).Data(), &p.decoded)
		fmt.Println(p.decoded)
		p.PayloadMatch()
	}
}

// StartHandle 开启数据包处理功能
func (p *PacketHandler) StartHandle() bool {
	if p.isStart == false {
		go p.HandlePacket()
		p.isStart = true
		return true
	}
	return false
}

// NewHandlerPacket 创建一个新的数据包处理对象
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
