package Initialize

import (
	"graduationDesign/Handler"
	"graduationDesign/Parser/ParserConfig"
	"graduationDesign/Parser/ParserRuler"
	"graduationDesign/structure/ruleStructure"
)

var ( //配置键值对，默认值
	MaxCoroutine                 = 512
	MinCoroutine                 = 256
	ConfigPath                   = "./conf/ids.conf"
	LogPath                      = "./log"
	AlertPath                    = "./log/alert.log"
	HardwareDataRefreshFrequency = 500
	SmtpHost                     string
	SmtpPort                     int
	ToEmail                      string
	AuthEmail                    string
	AuthPassword                 string
)

var ( //通用变量
	PacketHandlerArr          []*Handler.PacketHandler           //数据包解析结构体数组
	ConfigVars                map[string]*ParserConfig.ConfigVar //配置变量数组
	RulePaths                 []*string                          //路径变量数组
	ConfigMapNameToValue      map[string]string                  //配置变量键值对
	MapVariableNameToVariable = map[string]interface{}{          // 配置变量键值对，用interface{}任意类型表示
		"min_coroutine":        &MinCoroutine,                 //最大数据包解析协程
		"max_coroutine":        &MaxCoroutine,                 //最少数据包解析协程
		"config_path":          &ConfigPath,                   //配置文件路径
		"log_path":             &LogPath,                      //日志文件夹路径
		"alert_path":           &AlertPath,                    //告警日志文件路径
		"hardware_frequent":    &HardwareDataRefreshFrequency, //硬件数据刷新频率
		"mail_host":            &SmtpHost,                     //邮件服务器主机域名
		"mail_port":            &SmtpPort,                     //邮件服务器主机端口
		"mail_receive":         &ToEmail,                      //邮件接收人
		"mail_sender_name":     &AuthEmail,                    //邮件发送人账户
		"mail_sender_password": &AuthPassword,                 //邮件发送人密码
	}
	ConfigLexer     *ParserConfig.ConfigLexer     //配置文件词法分析器
	ConfigParser    *ParserConfig.Configparser    //配置文件语法分析器
	MainRulerParser *ParserRuler.RuleParser       //规则解析器
	ListRules       *ruleStructure.RuleLinkedList //规则链表
	mail            = &Handler.MailSetting{}      //邮箱结构体
)
