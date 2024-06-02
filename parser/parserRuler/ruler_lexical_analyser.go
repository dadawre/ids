package parserRuler

import (
	"fmt"
	"graduationDesign/parser/parserConfig"
	"graduationDesign/structure/ruleStructure"
	"regexp"
	"strconv"
	"strings"
)

// RuleParser 规则解析结构体
type RuleParser struct {
	curRule       int
	pos           int
	ruleStr       string
	Tokens        []string
	protocolRegex *regexp.Regexp
	fieldRegex    *regexp.Regexp
	operatorRegex *regexp.Regexp
	priorityRegex *regexp.Regexp
	actionRegex   *regexp.Regexp
	msgRegex      *regexp.Regexp
	stringRegex   *regexp.Regexp
	msgEndRegex   *regexp.Regexp
	ConfigVarsMap map[string]*parserConfig.ConfigVar
}

// RuleParserError 自定义错误类型
type RuleParserError string

func (e RuleParserError) Error() string {
	return string(e)
}

// NewRuleParser 创建并初始化一个新的 RuleParser
func NewRuleParser(ruleStr string, configVar map[string]*parserConfig.ConfigVar) *RuleParser {
	return &RuleParser{
		pos:           0,
		ruleStr:       ruleStr,
		Tokens:        strings.Fields(ruleStr),
		protocolRegex: regexp.MustCompile(`^(Eth|IPv4|IPv6|TCP|UDP)$`),
		fieldRegex:    regexp.MustCompile(`(\b)(\w+)\b (\b|\$)(\w+)\b`),
		operatorRegex: regexp.MustCompile(`[&|]`),
		priorityRegex: regexp.MustCompile(`priority (\d+)`),
		actionRegex:   regexp.MustCompile(`action (pass|record|alert|send|drop)`),
		msgRegex:      regexp.MustCompile(`msg:"([^"]*)"`),
		msgEndRegex:   regexp.MustCompile(`"$`),
		stringRegex:   regexp.MustCompile(`(\b|\#).+\b`),
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
func (rp *RuleParser) Parse() (*ruleStructure.Rule, error) {
	if rp.pos >= len(rp.Tokens) {
		return nil, RuleParserError("parsing completed")
	}
	result := &ruleStructure.Rule{}

	err := rp.parseProtocols(result)
	if err != nil {
		return nil, err
	}

	err = rp.parsePriorityAndAction(result)
	if err != nil {
		return nil, err
	}

	err = rp.parseMsg(result)
	if err != nil {
		return nil, err
	}

	err = rp.parseStrings(result)
	if err != nil {
		return nil, err
	}

	if rp.isValidRule(result) {
		result.ID = rp.curRule
		rp.curRule++
		return result, nil
	}
	fmt.Println(rp.pos, "     ", rp.Tokens[rp.pos])
	return nil, RuleParserError("something is invalid value")
}

// parseProtocols 解析协议部分
func (rp *RuleParser) parseProtocols(result *ruleStructure.Rule) error {
	for rp.protocolRegex.MatchString(rp.Tokens[rp.pos]) {
		protocol := rp.Tokens[rp.pos]
		p := ruleStructure.Protocol{NameId: mapStrToLayerType[protocol]}
		rp.pos++

		for rp.pos < len(rp.Tokens) && rp.fieldRegex.MatchString(rp.Tokens[rp.pos]+" "+rp.Tokens[rp.pos+1]) {

			if rp.priorityRegex.MatchString(rp.Tokens[rp.pos] + " " + rp.Tokens[rp.pos+1]) {
				break
			}
			field := MapLayerTypeTOParseType[p.NameId](rp.Tokens[rp.pos])
			tokenValue := rp.Tokens[rp.pos+1]

			if tokenValue[0] == '$' {
				tokenValue = tokenValue[1:]
				value := rp.ConfigVarsMap[tokenValue].VarValue
				if value != nil {
					p.Fields = append(p.Fields, &(ruleStructure.Field{NameId: field, ValuePtr: value}))
				}
			} else {
				value, err := MapLayerTypeToParseFunc[p.NameId](field, rp.Tokens[rp.pos+1])
				if err != nil {
					return err
				}
				p.Fields = append(p.Fields, &(ruleStructure.Field{NameId: field, ValuePtr: value}))
			}
			rp.pos += 2
		}
		result.Protocols = append(result.Protocols, &p)

		if rp.operatorRegex.MatchString(rp.Tokens[rp.pos]) {
			result.Operator = append(result.Operator, mapStrLogicToInt8[rp.Tokens[rp.pos]])
			rp.pos++
			continue
		}
		break
	}
	return nil
}

// parsePriorityAndAction 解析优先级和动作
func (rp *RuleParser) parsePriorityAndAction(result *ruleStructure.Rule) error {
	if rp.pos+2 < len(rp.Tokens) {
		if rp.priorityRegex.MatchString(strings.Join(rp.Tokens[rp.pos:rp.pos+2], " ")) {
			priority, _ := strconv.Atoi(rp.priorityRegex.FindStringSubmatch(strings.Join(rp.Tokens[rp.pos:rp.pos+2], " "))[1])
			result.Priority = priority
			rp.pos += 2
		}

		if rp.actionRegex.MatchString(strings.Join(rp.Tokens[rp.pos:rp.pos+2], " ")) {
			action := rp.actionRegex.FindStringSubmatch(strings.Join(rp.Tokens[rp.pos:rp.pos+2], " "))[1]
			result.Action = mapStringToAction[action]
			rp.pos += 2
		}
	}
	return nil
}

// parseMsg 解析消息
func (rp *RuleParser) parseMsg(result *ruleStructure.Rule) error {
	// 循环到以"结尾的字符串，代表消息结尾
	endPos := rp.pos
	for !rp.msgEndRegex.MatchString(rp.Tokens[endPos]) {
		endPos++
		if endPos >= len(rp.Tokens) {
			return RuleParserError("message parsing error")
		}
	}

	if rp.msgRegex.MatchString(strings.Join(rp.Tokens[rp.pos:endPos+1], " ")) {
		msg := rp.msgRegex.FindStringSubmatch(strings.Join(rp.Tokens[rp.pos:endPos+1], " "))[1]
		result.Msg = msg
		rp.pos = endPos + 1
	}
	return nil
}

// parseStrings 解析字符串
func (rp *RuleParser) parseStrings(result *ruleStructure.Rule) error {
	if rp.pos >= len(rp.Tokens) || (rp.pos < len(rp.Tokens) && rp.protocolRegex.MatchString(rp.Tokens[rp.pos])) {
		return nil
	}
	for rp.stringRegex.MatchString(rp.Tokens[rp.pos]) {
		//result.Strings = append(result.Strings, &rp.Tokens[rp.pos])

		result.Strings = append(result.Strings, int(MapMsgParaStringToProtocolFieldId[rp.Tokens[rp.pos]]))

		rp.pos++

		if rp.pos >= len(rp.Tokens) || (rp.pos < len(rp.Tokens) && rp.protocolRegex.MatchString(rp.Tokens[rp.pos])) {
			break
		}

	}
	return nil
}

// isValidRule 检查规则的有效性
func (rp *RuleParser) isValidRule(result *ruleStructure.Rule) bool {
	return (result.Msg != "" && result.Action != PassAction && result.Action != 0 || result.Action == PassAction) &&
		len(result.Protocols) == len(result.Operator)+1 &&
		result.Priority != 0
}

// GetAllRule 获取所有规则并插入到给定的链表中
func (rp *RuleParser) GetAllRule(r *ruleStructure.RuleLinkedList) error {
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
