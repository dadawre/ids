package main

import (
	"fmt"
	"testing"
)

func TestSprintf(t *testing.T) {
	num := 42
	ptr := &num

	// 使用 fmt.Sprintf 格式化指针的值
	result := fmt.Sprintf("Pointer address: %p", ptr)

	fmt.Println(result)
}
