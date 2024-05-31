package parserConfig

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Configparser struct {
	ConfigTokens         *[]Token
	curr                 int
	ConfigVars           map[string]*ConfigVar
	IncludePath          []*string
	ConfigMapNameToValue map[string]string
}

type ConfigVar struct {
	VarToken *Token
	VarName  string
	VarValue interface{}
}

// NewConfigparser 新建一个配置文件语法分析器
func NewConfigparser(tokens *[]Token) *Configparser {
	newConfigParser := Configparser{ConfigTokens: tokens}
	newConfigParser.IncludePath = make([]*string, 0, 0)
	newConfigParser.ConfigVars = make(map[string]*ConfigVar, 0)
	newConfigParser.ConfigMapNameToValue = make(map[string]string, 0)
	return &newConfigParser
}

// ParseError 语法分析错误对象
type ParseError string

func (s *ParseError) Error() string {
	return fmt.Sprintf("parse error : %s", string(*s))
}

// ParseSelect 选择对Token进行哪种语法分析
func (c *Configparser) ParseSelect() error {
	if c.curr >= len(*c.ConfigTokens) {
		parseError := ParseError("This is already the last token")
		return &parseError
	}
	//fmt.Println("curr:", c.curr, "  token:", (*c.ConfigTokens)[c.curr])
	if (*c.ConfigTokens)[c.curr].Type == TokenKeywordInclude {
		includeString := c.parseInclude()
		if includeString == nil {
			err := ParseError("something is wrong with parse include")
			return &err
		}
		c.IncludePath = append(c.IncludePath, includeString)
	} else if (*c.ConfigTokens)[c.curr].Type == TokenKeywordVar {
		configvar := c.parseKeywordVar()
		if configvar == nil {
			err := ParseError("something is wrong with parse var")
			return &err
		}
		//c.ConfigVars = append(c.ConfigVars, configvar)
		c.ConfigVars[configvar.VarName] = configvar
	} else if (*c.ConfigTokens)[c.curr].Type == TokenKeywordIPv4Var {
		IPValue := c.parserKeywordIp(true)
		if IPValue == nil {
			err := ParseError("something is wrong with parse ipv4")
			return &err
		}
		//c.ConfigVars = append(c.ConfigVars, IPValue)
		c.ConfigVars[IPValue.VarName] = IPValue
	} else if (*c.ConfigTokens)[c.curr].Type == TokenKeywordIPv6Var {
		IPValue := c.parserKeywordIp(false)
		if IPValue == nil {
			err := ParseError("something is wrong with parse ipv6")
			return &err
		}
		//c.ConfigVars = append(c.ConfigVars, IPValue)
		c.ConfigVars[IPValue.VarName] = IPValue
	} else if (*c.ConfigTokens)[c.curr].Type == TokenKeywordPortVar {
		portArray := c.parserKeywordPort()
		if portArray == nil {
			err := ParseError("something is wrong with parse port")
			return &err
		}
		//c.ConfigVars = append(c.ConfigVars, portArray)
		c.ConfigVars[portArray.VarName] = portArray
	} else if (*c.ConfigTokens)[c.curr].Type == TokenKeywordConfig {
		c.curr++
		var err error
		var haveConfPara bool
		for haveConfPara, err = c.parseConfigVar(); haveConfPara; haveConfPara, err = c.parseConfigVar() {
		}
		return err
	}
	return nil
}

// 解析配置文件中的环境参数
func (c *Configparser) parseConfigVar() (bool, error) {
	var Key string
	if c.curr < len(*c.ConfigTokens) && (*c.ConfigTokens)[c.curr].Type == TokenIdentifier {
		Key = (*c.ConfigTokens)[c.curr].Value
		c.curr++
	} else {
		return false, nil
	}
	if c.curr < len(*c.ConfigTokens) && (*c.ConfigTokens)[c.curr].Type == TokenStringValue {
		Value := (*c.ConfigTokens)[c.curr].Value
		c.ConfigMapNameToValue[Key] = Value
		c.curr++
	} else {
		err := ParseError("something is wrong with config")
		return false, &err
	}
	return true, nil
}

// 解析配置文件中一般关键字
func (c *Configparser) parseKeywordVar() *ConfigVar {
	NewConfigVar := ConfigVar{}
	NewConfigVar.VarToken = &((*c.ConfigTokens)[c.curr])

	c.curr++
	if c.curr < len(*c.ConfigTokens) && ((*c.ConfigTokens)[c.curr]).Type == TokenIdentifier {
		NewConfigVar.VarName = ((*c.ConfigTokens)[c.curr]).Value
	} else {
		return nil
	}
	c.curr++
	if c.curr < len(*c.ConfigTokens) && ((*c.ConfigTokens)[c.curr]).Type == TokenStringValue {
		NewConfigVar.VarValue = ((*c.ConfigTokens)[c.curr]).Value
	} else {
		return nil
	}
	c.curr++
	return &NewConfigVar
}

// 解析配置文件中规则路径
func (c *Configparser) parseInclude() *string {
	curString := ""
	if c.curr+1 < len(*c.ConfigTokens) {
		c.curr += 1
		for c.curr < len(*c.ConfigTokens) {
			currToken := (*c.ConfigTokens)[c.curr]
			if currToken.Type != TokenStringValue && currToken.Type != TokenVarUse {
				break
			}
			if currToken.Type == TokenStringValue {
				curString += currToken.Value
			} else {
				for _, currConfigVar := range c.ConfigVars {
					if currConfigVar.VarToken.Type == TokenKeywordVar && currConfigVar.VarName == currToken.Value {
						curString += currConfigVar.VarValue.(string)
					}
				}
			}
			c.curr++
		}
	}
	return &curString
}

// InvalidIPAddressError 解析IP错误类型
type InvalidIPAddressError string

func (e *InvalidIPAddressError) Error() string {
	return fmt.Sprintf("invalid IPv4 address: %s", string(*e))
}

// IPStringToIPNet 将字符串解析为net.IPNet,支持IPv6
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

// 解析IP
func (c *Configparser) parserKeywordIp(isIpv4 bool) *ConfigVar {
	NewConfigVar := ConfigVar{}
	NewConfigVar.VarToken = &((*c.ConfigTokens)[c.curr])

	c.curr++
	if c.curr < len(*c.ConfigTokens) && ((*c.ConfigTokens)[c.curr]).Type == TokenIdentifier {
		NewConfigVar.VarName = ((*c.ConfigTokens)[c.curr]).Value
	} else {
		return nil
	}
	c.curr++
	if c.curr < len(*c.ConfigTokens) && ((*c.ConfigTokens)[c.curr]).Type >= TokenIpv4Value && ((*c.ConfigTokens)[c.curr]).Type <= TokenIpv6NetworkValue {
		IPvalue, err := IPStringToIPNet(((*c.ConfigTokens)[c.curr]).Value, isIpv4)
		if err != nil {
			return nil
		}
		NewConfigVar.VarValue = IPvalue
	} else {
		return nil
	}
	c.curr++
	return &NewConfigVar

}

// StrToArr 解析字符串到数组
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

func (c *Configparser) parserKeywordPort() *ConfigVar {
	NewConfigVar := ConfigVar{}
	NewConfigVar.VarToken = &((*c.ConfigTokens)[c.curr])

	c.curr++
	if c.curr < len(*c.ConfigTokens) && ((*c.ConfigTokens)[c.curr]).Type == TokenIdentifier {
		NewConfigVar.VarName = ((*c.ConfigTokens)[c.curr]).Value
	} else {
		return nil
	}
	c.curr++
	if c.curr < len(*c.ConfigTokens) && ((*c.ConfigTokens)[c.curr]).Type == TokenStringValue {
		portArray := StrToArr(((*c.ConfigTokens)[c.curr]).Value)
		if portArray == nil {
			return nil
		}
		NewConfigVar.VarValue = portArray
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
