package parserConfig

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Configparser 表示配置文件语法分析器
type Configparser struct {
	ConfigTokens         *[]Token
	curr                 int
	ConfigVars           map[string]*ConfigVar
	IncludePath          []*string
	ConfigMapNameToValue map[string]string
}

// ConfigVar 表示配置变量
type ConfigVar struct {
	VarToken *Token
	VarName  string
	VarValue interface{}
}

// NewConfigparser 新建一个配置文件语法分析器
func NewConfigparser(tokens *[]Token) *Configparser {
	return &Configparser{
		ConfigTokens:         tokens,
		IncludePath:          []*string{},
		ConfigVars:           make(map[string]*ConfigVar),
		ConfigMapNameToValue: make(map[string]string),
	}
}

// ParseError 语法分析错误对象
type ParseError string

func (s ParseError) Error() string {
	return fmt.Sprintf("parse error: %s", string(s))
}

// ParseSelect 选择对Token进行哪种语法分析
func (c *Configparser) ParseSelect() error {
	if c.curr >= len(*c.ConfigTokens) {
		return ParseError("this is already the last token")
	}

	switch (*c.ConfigTokens)[c.curr].Type {
	case TokenKeywordInclude:
		if includeString := c.parseInclude(); includeString != nil {
			c.IncludePath = append(c.IncludePath, includeString)
		} else {
			return ParseError("something is wrong with parse include")
		}
	case TokenKeywordVar:
		if configvar := c.parseKeywordVar(); configvar != nil {
			c.ConfigVars[configvar.VarName] = configvar
		} else {
			return ParseError("something is wrong with parse var")
		}
	case TokenKeywordIPv4Var:
		if IPValue := c.parseKeywordIP(true); IPValue != nil {
			c.ConfigVars[IPValue.VarName] = IPValue
		} else {
			return ParseError("something is wrong with parse ipv4")
		}
	case TokenKeywordIPv6Var:
		if IPValue := c.parseKeywordIP(false); IPValue != nil {
			c.ConfigVars[IPValue.VarName] = IPValue
		} else {
			return ParseError("something is wrong with parse ipv6")
		}
	case TokenKeywordPortVar:
		if portArray := c.parseKeywordPort(); portArray != nil {
			c.ConfigVars[portArray.VarName] = portArray
		} else {
			return ParseError("something is wrong with parse port")
		}
	case TokenKeywordConfig:
		c.curr++
		for {
			if haveConfPara, err := c.parseConfigVar(); !haveConfPara {
				return err
			}
		}
	}
	return nil
}

// parseConfigVar 解析配置文件中的环境参数
func (c *Configparser) parseConfigVar() (bool, error) {
	if c.curr >= len(*c.ConfigTokens) || (*c.ConfigTokens)[c.curr].Type != TokenIdentifier {
		return false, nil
	}
	Key := (*c.ConfigTokens)[c.curr].Value
	c.curr++

	if c.curr < len(*c.ConfigTokens) && (*c.ConfigTokens)[c.curr].Type == TokenStringValue {
		c.ConfigMapNameToValue[Key] = (*c.ConfigTokens)[c.curr].Value
		c.curr++
		return true, nil
	}
	return false, ParseError("something is wrong with config")
}

// parseKeywordVar 解析配置文件中一般关键字
func (c *Configparser) parseKeywordVar() *ConfigVar {
	if c.curr+2 >= len(*c.ConfigTokens) {
		return nil
	}

	NewConfigVar := &ConfigVar{VarToken: &(*c.ConfigTokens)[c.curr]}
	c.curr++

	if (*c.ConfigTokens)[c.curr].Type == TokenIdentifier {
		NewConfigVar.VarName = (*c.ConfigTokens)[c.curr].Value
		c.curr++
	} else {
		return nil
	}

	if (*c.ConfigTokens)[c.curr].Type == TokenStringValue {
		NewConfigVar.VarValue = (*c.ConfigTokens)[c.curr].Value
		c.curr++
	} else {
		return nil
	}

	return NewConfigVar
}

// parseInclude 解析配置文件中规则路径
func (c *Configparser) parseInclude() *string {
	curString := ""
	c.curr++
	for c.curr < len(*c.ConfigTokens) {
		currToken := (*c.ConfigTokens)[c.curr]
		if currToken.Type != TokenStringValue && currToken.Type != TokenVarUse {
			break
		}

		if currToken.Type == TokenStringValue {
			curString += currToken.Value
		} else {
			if configVar, exists := c.ConfigVars[currToken.Value]; exists && configVar.VarToken.Type == TokenKeywordVar {
				curString += configVar.VarValue.(string)
			}
		}
		c.curr++
	}
	return &curString
}

// InvalidIPAddressError 解析IP错误类型
type InvalidIPAddressError string

func (e InvalidIPAddressError) Error() string {
	return fmt.Sprintf("invalid IP address: %s", string(e))
}

// IPStringToIPNet 将字符串解析为net.IPNet,支持IPv6
func IPStringToIPNet(ipStr string, isIpv4 bool) (*net.IPNet, error) {
	_, subnet, err := net.ParseCIDR(ipStr)
	if err == nil {
		return subnet, nil
	}

	maskStr := "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"
	if isIpv4 {
		maskStr = "255.255.255.255"
	}

	networkIP := net.ParseIP(ipStr)
	if networkIP == nil {
		return nil, fmt.Errorf("invalid network IP address: %s", ipStr)
	}

	maskIP := net.ParseIP(maskStr)
	if maskIP == nil {
		return nil, fmt.Errorf("invalid mask IP address: %s", maskStr)
	}

	mask := net.IPMask(maskIP.To16())
	if isIpv4 {
		mask = net.IPMask(maskIP.To4())
	}

	if len(mask) != 4 && len(mask) != 16 {
		return nil, fmt.Errorf("invalid mask length: %d", len(mask))
	}

	return &net.IPNet{IP: networkIP, Mask: mask}, nil
}

// parseKeywordIP 解析IP
func (c *Configparser) parseKeywordIP(isIpv4 bool) *ConfigVar {
	if c.curr+2 >= len(*c.ConfigTokens) {
		return nil
	}

	NewConfigVar := &ConfigVar{VarToken: &(*c.ConfigTokens)[c.curr]}
	c.curr++

	if (*c.ConfigTokens)[c.curr].Type == TokenIdentifier {
		NewConfigVar.VarName = (*c.ConfigTokens)[c.curr].Value
		c.curr++
	} else {
		return nil
	}

	if (*c.ConfigTokens)[c.curr].Type >= TokenIpv4Value && (*c.ConfigTokens)[c.curr].Type <= TokenIpv6NetworkValue {
		if IPvalue, err := IPStringToIPNet((*c.ConfigTokens)[c.curr].Value, isIpv4); err == nil {
			NewConfigVar.VarValue = IPvalue
			c.curr++
			return NewConfigVar
		}
	}
	return nil
}

// StrToArr 解析字符串到数组
func StrToArr(arrStr string) []int {
	arrStr = strings.Trim(arrStr, "[]")
	parts := strings.Split(arrStr, ",")
	result := make([]int, len(parts))

	for i, part := range parts {
		if number, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			result[i] = number
		} else {
			fmt.Println("Error converting string to integer:", err)
			return nil
		}
	}
	return result
}

// parseKeywordPort 解析端口关键字
func (c *Configparser) parseKeywordPort() *ConfigVar {
	if c.curr+2 >= len(*c.ConfigTokens) {
		return nil
	}

	NewConfigVar := &ConfigVar{VarToken: &(*c.ConfigTokens)[c.curr]}
	c.curr++

	if (*c.ConfigTokens)[c.curr].Type == TokenIdentifier {
		NewConfigVar.VarName = (*c.ConfigTokens)[c.curr].Value
		c.curr++
	} else {
		return nil
	}

	if (*c.ConfigTokens)[c.curr].Type == TokenStringValue {
		if portArray := StrToArr((*c.ConfigTokens)[c.curr].Value); portArray != nil {
			NewConfigVar.VarValue = portArray
			c.curr++
			return NewConfigVar
		}
	}
	return nil
}

// ParseAllTokens 解析所有Token
func (c *Configparser) ParseAllTokens() {
	for {
		if err := c.ParseSelect(); err != nil {
			break
		}
	}
}
