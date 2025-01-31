package ParserRuler

import (
	"graduationDesign/back_end/Parser/ParserConfig"
	"graduationDesign/back_end/Structure"
	"regexp"
	"strconv"
	"strings"
)

// RuleParser 规则解析结构体
type RuleParser struct {
	curRule       int                                //当前是第几个规则
	pos           int                                //token偏移量
	ruleStr       string                             //规则输入字符串
	Tokens        []string                           //输入字符串转化为token
	protocolRegex *regexp.Regexp                     //协议正则表达式
	fieldRegex    *regexp.Regexp                     //协议字段正则表达式
	operatorRegex *regexp.Regexp                     //逻辑表达式运算符正则表达式
	priorityRegex *regexp.Regexp                     //优先级正则表达式
	actionRegex   *regexp.Regexp                     //行为正则表达式
	nameRegex     *regexp.Regexp                     //规则名正则表达式
	banTimeRegex  *regexp.Regexp                     //封禁时间正则表达式
	msgRegex      *regexp.Regexp                     //输出字符串正则表达式
	stringRegex   *regexp.Regexp                     //string正则表达式
	msgEndRegex   *regexp.Regexp                     //消息结尾正则表达式
	ConfigVarsMap map[string]*ParserConfig.ConfigVar //配置变量映射

}

// RuleParserError 自定义错误类型
type RuleParserError string

func (e RuleParserError) Error() string {
	return string(e)
}

// NewRuleParser 创建并初始化一个新的 RuleParser
func NewRuleParser(ruleStr string, configVar map[string]*ParserConfig.ConfigVar) *RuleParser {
	return &RuleParser{
		pos:           0,
		ruleStr:       ruleStr,
		Tokens:        strings.Fields(ruleStr), //将字符串按空格分割开来存到数组中
		protocolRegex: regexp.MustCompile(`^(Eth|IPv4|IPv6|TCP|UDP|OS)$`),
		fieldRegex:    regexp.MustCompile(`(\b)(\w+)\b (\b|\$)(\w+)\b`),
		operatorRegex: regexp.MustCompile(`[&|()]`),
		priorityRegex: regexp.MustCompile(`priority (\d+)`),
		actionRegex:   regexp.MustCompile(`action (pass|record|alert|send|drop)`),
		banTimeRegex:  regexp.MustCompile(`^banTime.*[0-9]$`),
		nameRegex:     regexp.MustCompile(`^name.*`),
		msgRegex:      regexp.MustCompile(`msg:"([^"]*)"`),
		msgEndRegex:   regexp.MustCompile(`"$`),
		stringRegex:   regexp.MustCompile(`(\b|).+\b`),
		ConfigVarsMap: configVar,
	}
}

// Refresh 刷新 Tokens 并重新进行初始化
func (rp *RuleParser) Refresh(ruleStr string) {
	rp.ruleStr = ruleStr
	rp.Tokens = strings.Fields(ruleStr)
	rp.pos = 0
}

// Parse 解析规则
func (rp *RuleParser) Parse() (*Structure.RuleNodeValue, error) {

	if rp.pos >= len(rp.Tokens) {
		return nil, RuleParserError("parsing completed")
	}
	result := &Structure.RuleNodeValue{}

	//fmt.Println(rp.pos, "     ", rp.Tokens[rp.pos])
	err := rp.parseProtocols(result)
	if err != nil {
		return nil, err
	}

	//fmt.Println(rp.pos, "     ", rp.Tokens[rp.pos])
	err = rp.parsePriorityAndActionAndName(result)
	if err != nil {
		return nil, err
	}

	//fmt.Println(rp.pos, "     ", rp.Tokens[rp.pos])
	err = rp.parseMsg(result)
	if err != nil {
		return nil, err
	}

	//fmt.Println(rp.pos, "     ", rp.Tokens[rp.pos])
	err = rp.parseStrings(result)
	if err != nil {
		return nil, err
	}

	if rp.isValidRule(result) {
		result.ID = rp.curRule
		result.NetworkValues.ID = result.ID
		rp.curRule++
		return result, nil
	}

	return nil, RuleParserError("something is invalid value")
}

// parseProtocols 解析协议部分
func (rp *RuleParser) parseProtocols(result *Structure.RuleNodeValue) error {
	startIndex := rp.pos
	for true {
		if rp.protocolRegex.MatchString(rp.Tokens[rp.pos]) { //解析协议字符串
			//fmt.Println(rp.pos, "     ", rp.Tokens[rp.pos])
			result.Operators = append(result.Operators, BoolPlaceholder)

			protocol := rp.Tokens[rp.pos]
			p := Structure.Protocol{NameId: mapStrToLayerType[protocol]}
			rp.pos++

			for rp.pos < len(rp.Tokens) && rp.fieldRegex.MatchString(rp.Tokens[rp.pos]+" "+rp.Tokens[rp.pos+1]) {

				if rp.priorityRegex.MatchString(rp.Tokens[rp.pos] + " " + rp.Tokens[rp.pos+1]) {
					break
				}
				//fmt.Println(rp.pos, "   a  ", rp.Tokens[rp.pos])
				field := MapLayerTypeToParseType[p.NameId](rp.Tokens[rp.pos])
				tokenValue := rp.Tokens[rp.pos+1]

				if tokenValue[0] == '$' {
					tokenValue = tokenValue[1:]
					value := rp.ConfigVarsMap[tokenValue].VarValue
					if value != nil {
						p.Fields = append(p.Fields, Structure.Field{NameId: field, ValuePtr: value})
					}
				} else {
					if rp.priorityRegex.MatchString(rp.Tokens[rp.pos] + " " + rp.Tokens[rp.pos+1]) {
						break
					}

					value, err := MapLayerTypeToParseFunc[p.NameId](field, rp.Tokens[rp.pos+1])
					if err != nil {
						return err
					}
					p.Fields = append(p.Fields, Structure.Field{NameId: field, ValuePtr: value})
				}
				rp.pos += 2
			}
			result.Protocols = append(result.Protocols, p)
		} else if rp.operatorRegex.MatchString(rp.Tokens[rp.pos]) { //解析逻辑运算符字符串
			//fmt.Println(rp.pos, "   b  ", rp.Tokens[rp.pos])
			result.Operators = append(result.Operators, mapStrLogicToInt8[rp.Tokens[rp.pos]])
			//fmt.Println(result.Operators)
			rp.pos++
			continue
		} else {
			//fmt.Println(rp.pos, "   c  ", rp.Tokens[rp.pos])
			break
		}

	}
	result.NetworkValues.ProtocolsString = strings.Join(rp.Tokens[startIndex:rp.pos], " ")
	//fmt.Println("--------------------------------------")
	//fmt.Println(result.NetworkValues.ProtocolsString)
	return nil
}

// parsePriorityAndActionAndName 解析优先级和动作
func (rp *RuleParser) parsePriorityAndActionAndName(result *Structure.RuleNodeValue) error {
	//fmt.Println(rp.pos, "   a  ", rp.Tokens[rp.pos])

	if rp.pos+2 <= len(rp.Tokens) {
		if rp.priorityRegex.MatchString(strings.Join(rp.Tokens[rp.pos:rp.pos+2], " ")) {
			priority, _ := strconv.Atoi(rp.priorityRegex.FindStringSubmatch(strings.Join(rp.Tokens[rp.pos:rp.pos+2], " "))[1])
			result.Priority = priority
			result.NetworkValues.Priority = result.Priority
			rp.pos += 2
		}

		if rp.pos+2 <= len(rp.Tokens) && rp.actionRegex.MatchString(strings.Join(rp.Tokens[rp.pos:rp.pos+2], " ")) {
			action := rp.actionRegex.FindStringSubmatch(strings.Join(rp.Tokens[rp.pos:rp.pos+2], " "))[1]
			result.NetworkValues.ActionString = action
			result.Action = mapStringToAction[action]
			rp.pos += 2

			if result.Action == DropAction {
				err := rp.parseOption(result)
				if err != nil {
					return err
				}
			}
		}

		if rp.pos+2 <= len(rp.Tokens) {
			joinNameString := strings.Join(rp.Tokens[rp.pos:rp.pos+2], " ")

			if rp.nameRegex.MatchString(joinNameString) {
				name := strings.Split(joinNameString, " ")[1]
				result.RuleName = name
				rp.pos += 2
			}
		}
	}
	return nil
}

func (rp *RuleParser) parseOption(result *Structure.RuleNodeValue) error {
	result.Option = &Structure.RuleOption{}
	joinString := strings.Join(rp.Tokens[rp.pos:rp.pos+2], " ")

	if rp.banTimeRegex.MatchString(joinString) {
		banTime, _ := strconv.Atoi(strings.Split(joinString, " ")[1])
		result.Option.BanTime = banTime
		rp.pos += 2
	} else {
		result.Option.BanTime = 15
	}
	return nil
}

// parseMsg 解析消息
func (rp *RuleParser) parseMsg(result *Structure.RuleNodeValue) error {
	//fmt.Println(rp.pos, "   b  ", rp.Tokens[rp.pos])

	// 循环到以"结尾的字符串，代表消息结尾
	endPos := rp.pos
	for !rp.msgEndRegex.MatchString(rp.Tokens[endPos]) {
		endPos++
		if endPos >= len(rp.Tokens) {
			return RuleParserError("message parsing error")
		}
	}

	msgString := strings.Join(rp.Tokens[rp.pos:endPos+1], " ")
	result.NetworkValues.MsgString = msgString

	if rp.msgRegex.MatchString(msgString) {
		msg := rp.msgRegex.FindStringSubmatch(strings.Join(rp.Tokens[rp.pos:endPos+1], " "))[1]
		result.Msg = msg
		rp.pos = endPos + 1
	}
	return nil
}

// parseStrings 解析字符串
func (rp *RuleParser) parseStrings(result *Structure.RuleNodeValue) error {
	//fmt.Println(rp.pos, "   c  ", rp.Tokens[rp.pos])

	if rp.pos >= len(rp.Tokens) || (rp.pos < len(rp.Tokens) && rp.protocolRegex.MatchString(rp.Tokens[rp.pos])) {
		return nil
	}
	for rp.pos < len(rp.Tokens) && rp.stringRegex.MatchString(rp.Tokens[rp.pos]) {
		//result.StringIDs = append(result.StringIDs, &rp.Tokens[rp.pos])

		if rp.pos >= len(rp.Tokens) || (rp.pos < len(rp.Tokens) && rp.protocolRegex.MatchString(rp.Tokens[rp.pos])) {
			//fmt.Println(6666)
			break
		}
		result.StringIDs = append(result.StringIDs, int(MapMsgParaStringToProtocolFieldId[rp.Tokens[rp.pos]]))
		result.NetworkValues.StringParameterStrings = append(result.NetworkValues.StringParameterStrings, rp.Tokens[rp.pos])
		//fmt.Println(rp.pos, "   fff  ", rp.Tokens[rp.pos],"    ")
		rp.pos++

	}
	return nil
}

// isValidRule 检查规则的有效性
func (rp *RuleParser) isValidRule(result *Structure.RuleNodeValue) bool {
	//fmt.Println(rp.pos, "   d  ", rp.Tokens[rp.pos])
	return ((result.Msg != "" && result.Action != PassAction && result.Action != 0) || result.Action == PassAction) &&
		result.Priority != 0
}

// GetAllRule 获取所有规则并插入到给定的链表中
func (rp *RuleParser) GetAllRule(r *Structure.RuleLinkedList) error {
	for {
		rule, err := rp.Parse()
		if err != nil {
			if err.Error() == "parsing completed" {
				return nil
			}
			return err
		}
		r.Insert(rule)

	}
}
