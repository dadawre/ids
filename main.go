package main

import (
	"graduationDesign/parser/parserConfig"
	"graduationDesign/parser/parserRuler"
	"graduationDesign/structure/ruleStructure"
	"io/ioutil"
	_ "net/http/pprof"
)
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

var MaxCoroutine int = 1024
var PacketHandlerArr []*handler.PacketHandler

var ConfigPath = "./conf/ids.conf"
var ConfigVars map[string]*parserConfig.ConfigVar
var IncludePath []*string
var ConfigMapNameToValue map[string]string
var LogPath = "./log"
var MapVariableNameToVariable = map[string]interface{}{
	"max_coroutine": &MaxCoroutine,
}

var LISTRules *ruleStructure.RuleLinkedList

func InitializeConfig() {
	myConfigLexer := parserConfig.NewConfigLexer(ConfigPath)
	confTokens := myConfigLexer.GetAllToken()

	myConfigparser := parserConfig.NewConfigparser(&confTokens)

	myConfigparser.ParseAllTokens()

	ConfigVars = myConfigparser.TconfigVars

	IncludePath = myConfigparser.TincludePath

	ConfigMapNameToValue = myConfigparser.TConfigMapNameToValue
}

func InitializeRules() {
	LISTRules = &ruleStructure.RuleLinkedList{}
	parser := parserRuler.NewRuleParser("", ConfigVars)
	for index := range IncludePath {
		content, err := ioutil.ReadFile(*(IncludePath[index]))
		if err != nil {
			log.Fatal("Error reading file:", err)
			//return nil
		}

		rule := string(content)
		parser.Refresh(rule)
		err = parser.GetAllRule(LISTRules)

		if err != nil {
			fmt.Println("Error parsing rule:", err)
			return
		}
	}
}

func InitializeVariable() {
	for key, value := range ConfigMapNameToValue {
		//fmt.Println(key, "   ", value)
		if MapVariableNameToVariable[key] != nil {
			// 使用 strconv.Atoi 将字符串转换为 int
			num, err := strconv.Atoi(value)
			if err == nil {
				*(MapVariableNameToVariable[key].(*int)) = num
			} else {
				*(MapVariableNameToVariable[key]).(*string) = value
			}
		}
	}
	handler.LogPath = LogPath
}

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
	InitializeConfig()
	InitializeRules()
	InitializeVariable()
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
		PacketHandlerArr[i] = handler.NewHandlerPacket(Channels[i], "coroutine"+strconv.Itoa(i), LISTRules)
		PacketHandlerArr[i].StartHandle()
	}

	for packet := range sourcePacket.Packets() {
		//fmt.Println(packet)
		networkLayerType := packet.Layer(layers.LayerTypeIPv4) != nil || packet.Layer(layers.LayerTypeIPv6) != nil
		if networkLayerType != false {
			IPFlow := packet.NetworkLayer().NetworkFlow()
			IPFlowHash := IPFlow.FastHash()
			Channels[IPFlowHash%uint64(MaxCoroutine)] <- &packet
		}
	}

}
