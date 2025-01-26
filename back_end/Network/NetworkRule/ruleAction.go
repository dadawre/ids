package NetworkRule

import (
	"encoding/json"
	"graduationDesign/back_end/Initialize"
	"graduationDesign/back_end/structure/ruleStructure"
	"net/http"
	"strconv"
	"strings"
)

func RuleStringInsert(rule ruleStructure.RuleNodeStringValue) error {
	ruleString := ""
	ruleString += rule.ProtocolsString
	ruleString += " priority " + strconv.Itoa(rule.Priority)
	ruleString += " action " + rule.ActionString
	ruleString += " " + rule.MsgString
	ruleString += " " + strings.Join(rule.StringParameterStrings, " ")
	//fmt.Println(ruleString)
	Initialize.MainRulerParser.Refresh(ruleString)
	err := Initialize.MainRulerParser.GetAllRule(Initialize.ListRules)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	//Initialize.ListRules.Traverse()
	return nil
}

func RuleStringChange(rule ruleStructure.RuleNodeStringValue) error {
	err := Initialize.ListRules.Delete(rule.ID)
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
	var rule ruleStructure.RuleNodeStringValue
	err := json.NewDecoder(r.Body).Decode(&rule)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	//fmt.Println(rule)
	mu.Lock()
	defer mu.Unlock()
	Initialize.ArrayRuleStrings = append(Initialize.ArrayRuleStrings, rule)
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
	for i, rule := range Initialize.ArrayRuleStrings {
		if rule.ID == req.ID {
			Initialize.ArrayRuleStrings = append(Initialize.ArrayRuleStrings[:i], Initialize.ArrayRuleStrings[i+1:]...)
			err := Initialize.ListRules.Delete(req.ID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			break
		}
	}
	w.WriteHeader(http.StatusOK)
}

func RuleEdit(w http.ResponseWriter, r *http.Request) {
	var updatedRule ruleStructure.RuleNodeStringValue
	err := json.NewDecoder(r.Body).Decode(&updatedRule)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	//fmt.Println(updatedRule)
	mu.Lock()
	defer mu.Unlock()
	for i, rule := range Initialize.ArrayRuleStrings {
		if rule.ID == updatedRule.ID {
			Initialize.ArrayRuleStrings[i] = updatedRule
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
