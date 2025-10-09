package ioc

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// InitDB 初始化数据库连接
func InitDB() *sql.DB {

	type Config struct {
		Dsn string `yaml:"dsn"`
	}
	var cfg Config = Config{
		Dsn: "root:@tcp(127.0.0.1:3306)/admin_hub",
	}
	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}

	// 打开数据库连接
	db, err := sql.Open("mysql", cfg.Dsn)
	if err != nil {
		// 如果连接失败，抛出错误
		panic(err)
	}

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		// 如果连接失败，抛出错误
		panic(err)
	}

	// 打印连接成功信息
	fmt.Println("Successfully connected to the database!")
	return db
}
