package main

import (
	"fmt"
	"log"

	conf "github.com/MagicPig9898/familychefassistant_server/conf"
	dbconfig "github.com/MagicPig9898/familychefassistant_server/config/db_config"
	router_config "github.com/MagicPig9898/familychefassistant_server/router"
)

func main() {
	// 1. 加载配置文件（失败直接终止）
	err := conf.MustLoad("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 3. 初始化数据库连接池
	if err := dbconfig.NewDbConfig(); err != nil {
		log.Fatalf("Failed to initialize database config: %v", err)
	}
	defer func() {
		dbconfig.Close()
	}()

	// 4. 注册路由
	r := router_config.NewRouter()

	// 5. 启动服务
	port := fmt.Sprintf(":%d", conf.Cfg.Server.Port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
