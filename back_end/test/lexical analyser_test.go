package main

import (
	"fmt"
	"graduationDesign/back_end/Parser/ParserConfig"
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	//reader := bufio.NewReader(file)
	// 创建一个bufio.Reader

	myConfigLexer := ParserConfig.NewConfigLexer("../conf/ids.conf")
	//myTokens := make([]parser.Token, 0, 0)

	//for {
	//	myToken := myConfigLexer.NextToken()
	//	if myToken.Type == 0 {
	//		break
	//	}
	//	myTokens = append(myTokens, myToken)
	//}
	//fmt.Println(myConfigLexer.GetAllToken())
	//fmt.Println(myConfigLexer.GetAllToken())
	//fmt.Println(myConfigLexer.GetAllToken())
	//fmt.Println(myConfigLexer.GetAllToken())
	//fmt.Println(myConfigLexer.GetAllToken())
	fmt.Println(myConfigLexer.GetAllToken())

	//fmt.Println(myConfigLexer.CurrTokens)
}
