package main

import (
	appinit "github.com/Light-ink-yht/admin-hub/init"
	"github.com/spf13/viper"
)

func main() {
	// 初始化 viper
	appinit.Viper()

	// 初始化 Web 服务器（通过 wire 自动完成日志初始化）
	server := InitWebServer()

	// 启动 Web 服务器
	port := viper.GetString("server.port")
	err := server.Run(port)
	if err != nil {
		// 如果启动失败，打印错误信息并终止程序
		panic(err)
	}
}
