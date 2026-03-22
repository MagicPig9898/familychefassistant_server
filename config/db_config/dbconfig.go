package db_config

import (
	"fmt"
	"sync"

	mysql "github.com/MagicPig9898/easy_db/mysql"
)

// DbConfig 数据库配置
type dbConfig struct {
	MysqlCli *mysql.Client
}

// 定义变量
var (
	initOnce sync.Once
	initErr  error
)

var db = &dbConfig{
	MysqlCli: nil,
	// todo: 添加其他数据库客户端
}

// NewDbConfig 初始化数据库配置
func NewDbConfig() error {
	initOnce.Do(func() {
		fmt.Printf("DB Config init\n")
		db.MysqlCli, initErr = mysql.NewClient(
			"localhost",
			3306,
			"root",
			"123456",
			"lhs",
		)
		// todo: 初始化其他数据库客户端
	})
	return initErr
}

// GetDbConfig 获取数据库客户端
func GetDbConfig() (*dbConfig, error) {
	if initErr != nil {
		return nil, initErr
	}
	if db.MysqlCli == nil {
		return nil, fmt.Errorf("db.MysqlCli not initialized")
	}
	// todo: 检查其他数据库客户端是否初始化
	return db, nil
}

// Close 关闭数据库客户端
func Close() error {
	if db.MysqlCli != nil {
		db.MysqlCli.Close()
	}
	// todo: 关闭其他数据库客户端
	return nil
}
