package ParserRuler

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/gopacket/layers"
	"graduationDesign/back_end/Parser/ParserConfig"
	"net"
	"strconv"
	"strings"
)

// ParseFieldError 定义解析错误类型
type ParseFieldError string

func (p ParseFieldError) Error() string {
	return string(p)
}

func parseOSType(fieldName string) int8 {
	switch fieldName {
	case "cpu":
		return OSCpuUsageRateField
	case "gpu":
		return OSGpuUsageRateField
	case "memory":
		return OSMemoryUsageRateField
	default:
		return -1
	}
}

func parseOSFieldValue(fieldType int8, str string) (interface{}, error) {
	// 去掉百分号
	str = strings.TrimSuffix(str, "%")

	// 将字符串解析为 float64
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("解析错误:", err)
		return nil, err
	}

	switch fieldType {
	case OSCpuUsageRateField:
		// 转换为小数形式
		return f, nil
	case OSGpuUsageRateField:
		return f, nil
	case OSMemoryUsageRateField:
		return f, nil
	default:
		return nil, err
	}

}

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
		return &mac, nil
	case EthLengthField:
		// 将整数字符串转换为 int 类型
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num16 := uint16(num)
		return num16, nil
	case EthTypeField:
		// 将整数字符串转换为 int 类型
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num16 := uint16(num)
		ethType := layers.EthernetType(num16)
		return ethType, nil
	default:
		return nil, ParseFieldError("ETH fieldTypeError")
	}
}

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
	switch fieldType {
	case IPv4SrcIPField, IPv4DstIPField:
		ipNet, err := ParserConfig.IPStringToIPNet(str, false)
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
		num, err := stringToUint32(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		//num32 := uint32(num)
		return num, nil

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
		return nil, ParseFieldError("TCP fieldTypeError")
	}
}

func parseUDPType(fieldName string) int8 {
	switch fieldName {
	case "src":
		return UDPSrcField
	case "dst":
		return UDPDstField
	case "length":
		return UDPLengthField
	case "checksum":
		return UDPChecksumField
	default:
		return -1
	}
}

func parseUDPFieldValue(fieldType int8, str string) (interface{}, error) {
	switch fieldType {
	case UDPSrcField, UDPDstField:
		port, err := strconv.Atoi(str)
		if err != nil {
			return 0, errors.New("invalid port: not a number")
		}
		udpPort := layers.UDPPort(port)
		if err = Validate(layers.TCPPort(udpPort)); err != nil {
			return nil, err
		}
		return udpPort, nil
	case UDPLengthField, UDPChecksumField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num16 := uint16(num)
		return num16, nil
	default:
		return nil, ParseFieldError("UDP fieldTypeError")
	}
}

func parseIPv6Type(fieldName string) int8 {
	switch fieldName {
	case "dst":
		return IPv6DstField
	case "src":
		return IPv6SrcField
	case "hlimit":
		return IPv6HopLimitField
	case "nexthead":
		return IPv6NextHeadField
	case "length":
		return IPv6LengthField
	case "label":
		return IPv6FlowLabelField
	case "class":
		return IPv6TrafficClassField
	case "version":
		return IPv6VersionField
	default:
		return -1
	}
}

func parseIPv6FieldValue(fieldType int8, str string) (interface{}, error) {
	switch fieldType {
	case IPv6DstField, IPv6SrcField:
		ipNet, err := ParserConfig.IPStringToIPNet(str, false)
		if err != nil {
			return nil, err
		}
		return ipNet, nil

	case IPv6VersionField, IPv6TrafficClassField, IPv6HopLimitField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num8 := uint8(num)
		return num8, nil

	case IPv6FlowLabelField:
		num, err := stringToUint32(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		//num32 := uint32(num)
		return num, nil

	case IPv6LengthField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num16 := uint16(num)
		return num16, nil

	case IPv6NextHeadField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, err
		}
		// 将 int 类型转换为 int8 类型
		num8 := uint8(num)
		ipv4protocol := layers.IPProtocol(num8)
		return ipv4protocol, nil

	default:
		return nil, ParseFieldError("UDP fieldTypeError")
	}
}

func parseARPType(fieldName string) int8 {
	switch fieldName {
	case "type":
		return ARPTypeField
	case "protocol":
		return ARPProtocolField
	case "hsize":
		return ARPHwAddressSizeField
	case "psize":
		return ARPProtocolAddressSizeField
	case "option":
		return ARPOptionField
	case "shaddress":
		return ARPSourceHwAddressField
	case "spaddress":
		return ARPSourceProtocolAddressField
	case "dhaddress":
		return ARPDestHwAddressField
	case "dpaddress":
		return ARPDestProtocolAddressField
	default:
		return -1
	}
}

// ParseIPString parses an IP string and returns the corresponding byte array.
func ParseIPString(ipStr string) ([]byte, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}
	// IPv4 addresses are 4 bytes, IPv6 are 16 bytes.
	if ipv4 := ip.To4(); ipv4 != nil {
		return ipv4, nil
	}
	return ip.To16(), nil
}

// ParseMACString parses a MAC address string and returns the corresponding byte array.
func ParseMACString(macStr string) ([]byte, error) {
	macStr = strings.ToLower(macStr)
	macStr = strings.ReplaceAll(macStr, ":", "")
	macStr = strings.ReplaceAll(macStr, "-", "")
	mac, err := hex.DecodeString(macStr)
	if err != nil || len(mac) != 6 {
		return nil, fmt.Errorf("invalid MAC address: %s", macStr)
	}
	return mac, nil
}

func parseARPFieldValue(fieldType int8, str string) (interface{}, error) {
	switch fieldType {
	case ARPTypeField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num8 := uint8(num)
		linkType := layers.LinkType(num8)
		return linkType, nil

	case ARPProtocolField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num16 := uint16(num)
		EthType := layers.EthernetType(num16)
		return EthType, nil

	case ARPHwAddressSizeField, ARPProtocolAddressSizeField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num8 := uint8(num)
		return num8, nil

	case ARPOptionField:
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to convert string to int: %v\n", err)
			return nil, err
		}
		num16 := uint16(num)
		return num16, nil
	case ARPSourceHwAddressField, ARPDestHwAddressField:
		mac, err := ParseMACString(str)
		if err != nil {
			return nil, err
		}
		return mac, nil

	case ARPSourceProtocolAddressField, ARPDestProtocolAddressField:
		ip, err := ParseIPString(str)
		if err != nil {
			return nil, err
		}
		return ip, nil

	default:
		return nil, ParseFieldError("UDP fieldTypeError")
	}
}

// 将字符串转换为uint64，第二个参数是进制，第三个参数是位大小
func stringToUint32(s string) (uint32, error) {

	u64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	// 将uint64转换为uint32
	return uint32(u64), nil
}
