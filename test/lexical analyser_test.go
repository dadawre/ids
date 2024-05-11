package main

import (
	"fmt"
	"graduationDesign/parser"
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

	myConfigLexer := parser.NewConfigLexer("../conf/ids.conf")
	myTokens := make([]parser.Token, 1, 1)
	myTokens = append(myTokens, myConfigLexer.NextToken())
	myTokens = append(myTokens, myConfigLexer.NextToken())
	myTokens = append(myTokens, myConfigLexer.NextToken())
	myTokens = append(myTokens, myConfigLexer.NextToken())
	myTokens = append(myTokens, myConfigLexer.NextToken())

	fmt.Println(myTokens)
	//fmt.Println("gawgagagwagwagwa")
	//fmt.Println(myConfigLexer)

	//fmt.Println(m, "  值", myConfigLexer.NextToken().Value)
	//for myConfigLexer.NextToken().Type != parser.TokenEOF {
	//myToken = append(myToken, myConfigLexer.NextToken())
	//}

}
