// Package console 命令行辅助方法
package console

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"api/pkg/logger"
)

// Success 打印一条成功消息，绿色输出
func Success(msg string) {
	color.NoColor = false
	_, err := fmt.Fprintln(color.Output, color.GreenString(msg))
	logger.LogIf(err)
}

// Error 打印一条报错消息，红色输出
func Error(msg string) {
	color.NoColor = false
	_, err := fmt.Fprintln(color.Output, color.RedString(msg))
	logger.LogIf(err)
}

// Warning 打印一条提示消息，黄色输出
func Warning(msg string) {
	color.NoColor = false
	_, err := fmt.Fprintln(color.Output, color.YellowString(msg))
	logger.LogIf(err)
}

// Exit 打印一条报错消息，并退出 os.Exit(1)
func Exit(msg string) {
	Error(msg)
	os.Exit(1)
}

// ExitIf 语法糖，自带 err != nil 判断
func ExitIf(err error) {
	if err != nil {
		Exit(err.Error())
	}
}
