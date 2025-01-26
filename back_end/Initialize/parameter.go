package Initialize

import (
	ParserConfig2 "graduationDesign/back_end/Parser/ParserConfig"
	"graduationDesign/back_end/Parser/ParserRuler"
	"graduationDesign/back_end/structure/ruleStructure"
)

var ( //通用变量

	ConfigVars                    map[string]*ParserConfig2.ConfigVar              //配置变量数组
	RulePaths                     []*string                                        //路径变量数组
	ConfigMapNameToValue          map[string]string                                //配置变量键值对
	ConfigLexer                   *ParserConfig2.ConfigLexer                       //配置文件词法分析器
	ConfigParser                  *ParserConfig2.Configparser                      //配置文件语法分析器
	MainRulerParser               *ParserRuler.RuleParser                          //规则解析器
	ListRules                     *ruleStructure.RuleLinkedList                    //规则链表
	ArrayRuleStrings              []ruleStructure.RuleNodeStringValue              //规则字符串结构体数组
	Mail                          = &MailSetting{}                                 //邮箱结构体
	ActionDescribeTriggerLogArray = make([]ActionDescribeTriggerLog, 0)            //Action执行日志
	BanBanIPSet                   = BanIPSet{data: make(map[string]BanIPDescribe)} //封禁IP集合
	MapVariableNameToVariable     = map[string]interface{}{                        // 配置变量键值对，用interface{}任意类型表示
		"min_coroutine":     &MinCoroutine,                 //最大数据包解析协程
		"max_coroutine":     &MaxCoroutine,                 //最少数据包解析协程
		"config_path":       &ConfigPath,                   //配置文件路径
		"log_path":          &LogPath,                      //日志文件夹路径
		"alert_path":        &AlertPath,                    //告警日志文件路径
		"hardware_frequent": &HardwareDataRefreshFrequency, //硬件数据刷新频率

		"mail_host":            &SmtpHost,     //邮件服务器主机域名
		"mail_port":            &SmtpPort,     //邮件服务器主机端口
		"mail_receive":         &ToEmail,      //邮件接收人
		"mail_sender_name":     &AuthEmail,    //邮件发送人账户
		"mail_sender_password": &AuthPassword, //邮件发送人密码

		"init_ban_ip_sh_path":       &InitBanIpShPath,     //ipset的ban_ip集合初始化脚本文件路径
		"add_ip_to_ban_sh_path":     &AddIpToBanShPath,    //添加IP到ban_ip中的脚本文件路径
		"ban_ip_delete_from_ban_sh": &DeleteIpFromShPath,  //从ban_ip中删除IP的脚本文件路径
		"save_ban_ip_sh_path":       &SaveBanIpShPath,     //保存ban_ip ipset配置到指定文件的脚本文件路径
		"ban_ip_save_path":          &BanIpConfigSavePath, //ban_ip配置文件路径

		"init_ban_ip_network_sh_path":       &InitBanIpNetWorkShPath,     //ipset的ban_ip_network集合初始化脚本文件路径
		"add_ip_network_to_ban_sh_path":     &AddIpNetworkToBanShPath,    //添加IP到ban_ip_network中的脚本文件路径
		"ban_ip_network_delete_from_ban_sh": &DeleteIpNetworkFromShPath,  //从ban_ip中删除IP的脚本文件路径
		"save_ban_ip_network_sh_path":       &SaveBanIpNetworkShPath,     //保存ban_ip_network ipset配置到指定文件的脚本文件路径
		"ban_ip_network_save_path":          &BanIpNetworkConfigSavePath, //ban_ip_network配置文件路径

		"save_iptables_sh_path": &SaveIptablesShPath,     //保存iptables配置到指定文件的脚本文件路径
		"iptables_save_path":    &IptablesConfigSavePath, //iptables配置文件路径

		"ids_mode": &IdsMode, //ids模式
		"ids_name": &IdsName, //ids主机名
		"ids_id":   &IdsID,   //ids标识符
	}
)

var ( //配置键值对，默认值
	MaxCoroutine                 = 512
	MinCoroutine                 = 256
	ConfigPath                   = "./conf/ids.conf"
	LogPath                      = "./log"
	AlertPath                    = "./log/alert.log"
	HardwareDataRefreshFrequency = 500

	SmtpHost     string
	SmtpPort     int
	ToEmail      string
	AuthEmail    string
	AuthPassword string

	InitBanIpShPath     string
	AddIpToBanShPath    string
	DeleteIpFromShPath  string
	SaveBanIpShPath     string
	BanIpConfigSavePath string

	InitBanIpNetWorkShPath     string
	AddIpNetworkToBanShPath    string
	DeleteIpNetworkFromShPath  string
	SaveBanIpNetworkShPath     string
	BanIpNetworkConfigSavePath string

	SaveIptablesShPath     string
	IptablesConfigSavePath string

	IdsMode = 1
	IdsName = "node1"
	IdsID   = 1
)
