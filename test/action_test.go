package main

import (
	"fmt"
	"graduationDesign/Handler"
	"log"
	"sync"
	"testing"
)

func TestSendMail(t *testing.T) {
	mail := &(Handler.MailSetting{})
	mail = mail.SetSmtpHost("smtp.qq.com").SetSmtpPort("587").SetToEmail("3555923360@qq.com")
	mail = mail.SetAuthEmail("112724772@qq.com").SetAuthPassword("bcernpppsoejbifi")
	mail = mail.SetSubject("Subject: 测试邮件\n").SetBody(`这是通过Go语言发送的测试邮件。<br><a href="http://192.168.190.42/pikachu/vul/xss/xss_stored.php">点击这里访问链接</a>`)
	err := mail.SendMail()
	if err != nil {
		err.Error()
	}
}

func TestWriteFile(t *testing.T) {
	err := Handler.AppendWriteFile("./test.txt", "大大大大大大发发生飞洒飞洒发生大2")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestTCP(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	// Example usage
	go func() {
		if err := Handler.ReceiveData("9999"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	data := map[string]interface{}{
		"int":    123,
		"string": "Hello, World!",
		"slice":  []int{1, 2, 3, 4, 5},
		"float":  3.14,
	}

	if err := Handler.SendData("192.168.0.209:9999", data); err != nil {
		log.Fatalf("Client error: %v", err)
	}

	wg.Wait() // Wait for all goroutines to complete
}
