package main

import (
	"graduationDesign/OSController"
	"testing"
	"time"
)

func TestGetOSType(t *testing.T) {
	OSController.GetOSStates(2 * time.Second)
}
