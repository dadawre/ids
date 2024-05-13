package parser

type Configparser struct {
	Ttokens     *[]Token
	curr        int
	configVars  []*configVar
	includePath []*string
}

type configVar struct {
	token *Token
	name  string
	value interface{}
}

func NewConfigparser(tokens *[]Token) *Configparser {
	newConfigParser := Configparser{Ttokens: tokens}
	newConfigParser.includePath = make([]*string, 0, 0)
	newConfigParser.configVars = make([]*configVar, 0, 0)
	return &newConfigParser
}

func (c *Configparser) ParseSelect() {
	if c.curr >= len(*c.Ttokens) {
		return
	}
	if (*c.Ttokens)[c.curr].Type == TokenKeywordInclude {
		c.includePath = append(c.includePath, c.parseInclude())
	} else if (*c.Ttokens)[c.curr].Type == TokenKeywordVar {
		c.configVars = append(c.configVars, c.parseKeywordVar())
	}
}

func (c *Configparser) parseKeywordVar() *configVar {
	NewConfigVar := configVar{}
	NewConfigVar.token = &((*c.Ttokens)[c.curr])

	c.curr++
	if c.curr <= len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type == TokenIdentifier {
		NewConfigVar.name = ((*c.Ttokens)[c.curr]).Value
	}
	c.curr++
	if c.curr <= len(*c.Ttokens) && ((*c.Ttokens)[c.curr]).Type == TokenStringValue {
		NewConfigVar.value = ((*c.Ttokens)[c.curr]).Value
	}
	return &NewConfigVar
}

func (c *Configparser) parseInclude() *string {
	curString := ""
	if c.curr+1 < len(*c.Ttokens) {
		c.curr += 1
		for currToken := (*c.Ttokens)[c.curr]; c.curr <= len(*c.Ttokens) && (currToken.Type == TokenStringValue || currToken.Type == TokenVarUse); currToken = (*c.Ttokens)[c.curr] {
			if currToken.Type == TokenStringValue {
				curString += currToken.Value
			} else {
				for _, currConfigVar := range c.configVars {
					if currConfigVar.token.Type == TokenKeywordVar || currConfigVar.token.Value == currToken.Value {
						curString += currConfigVar.value.(string)
					}
				}
			}
		}
		c.curr++
	}
	return &curString
}

//
//func (c *Configparser[T]) parseIpv4() configVar[net.IP] {
//	ipv4CIDRStr := "192.168.1.0"
//	ip, ipnet, err := net.ParseCIDR(ipv4CIDRStr)
//
//	if err != nil {
//		fmt.Println("Error parsing CIDR:", err)
//		return configVar[net.IP]{}
//	}
//
//	// Get the IP address
//	ipv4 := ip.To4()
//	if ipv4 == nil {
//		fmt.Println("Not an IPv4 address")
//		return configVar[net.IP]{}
//	}
//
//	// Print IPv4 address and subnet mask
//	fmt.Println(ip)
//	fmt.Println("IPv4 Address:", ipv4)
//	fmt.Println("Subnet Mask:", ipnet.Mask)
//	fmt.Println(ipv4[0])
//	fmt.Println(ipv4[1])
//	fmt.Println(ipv4[2])
//	fmt.Println(ipv4[3])
//	fmt.Println(ipnet.Contains(ipv4))
//	return configVar[net.IP]{}
//}
