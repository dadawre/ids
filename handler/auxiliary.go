package handler

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
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

func CreateFile(dirPath string) error {
	// 创建文件夹
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func AppendWriteFile(path, content string) error {
	// 打开或创建文件
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("关闭文件时发生错误: %v", err)
		}
	}(file)

	// 创建一个带缓冲的写入器
	writer := bufio.NewWriter(file)

	// 写入数据到文件
	_, err = writer.WriteString(content)
	if err != nil {
		return fmt.Errorf("无法写入数据: %w", err)
	}

	// 刷新缓冲区
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("无法刷新缓冲区: %w", err)
	}

	return nil
}

// SendData sends data to a specified address
func SendData(address string, data map[string]interface{}) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	log.Println("Data sent successfully")
	return nil
}

// ReceiveData listens on the specified port and handles incoming connections
func ReceiveData(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port: %w", err)
	}
	defer listener.Close()

	log.Println("Server started, waiting for connections...")

	for true {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
	return nil
}

// handleConnection handles a single connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected")

	decoder := gob.NewDecoder(conn)
	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		log.Printf("Failed to receive data: %v", err)
		return
	}

	fmt.Printf("Received data: %v\n", data)
}
