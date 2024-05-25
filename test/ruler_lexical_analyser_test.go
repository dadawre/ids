package main

import (
	"fmt"
	"graduationDesign/parser/parserRuler"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func TestLexical(t *testing.T) {
	content, err := ioutil.ReadFile("../rulers/myruler")
	if err != nil {
		log.Fatal("Error reading file:", err)
		//return nil
	}

	rule := string(content)
	//strField := strings.Fields(rule)

	//rule := `IPv4 src 192.168.1.1 dest 192.168.80.0 TTL 128 & TCP dst 80 | UDP dst 8080 priority 10 action pass msg:"Connection allowed%c%c" note1 note2`
	parser := parserRuler.NewRuleParser(rule)
	result, err := parser.GetAllRule()
	if err != nil {
		fmt.Println("Error parsing rule:", err)
		return
	}
	//fmt.Println((*result))
	for index, rule := range *result {
		fmt.Println(index)
		fmt.Println(rule.Msg)
		fmt.Println(rule.Action)
		fmt.Println(rule.Priority)
		fmt.Println(rule.Operator)
		fmt.Println(rule.Protocols)
		fmt.Println(rule.Strings)
		fmt.Println("--------------------------------")
	}

	//fmt.Printf("Parsed Rule: %+v\n", result)
}

func TestField(t *testing.T) {

	// 打开文件
	//file, err := os.Open("../rulers/myruler")
	//if err != nil {
	//	fmt.Println("Error opening file:", err)
	//	//return nil
	//}

	// 创建bufio.Reader之前检查文件打开是否成功
	//info, err := file.Stat()
	//if err != nil {
	//	fmt.Println("Error getting file info:", err)
	//	//return nil
	//}
	//fmt.Println(info)

	content, err := ioutil.ReadFile("../rulers/myruler")
	if err != nil {
		log.Fatal("Error reading file:", err)
		//return nil
	}

	str := string(content)
	strField := strings.Fields(str)
	for index, iteratorStr := range strField {
		fmt.Println(index, "  ", iteratorStr)
	}
}
