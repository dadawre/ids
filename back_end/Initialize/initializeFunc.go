package Initialize

import (
	"fmt"
	Structure2 "graduationDesign/back_end/Initialize/InitializeStructure"
	"graduationDesign/back_end/Initialize/Parameter"
	"graduationDesign/back_end/Parser/ParserConfig"
	"graduationDesign/back_end/Parser/ParserRuler"
	"graduationDesign/back_end/Structure"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// MainInitialize ids初始化总函数
func MainInitialize() {
	ConfigInitialize()
	RulesInitialize()
	VariableInitialize()
	TerminateParameterInitialize()
	BashScriptInitialize()
	EndListenInitialize()
	MailInitialize()
}

// EndConduction 程序结束时执行的函数
func EndConduction() {
	//保存ip_ban集合的ipset配置
	cmd := exec.Command("/bin/bash", Parameter.SaveBanIpShPath, Parameter.BanIpConfigSavePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 打印脚本的输出
	fmt.Println("ip_ban save output:", string(output))

	//保存ip_ban集合的ipset配置
	cmd = exec.Command("/bin/bash", Parameter.SaveBanIpNetworkShPath, Parameter.BanIpNetworkConfigSavePath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 打印脚本的输出
	fmt.Println("ip_ban save output:", string(output))

	//保存iptables集合的配置
	cmd = exec.Command("/bin/bash", Parameter.SaveIptablesShPath, Parameter.IptablesConfigSavePath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 打印脚本的输出
	fmt.Println("iptables output:", string(output))
}

// EndListenInitialize 程序结束监听函数初始化
func EndListenInitialize() {
	// 创建一个用于接收信号的 channel
	signalChan := make(chan os.Signal, 1)

	// 捕获 SIGINT 和 SIGTERM 信号
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个 goroutine 监听信号
	go func() {
		// 阻塞直到接收到信号
		sigReceived := <-signalChan
		fmt.Printf("接收到信号: %s，开始执行退出清理...\n", sigReceived)

		// 在程序关闭前执行清理函数
		EndConduction()

		// 程序退出
		os.Exit(0)
	}()

}

func BashScriptInitialize() {
	cmd := exec.Command("/bin/bash", Parameter.InitBanIpShPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	cmd = exec.Command("/bin/bash", Parameter.InitBanIpNetWorkShPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 打印脚本的输出
	fmt.Println("Output:", string(output))

}

func ConfigInitialize() {
	Structure2.ConfigLexer = ParserConfig.NewConfigLexer(Parameter.ConfigPath) //词法分析器
	confTokens := Structure2.ConfigLexer.GetAllToken()                         //Token标志数组

	Structure2.ConfigParser = ParserConfig.NewConfigparser(&confTokens) //语法分析器
	Structure2.ConfigParser.ParseAllTokens()                            //语法解析

	Structure2.ConfigVars = Structure2.ConfigParser.ConfigVars //变量

	Structure2.RulePaths = Structure2.ConfigParser.RulePaths //规则路径变量

	Structure2.ConfigMapNameToValue = Structure2.ConfigParser.ConfigMapNameToValue //配置变量键值对

}

// RulesInitialize 初始化规则链表和数组
func RulesInitialize() {
	Structure2.ListRules = &Structure.RuleLinkedList{}
	Structure2.MainRulerParser = ParserRuler.NewRuleParser("", Structure2.ConfigVars)

	for _, path := range Structure2.RulePaths {
		content, err := os.ReadFile(*path)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		rule := string(content)
		Structure2.MainRulerParser.Refresh(rule)
		err = Structure2.MainRulerParser.GetAllRule(Structure2.ListRules)

		if err != nil {
			log.Printf("Error parsing rule: %v", err)
		}
	}
	//fmt.Println("---------------------------------")
	//fmt.Println(ListRules.GetHead().RuleValue)
	Structure2.ArrayRuleStrings = Structure2.ListRules.GetStringArray()
	//fmt.Println(ArrayRuleStrings)
}

func VariableInitialize() {
	for key, value := range Structure2.ConfigMapNameToValue {
		if varPtr, exists := Structure2.MapVariableNameToVariable[key]; exists {
			num, err := strconv.Atoi(value)
			if err == nil {
				*(varPtr.(*int)) = num
			} else {
				*(varPtr.(*string)) = value
			}
		}
	}
	//fmt.Println(MaxCoroutine)
	//Handler.LogPath = LogPath
}

func MailInitialize() {
	Structure2.Mail = Structure2.Mail.SetSmtpHost(Parameter.SmtpHost).
		SetSmtpPort(strconv.Itoa(Parameter.SmtpPort)).
		SetToEmail(Parameter.ToEmail).
		SetAuthEmail(Parameter.AuthEmail).
		SetAuthPassword(Parameter.AuthPassword).
		SetSubject("Subject: 测试邮件\n").
		SetBody("这是通过Go语言发送的测试邮件。")
	//Handler.RuleMailSetting = Mail
	//fmt.Println(Mail)
}

func TerminateParameterInitialize() {
	// 手动解析命令行参数
	args := os.Args[1:]
	for index, arg := range args {
		fmt.Println(index, "  ", arg)
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-l", "--log":
			if i+1 < len(args) {
				Parameter.LogPath = args[i+1]
				i++
			} else {
				fmt.Println("Error: missing value for", arg)
				return
			}
		case "-c", "--config":
			if i+1 < len(args) {
				Parameter.ConfigPath = args[i+1]
				i++
			} else {
				fmt.Println("Error: missing value for", arg)
				return
			}
		default:
			if strings.HasPrefix(arg, "--log=") {
				Parameter.LogPath = strings.SplitN(arg, "=", 2)[1]
			} else if strings.HasPrefix(arg, "--config=") {
				Parameter.ConfigPath = strings.SplitN(arg, "=", 2)[1]
			} else {
				fmt.Println("Unknown argument:", arg)
				return
			}
		}
	}

	// 输出解析后的参数
	//fmt.Printf("Log file path: %s\n", LogPath)
	//fmt.Printf("Config file path: %s\n", ConfigPath)
	//fmt.Println("硬件刷新频率", HardwareDataRefreshFrequency)
}
