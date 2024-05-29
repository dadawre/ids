package main

import (
	"fmt"
	"graduationDesign/parser/parserConfig"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TestConfigParser(t *testing.T) {

	myTokens := make([]parserConfig.Token, 0, 0)
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenKeywordVar, Value: "var"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenIdentifier, Value: "homePath"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenStringValue, Value: "/home/dadawre"})

	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenKeywordInclude, Value: "include"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenStringValue, Value: "/first"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenVarUse, Value: "homePath"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenStringValue, Value: "/home1111"})

	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenKeywordIPv4Var, Value: "ipv4var"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenIdentifier, Value: "thisNetworkIpv4"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenIpv4Value, Value: "192.168.80.1"})

	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenKeywordIPv6Var, Value: "ipv6var"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenIdentifier, Value: "thisNetworkIpv6"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenIpv6NetworkValue, Value: "AA22:BB11:1122:CDEF:1234:AA99:7654:7410/64"})

	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenKeywordPortVar, Value: "portvar"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenIdentifier, Value: "httpPort"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenStringValue, Value: "[80,433]"})

	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenKeywordConfig, Value: "config"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenIdentifier, Value: "max_coroutine"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenStringValue, Value: "1024"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenIdentifier, Value: "min_coroutine"})
	myTokens = append(myTokens, parserConfig.Token{Type: parserConfig.TokenStringValue, Value: "512"})

	myConfigparser := parserConfig.NewConfigparser(&myTokens)
	fmt.Println(*(myConfigparser.Ttokens))

	fmt.Println("--------------------------------------------")
	myConfigparser.ParseSelect()
	fmt.Println(myConfigparser.TconfigVars[0].Ttoken)
	fmt.Println(myConfigparser.TconfigVars[0].Tname)
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

	fmt.Println("--------------------------------------------")
	myConfigparser.ParseSelect()
	fmt.Println((myConfigparser.TconfigVars)[3].Ttoken.Type)
	fmt.Println((myConfigparser.TconfigVars)[3].Ttoken.Value)
	fmt.Println((myConfigparser.TconfigVars)[3].Tname)
	fmt.Println((myConfigparser.TconfigVars)[3].Tvalue)
	myConfigparser.ParseSelect()
	//myConfigparser.ParseSelect()
	fmt.Println("max_coroutine : ", myConfigparser.TConfigMapNameToValue["max_coroutine"])
	fmt.Println("min_coroutine : ", myConfigparser.TConfigMapNameToValue["min_coroutine"])
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

func TestStrToArr(t *testing.T) {
	// 给定的字符串
	input := "[80,433,555]"

	// 去除前后的方括号
	input = strings.Trim(input, "[]")

	// 使用逗号分割字符串
	parts := strings.Split(input, ",")

	// 创建一个整数数组存储结果
	var result []int

	// 遍历分割后的字符串，将每个部分转换为整数并添加到数组中
	for _, part := range parts {
		number, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}
		result = append(result, number)
	}

	// 输出结果数组
	fmt.Println("Result array:", result)
}

func TestGrammarLexical(t *testing.T) {
	// 打开文件
	//file, err := os.Open("../conf/ids.conf")
	//if err != nil {
	//	fmt.Println("Error opening file:", err)
	//	return
	//}
	//defer func(file *os.File) {
	//	err := file.Close()
	//	if err != nil {
	//		err.Error()
	//	}
	//}(file)

	//reader := bufio.NewReader(file)
	// 创建一个bufio.Reader

	myConfigLexer := parserConfig.NewConfigLexer("../conf/ids.conf")
	confTokens := myConfigLexer.GetAllToken()
	//fmt.Println(myConfigLexer.GetAllToken())

	myConfigparser := parserConfig.NewConfigparser(&confTokens)
	fmt.Println(*(myConfigparser.Ttokens))
	myConfigparser.ParseAllTokens()
	for _, configvar := range myConfigparser.TconfigVars {
		fmt.Println("====================================================")
		fmt.Println(*configvar.Ttoken)
		fmt.Println(configvar.Tname)
		fmt.Println(configvar.Tvalue)
	}

	for _, IncludeString := range myConfigparser.TincludePath {
		fmt.Println("====================================================")
		fmt.Println(*IncludeString)
	}

	for key, value := range myConfigparser.TConfigMapNameToValue {
		fmt.Println("====================================================")

		fmt.Println(key, " :", value)
	}
}
