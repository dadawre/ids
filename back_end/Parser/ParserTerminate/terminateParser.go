package ParserTerminate

import (
	"fmt"
	"os"
	"strings"
)

func TerminateParameter() {
	var logPath string
	var configPath string

	// 手动解析命令行参数
	args := os.Args[1:]
	for index, arg := range args {
		fmt.Println(index, "  ", arg)
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-l", "--log":
			if i+1 < len(args) {
				logPath = args[i+1]
				i++
			} else {
				fmt.Println("Error: missing value for", arg)
				return
			}
		case "-c", "--config":
			if i+1 < len(args) {
				configPath = args[i+1]
				i++
			} else {
				fmt.Println("Error: missing value for", arg)
				return
			}
		default:
			if strings.HasPrefix(arg, "--log=") {
				logPath = strings.SplitN(arg, "=", 2)[1]
			} else if strings.HasPrefix(arg, "--config=") {
				configPath = strings.SplitN(arg, "=", 2)[1]
			} else {
				fmt.Println("Unknown argument:", arg)
				return
			}
		}
	}

	// 输出解析后的参数
	fmt.Printf("Log file path: %s\n", logPath)
	fmt.Printf("Config file path: %s\n", configPath)
}
