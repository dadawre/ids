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
	layersContainer                 gopacket.DecodingLayerContainer
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
	mapLayerTypeToProtocolStructure map[gopacket.LayerType]interface{}
	ruleBools                       []bool
	// Action标识符到Action函数的映射
	packetActionMap map[int8]func(rule *ruleStructure.Rule)
}

// initialContainer 初始化协议层容器和解码器
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

	p.layersContainer = gopacket.DecodingLayerContainer(gopacket.DecodingLayerMap{})
	p.layersContainer = p.layersContainer.Put(p.eth)
	p.layersContainer = p.layersContainer.Put(p.arp)
	p.layersContainer = p.layersContainer.Put(p.ipv4)
	p.layersContainer = p.layersContainer.Put(p.ipv6)
	p.layersContainer = p.layersContainer.Put(p.icmpv4)
	p.layersContainer = p.layersContainer.Put(p.icmpv6)
	p.layersContainer = p.layersContainer.Put(p.udp)
	p.layersContainer = p.layersContainer.Put(p.tcp)
	p.layersContainer = p.layersContainer.Put(p.dns)

	p.decoder = p.layersContainer.LayersDecoder(layers.LayerTypeEthernet, gopacket.NilDecodeFeedback)
	p.decoded = make([]gopacket.LayerType, 0, 20)

	p.mapLayerTypeToProtocolStructure = map[gopacket.LayerType]interface{}{
		layers.LayerTypeEthernet: p.eth,
		layers.LayerTypeARP:      p.arp,
		layers.LayerTypeIPv4:     p.ipv4,
		layers.LayerTypeIPv6:     p.ipv6,
		layers.LayerTypeTCP:      p.tcp,
		layers.LayerTypeUDP:      p.udp,
	}
	p.ruleBools = make([]bool, 10)

	p.packetActionMap = map[int8]func(rule *ruleStructure.Rule){
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
			for _, decodedLayer := range p.decoded {
				if protocol.NameId == decodedLayer {
					matchFunc, exists := ruleMatch.MapLayerTypeToMatchFunc[protocol.NameId]
					if !exists {
						fmt.Printf("No match function for protocol: %v\n", protocol.NameId)
						continue
					}
					p.ruleBools[index] = matchFunc(protocol, p.mapLayerTypeToProtocolStructure[protocol.NameId])
				}
			}
		}

		logic, err := ruleMatch.ExecuteLogic(p.ruleBools, node.RuleValue.Operator)
		if err != nil {
			fmt.Printf("Error executing logic: %v\n", err)
			continue
		}

		if logic {
			actionFunc, exists := p.packetActionMap[node.RuleValue.Action]
			if exists {

				actionFunc(node.RuleValue)
			} else {
				fmt.Printf("No action function for action: %d\n", node.RuleValue.Action)
			}
			return
		}
	}
	p.PassPacket(nil)
}

// HandlePacket 一个协程函数负责接收数据包并进行解包
func (p *PacketHandler) HandlePacket() {
	for packet := range p.packetChan {
		_, err := p.decoder((*packet).Data(), &p.decoded)
		if err != nil {
			fmt.Printf("Error decoding packet: %v\n", err)
			continue
		}

		fmt.Println(p.decoded)
		p.PayloadMatch()
	}
}

// StartHandle 开启数据包处理功能
func (p *PacketHandler) StartHandle() bool {
	if !p.isStart {
		go p.HandlePacket()
		p.isStart = true
		return true
	}
	return false
}

// NewHandlerPacket 创建一个新的数据包处理对象
func NewHandlerPacket(packetChan chan *gopacket.Packet, handlerName string, ruleList *ruleStructure.RuleLinkedList) *PacketHandler {
	newHandlerPacket := &PacketHandler{
		packetChan:     packetChan,
		handlerName:    handlerName,
		isStart:        false,
		ruleLinkedList: ruleList,
	}
	newHandlerPacket.initialContainer()
	return newHandlerPacket
}
