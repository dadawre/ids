package parserRuler

import (
	"graduationDesign/structure/ruleStructure"
	"regexp"
	"strconv"
	"strings"
)

type RuleParser struct {
	curRule       int
	RuleLIST      ruleStructure.RuleLinkedList
	pos           int
	ruleStr       string
	tokens        []string
	protocolRegex *regexp.Regexp
	fieldRegex    *regexp.Regexp
	operatorRegex *regexp.Regexp
	priorityRegex *regexp.Regexp
	actionRegex   *regexp.Regexp
	msgRegex      *regexp.Regexp
	stringRegex   *regexp.Regexp
}

type RuleParserError string

func (e RuleParserError) Error() string {
	return string(e)
}

func NewRuleParser(ruleStr string) *RuleParser {
	return &RuleParser{
		pos:           0,
		ruleStr:       ruleStr,
		tokens:        strings.Fields(ruleStr),
		protocolRegex: regexp.MustCompile(`\b(Eth|IPv4|IPv6|TCP|UDP)\b`),
		fieldRegex:    regexp.MustCompile(`\b(\w+)\b \b(\w+)\b`),
		operatorRegex: regexp.MustCompile(`[&|]`),
		priorityRegex: regexp.MustCompile(`priority (\d+)`),
		actionRegex:   regexp.MustCompile(`action (pass|record|alert|send|drop)`),
		msgRegex:      regexp.MustCompile(`msg:"([^"]*)"`),
		stringRegex:   regexp.MustCompile(`\b\w+\b`),
		//RuleLIST:      make([]*Rule, 0, 0),
	}
}

func (rp *RuleParser) Parse() (*ruleStructure.Rule, error) {
	var err error
	if rp.pos >= len(rp.tokens) {
		err = RuleParserError("parsing completed")
		return nil, err
	}

	//tokens := strings.Fields(rp.ruleStr)
	//for index, token := range rp.tokens {
	//	fmt.Println(index, "  ", token)
	//}
	result := &ruleStructure.Rule{}

	//for rp.pos < len(rp.tokens) {

	for rp.protocolRegex.MatchString(rp.tokens[rp.pos]) {
		protocol := rp.tokens[rp.pos]
		p := ruleStructure.Protocol{Name: mapStrToLayerType[protocol]}
		rp.pos++
		for rp.pos < len(rp.tokens) && rp.fieldRegex.MatchString(rp.tokens[rp.pos]+" "+rp.tokens[rp.pos+1]) {
			if rp.priorityRegex.MatchString(rp.tokens[rp.pos] + " " + rp.tokens[rp.pos+1]) {
				break
			}
			field := MapLayerTypeTOParseType[p.Name](rp.tokens[rp.pos])
			value, err := MapLayerTypeToParseFunc[p.Name](field, rp.tokens[rp.pos+1])
			if err != nil {
				return nil, err
			}
			p.Fields = append(p.Fields, ruleStructure.Field{Name: field, Value: value})
			rp.pos += 2
		}
		result.Protocols = append(result.Protocols, p)

		if rp.operatorRegex.MatchString(rp.tokens[rp.pos]) {
			result.Operator = append(result.Operator, mapStrLogicToInt8[rp.tokens[rp.pos]])
			rp.pos++
			continue
		}
		break
	}

	if rp.pos+2 < len(rp.tokens) {
		if rp.priorityRegex.MatchString(strings.Join(rp.tokens[rp.pos:rp.pos+2], " ")) {
			priority, _ := strconv.Atoi(rp.priorityRegex.FindStringSubmatch(strings.Join(rp.tokens[rp.pos:rp.pos+2], " "))[1])
			result.Priority = priority
			rp.pos += 2
			//continue
		}

		if rp.actionRegex.MatchString(strings.Join(rp.tokens[rp.pos:rp.pos+2], " ")) {
			action := rp.actionRegex.FindStringSubmatch(strings.Join(rp.tokens[rp.pos:rp.pos+2], " "))[1]
			result.Action = mapStringToAction[action]
			rp.pos += 2
			//continue
		}

		if rp.msgRegex.MatchString(strings.Join(rp.tokens[rp.pos:rp.pos+2], " ")) {
			msg := rp.msgRegex.FindStringSubmatch(strings.Join(rp.tokens[rp.pos:rp.pos+2], " "))[1]
			result.Msg = msg
			rp.pos += 2
			//continue
		}
	}

	for rp.stringRegex.MatchString(rp.tokens[rp.pos]) {
		result.Strings = append(result.Strings, rp.tokens[rp.pos])
		rp.pos++
		if rp.pos >= len(rp.tokens) || (rp.pos < len(rp.tokens) && rp.protocolRegex.MatchString(rp.tokens[rp.pos])) {
			if ((result.Msg != "" && result.Action != PassAction && result.Action != 0) || result.Action == PassAction) && len(result.Protocols) == len(result.Operator)+1 && result.Priority != 0 {
				result.ID = rp.curRule
				rp.curRule++
				return result, nil
			} else {
				err = RuleParserError("something is invalid value")
				return nil, err
			}
		}
		//continue
	}

	err = RuleParserError("RuleParseError")
	return nil, err
	//}

	//return result, nil
}

func (rp *RuleParser) GetAllRule() (*ruleStructure.RuleLinkedList, error) {
	var err error
	var rule *ruleStructure.Rule
	for rule, err = rp.Parse(); err == nil; rule, err = rp.Parse() {
		//rp.RuleLIST = append(rp.RuleLIST, rule)
		//fmt.Println(rule)
		rp.RuleLIST.Insert(rule)
	}
	if err == nil || (err != nil && err.Error() == "parsing completed") {
		return &rp.RuleLIST, nil
	}
	return nil, err
}
