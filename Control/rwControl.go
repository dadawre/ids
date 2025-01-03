package Control

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// ReadConfig 读配置文件
func ReadConfig(configPath string) {
	// 打开文件
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 创建一个新的Reader对象
	reader := bufio.NewReader(file)

	readString, err := reader.ReadString('\n')
	print(readString)

	if err != nil {
		return
	}

	// 如果遇到文件末尾，退出循环
	if err != nil && err.Error() == "EOF" {
		fmt.Println(err)
	}

	// 如果读取过程中出现错误，报错并退出
	if err != nil {
		log.Fatalf("读取文件出错: %v", err)
	}

}

// WriteToFile 写文件
func WriteToFile(filename string, content string) error {
	// 将字符串内容转换为字节切片
	data := []byte(content)

	// 使用 ioutil.WriteFile 写入文件
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("文件 %s 写入成功\n", filename)
	return nil
}
