package parserRuler

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"graduationDesign/parser/parserConfig"
	"net"
	"strconv"
)

var mapStrToLayerType = map[string]gopacket.LayerType{
	"Eth":  layers.LayerTypeEthernet,
	"IPv4": layers.LayerTypeIPv4,
	"IPv6": layers.LayerTypeIPv6,
	"TCP":  layers.LayerTypeTCP,
	"UDP":  layers.LayerTypeUDP,
}

const (
	AndOperator = iota + 1
	ArOperator
)

var mapStrLogicToInt8 = map[string]int8{
	"&": AndOperator,
	"|": ArOperator,
}

const (
	PassAction = iota + 1
	RecordAction
	AlertAction
	SendAction
	DropAction
)

var mapStringToAction = map[string]int8{
	"pass":   PassAction,
	"record": RecordAction,
	"alert":  AlertAction,
	"send":   SendAction,
	"drop":   DropAction,
}

// MapLayerTypeToParseFunc 定义协议标识符到字段标识符的映射
var MapLayerTypeToParseFunc = map[gopacket.LayerType]func(fieldType int8, str string) (interface{}, error){
	layers.LayerTypeEthernet: parseEthFieldValue,
	layers.LayerTypeIPv4:     parseIPv4FieldValue,
	layers.LayerTypeTCP:      parseTCPFieldValue,
}

var MapLayerTypeTOParseType = map[gopacket.LayerType]func(fieldName string) int8{
	layers.LayerTypeEthernet: parseEthType,
	layers.LayerTypeIPv4:     parseIPv4Type,
	layers.LayerTypeTCP:      parseTcpType,
}

type ParseFieldError string

func (p ParseFieldError) Error() string {
	return string(p)
}

// 定义Ethernet字段标识符
const (
	EthSrcMACField = iota + 1
	EthDstMACField
	EthLengthField
	EthTypeField
)

func parseEthType(fieldName string) int8 {
	switch fieldName {
	case "src":
		return EthSrcMACField
	case "dst":
		return EthDstMACField
	case "length":
		return EthLengthField
	case "type":
		return EthTypeField
	default:
		return -1
	}
}

func parseEthFieldValue(fieldType int8, str string) (interface{}, error) {
	switch fieldType {
	case EthSrcMACField, EthDstMACField:
		// 将 MAC 地址字符串转换为 net.HardwareAddr
		mac, err := net.ParseMAC(str)
		if err != nil {
			fmt.Printf("Failed to parse MAC address: %v\n", err)
			return nil, err
		}
		return mac, nil
	case EthTypeField, EthLengthField:
		// 将整数字符串转换为 int 类型
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		return num, nil
	default:
		return nil, ParseFieldError("ETH fieldTypeError")
	}
}

// 定义IPv4字段标识符
const (
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
)

func parseIPv4Type(fieldName string) int8 {
	switch fieldName {
	case "src":
		return IPv4SrcIPField
	case "dst":
		return IPv4DstIPField
	case "checksum":
		return IPv4ChecksumField
	case "protocol":
		return IPv4ProtocolField
	case "TTL":
		return IPv4TTLField
	case "offset":
		return IPv4OffsetField
	case "flags":
		return IPv4FlagField
	case "id":
		return IPv4IDField
	case "length":
		return IPv4TotalLengthField
	case "TOS":
		return IPv4TOSField
	case "IHL":
		return IPv4IHLField
	case "version":
		return IPv4VersionField
	default:
		return -1
	}
}

func parseIPv4FieldValue(fieldType int8, str string) (interface{}, error) {
	fmt.Println(fieldType, "   ", str)
	switch fieldType {
	case IPv4SrcIPField, IPv4DstIPField:
		ipNet, err := parserConfig.IPStringToIPNet(str, false)
		if err != nil {
			return nil, err
		}
		return ipNet, nil
	case IPv4VersionField, IPv4IHLField, IPv4TOSField, IPv4TTLField:
		// 将字符串转换为 int 类型
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, err
		}

		// 将 int 类型转换为 int8 类型
		num8 := uint8(num)
		return num8, nil
	case IPv4FlagField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, err
		}
		// 将 int 类型转换为 int8 类型
		num8 := uint8(num)
		ipv4flag := layers.IPv4Flag(num8)
		return ipv4flag, nil

	case IPv4ProtocolField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, err
		}
		// 将 int 类型转换为 int8 类型
		num8 := uint8(num)
		ipv4protocol := layers.IPProtocol(num8)
		return ipv4protocol, nil

	case IPv4TotalLengthField, IPv4IDField, IPv4OffsetField, IPv4ChecksumField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, err
		}
		// 将 int 类型转换为 int8 类型
		num16 := uint16(num)
		return num16, nil
	default:
		err := ParseFieldError("IPv4 parseFailed")
		return nil, err
	}

}

const (
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
)

func parseTcpType(fieldName string) int8 {
	switch fieldName {
	case "src":
		return TCPSrcPortField
	case "dst":
		return TCPDstPortField
	case "seq":
		return TCPSequenceField
	case "seqack":
		return TCPAckSeqField
	case "offset":
		return TCPOffsetField
	case "URG":
		return TCPUrgField
	case "ACK":
		return TCPAckField
	case "PSH":
		return TCPPstField
	case "RST":
		return TCPRstField
	case "SYN":
		return TCPSynField
	case "FIN":
		return TCPFinField
	case "window":
		return TCPWindowField
	case "checksum":
		return TCPChecksumField
	case "urgptr":
		return TCPUrgPtrField
	default:
		return -1
	}
}

// Validate 验证端口号是否在有效范围内
func Validate(p layers.TCPPort) error {
	if p < 1 || p > 65535 {
		return errors.New("invalid port: out of range")
	}
	return nil
}

func parseTCPFieldValue(fieldType int8, str string) (interface{}, error) {
	switch fieldType {
	case TCPSrcPortField, TCPDstPortField:
		port, err := strconv.Atoi(str)
		if err != nil {
			return 0, errors.New("invalid port: not a number")
		}
		tcpPort := layers.TCPPort(port)
		if err = Validate(tcpPort); err != nil {
			return nil, err
		}
		return tcpPort, nil

	case TCPAckSeqField, TCPSequenceField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num32 := uint32(num)
		return num32, nil

	case TCPAckField, TCPSynField, TCPUrgField, TCPPstField, TCPFinField, TCPRstField:
		if str == "1" {
			return true, nil
		} else if str == "0" {
			return false, nil
		}
		return nil, ParseFieldError("TCP Parser Error")

	case TCPWindowField, TCPChecksumField, TCPUrgPtrField:
		// 将整数字符串转换为 int 类型
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num16 := uint16(num)
		return num16, nil

	case TCPOffsetField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num8 := uint8(num)
		return num8, nil

	default:
		return nil, ParseFieldError("ETH fieldTypeError")
	}
}
