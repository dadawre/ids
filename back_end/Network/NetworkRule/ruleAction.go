package NetworkRule

import (
	"encoding/json"
	Structure2 "graduationDesign/back_end/Initialize/InitializeStructure"
	"graduationDesign/back_end/Structure"
	"net/http"
	"strconv"
	"strings"
)

func RuleStringInsert(rule Structure.RuleNodeStringValue) error {
	ruleString := ""
	ruleString += rule.ProtocolsString
	ruleString += " priority " + strconv.Itoa(rule.Priority)
	ruleString += " action " + rule.ActionString
	ruleString += " " + rule.MsgString
	ruleString += " " + strings.Join(rule.StringParameterStrings, " ")
	//fmt.Println(ruleString)
	Structure2.MainRulerParser.Refresh(ruleString)
	err := Structure2.MainRulerParser.GetAllRule(Structure2.ListRules)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	//Initialize.ListRules.Traverse()
	return nil
}

func RuleStringChange(rule Structure.RuleNodeStringValue) error {
	err := Structure2.ListRules.Delete(rule.ID)
	if err != nil {
		return err
	}
	err = RuleStringInsert(rule)
	if err != nil {
		return err
	}
	return err
}

func RuleAdd(w http.ResponseWriter, r *http.Request) {
	var rule Structure.RuleNodeStringValue
	err := json.NewDecoder(r.Body).Decode(&rule)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	//fmt.Println(rule)
	mu.Lock()
	defer mu.Unlock()
	Structure2.ArrayRuleStrings = append(Structure2.ArrayRuleStrings, rule)
	err = RuleStringInsert(rule)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func RuleDelete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID int `json:"ID"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	mu.Lock()
	defer mu.Unlock()
	for i, rule := range Structure2.ArrayRuleStrings {
		if rule.ID == req.ID {
			Structure2.ArrayRuleStrings = append(Structure2.ArrayRuleStrings[:i], Structure2.ArrayRuleStrings[i+1:]...)
			err := Structure2.ListRules.Delete(req.ID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			break
		}
	}
	w.WriteHeader(http.StatusOK)
}

func RuleEdit(w http.ResponseWriter, r *http.Request) {
	var updatedRule Structure.RuleNodeStringValue
	err := json.NewDecoder(r.Body).Decode(&updatedRule)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	//fmt.Println(updatedRule)
	mu.Lock()
	defer mu.Unlock()
	for i, rule := range Structure2.ArrayRuleStrings {
		if rule.ID == updatedRule.ID {
			Structure2.ArrayRuleStrings[i] = updatedRule
			//RuleStringInsert(updatedRule)
			err := RuleStringChange(updatedRule)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			break
		}
	}
	w.WriteHeader(http.StatusOK)

}
