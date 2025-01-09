package main

import (
	"graduationDesign/Control"
	"graduationDesign/Initialize"
	"graduationDesign/Mode"
	"graduationDesign/OSController"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	Initialize.MainInitialize()
	Control.RulerInputListen()
	go func() { //程序状态监控
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	go func() { //硬件利用率刷新
		OSController.GetOSStates(time.Duration(Initialize.HardwareDataRefreshFrequency) * time.Millisecond) //多少ms刷新一次硬件数据
	}()

	if Initialize.IdsMode == 1 {
		Mode.IDSMode(Mode.NewPacketSource())
	} else if Initialize.IdsMode == 2 {
		Mode.PacketSniffMode(Mode.NewPacketSource())
	}

}
