package Mode

import (
	"github.com/google/gopacket"
	"graduationDesign/back_end/Handler"
	"graduationDesign/back_end/Initialize/InitializeStructure"
	"graduationDesign/back_end/Initialize/Parameter"
	"strconv"
)

// IDSMode IDS模式
func IDSMode(packetSource *gopacket.PacketSource) {
	defer inactiveHandle.CleanUp()
	defer handle.Close()

	//定义通道数组，将数据包源的数据包分发给各个数据包解析结构体
	Channels := make([]chan *gopacket.Packet, Parameter.MaxCoroutine)
	Handler.PacketHandlerArr = make([]*Handler.PacketHandler, Parameter.MaxCoroutine)

	//通道数组为各个数据包解析结构体提供数据包
	for i := 0; i < Parameter.MaxCoroutine; i++ {
		Channels[i] = make(chan *gopacket.Packet)
		Handler.PacketHandlerArr[i] = Handler.NewHandlerPacket(Channels[i], "coroutine"+strconv.Itoa(i), InitializeStructure.ListRules)
		Handler.PacketHandlerArr[i].StartHandle()
	}

	//将数据包按照ip流生成的hash值分发给不同的数据包处理结构体
	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() != nil {
			ipFlow := packet.NetworkLayer().NetworkFlow()
			ipFlowHash := ipFlow.FastHash()
			Channels[ipFlowHash%uint64(Parameter.MaxCoroutine)] <- &packet
		} else {
			Channels[Parameter.MaxCoroutine-1] <- &packet
		}
	}
}
