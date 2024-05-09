package main

import "testing"

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
