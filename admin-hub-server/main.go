package main

import (
	appinit "github.com/Light-ink-yht/admin-hub/init"
	"github.com/Light-ink-yht/admin-hub/ioc"
)

func main() {
	// 初始化 viper
	appinit.Viper()

	// 初始化邮件配置表和默认数据
	db := ioc.InitDB()
	ioc.InitEmailTables(db)
	ioc.InitDefaultEmailConfig(db)
	ioc.InitDefaultEmailTemplates(db)

	// 初始化 Web 服务器（通过 wire 自动完成日志初始化）
	server := InitWebServer()

	// 启动 Web 服务器，监听 8080 端口
	err := server.Run(":8081")
	if err != nil {
		// 如果启动失败，打印错误信息并终止程序
		panic(err)
	}
}
