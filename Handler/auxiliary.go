package Handler

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
)

func CreateFile(dirPath string) error { //创建文件
	// 创建文件夹
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// AppendWriteFile 在指定文件中以append的方式写入字符串
func AppendWriteFile(path, content string) error { //写入文件
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

// SendData 发送数据到特定IP和端口
func SendData(address string, data map[string]interface{}) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	log.Println("Data sent successfully")
	return nil
}

// ReceiveData 监听端口
func ReceiveData(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port: %w", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(listener)

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
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	fmt.Println("Client connected")

	decoder := gob.NewDecoder(conn)
	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		log.Printf("Failed to receive data: %v", err)
		return
	}

	fmt.Printf("Received data: %v\n", data)
}
