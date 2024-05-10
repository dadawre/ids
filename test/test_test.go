package main

import (
	"fmt"
	"testing"
)

// Add 要测试的函数，这里简单示范一个加法函数
func Add(a, b int) int {
	return a + b
}

// TestAdd 测试函数，函数名必须以 Test 开头，并接受一个 *testing.T 类型的参数
func TestAdd(t *testing.T) {
	// 调用被测试的函数
	result := Add(2, 3)
	// 检查结果是否符合预期
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; want 5", result)
	}
}

func TestString(t *testing.T) {

	str := "Hello\n"
	for i, char := range str {
		// 修改 char 的值
		char = 'X'
		fmt.Printf("字符 %d: %c    %c\n", i, str[i], char)
		fmt.Println(str[i] == '\n')
	}
	// 打印原始字符串
	fmt.Println("原始字符串:", str) // 输出 "Hello"

}
