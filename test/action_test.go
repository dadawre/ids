package main

import (
	"graduationDesign/action"
	"log"
	"sync"
	"testing"
)

func TestSendMail(t *testing.T) {
	mail := &(action.MailSetting{})
	mail = mail.SetSmtpHost("smtp.qq.com").SetSmtpPort("587").SetToEmail("2939188260@qq.com")
	mail = mail.SetAuthEmail("112724772@qq.com").SetAuthPassword("bcernpppsoejbifi")
	mail = mail.SetSubject("Subject: 测试邮件\n").SetBody("这是通过Go语言发送的测试邮件。")
	err := mail.SendMail()
	if err != nil {
		err.Error()
	}
}

func TestWriteFile(t *testing.T) {
	action.AppendWriteFile("./test.txt", "大大大大大大发发生飞洒飞洒发生大2")
}

func TestTCP(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	// Example usage
	go func() {
		if err := action.ReceiveData("9999"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	data := map[string]interface{}{
		"int":    123,
		"string": "Hello, World!",
		"slice":  []int{1, 2, 3, 4, 5},
		"float":  3.14,
	}

	if err := action.SendData("192.168.0.209:9999", data); err != nil {
		log.Fatalf("Client error: %v", err)
	}

	wg.Wait() // Wait for all goroutines to complete
}
