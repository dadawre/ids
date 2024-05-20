package main

import (
	"fmt"
	"graduationDesign/parser"
	"net"
	"testing"
)

func TestConfigParser(t *testing.T) {

	myTokens := make([]parser.Token, 0, 0)
	myTokens = append(myTokens, parser.Token{Type: parser.TokenKeywordVar, Value: "var"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenIdentifier, Value: "homePath"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenStringValue, Value: "/home/dadawre"})

	myTokens = append(myTokens, parser.Token{Type: parser.TokenKeywordInclude, Value: "include"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenStringValue, Value: "/first"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenVarUse, Value: "homePath"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenStringValue, Value: "/home1111"})

	myTokens = append(myTokens, parser.Token{Type: parser.TokenKeywordIPv4Var, Value: "ipv4var"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenIdentifier, Value: "thisNetworkIpv4"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenIpv4Value, Value: "192.168.80.1"})

	myTokens = append(myTokens, parser.Token{Type: parser.TokenKeywordIPv6Var, Value: "ipv6var"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenIdentifier, Value: "thisNetworkIpv6"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenIpv6NetworkValue, Value: "AA22:BB11:1122:CDEF:1234:AA99:7654:7410/64"})

	myConfigparser := parser.NewConfigparser(&myTokens)
	fmt.Println(*(myConfigparser.Ttokens))

	fmt.Println("--------------------------------------------")
	myConfigparser.ParseSelect()
	fmt.Println(myConfigparser.TconfigVars[0].Ttoken)
	fmt.Println(myConfigparser.TconfigVars[0].Tvalue.(string))

	fmt.Println("--------------------------------------------")
	myConfigparser.ParseSelect()
	fmt.Println(*(myConfigparser.TincludePath)[0])
	fmt.Println()

	fmt.Println("--------------------------------------------")
	myConfigparser.ParseSelect()
	fmt.Println((myConfigparser.TconfigVars)[1].Ttoken.Type)
	fmt.Println((myConfigparser.TconfigVars)[1].Ttoken.Value)
	fmt.Println((myConfigparser.TconfigVars)[1].Tname)
	getIPnet := (myConfigparser.TconfigVars)[1].Tvalue.(*net.IPNet)
	fmt.Println(getIPnet.IP)
	fmt.Println(getIPnet.Mask)
	hostIP := net.ParseIP("192.168.80.1")
	fmt.Println(getIPnet.Contains(hostIP))

	fmt.Println("--------------------------------------------")
	myConfigparser.ParseSelect()
	fmt.Println((myConfigparser.TconfigVars)[2].Ttoken.Type)
	fmt.Println((myConfigparser.TconfigVars)[2].Ttoken.Value)
	fmt.Println((myConfigparser.TconfigVars)[2].Tname)
	getIPnet = (myConfigparser.TconfigVars)[2].Tvalue.(*net.IPNet)
	fmt.Println(getIPnet.IP)
	fmt.Println(getIPnet.Mask)
	hostIP = net.ParseIP("AA22:BB11:1122:CDEF:1234:AA99:7654:7410")
	fmt.Println(getIPnet.Contains(hostIP))

}
func TestIPv6(t *testing.T) {
	// 定义一个带子网掩码的冒号分隔16进制IPv6地址段字符串
	ipv6CIDR := "2001:0db8:85a3:0000:0000:8a2e:0370:7334/128"

	// 使用net.ParseCIDR将IPv6地址段字符串解析为IP和IPNet
	ip, ipNet, err := net.ParseCIDR(ipv6CIDR)
	if err != nil {
		fmt.Println("Invalid IPv6 CIDR notation:", err)
		return
	}

	fmt.Printf("IP address: %v\n", ip)
	fmt.Printf("Network: %v\n", ipNet)

	// 如果只需要IP地址部分
	fmt.Printf("IP address only: %v\n", ip)

	// 如果只需要子网掩码部分
	fmt.Printf("Network only: %v\n", ipNet.Mask)
}
