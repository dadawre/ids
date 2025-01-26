package Handler

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"graduationDesign/back_end/OSController"
	RuleMatch2 "graduationDesign/back_end/RuleMatch"
	"graduationDesign/back_end/structure/ruleStructure"
)

var PacketHandlerArr []*PacketHandler //数据包解析结构体数组

// PacketHandler 进行数据包处理结构体，功能包括提取数据包各层协议并与规则进行匹配
type PacketHandler struct {
	handlerName string //数据处理器名称
	isStart     bool   //是否启动

	packetChan chan *gopacket.Packet //接收数据的管程
	//parser          *gopacket.DecodingLayerParser
	layersContainer gopacket.DecodingLayerContainer //数据栈解析容器
	decoder         gopacket.DecodingLayerFunc      //数据栈解析函数
	decoded         []gopacket.LayerType            //数据栈各层使用协议

	ruleLinkedList                  *ruleStructure.RuleLinkedList                          //规则头指针
	mapLayerTypeToProtocolStructure map[gopacket.LayerType]interface{}                     //协议标识符到协议结构体映射
	ruleBools                       []bool                                                 //规则匹配bool
	packetActionMap                 map[int8]func(rule *ruleStructure.RuleNodeValue) error // Action标识符到Action函数的映射
	eth                             *layers.Ethernet                                       //以太网结构体
	arp                             *layers.ARP                                            //ARP结构体
	ipv4                            *layers.IPv4                                           //IPV4结构体
	ipv6                            *layers.IPv6                                           //IPV6结构体
	icmpv4                          *layers.ICMPv4                                         //ICMPv4结构体
	icmpv6                          *layers.ICMPv6                                         //ICMPv6结构体
	tcp                             *layers.TCP                                            //TCP结构体
	udp                             *layers.UDP                                            //UDP结构体
	dns                             *layers.DNS                                            //DNS结构体
	payload                         *gopacket.Payload                                      //载荷结构体
}

// initialContainer 初始化协议层容器和解码器
func (p *PacketHandler) initialContainer() {
	// 创建协议层
	p.eth = &layers.Ethernet{}
	p.arp = &layers.ARP{}
	p.ipv4 = &layers.IPv4{}
	p.ipv6 = &layers.IPv6{}
	p.icmpv4 = &layers.ICMPv4{}
	p.icmpv6 = &layers.ICMPv6{}
	p.tcp = &layers.TCP{}
	p.udp = &layers.UDP{}
	p.dns = &layers.DNS{}
	p.payload = &gopacket.Payload{}

	// 创建一个新的 DecodingLayerMap 实例
	p.layersContainer = gopacket.DecodingLayerMap{}

	// 将各个层次加入到容器中，链式调用方式
	p.layersContainer = p.layersContainer.Put(p.eth).
		Put(p.arp).
		Put(p.ipv4).
		Put(p.ipv6).
		Put(p.icmpv4).
		Put(p.icmpv6).
		Put(p.udp).
		Put(p.tcp).
		Put(p.dns)

	// 初始化 decoder
	p.decoder = p.layersContainer.LayersDecoder(layers.LayerTypeEthernet, gopacket.NilDecodeFeedback)

	// 初始化解码的协议层类型切片
	p.decoded = make([]gopacket.LayerType, 0, 20)

	p.mapLayerTypeToProtocolStructure = map[gopacket.LayerType]interface{}{
		layers.LayerTypeEthernet: p.eth,
		layers.LayerTypeARP:      p.arp,
		layers.LayerTypeIPv4:     p.ipv4,
		layers.LayerTypeIPv6:     p.ipv6,
		layers.LayerTypeTCP:      p.tcp,
		layers.LayerTypeUDP:      p.udp,
		layers.LayerTypeICMPv4:   p.icmpv4,
		layers.LayerTypeICMPv6:   p.icmpv6,
		layers.LayerTypeDNS:      p.dns,
		10000:                    OSController.State,
	}

	p.ruleBools = make([]bool, 20)

	p.packetActionMap = map[int8]func(rule *ruleStructure.RuleNodeValue) error{
		PacketPassAction:   p.PassPacket,
		PacketRecordAction: p.RecordPacket,
		PacketAlertAction:  p.AlertPacket,
		PacketSendAction:   p.SendPacket,
		PacketDropAction:   p.DropIPv4Packets,
	}
}

// PayloadMatch 进行匹配操作
func (p *PacketHandler) PayloadMatch() {
	for node := p.ruleLinkedList.GetHead(); node != nil; node = node.Next { //遍历规则链表
		for index, protocol := range node.RuleValue.Protocols { //和规则的各个部分进行匹配
			for _, decodedLayer := range p.decoded { //每层协议
				if protocol.NameId == decodedLayer {
					matchFunc, exists := RuleMatch2.MapLayerTypeToMatchFunc[protocol.NameId]
					if !exists {
						fmt.Printf("No match function for protocol: %v\n", protocol.NameId)
						continue
					}
					p.ruleBools[index] = matchFunc(&protocol, p.mapLayerTypeToProtocolStructure[protocol.NameId])
				} else if protocol.NameId == 10000 {
					//fmt.Println("OS执行", p.mapLayerTypeToProtocolStructure[10000])
					OSMatchFunc, exists := RuleMatch2.MapLayerTypeToMatchFunc[10000]
					if !exists {
						fmt.Printf("No match function for OS: %v\n", protocol.NameId)
						break
					}
					p.ruleBools[index] = OSMatchFunc(&protocol, p.mapLayerTypeToProtocolStructure[10000])
					break
				}
			}
		}

		result, err := RuleMatch2.EvaluateExpression(p.ruleBools, node.RuleValue.Operators)
		//if err != nil {
		//	fmt.Printf("Error executing result: %v\n", err)
		//	continue
		//}

		if err != nil {
			fmt.Println("逻辑表达式有错误:", err)
			continue
		} //else {
		//	fmt.Println("逻辑表达式结果:", result)
		//}

		if result {
			actionFunc, exists := p.packetActionMap[node.RuleValue.Action]
			if exists {
				actionFunc(node.RuleValue)
			} else {
				fmt.Printf("No action function for action: %d\n", node.RuleValue.Action)
			}
			return
		}
	}
	p.defaultPass()
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
