package main

import (
	"fmt"
	"graduationDesign/parser"
	"testing"
)

func TestConfigParser(t *testing.T) {
	myTokens := make([]parser.Token, 0, 0)
	myTokens = append(myTokens, parser.Token{Type: parser.TokenKeywordVar, Value: "var"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenIdentifier, Value: "homePath"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenStringValue, Value: "/home/dadawre"})

	myTokens = append(myTokens, parser.Token{Type: parser.TokenKeywordInclude, Value: "include"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenStringValue, Value: "/agwagwagwag"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenVarUse, Value: "homePath"})
	myTokens = append(myTokens, parser.Token{Type: parser.TokenStringValue, Value: "home1111"})
	myConfigparser := parser.NewConfigparser(&myTokens)
	myConfigparser.ParseSelect()
	fmt.Println(myConfigparser)
	fmt.Println()
	//myConfigparser.ParseIpv4()
}
