package main

import (
	"fmt"
	"graduationDesign/Control"
	"testing"
)

func TestReadConf(t *testing.T) {
	Control.ReadConfig("F:\\go_learn\\graduationDesign\\conf\\ids.conf")
	fmt.Println(66666)
}

func TestWriteConf(t *testing.T) {
	err := Control.WriteToFile("g", "dafafa")
	if err != nil {
		fmt.Println("error")
	}
}
