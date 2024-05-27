package main

import (
	"fmt"
	"graduationDesign/parser/parserRuler"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

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
	for node := result.GetHead(); node != nil; node = node.Next {
		fmt.Println(node.RuleValue.ID)
		fmt.Println(node.RuleValue.Msg)
		fmt.Println(node.RuleValue.Action)
		fmt.Println(node.RuleValue.Priority)
		fmt.Println(node.RuleValue.Operator)
		fmt.Println(node.RuleValue.Protocols)
		fmt.Println(node.RuleValue.Strings)
		fmt.Println("--------------------------------")
	}
	//fmt.Println((*result))
	//for index, rule := range *result {
	//	fmt.Println(index)
	//	fmt.Println(rule.Msg)
	//	fmt.Println(rule.Action)
	//	fmt.Println(rule.Priority)
	//	fmt.Println(rule.Operator)
	//	fmt.Println(rule.Protocols)
	//	fmt.Println(rule.Strings)
	//	fmt.Println("--------------------------------")
	//}

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
