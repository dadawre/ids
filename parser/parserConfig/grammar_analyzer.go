package parserConfig

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Configparser struct {
	Ttokens               *[]Token
	curr                  int
	TconfigVars           []*configVar
	TincludePath          []*string
	TConfigMapNameToValue map[string]interface{}
}

type configVar struct {
	Ttoken *Token
	Tname  string
	Tvalue interface{}
}

func NewConfigparser(tokens *[]Token) *Configparser {
	newConfigParser := Configparser{Ttokens: tokens}
	newConfigParser.TincludePath = make([]*string, 0, 0)
	newConfigParser.TconfigVars = make([]*configVar, 0, 0)
	newConfigParser.TConfigMapNameToValue = make(map[string]interface{}, 0)
	return &newConfigParser
}

type ParseError string

func (s *ParseError) Error() string {
	return fmt.Sprintf("parse error : %s", string(*s))
}

func (c *Configparser) ParseSelect() error {
	if c.curr >= len(*c.Ttokens) {
		parseError := ParseError("This is already the last token")
		return &parseError
	}
	fmt.Println("curr:", c.curr, "  token:", (*c.Ttokens)[c.curr])
	if (*c.Ttokens)[c.curr].Type == TokenKeywordInclude {
		includeString := c.parseInclude()
		if includeString == nil {
			err := ParseError("something is wrong with parse include")
			return &err
		}
		c.TincludePath = append(c.TincludePath, includeString)
	} else if (*c.Ttokens)[c.curr].Type == TokenKeywordVar {
		configvar := c.parseKeywordVar()
		if configvar == nil {
			err := ParseError("something is wrong with parse var")
			return &err
		}
		c.TconfigVars = append(c.TconfigVars, configvar)
	} else if (*c.Ttokens)[c.curr].Type == TokenKeywordIPv4Var {
		IPValue := c.parserKeywordIp(true)
		if IPValue == nil {
			err := ParseError("something is wrong with parse ipv4")
			return &err
		}
		c.TconfigVars = append(c.TconfigVars, IPValue)
	} else if (*c.Ttokens)[c.curr].Type == TokenKeywordIPv6Var {
		IPValue := c.parserKeywordIp(false)
		if IPValue == nil {
			err := ParseError("something is wrong with parse ipv6")
			return &err
		}
		c.TconfigVars = append(c.TconfigVars, IPValue)
	} else if (*c.Ttokens)[c.curr].Type == TokenKeywordPortVar {
		portArray := c.parserKeywordPort()
		if portArray == nil {
			err := ParseError("something is wrong with parse port")
			return &err
		}
		c.TconfigVars = append(c.TconfigVars, portArray)
	} else if (*c.Ttokens)[c.curr].Type == TokenKeywordConfig {
		c.curr++
		var err error
		var haveConfPara bool
		for haveConfPara, err = c.parseConfig(); haveConfPara; haveConfPara, err = c.parseConfig() {
		}
		return err
	}
	return nil
}

func (c *Configparser) parseConfig() (bool, error) {
	var Key string
	if c.curr < len(*c.Ttokens) && (*c.Ttokens)[c.curr].Type == TokenIdentifier {
		Key = (*c.Ttokens)[c.curr].Value
		c.curr++
	} else {
		return false, nil
	}
	if c.curr < len(*c.Ttokens) && (*c.Ttokens)[c.curr].Type == TokenStringValue {
		Value := (*c.Ttokens)[c.curr].Value
		c.TConfigMapNameToValue[Key] = Value
		c.curr++
	} else {
		err := ParseError("something is wrong with config")
		return false, &err
	}
	return true, nil
}

func (c *Configparser) parseKeywordVar() *configVar {
	NewConfigVar := configVar{}
	NewConfigVar.Ttoken = &((*c.Ttokens)[c.curr])

	c.curr++
	if c.curr < len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type == TokenIdentifier {
		NewConfigVar.Tname = ((*c.Ttokens)[c.curr]).Value
	} else {
		return nil
	}
	c.curr++
	if c.curr < len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type == TokenStringValue {
		NewConfigVar.Tvalue = ((*c.Ttokens)[c.curr]).Value
	} else {
		return nil
	}
	c.curr++
	return &NewConfigVar
}

func (c *Configparser) parseInclude() *string {
	curString := ""
	if c.curr+1 < len(*c.Ttokens) {
		c.curr += 1
		for c.curr < len(*c.Ttokens) {
			currToken := (*c.Ttokens)[c.curr]
			if currToken.Type != TokenStringValue && currToken.Type != TokenVarUse {
				break
			}
			if currToken.Type == TokenStringValue {
				curString += currToken.Value
			} else {
				for _, currConfigVar := range c.TconfigVars {
					if currConfigVar.Ttoken.Type == TokenKeywordVar || currConfigVar.Ttoken.Value == currToken.Value {
						curString += currConfigVar.Tvalue.(string)
					}
				}
			}
			c.curr++
		}
	}
	return &curString
}

// InvalidIPAddressError 自定义错误类型
type InvalidIPAddressError string

func (e *InvalidIPAddressError) Error() string {
	return fmt.Sprintf("invalid IPv4 address: %s", string(*e))
}

func IPStringToIPNet(ipStr string, isIpv4 bool) (*net.IPNet, error) {
	_, subnet, err := net.ParseCIDR(ipStr)
	if err == nil {
		return subnet, nil
	}

	var maskStr string
	if isIpv4 {
		maskStr = "255.255.255.255"
	} else {
		maskStr = "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"
	}

	networkIP := net.ParseIP(ipStr)

	if networkIP == nil {
		return nil, fmt.Errorf("invalid network IP address: %s", ipStr)
	}
	//fmt.Println(networkIP)

	maskIP := net.ParseIP(maskStr)
	if maskIP == nil {
		return nil, fmt.Errorf("invalid mask IP address: %s", maskStr)
	}
	//fmt.Println(maskIP)
	// Convert mask IP to net.IPMask
	var mask net.IPMask
	if isIpv4 {
		mask = net.IPMask(maskIP.To4())
	} else {
		mask = net.IPMask(maskIP.To16())
	}
	//fmt.Println(mask)

	if !(len(mask) == 4) && !(len(mask) == 16) {
		return nil, fmt.Errorf("invalid mask length: %d", len(mask))
	}

	// Create an IPNet from the network IP and mask
	ipNet := &net.IPNet{
		IP:   networkIP,
		Mask: mask,
	}
	return ipNet, nil
}

func (c *Configparser) parserKeywordIp(isIpv4 bool) *configVar {
	NewConfigVar := configVar{}
	NewConfigVar.Ttoken = &((*c.Ttokens)[c.curr])

	c.curr++
	if c.curr < len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type == TokenIdentifier {
		NewConfigVar.Tname = ((*c.Ttokens)[c.curr]).Value
	} else {
		return nil
	}
	c.curr++
	if c.curr < len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type >= TokenIpv4Value && ((*c.Ttokens)[c.curr]).Type <= TokenIpv6NetworkValue {
		IPvalue, err := IPStringToIPNet(((*c.Ttokens)[c.curr]).Value, isIpv4)
		if err != nil {
			return nil
		}
		NewConfigVar.Tvalue = IPvalue
	} else {
		return nil
	}
	c.curr++
	return &NewConfigVar

}

func StrToArr(arrStr string) []int {
	// 去除前后的方括号
	arrStr = strings.Trim(arrStr, "[]")

	// 使用逗号分割字符串
	parts := strings.Split(arrStr, ",")

	// 创建一个整数数组存储结果
	var result []int

	// 遍历分割后的字符串，将每个部分转换为整数并添加到数组中
	for _, part := range parts {
		number, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return nil
		}
		result = append(result, number)
	}

	// 输出结果数组
	//fmt.Println("Result array:", result)
	return result
}

func (c *Configparser) parserKeywordPort() *configVar {
	NewConfigVar := configVar{}
	NewConfigVar.Ttoken = &((*c.Ttokens)[c.curr])

	c.curr++
	if c.curr < len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type == TokenIdentifier {
		NewConfigVar.Tname = ((*c.Ttokens)[c.curr]).Value
	} else {
		return nil
	}
	c.curr++
	if c.curr < len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type == TokenStringValue {
		portArray := StrToArr(((*c.Ttokens)[c.curr]).Value)
		if portArray == nil {
			return nil
		}
		NewConfigVar.Tvalue = portArray
	} else {
		return nil
	}
	c.curr++
	return &NewConfigVar
}

func (c *Configparser) ParseAllTokens() {
	var err error
	for err == nil {
		err = c.ParseSelect()
	}
}
