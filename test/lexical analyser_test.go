package main

import (
	"fmt"
	"graduationDesign/parser/parserConfig"
	"os"
	"testing"
)

func TestConfigLexer(t *testing.T) {
	// 打开文件
	file, err := os.Open("../conf/ids.conf")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	//reader := bufio.NewReader(file)
	// 创建一个bufio.Reader

	myConfigLexer := parserConfig.NewConfigLexer("../conf/ids.conf")
	//myTokens := make([]parser.Token, 0, 0)

	//for {
	//	myToken := myConfigLexer.NextToken()
	//	if myToken.Type == 0 {
	//		break
	//	}
	//	myTokens = append(myTokens, myToken)
	//}
	fmt.Println(myConfigLexer.GetAllToken())
	fmt.Println(myConfigLexer.GetAllToken())
	fmt.Println(myConfigLexer.GetAllToken())
	fmt.Println(myConfigLexer.GetAllToken())
	fmt.Println(myConfigLexer.GetAllToken())

	//fmt.Println(myConfigLexer.CurrTokens)
}
