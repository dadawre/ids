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

// Refresh 刷新Tokens并重新进行初始化
func (rp *RuleParser) Refresh(ruleStr string) {
	rp.ruleStr = ruleStr
	rp.Tokens = strings.Fields(ruleStr)
	rp.pos = 0
}

type RuleParserError string

func (e RuleParserError) Error() string {
	return string(e)
}

func NewRuleParser(ruleStr string, configVar map[string]*parserConfig.ConfigVar) *RuleParser {
	return &RuleParser{
		pos:           0,
		ruleStr:       ruleStr,
		Tokens:        strings.Fields(ruleStr),
		protocolRegex: regexp.MustCompile(`\b(Eth|IPv4|IPv6|TCP|UDP)\b`),
		fieldRegex:    regexp.MustCompile(`(\b)(\w+)\b (\b|\$)(\w+)\b`),
		operatorRegex: regexp.MustCompile(`[&|]`),
		priorityRegex: regexp.MustCompile(`priority (\d+)`),
		actionRegex:   regexp.MustCompile(`action (pass|record|alert|send|drop)`),
		msgRegex:      regexp.MustCompile(`msg:"([^"]*)"`),
		msgEndRegex:   regexp.MustCompile(`"$`),
		stringRegex:   regexp.MustCompile(`\b\w+\b`),
		ConfigVarsMap: configVar,
	}
}

// Parse 解析规则
func (rp *RuleParser) Parse() (*ruleStructure.Rule, error) {
	var err error
	if rp.pos >= len(rp.Tokens) {
		err = RuleParserError("parsing completed")
		return nil, err
	}

	result := &ruleStructure.Rule{}

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
				//fmt.Println("$value is : ", value)
				if value != nil {
					p.Fields = append(p.Fields, &(ruleStructure.Field{NameId: field, ValuePtr: value}))
				}

			} else {
				value, err := MapLayerTypeToParseFunc[p.NameId](field, rp.Tokens[rp.pos+1])
				//fmt.Println(p.NameId, "   ", rp.Tokens[rp.pos+1])
				//fmt.Println(field, "  value is :  ", value)
				if err != nil {
					return nil, err
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

		endPos := rp.pos
		for rp.msgEndRegex.MatchString(rp.Tokens[endPos]) == false {
			endPos++
		}

		if rp.msgRegex.MatchString(strings.Join(rp.Tokens[rp.pos:endPos+1], " ")) {
			msg := rp.msgRegex.FindStringSubmatch(strings.Join(rp.Tokens[rp.pos:endPos+1], " "))[1]
			result.Msg = msg
			rp.pos = endPos
		}
	}

	for rp.stringRegex.MatchString(rp.Tokens[rp.pos]) {
		result.Strings = append(result.Strings, &rp.Tokens[rp.pos])
		rp.pos++
		if rp.pos >= len(rp.Tokens) || (rp.pos < len(rp.Tokens) && rp.protocolRegex.MatchString(rp.Tokens[rp.pos])) {

			if ((result.Msg != "" && result.Action != PassAction && result.Action != 0) || result.Action == PassAction) && len(result.Protocols) == len(result.Operator)+1 && result.Priority != 0 {
				result.ID = rp.curRule
				rp.curRule++
				return result, nil
			} else {
				err = RuleParserError("something is invalid value")
				return nil, err
			}
		}
	}

	err = RuleParserError("RuleParseError")
	return nil, err
}

// GetAllRule 获取所有规则并插入到给定的链表中
func (rp *RuleParser) GetAllRule(r *ruleStructure.RuleLinkedList) error {
	var err error
	var rule *ruleStructure.Rule
	for rule, err = rp.Parse(); err == nil; rule, err = rp.Parse() {
		//rp.RuleLIST = append(rp.RuleLIST, rule)
		//fmt.Println(rule)
		//rp.RuleLIST.Insert(rule)
		r.Insert(rule)
	}
	if err == nil || (err != nil && err.Error() == "parsing completed") {
		return nil
	}
	fmt.Println("gwagawapogjapaw      ", err)
	return err
}
