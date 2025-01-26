package main

import (
	"graduationDesign/back_end/OSController"
	"testing"
	"time"
)

func TestGetOSType(t *testing.T) {
	OSController.GetOSStates(2 * time.Second)
}
