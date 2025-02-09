package Initialize

import (
	"fmt"
	"graduationDesign/back_end/DB"
	"graduationDesign/back_end/Initialize/InitializeStructure"
	"graduationDesign/back_end/Initialize/Parameter"
	"graduationDesign/back_end/List"
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
	"time"
)

// MainInitialize ids初始化总函数
func MainInitialize() {
	ConfigInitialize()             //配置文件解析
	RulesInitialize()              //规矩解析
	VariableInitialize()           //配置变量解析
	TerminateParameterInitialize() //终端参数解析
	BashScriptInitialize()         //bash解析
	EndListenInitialize()          //结束监听
	EmailInitialize()              //邮件发送
	WhiteAndBlackListInitialize()  //黑白名单
	SQLInitialize()                //SQL管理
}

func ConfigInitialize() {
	InitializeStructure.ConfigLexer = ParserConfig.NewConfigLexer(Parameter.ConfigPath) //词法分析器
	confTokens := InitializeStructure.ConfigLexer.GetAllToken()                         //Token标志数组

	InitializeStructure.ConfigParser = ParserConfig.NewConfigparser(&confTokens) //语法分析器
	InitializeStructure.ConfigParser.ParseAllTokens()                            //语法解析

	InitializeStructure.ConfigVars = InitializeStructure.ConfigParser.ConfigVars //变量

	InitializeStructure.RulePaths = InitializeStructure.ConfigParser.RulePaths //规则路径变量

	InitializeStructure.ConfigMapNameToValue = InitializeStructure.ConfigParser.ConfigMapNameToValue //配置变量键值对

}

// RulesInitialize 初始化规则链表和数组
func RulesInitialize() {
	InitializeStructure.ListRules = &Structure.RuleLinkedList{}
	InitializeStructure.MainRulerParser = ParserRuler.NewRuleParser("", InitializeStructure.ConfigVars)

	for _, path := range InitializeStructure.RulePaths {
		content, err := os.ReadFile(*path)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		rule := string(content)
		InitializeStructure.MainRulerParser.Refresh(rule)
		err = InitializeStructure.MainRulerParser.GetAllRule(InitializeStructure.ListRules)

		if err != nil {
			log.Printf("Error parsing rule: %v", err)
		}
	}
	//fmt.Println("---------------------------------")
	//fmt.Println(InitializeStructure.ListRules.GetHead().RuleValue)
	InitializeStructure.ArrayRuleStrings = InitializeStructure.ListRules.GetStringArray()
	//fmt.Println(ArrayRuleStrings)
}

func VariableInitialize() {
	for key, value := range InitializeStructure.ConfigMapNameToValue {
		if varPtr, exists := InitializeStructure.MapVariableNameToVariable[key]; exists {
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

func BashScriptInitialize() {
	cmd := exec.Command("/bin/bash", Parameter.InitBanIpShPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Output:", string(output))

	cmd = exec.Command("/bin/bash", Parameter.InitBanIpNetWorkShPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Output:", string(output))

	cmd = exec.Command("/bin/bash", Parameter.InitWhiteListShPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Initialize White Set Error:", err)
	}
	fmt.Println("Output:", string(output))

	cmd = exec.Command("/bin/bash", Parameter.InitWhiteListNetworkShPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Initialize White Network Set  Error:", err)
	}
	// 打印脚本的输出
	fmt.Println("Output:", string(output))
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
	fmt.Println("ip_network_ban save output:", string(output))

	//保存ip_ban集合的ipset配置
	cmd = exec.Command("/bin/bash", Parameter.SaveWhiteListShPath, Parameter.WhiteListSavePath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 打印脚本的输出
	fmt.Println("white list save output:", string(output))

	//保存ip_ban集合的ipset配置
	cmd = exec.Command("/bin/bash", Parameter.SaveWhiteListNetworkShPath, Parameter.WhiteListNetworkSavePath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 打印脚本的输出
	fmt.Println("white list network save output:", string(output))

	//保存iptables集合的配置
	cmd = exec.Command("/bin/bash", Parameter.SaveIptablesShPath, Parameter.IptablesConfigSavePath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 打印脚本的输出
	fmt.Println("iptables output:", string(output))

	//ipset清空
	cmd = exec.Command("/bin/bash", "ipset flush")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 打印脚本的输出
	fmt.Println("ipset flush output:", string(output))
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

func EmailInitialize() {
	InitializeStructure.Mail = InitializeStructure.Mail.SetSmtpHost(Parameter.SmtpHost).
		SetSmtpPort(strconv.Itoa(Parameter.SmtpPort)).
		SetToEmail(Parameter.ToEmail).
		SetAuthEmail(Parameter.AuthEmail).
		SetAuthPassword(Parameter.AuthPassword).
		SetSubject(fmt.Sprintf("IDS告警 - ID:%d 名称:%s", Parameter.IdsID, Parameter.IdsName))

	InitializeStructure.Mail.AddStringEvent("这是通过Go语言发送的测试邮件。")
	InitializeStructure.Mail.StartTimer(time.Duration(Parameter.EmailSendInterval))
	//Handler.RuleMailSetting = Mail
	//fmt.Println(Mail)
}

func WhiteAndBlackListInitialize() {
	// 示例用法
	_, _ = List.LoadAccessList(Parameter.WhiteListPath, true, InitializeStructure.BlackWhiteList)
	_, _ = List.LoadAccessList(Parameter.BlackListPath, false, InitializeStructure.BlackWhiteList)
	_, _ = List.LoadAccessList(Parameter.BlackPortListPath, false, InitializeStructure.BlackWhiteList)

	fmt.Println(InitializeStructure.BlackWhiteList)
	fmt.Println(len(InitializeStructure.BlackWhiteList.Addrs))
}

func SQLInitialize() {
	Parameter.SqlRecordTableFullName = Parameter.SqlSchemaName + "." + Parameter.SqlRecordTableName
	// MySQL 连接信息 (请替换成你的 MySQL 配置)
	dsn := DB.ConnectAddrStr(Parameter.SqlUser, Parameter.SqlPassword, Parameter.SqlHost)

	// 创建 MySQL 管理实例
	var err error
	InitializeStructure.LogSQLManager, err = DB.NewMySQLManager(dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 确保数据库存在
	if err := InitializeStructure.LogSQLManager.EnsureDatabase(Parameter.SqlSchemaName); err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}

	// 切换到 `IDSdata` 数据库
	InitializeStructure.LogSQLManager.DB = InitializeStructure.LogSQLManager.DB.Exec("USE IDSdata")
	fullPathTableName := Parameter.SqlSchemaName + "." + Parameter.SqlRecordTableName

	// 检查表是否存在
	exists, err := InitializeStructure.LogSQLManager.TableExists(Parameter.SqlSchemaName, Parameter.SqlRecordTableName)
	if err != nil {
		log.Fatalf("检查表失败: %v", err)
	}
	if !exists {
		log.Println("表不存在，正在创建...")
		if err := InitializeStructure.LogSQLManager.CreateTable(fullPathTableName); err != nil {
			log.Fatalf("创建表失败: %v", err)
		}
	}
}
