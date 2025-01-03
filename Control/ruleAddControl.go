package Control

import (
	"fmt"
	"graduationDesign/Initialize"
)

func RuleInsert(ruleString string) bool {
	Initialize.MainRulerParser.Refresh(ruleString)

	if err := Initialize.MainRulerParser.GetAllRule(Initialize.ListRules); err != nil && err.Error() != "parsing completed" {
		//fmt.Println(Initialize.MainRulerParser.Tokens)
		fmt.Printf("Error parsing rule: %v", err)
		return false
	}

	return true
}
