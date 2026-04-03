package main

import (
	"fmt"

	conf "github.com/MagicPig9898/familychefassistant_server/conf"
	dbconfig "github.com/MagicPig9898/familychefassistant_server/config/db_config"
	log "github.com/MagicPig9898/familychefassistant_server/config/log_config"
	router_config "github.com/MagicPig9898/familychefassistant_server/router"
)

func main() {
	// 1. 加载配置文件（失败直接终止）
	err := conf.MustLoad("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 2. 初始化日志
	if err := log.NewLogConfig(); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Close()

	// 3. 初始化数据库连接池
	if err := dbconfig.NewDbConfig(); err != nil {
		log.Fatalf("Failed to initialize database config: %v", err)
	}
	defer dbconfig.Close()

	// 4. 注册路由
	r := router_config.NewRouter()

	// 5. 启动服务
	port := fmt.Sprintf(":%d", conf.Cfg.Server.Port)
	log.Infof("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
