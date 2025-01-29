package Control

import (
	"fmt"
	"graduationDesign/back_end/Initialize/InitializeStructure"
)

func RuleInsert(ruleString string) bool {
	InitializeStructure.MainRulerParser.Refresh(ruleString)

	if err := InitializeStructure.MainRulerParser.GetAllRule(InitializeStructure.ListRules); err != nil && err.Error() != "parsing completed" {
		//fmt.Println(Initialize.MainRulerParser.Tokens)
		fmt.Printf("Error parsing rule: %v", err)
		return false
	}

	return true
}
