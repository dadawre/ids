package InitializeStructure

import (
	"graduationDesign/back_end/DB"
	"graduationDesign/back_end/Initialize/Parameter"
	"graduationDesign/back_end/List"
	"graduationDesign/back_end/Parser/ParserConfig"
	"graduationDesign/back_end/Parser/ParserRuler"
	"graduationDesign/back_end/Structure"
	"time"
)

var ( //通用变量

	ConfigVars           map[string]*ParserConfig.ConfigVar //配置变量数组
	RulePaths            []*string                          //路径变量数组
	ConfigMapNameToValue map[string]string                  //配置变量键值对
	ConfigLexer          *ParserConfig.ConfigLexer          //配置文件词法分析器
	ConfigParser         *ParserConfig.Configparser         //配置文件语法分析器
	MainRulerParser      *ParserRuler.RuleParser            //规则解析器
	ListRules            *Structure.RuleLinkedList          //规则链表
	ArrayRuleStrings     []Structure.RuleNodeStringValue    //规则字符串结构体数组
	Mail                 = &Structure.MailSetting{}         //邮箱结构体
	//ActionDescribeTriggerLogArray = make([]Structure.ActionDescribeTriggerLog, 0) //Action执行日志
	MainBanIPSet              = Structure.NewBanIPSet(time.Duration(Parameter.IPBanAllowedDiff * 1000000))     //封禁IP集合
	BlackWhiteList            = &List.AccessList{Addrs: make(map[string]List.WBValueStructure, 0), IDindex: 0} //黑白名单
	LogSQLManager             *DB.RecordSQLManager                                                             //日志SQL驱动器
	MapVariableNameToVariable = map[string]interface{}{                                                        // 配置变量键值对，用interface{}任意类型表示
		"min_coroutine":       &Parameter.MinCoroutine,                 //最大数据包解析协程
		"max_coroutine":       &Parameter.MaxCoroutine,                 //最少数据包解析协程
		"config_path":         &Parameter.ConfigPath,                   //配置文件路径
		"log_path":            &Parameter.LogPath,                      //日志文件夹路径
		"alert_path":          &Parameter.AlertPath,                    //告警日志文件路径
		"hardware_frequent":   &Parameter.HardwareDataRefreshFrequency, //硬件数据刷新频率
		"ban_allowed_diff":    &Parameter.IPBanAllowedDiff,             //解禁ip允许时间误差
		"email_send_interval": &Parameter.EmailSendInterval,            //短信发送间隔

		"sql_user":         &Parameter.SqlUser,
		"sql_password":     &Parameter.SqlPassword,
		"sql_host":         &Parameter.SqlHost,
		"sql_schema":       &Parameter.SqlSchemaName,
		"sql_record_table": &Parameter.SqlRecordTableName,

		"white_list_path":      &Parameter.WhiteListPath,     //白名单文件路径
		"black_list_path":      &Parameter.BlackListPath,     //黑名单文件路径
		"port_black_list_path": &Parameter.BlackPortListPath, //端口黑名单

		"ids_mode": &Parameter.IdsMode, //ids模式
		"ids_name": &Parameter.IdsName, //ids主机名
		"ids_id":   &Parameter.IdsID,   //ids标识符

		"account":  &Parameter.Account,  //账户
		"password": &Parameter.Password, //密码

		"mail_host":            &Parameter.SmtpHost,     //邮件服务器主机域名
		"mail_port":            &Parameter.SmtpPort,     //邮件服务器主机端口
		"mail_receive":         &Parameter.ToEmail,      //邮件接收人
		"mail_sender_name":     &Parameter.AuthEmail,    //邮件发送人账户
		"mail_sender_password": &Parameter.AuthPassword, //邮件发送人密码

		"init_ban_ip_sh_path":       &Parameter.InitBanIpShPath,     //ipset的ban_ip集合初始化脚本文件路径
		"add_ip_to_ban_sh_path":     &Parameter.AddIpToBanShPath,    //添加IP到ban_ip中的脚本文件路径
		"ban_ip_delete_from_ban_sh": &Parameter.DeleteIpFromShPath,  //从ban_ip中删除IP的脚本文件路径
		"save_ban_ip_sh_path":       &Parameter.SaveBanIpShPath,     //保存ban_ip ipset配置到指定文件的脚本文件路径
		"ban_ip_save_path":          &Parameter.BanIpConfigSavePath, //ban_ip配置文件路径

		"init_ban_ip_network_sh_path":       &Parameter.InitBanIpNetWorkShPath,     //ipset的ban_ip_network集合初始化脚本文件路径
		"add_ip_network_to_ban_sh_path":     &Parameter.AddIpNetworkToBanShPath,    //添加IP到ban_ip_network中的脚本文件路径
		"ban_ip_network_delete_from_ban_sh": &Parameter.DeleteIpNetworkFromShPath,  //从ban_ip中删除IP的脚本文件路径
		"save_ban_ip_network_sh_path":       &Parameter.SaveBanIpNetworkShPath,     //保存ban_ip_network ipset配置到指定文件的脚本文件路径
		"ban_ip_network_save_path":          &Parameter.BanIpNetworkConfigSavePath, //ban_ip_network配置文件路径

		"init_white_list_sh_path":      &Parameter.InitWhiteListShPath,
		"add_ip_to_white_list_sh_path": &Parameter.AddIpToWhiteListShPath,
		"delete_white_ip_sh":           &Parameter.DeleteWhiteIpSh,
		"save_white_list_sh_path":      &Parameter.SaveWhiteListShPath,
		"white_list_save_path":         &Parameter.WhiteListSavePath,

		"init_white_list_network_sh_path":      &Parameter.InitWhiteListNetworkShPath,
		"add_ip_to_white_list_network_sh_path": &Parameter.AddIpToWhiteListNetworkShPath,
		"delete_white_ip_network_sh":           &Parameter.DeleteWhiteIpNetworkSh,
		"save_white_list_network_sh_path":      &Parameter.SaveWhiteListNetworkShPath,
		"white_list_network_save_path":         &Parameter.WhiteListNetworkSavePath,

		"init_ban_port_sh_path":       &Parameter.InitBanPortShPath,
		"add_port_to_ban_sh_path":     &Parameter.AddPortToBanShPath,
		"ban_port_delete_from_ban_sh": &Parameter.DeleteBanPortFromBanSh,
		"save_ban_port_sh_path":       &Parameter.SaveBanPortShPath,
		"ban_port_save_path":          &Parameter.BanPortSavePath,

		"save_iptables_sh_path": &Parameter.SaveIptablesShPath,     //保存iptables配置到指定文件的脚本文件路径
		"iptables_save_path":    &Parameter.IptablesConfigSavePath, //iptables配置文件路径

	}
)
