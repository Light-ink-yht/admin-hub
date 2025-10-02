package main

import (
	"github.com/Light-ink-yht/admin-hub/init"
)

func main() {
	// 初始化 viper
	init.Viper()

	// 初始化 logger
	init.Logger()

	// 初始化 Web 服务器
	server := InitWebServer()

	// 启动 Web 服务器，监听 8080 端口
	err := server.Run(":8081")
	if err != nil {
		// 如果启动失败，打印错误信息并终止程序
		panic(err)
	}
}
