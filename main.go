package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"graduationDesign/Control"
	"graduationDesign/Handler"
	"graduationDesign/Initialize"
	"graduationDesign/OSController"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

// NetworkDevice 网络设备结构体
type NetworkDevice struct {
	Name        string
	Description string //网络设备描述
}

func getAllNetworkDev() ([]NetworkDevice, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, fmt.Errorf("error finding devices: %v", err)
	}

	var networkDevices []NetworkDevice
	for _, device := range devices {
		networkDevices = append(networkDevices, NetworkDevice{
			Name:        device.Name,
			Description: device.Description,
		})
	}
	return networkDevices, nil
}

func main() {
	Initialize.ConfigInitialize()
	Initialize.RulesInitialize()
	Initialize.VariableInitialize()
	Initialize.TerminateParameterInitialize()
	Control.RulerInputListen()

	go func() { //程序状态监控
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	devices, err := getAllNetworkDev()
	for index, device := range devices {
		fmt.Println(index, "  ", device.Name, "   ", device.Description)
	}

	var userChooseDevice int
	fmt.Println("请输入想要监听的设备ID")
	_, err = fmt.Scan(&userChooseDevice)

	if err != nil {
		fmt.Println("输入设备 ID 有误")
		return
	}
	fmt.Println("开始监听设备", devices[userChooseDevice].Description)
	inactiveHandle, err := pcap.NewInactiveHandle(devices[userChooseDevice].Name)
	if err != nil {
		log.Fatalf("Error creating inactive handle: %v", err)
	}
	defer inactiveHandle.CleanUp()

	go func() { //硬件利用率刷新
		OSController.GetOSStates(time.Duration(Initialize.HardwareDataRefreshFrequency) * time.Millisecond) //多少ms刷新一次硬件数据
	}()

	if err := inactiveHandle.SetTimeout(time.Second); err != nil {
		log.Fatalf("Error setting timeout: %v", err)
	}

	if err := inactiveHandle.SetPromisc(true); err != nil {
		log.Fatalf("Error setting promiscuous mode: %v", err)
	}

	handle, err := inactiveHandle.Activate()
	if err != nil {
		log.Fatalf("Error activating handle: %v", err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	Channels := make([]chan *gopacket.Packet, Initialize.MaxCoroutine)
	Initialize.PacketHandlerArr = make([]*Handler.PacketHandler, Initialize.MaxCoroutine)

	for i := 0; i < Initialize.MaxCoroutine; i++ {
		Channels[i] = make(chan *gopacket.Packet)
		Initialize.PacketHandlerArr[i] = Handler.NewHandlerPacket(Channels[i], "coroutine"+strconv.Itoa(i), Initialize.ListRules)
		Initialize.PacketHandlerArr[i].StartHandle()
	}

	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() != nil {
			ipFlow := packet.NetworkLayer().NetworkFlow()
			ipFlowHash := ipFlow.FastHash()
			Channels[ipFlowHash%uint64(Initialize.MaxCoroutine)] <- &packet
		} else if packet.Layer(layers.LayerTypeARP) != nil {
			Channels[Initialize.MaxCoroutine-1] <- &packet
		}
	}
}
