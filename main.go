package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"graduationDesign/OSController"
	"graduationDesign/handler"
	"graduationDesign/parser/parserConfig"
	"graduationDesign/parser/parserRuler"
	"graduationDesign/structure/ruleStructure"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

var (
	MaxCoroutine = 1024
	ConfigPath   = "./conf/ids.conf"
	LogPath      = "./log"
)

var (
	PacketHandlerArr          []*handler.PacketHandler
	ConfigVars                map[string]*parserConfig.ConfigVar
	IncludePaths              []*string
	ConfigMapNameToValue      map[string]string
	MapVariableNameToVariable = map[string]interface{}{"max_coroutine": new(int)}
	LISTRules                 *ruleStructure.RuleLinkedList
)

func InitializeConfig() {
	configLexer := parserConfig.NewConfigLexer(ConfigPath)
	confTokens := configLexer.GetAllToken()

	configParser := parserConfig.NewConfigparser(&confTokens)
	configParser.ParseAllTokens()

	ConfigVars = configParser.ConfigVars
	IncludePaths = configParser.IncludePath
	ConfigMapNameToValue = configParser.ConfigMapNameToValue
}

func InitializeRules() {
	LISTRules = &ruleStructure.RuleLinkedList{}
	parser := parserRuler.NewRuleParser("", ConfigVars)
	for _, path := range IncludePaths {
		content, err := ioutil.ReadFile(*path)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		rule := string(content)
		parser.Refresh(rule)
		err = parser.GetAllRule(LISTRules)
		if err != nil {
			log.Printf("Error parsing rule: %v", err)
		}
	}
}

func InitializeVariable() {
	for key, value := range ConfigMapNameToValue {
		if varPtr, exists := MapVariableNameToVariable[key]; exists {
			num, err := strconv.Atoi(value)
			if err == nil {
				*(varPtr.(*int)) = num
			} else {
				*(varPtr.(*string)) = value
			}
		}
	}
	InitializeMail()
	handler.LogPath = LogPath
}

func InitializeMail() {
	mail := &handler.MailSetting{}
	mail = mail.SetSmtpHost("smtp.qq.com").
		SetSmtpPort("587").
		SetToEmail("2939188260@qq.com").
		SetAuthEmail("112724772@qq.com").
		SetAuthPassword("bcernpppsoejbifi").
		SetSubject("Subject: 测试邮件\n").
		SetBody("这是通过Go语言发送的测试邮件。")
	handler.RuleMailSetting = mail
}

type NetworkDevice struct {
	Name        string
	Description string
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
	OSController.GetOSStates(1 * time.Second)
	InitializeConfig()
	InitializeRules()
	InitializeVariable()

	inactiveHandle, err := pcap.NewInactiveHandle("\\Device\\NPF_{8809DE90-DC4E-4E82-8495-740651F390A5}")
	if err != nil {
		log.Fatalf("Error creating inactive handle: %v", err)
	}
	defer inactiveHandle.CleanUp()

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

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	Channels := make([]chan *gopacket.Packet, MaxCoroutine)
	PacketHandlerArr = make([]*handler.PacketHandler, MaxCoroutine)

	for i := 0; i < MaxCoroutine; i++ {
		Channels[i] = make(chan *gopacket.Packet)
		PacketHandlerArr[i] = handler.NewHandlerPacket(Channels[i], "coroutine"+strconv.Itoa(i), LISTRules)
		PacketHandlerArr[i].StartHandle()
	}

	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() != nil {
			ipFlow := packet.NetworkLayer().NetworkFlow()
			ipFlowHash := ipFlow.FastHash()
			Channels[ipFlowHash%uint64(MaxCoroutine)] <- &packet
		}
	}
}
