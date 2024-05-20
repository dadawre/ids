package parser

import (
	"fmt"
	"net"
)

type Configparser struct {
	Ttokens      *[]Token
	curr         int
	TconfigVars  []*configVar
	TincludePath []*string
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
	return &newConfigParser
}

func (c *Configparser) ParseSelect() {
	fmt.Println("curr:", c.curr, "  token:", (*c.Ttokens)[c.curr])
	if c.curr >= len(*c.Ttokens) {
		return
	}
	if (*c.Ttokens)[c.curr].Type == TokenKeywordInclude {
		c.TincludePath = append(c.TincludePath, c.parseInclude())
	} else if (*c.Ttokens)[c.curr].Type == TokenKeywordVar {
		c.TconfigVars = append(c.TconfigVars, c.parseKeywordVar())
	} else if (*c.Ttokens)[c.curr].Type == TokenKeywordIPv4Var {
		c.TconfigVars = append(c.TconfigVars, c.parserKeywordIp(true))
	} else if (*c.Ttokens)[c.curr].Type == TokenKeywordIPv6Var {
		c.TconfigVars = append(c.TconfigVars, c.parserKeywordIp(false))
	}
}

func (c *Configparser) parseKeywordVar() *configVar {
	NewConfigVar := configVar{}
	NewConfigVar.Ttoken = &((*c.Ttokens)[c.curr])

	c.curr++
	if c.curr < len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type == TokenIdentifier {
		NewConfigVar.Tname = ((*c.Ttokens)[c.curr]).Value
	}
	c.curr++
	if c.curr < len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type == TokenStringValue {
		NewConfigVar.Tvalue = ((*c.Ttokens)[c.curr]).Value
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

func (e InvalidIPAddressError) Error() string {
	return fmt.Sprintf("invalid IPv4 address: %s", string(e))
}

func IPStringToIPNet(ipStr string, isIpv4 bool) (*net.IPNet, error) {
	_, subnet, err := net.ParseCIDR(ipStr)
	if err == nil {
		return subnet, nil
	}
	fmt.Println("ggggggggggggggggggggggggggggggggggggggggg")
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
	fmt.Println(networkIP)

	maskIP := net.ParseIP(maskStr)
	if maskIP == nil {
		return nil, fmt.Errorf("invalid mask IP address: %s", maskStr)
	}
	fmt.Println(maskIP)
	// Convert mask IP to net.IPMask
	var mask net.IPMask
	if isIpv4 {
		mask = net.IPMask(maskIP.To4())
	} else {
		mask = net.IPMask(maskIP.To16())
	}
	fmt.Println(mask)

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
	}
	c.curr++
	if c.curr < len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type >= TokenIpv4Value && ((*c.Ttokens)[c.curr]).Type <= TokenIpv6NetworkValue {
		NewConfigVar.Tvalue, _ = IPStringToIPNet(((*c.Ttokens)[c.curr]).Value, isIpv4)
	}
	c.curr++
	return &NewConfigVar

}
