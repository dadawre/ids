package main

import (
	"fmt"
	"graduationDesign/ruleMatch"
	"net"
	"testing"
)

func TestSlice(t *testing.T) {

	macAddr1, _ := net.ParseMAC("01:23:45:67:89:ab")
	macAddr2 := macAddr1

	// 修改 macAddr2
	macAddr2[0] = 0xff

	fmt.Println("macAddr1:", macAddr1)
	fmt.Println("macAddr2:", macAddr2)

}
func TestLogicOperator(t *testing.T) {
	// Example usage
	bools := []bool{true, false, false, true}
	//operators := []string{"&&", "||", "||"}
	operators := []int8{0, 1, 0}
	result, err := ruleMatch.ExecuteLogic(bools, operators)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
