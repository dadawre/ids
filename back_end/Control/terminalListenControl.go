package Control

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// openNewTerminal 创建一个新的终端并运行监听程序
func openNewTerminal() error {
	switch runtime.GOOS {
	case "windows":
		// Windows 使用 cmd 启动新终端
		return exec.Command("cmd", "/C", "start", "cmd", "/K", os.Args[0], "listen").Start()
	case "darwin":
		// macOS 使用 open -a Terminal 启动新终端
		return exec.Command("open", "-a", "Terminal", os.Args[0], "listen").Start()
	case "linux":
		// Linux 使用 xterm 启动新终端
		return exec.Command("xterm", "-e", os.Args[0], "listen").Start()
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

// listenInput 在当前终端中监听用户输入
func listenInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("监听已启动,可以输入规则动态向规则链表添加规则,输入exit退出")

	for {
		fmt.Print("请输入想要插入的规则: ")
		input, err := reader.ReadString('\n') //读取每一行输入
		if err != nil {
			fmt.Println("读取输入时出错:", err)
			continue
		}

		if input == "" {
			continue
		}

		// 移除输入中的换行符
		input = input[:len(input)-1]
		if input == "exit" {
			fmt.Println("监听已退出。")
			break
		}

		if RuleInsert(input) == true {
			fmt.Println("规则插入成功")
		} else {
			fmt.Println("规则插入失败,请重试")
		}
		//fmt.Println("您输入的是:", input)
	}
}

func RulerInputListen() {
	// 判断是否进入监听模式
	if len(os.Args) > 1 && os.Args[1] == "listen" {
		listenInput()
		return
	}

	// 如果不是监听模式，则创建新的终端
	err := openNewTerminal()
	if err != nil {
		fmt.Println("无法创建新终端:", err)
	} else {
		fmt.Println("新终端已创建，用于监听用户输入。")
	}
}
