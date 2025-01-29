package main

import (
	"fmt"
	ParserConfig2 "graduationDesign/back_end/Parser/ParserConfig"
	"graduationDesign/back_end/Parser/ParserRuler"
	"graduationDesign/back_end/Structure"
	"os"
	"strings"
	"testing"
)

// readFileContent 读取文件内容并返回字符串
func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", path, err)
	}
	return string(content), nil
}

// TestLexical 测试词法解析
func TestLexical(t *testing.T) {
	// 读取配置文件
	configPath := "../../conf/ids.conf"
	file, err := os.Open(configPath)
	if err != nil {
		t.Fatalf("Error opening file %s: %v", configPath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf(err.Error())
		}
	}(file)

	// 词法分析
	myConfigLexer := ParserConfig2.NewConfigLexer(configPath)
	tokens := myConfigLexer.GetAllToken()

	// 语法分析
	myConfigparser := ParserConfig2.NewConfigparser(&tokens)
	myConfigparser.ParseAllTokens()
	//if err := myConfigparser.ParseAllTokens(); err != nil {
	//	t.Fatalf("Error parsing tokens: %v", err)
	//}

	// 读取规则文件
	rulePath := "../rulers/myruler"
	rule, err := readFileContent(rulePath)
	if err != nil {
		t.Fatal(err)
	}

	// 解析规则
	parser := ParserRuler.NewRuleParser(rule, myConfigparser.ConfigVars)

	//for index, token := range parser.Tokens {
	//	fmt.Println(index, "   ", token)
	//}

	result := Structure.RuleLinkedList{}
	err = parser.GetAllRule(&result)
	if err != nil {
		return
	}
	// 打印词法分析结果
	for index, token := range parser.Tokens {
		t.Logf("Token %d: %s", index, token)
	}

	fmt.Println("----------------------------------------------")

	// 读取规则文件
	rulePath = "../rulers/backrule"
	rule, err = readFileContent(rulePath)
	if err != nil {
		t.Fatal(err)
	}
	parser.Refresh(rule)

	if err := parser.GetAllRule(&result); err != nil && err.Error() != "parsing completed" {
		fmt.Println(parser.Tokens)
		t.Fatalf("Error parsing rule: %v", err)
	}

	// 再次打印词法分析结果
	for index, token := range parser.Tokens {
		t.Logf("Token %d: %s", index, token)
	}

	// 打印解析结果
	for node := result.GetHead(); node != nil; node = node.Next {
		t.Logf("\nID: %d\nMsg: %s\nAction: %d\nPriority: %d\nOperator: %v\nProtocols: %v\nStrings: %v\n--------------------------------",
			node.RuleValue.ID, node.RuleValue.Msg, node.RuleValue.Action, node.RuleValue.Priority, node.RuleValue.Operators, node.RuleValue.Protocols, node.RuleValue.StringIDs)
		for _, protocol := range node.RuleValue.Protocols {
			fmt.Print(protocol.NameId, ":   ")
			for _, field := range protocol.Fields {
				fmt.Print(field, "  ")
			}
			fmt.Print(" | ")
		}
		fmt.Println()
	}
}

// TestField 测试字段解析
func TestField(t *testing.T) {
	// 读取规则文件
	rulePath := "../rulers/myruler"
	str, err := readFileContent(rulePath)
	if err != nil {
		t.Fatal(err)
	}

	// 解析字段
	strField := strings.Fields(str)
	for index, iteratorStr := range strField {
		t.Logf("Field %d: %s", index, iteratorStr)
	}
}
func TestAuxiliaryTransform(t *testing.T) {
	// 处理 int 类型指针
	//var a int = 10
	//parserRuler.ProcessPointer(&a)
	//fmt.Printf("Updated int value: %d\n", a)
	//
	// 处理 string 类型指针
	//var b string = "initial"
	//parserRuler.ProcessPointer(&b)
	//fmt.Printf("Updated string value: %s\n", b)

	// 处理结构体类型指针
	//person := &parserRuler.Person{NameId: "Alice", Age: 25}
	//parserRuler.ProcessPointer(person)
	//fmt.Printf("Updated person: %+v\n", person)
}

func TestAuxiliaryT(t *testing.T) {
	// 示例使用
	//var a int = 42
	//var b string = "Hello"
	//var c []int = []int{1, 2, 3}
	//var d map[string]int = map[string]int{"one": 1, "two": 2, "three": 3}
	//
	//// 获取接口的原始类型
	//typA := parserRuler.GetOriginalType(a)
	//typB := parserRuler.GetOriginalType(b)
	//typC := parserRuler.GetOriginalType(c)
	//typD := parserRuler.GetOriginalType(d)
	//
	//// 转换为原始类型
	//originalA := parserRuler.ConvertToOriginalType(typA)
	//originalB := parserRuler.ConvertToOriginalType(typB)
	//originalC := parserRuler.ConvertToOriginalType(typC)
	//originalD := parserRuler.ConvertToOriginalType(typD)
	//
	//// 打印转换后的值
	//fmt.Printf("Original type of a: %T\n", originalA)
	//fmt.Println(originalA)
	//fmt.Printf("Original type of b: %T\n", originalB)
	//fmt.Printf("Original type of c: %T\n", originalC)
	//fmt.Printf("Original type of d: %T\n", originalD)
}
