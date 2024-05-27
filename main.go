package main

import _ "net/http/pprof"
import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"graduationDesign/handler"
	"log"
	"net/http"
	"strconv"
	"time"
)

var MaxCoroutine uint64 = 1024
var PacketHandlerArr []*handler.PacketHandler

func getAllNetworkDev() {
	devices, err := pcap.FindAllDevs() //获取所有网络接口
	if err != nil {
		log.Fatal(err)
	}

	// 打印可用的网络接口
	for _, device := range devices {
		fmt.Println("NameId:", device.Name)
		fmt.Println("Description:", device.Description)
		fmt.Println("Addresses:", device.Addresses)
		fmt.Println()
	}
}

func main() {
	//go getOSStates()
	//getAllNetworkDev()
	//handle, err := pcap.OpenLive("\\Device\\NPF_{8809DE90-DC4E-4E82-8495-740651F390A5}", 65536, true, pcap.BlockForever)
	inactiveHandle, err := pcap.NewInactiveHandle("\\Device\\NPF_{8809DE90-DC4E-4E82-8495-740651F390A5}")
	if err != nil {
		fmt.Println("Error creating inactive handle:", err)
		return
	}
	defer inactiveHandle.CleanUp() // 确保在函数结束时关闭未激活的 pcap 句柄

	err = inactiveHandle.SetTimeout(time.Second)
	if err != nil {
		log.Fatal(err)
	}

	if err := inactiveHandle.SetPromisc(true); err != nil {
		fmt.Println("Error setting promisc mode:", err)
	}

	// 激活 pcap 句柄
	handle, err := inactiveHandle.Activate()
	if err != nil {
		fmt.Println("Error activating handle:", err)
		return
	}
	defer handle.Close() // 确保在函数结束时关闭 pcap 句柄

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	sourcePacket := gopacket.NewPacketSource(handle, handle.LinkType())

	Channels := make([]chan *gopacket.Packet, MaxCoroutine+1, MaxCoroutine+1)
	PacketHandlerArr = make([]*handler.PacketHandler, MaxCoroutine+1, MaxCoroutine+1)

	for i, _ := range Channels {
		Channels[i] = make(chan *gopacket.Packet)
		PacketHandlerArr[i] = handler.NewHandlerPacket(Channels[i], "coroutine"+strconv.Itoa(i))
		PacketHandlerArr[i].StartHandle()
	}

	for packet := range sourcePacket.Packets() {
		//fmt.Println(packet)
		networkLayerType := packet.Layer(layers.LayerTypeIPv4) != nil || packet.Layer(layers.LayerTypeIPv6) != nil
		if networkLayerType != false {
			IPFlow := packet.NetworkLayer().NetworkFlow()
			IPFlowHash := IPFlow.FastHash()
			Channels[IPFlowHash%MaxCoroutine] <- &packet
		}
	}

}
