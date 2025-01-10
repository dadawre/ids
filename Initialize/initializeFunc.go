package Initialize

import (
	"fmt"
	"graduationDesign/Parser/ParserConfig"
	"graduationDesign/Parser/ParserRuler"
	"graduationDesign/structure/ruleStructure"
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
}

// ConductionConduction 程序结束时执行的函数
func ConductionConduction() {
	//保存ip_ban集合的ipset配置
	cmd := exec.Command("/bin/bash", SaveBanIpShPath,BanIpConfigPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 打印脚本的输出
	fmt.Println("ip_ban save output:", string(output))

	//保存iptables集合的配置
	cmd = exec.Command("/bin/bash", SaveIptablesShPath,IptablesConfigPath)
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
		ConductionConduction()

		// 程序退出
		os.Exit(0)
	}()

}

func BashScriptInitialize() {
	cmd := exec.Command("/bin/bash", InitBanIpShPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 打印脚本的输出
	fmt.Println("Output:", string(output))

}

func ConfigInitialize() {
	ConfigLexer = ParserConfig.NewConfigLexer(ConfigPath) //词法分析器
	confTokens := ConfigLexer.GetAllToken()               //Token标志数组

	ConfigParser = ParserConfig.NewConfigparser(&confTokens) //语法分析器
	ConfigParser.ParseAllTokens()                            //语法解析
	//解析所有语法
	//fmt.Println("--------------------------------")
	//fmt.Println(ConfigParser.ConfigTokens)
	//fmt.Println("--------------------------------")
	//fmt.Println(ConfigParser.ConfigVars)
	//fmt.Println("--------------------------------")
	//fmt.Println(ConfigParser.RulePaths)
	//fmt.Println("--------------------------------")
	//fmt.Println(ConfigParser.ConfigMapNameToValue)
	//fmt.Println("--------------------------------")

	ConfigVars = ConfigParser.ConfigVars //变量

	RulePaths = ConfigParser.RulePaths //规则路径变量

	ConfigMapNameToValue = ConfigParser.ConfigMapNameToValue //配置变量键值对

}

func RulesInitialize() { //初始化规则链表
	ListRules = &ruleStructure.RuleLinkedList{}
	MainRulerParser = ParserRuler.NewRuleParser("", ConfigVars)
	for _, path := range RulePaths {
		content, err := os.ReadFile(*path)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		rule := string(content)
		MainRulerParser.Refresh(rule)
		err = MainRulerParser.GetAllRule(ListRules)
		if err != nil {
			log.Printf("Error parsing rule: %v", err)
		}
	}
}

func VariableInitialize() {
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
	fmt.Println(MaxCoroutine)
	MailInitialize()
	//Handler.LogPath = LogPath
}

func MailInitialize() {
	Mail = Mail.SetSmtpHost(SmtpHost).
		SetSmtpPort(strconv.Itoa(SmtpPort)).
		SetToEmail(ToEmail).
		SetAuthEmail(AuthEmail).
		SetAuthPassword(AuthPassword).
		SetSubject("Subject: 测试邮件\n").
		SetBody("这是通过Go语言发送的测试邮件。")
	//Handler.RuleMailSetting = Mail
	fmt.Println(Mail)
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
				LogPath = args[i+1]
				i++
			} else {
				fmt.Println("Error: missing value for", arg)
				return
			}
		case "-c", "--config":
			if i+1 < len(args) {
				ConfigPath = args[i+1]
				i++
			} else {
				fmt.Println("Error: missing value for", arg)
				return
			}
		default:
			if strings.HasPrefix(arg, "--log=") {
				LogPath = strings.SplitN(arg, "=", 2)[1]
			} else if strings.HasPrefix(arg, "--config=") {
				ConfigPath = strings.SplitN(arg, "=", 2)[1]
			} else {
				fmt.Println("Unknown argument:", arg)
				return
			}
		}
	}

	// 输出解析后的参数
	fmt.Printf("Log file path: %s\n", LogPath)
	fmt.Printf("Config file path: %s\n", ConfigPath)
	fmt.Println("硬件刷新频率", HardwareDataRefreshFrequency)
}
