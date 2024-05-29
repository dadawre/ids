package main

import (
	"graduationDesign/ruleAction"
	"testing"
)

func TestRuleAction(t *testing.T) {
	mail := &(ruleAction.MailSetting{})
	mail = mail.SetSmtpHost("smtp.qq.com").SetSmtpPort("587").SetToEmail("2939188260@qq.com")
	mail = mail.SetAuthEmail("112724772@qq.com").SetAuthPassword("bcernpppsoejbifi")
	mail = mail.SetSubject("Subject: 测试邮件\n").SetBody("这是通过Go语言发送的测试邮件。")
	err := mail.SendMail()
	if err != nil {
		err.Error()
	}
}
