package main

import (
	"graduationDesign/back_end/Control"
	"graduationDesign/back_end/Initialize"
	"graduationDesign/back_end/Initialize/Parameter"
	"graduationDesign/back_end/Mode"
	"graduationDesign/back_end/Network"
	"graduationDesign/back_end/OSController"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	Initialize.MainInitialize()
	Control.RulerInputListen()

	go func() { //前端请求监控
		Network.InitializeNetwork()
	}()

	go func() { //程序状态监控
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	go func() { //硬件利用率刷新
		OSController.GetOSStates(time.Duration(Parameter.HardwareDataRefreshFrequency) * time.Millisecond) //多少ms刷新一次硬件数据
	}()

	if Parameter.IdsMode == 1 {
		Mode.IDSMode(Mode.NewPacketSource())
	} else if Parameter.IdsMode == 2 {
		Mode.PacketSniffMode(Mode.NewPacketSource())
	}

}
