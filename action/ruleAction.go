package action

import (
	"fmt"
	"net/smtp"
)

type MailSetting struct {
	smtpHost     string
	smtpPort     string
	authEmail    string
	authPassword string
	toEmail      string
	subject      string
	body         string
}

func (s *MailSetting) SetSmtpHost(smtpHost string) *MailSetting {
	s.smtpHost = smtpHost
	return s
}

func (s *MailSetting) SetSmtpPort(smtpPort string) *MailSetting {
	s.smtpPort = smtpPort
	return s
}

func (s *MailSetting) SetAuthEmail(authEmail string) *MailSetting {
	s.authEmail = authEmail
	return s
}

func (s *MailSetting) SetAuthPassword(authPassword string) *MailSetting {
	s.authPassword = authPassword
	return s
}

func (s *MailSetting) SetToEmail(toEmail string) *MailSetting {
	s.toEmail = toEmail
	return s
}

func (s *MailSetting) SetSubject(subject string) *MailSetting {
	s.subject = subject
	return s
}

func (s *MailSetting) SetBody(body string) *MailSetting {
	s.body = body
	return s
}

func (s *MailSetting) SendMail() error {
	// QQ邮箱SMTP服务器配置
	//smtpHost := "smtp.qq.com"
	//smtpPort := "587"
	//authEmail := "your-email@qq.com" // 替换为你的QQ邮箱
	//authPassword := "your-smtp-auth-code" // 替换为你的SMTP授权码

	// 收件人地址
	//toEmail := "112724772@qq.com"

	// 邮件头和内容
	//subject := "Subject: 测试邮件\n"
	from := fmt.Sprintf("From: %s\n", s.authEmail)
	to := fmt.Sprintf("To: %s\n", s.toEmail)
	//body := "这是通过Go语言发送的测试邮件。"

	message := []byte(from + to + s.subject + "\n" + s.body)

	// 认证
	auth := smtp.PlainAuth("", s.authEmail, s.authPassword, s.smtpHost)

	// 发送邮件
	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.authEmail, []string{s.toEmail}, message)
	if err != nil {
		fmt.Println("邮件发送失败:", err)
		return err
	}

	fmt.Println("邮件发送成功")
	return nil
}
