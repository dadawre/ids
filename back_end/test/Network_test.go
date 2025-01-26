package main

import (
	"graduationDesign/back_end/Network"
	"graduationDesign/back_end/Network/NetworkConfig"
	"testing"
)

func TestConfigData(t *testing.T) {

	NetworkConfig.InitializeConfigData()
	Network.NewStartWebRoute()
}

func TestRoute(t *testing.T) {
	NetworkConfig.InitializeConfigData()
	Network.NewStartWebRoute()
}
