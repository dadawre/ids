package Structure

import (
	"fmt"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

type MailSetting struct {
	smtpHost     string //邮箱类型
	smtpPort     string //邮箱端口
	authEmail    string //邮箱发送账户
	authPassword string //邮箱发送账户密码
	toEmail      string //接收账户
	subject      string //邮件主题
	body         string //邮件内容

	events map[string]*EventEntry //字符串触发事件映射
	mutex  sync.Mutex             // 互斥锁
}

// EventEntry 触发事件结构体
type EventEntry struct {
	FirstTriggerTime time.Time
	Count            int
}

// AddStringEvent 字符串事件触发添加
func (s *MailSetting) AddStringEvent(str string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.events == nil {
		s.events = make(map[string]*EventEntry)
	}

	if entry, exists := s.events[str]; exists {
		entry.Count++
	} else {
		s.events[str] = &EventEntry{
			FirstTriggerTime: time.Now(),
			Count:            1,
		}
	}
}

func (s *MailSetting) StartTimer(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			s.sendCollectedEvents()
		}
	}()
}

func (s *MailSetting) sendCollectedEvents() error {
	s.mutex.Lock()
	events := s.events
	s.events = make(map[string]*EventEntry)
	s.mutex.Unlock()
	if len(events) == 0 {
		return nil
	}
	fmt.Println("正在发送警告日志")
	var bodyBuilder strings.Builder
	for str, entry := range events {
		bodyBuilder.WriteString(fmt.Sprintf("首次触发时间: %s, 字符串: %s, 触发次数: %d\n",
			entry.FirstTriggerTime.Format("2006-01-02 15:04:05"), str, entry.Count))
	}

	// 构造邮件消息
	message := []byte(
		fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
			s.authEmail,
			s.toEmail,
			s.subject,
			bodyBuilder.String(),
		))

	auth := smtp.PlainAuth("", s.authEmail, s.authPassword, s.smtpHost)
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(addr, auth, s.authEmail, []string{s.toEmail}, message)
	if err != nil {
		fmt.Printf("定时邮件发送失败: %v\n", err)
		return err
	}
	fmt.Println("定时邮件发送成功")
	return nil
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
