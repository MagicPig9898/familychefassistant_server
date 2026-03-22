package main

import (
	"log"

	dbconfig "github.com/MagicPig9898/familychefassistant_server/config/db_config"
	"github.com/MagicPig9898/familychefassistant_server/config/router_config"
)

func main() {
	// 1. 初始化数据库配置
	if err := dbconfig.NewDbConfig(); err != nil {
		log.Fatalf("Failed to initialize database config: %v", err)
	}
	defer func() {
		_ = dbconfig.Close()
	}()

	// 注册路由
	r := router_config.NewRouter()

	// 4. 启动服务
	port := ":8080"
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
