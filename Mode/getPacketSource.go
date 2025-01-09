package Mode

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

// NetworkDevice 网络设备结构体
type NetworkDevice struct {
	Name        string //网络设备名
	Description string //网络设备描述
}

func GetAllNetworkDev() ([]NetworkDevice, error) {
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

var packetSource *gopacket.PacketSource
var inactiveHandle *pcap.InactiveHandle
var handle *pcap.Handle

func NewPacketSource() *gopacket.PacketSource {
	devices, err := GetAllNetworkDev()
	for index, device := range devices {
		fmt.Println(index, "  ", device.Name, "   ", device.Description)
	}

	var userChooseDevice int
	fmt.Println("请输入想要监听的设备ID")
	_, err = fmt.Scan(&userChooseDevice)
	if err != nil {
		fmt.Println("输入设备 ID 有误")
		return nil
	}
	fmt.Println("开始监听设备", devices[userChooseDevice].Description)

	//根据设备名获取设备句柄
	inactiveHandle, err = pcap.NewInactiveHandle(devices[userChooseDevice].Name)
	if err != nil {
		log.Fatalf("Error creating inactive handle: %v", err)
	}
	//defer inactiveHandle.CleanUp()

	//设置句柄读取超时时间
	if err := inactiveHandle.SetTimeout(time.Second); err != nil {
		log.Fatalf("Error setting timeout: %v", err)
	}

	//设计是否为混杂模式
	if err := inactiveHandle.SetPromisc(true); err != nil {
		log.Fatalf("Error setting promiscuous mode: %v", err)
	}

	//获取激活的设备句柄
	handle, err = inactiveHandle.Activate()
	if err != nil {
		log.Fatalf("Error activating handle: %v", err)
	}
	//defer handle.Close()

	//通过设备句柄获取数据包源
	packetSource = gopacket.NewPacketSource(handle, handle.LinkType())
	return packetSource
}
