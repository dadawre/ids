package Initialize

import (
	"fmt"
	"graduationDesign/Handler"
	"graduationDesign/Parser/ParserConfig"
	"graduationDesign/Parser/ParserRuler"
	"graduationDesign/structure/ruleStructure"
	"log"
	"os"
	"strconv"
	"strings"
)

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
	Handler.LogPath = LogPath
}

func MailInitialize() {
	mail = mail.SetSmtpHost(SmtpHost).
		SetSmtpPort(strconv.Itoa(SmtpPort)).
		SetToEmail(ToEmail).
		SetAuthEmail(AuthEmail).
		SetAuthPassword(AuthPassword).
		SetSubject("Subject: 测试邮件\n").
		SetBody("这是通过Go语言发送的测试邮件。")
	Handler.RuleMailSetting = mail
	fmt.Println(mail)
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
