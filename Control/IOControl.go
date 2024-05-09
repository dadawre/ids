package Control

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// ReadConfig 读配置文件
func ReadConfig(configPath string) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 创建一个 Scanner 来逐行读取文件内容
	scanner := bufio.NewScanner(file)

	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	// 检查是否有错误发生
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// 写文件
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
