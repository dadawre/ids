package ParserRuler

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// 逻辑运算常量
const (
	BoolPlaceholder = iota //bool值占位符
	AndOperator
	OrOperator
	OpenParenOperator
	CloseParenOperator
)

// 行为常量标识符
const (
	PassAction = iota + 1
	RecordAction
	AlertAction
	SendAction
	DropAction
)

// mapStringToAction 行动字符到行动标识符的映射
var mapStringToAction = map[string]int8{
	"pass":   PassAction,
	"record": RecordAction,
	"alert":  AlertAction,
	"send":   SendAction,
	"drop":   DropAction,
}

// mapStrLogicToInt8 逻辑运算符字符串到int8的映射
var mapStrLogicToInt8 = map[string]int8{
	"&": AndOperator,
	"|": OrOperator,
	"(": OpenParenOperator,  // 表示左括号
	")": CloseParenOperator, // 表示右括号
}

// mapStrToLayerType 协议字符串到协议标识符的映射
var mapStrToLayerType = map[string]gopacket.LayerType{
	"Eth":  layers.LayerTypeEthernet,
	"IPv4": layers.LayerTypeIPv4,
	"IPv6": layers.LayerTypeIPv6,
	"TCP":  layers.LayerTypeTCP,
	"UDP":  layers.LayerTypeUDP,
	"OS":   10000,
}

// MapLayerTypeToParseType 协议标识符到协议字段标识符解析函数的映射
var MapLayerTypeToParseType = map[gopacket.LayerType]func(fieldName string) int8{
	layers.LayerTypeEthernet: parseEthType,
	layers.LayerTypeIPv4:     parseIPv4Type,
	layers.LayerTypeTCP:      parseTcpType,
	layers.LayerTypeUDP:      parseUDPType,
	layers.LayerTypeIPv6:     parseIPv6Type,
	layers.LayerTypeARP:      parseARPType,
	10000:                    parseOSType,
}

// MapLayerTypeToParseFunc 定义协议标识符到协议字段解析函数的映射
var MapLayerTypeToParseFunc = map[gopacket.LayerType]func(fieldType int8, str string) (interface{}, error){
	layers.LayerTypeEthernet: parseEthFieldValue,
	layers.LayerTypeIPv4:     parseIPv4FieldValue,
	layers.LayerTypeTCP:      parseTCPFieldValue,
	layers.LayerTypeUDP:      parseUDPFieldValue,
	layers.LayerTypeIPv6:     parseIPv6FieldValue,
	layers.LayerTypeARP:      parseARPFieldValue,
	10000:                    parseOSFieldValue, //操作系统参数
}

// MapMsgParaStringToProtocolFieldId 消息参数字段字符串转参数标识符映射
var MapMsgParaStringToProtocolFieldId = map[string]int8{
	"Eth.src":    EthSrcMACField,
	"Eth.dst":    EthDstMACField,
	"Eth.length": EthLengthField,
	"Eth.type":   EthTypeField,

	"ARP.type":      ARPTypeField,
	"ARP.protocol":  ARPProtocolField,
	"ARP.hsize":     ARPHwAddressSizeField,
	"ARP.psize":     ARPProtocolAddressSizeField,
	"ARP.option":    ARPOptionField,
	"ARP.shaddress": ARPSourceHwAddressField,
	"ARP.spaddress": ARPSourceProtocolAddressField,
	"ARP.dhaddress": ARPDestHwAddressField,
	"ARP.dpaddress": ARPDestProtocolAddressField,

	"IPv4.src":      IPv4SrcIPField,
	"IPv4.dst":      IPv4DstIPField,
	"IPv4.checksum": IPv4ChecksumField,
	"IPv4.protocol": IPv4ProtocolField,
	"IPv4.TTL":      IPv4TTLField,
	"IPv4.offset":   IPv4OffsetField,
	"IPv4.flags":    IPv4FlagField,
	"IPv4.id":       IPv4IDField,
	"IPv4.length":   IPv4TotalLengthField,
	"IPv4.TOS":      IPv4TOSField,
	"IPv4.IHL":      IPv4IHLField,
	"IPv4.version":  IPv4VersionField,

	"IPv6.dst":      IPv6DstField,
	"IPv6.src":      IPv6SrcField,
	"IPv6.hlimit":   IPv6HopLimitField,
	"IPv6.nexthead": IPv6NextHeadField,
	"IPv6.length":   IPv6LengthField,
	"IPv6.label":    IPv6FlowLabelField,
	"IPv6.class":    IPv6TrafficClassField,
	"IPv6.version":  IPv6VersionField,

	"TCP.src":      TCPSrcPortField,
	"TCP.dst":      TCPDstPortField,
	"TCP.seq":      TCPSequenceField,
	"TCP.seqack":   TCPAckSeqField,
	"TCP.offset":   TCPOffsetField,
	"TCP.URG":      TCPUrgField,
	"TCP.ACK":      TCPAckField,
	"TCP.PSH":      TCPPstField,
	"TCP.RST":      TCPRstField,
	"TCP.SYN":      TCPSynField,
	"TCP.FIN":      TCPFinField,
	"TCP.window":   TCPWindowField,
	"TCP.checksum": TCPChecksumField,
	"TCP.urgptr":   TCPUrgPtrField,

	"UDP.src":      UDPSrcField,
	"UDP.dst":      UDPDstField,
	"UDP.length":   UDPLengthField,
	"UDP.checksum": UDPChecksumField,

	"OS.cpu":    OSCpuUsageRateField,
	"OS.gpu":    OSGpuUsageRateField,
	"OS.memory": OSMemoryUsageRateField,
}

// 协议字段标识符
const (
	EthSrcMACField = iota + 1
	EthDstMACField
	EthLengthField
	EthTypeField

	ARPTypeField = iota + 1
	ARPProtocolField
	ARPHwAddressSizeField
	ARPProtocolAddressSizeField
	ARPOptionField
	ARPSourceHwAddressField
	ARPSourceProtocolAddressField
	ARPDestHwAddressField
	ARPDestProtocolAddressField

	IPv4SrcIPField = iota + 1
	IPv4DstIPField
	IPv4ChecksumField
	IPv4ProtocolField
	IPv4TTLField
	IPv4OffsetField
	IPv4FlagField
	IPv4IDField
	IPv4TotalLengthField
	IPv4TOSField
	IPv4IHLField
	IPv4VersionField

	IPv6DstField = iota + 1
	IPv6SrcField
	IPv6HopLimitField
	IPv6NextHeadField
	IPv6LengthField
	IPv6FlowLabelField
	IPv6TrafficClassField
	IPv6VersionField

	TCPSrcPortField = iota + 1
	TCPDstPortField
	TCPSequenceField
	TCPAckSeqField
	TCPOffsetField
	TCPUrgField
	TCPAckField
	TCPPstField
	TCPRstField
	TCPSynField
	TCPFinField
	TCPWindowField
	TCPChecksumField
	TCPUrgPtrField

	UDPSrcField = iota + 1
	UDPDstField
	UDPLengthField
	UDPChecksumField

	OSCpuUsageRateField
	OSGpuUsageRateField
	OSMemoryUsageRateField
)
