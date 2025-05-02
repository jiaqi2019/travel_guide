package main

import (
	"fmt"
	"log"
	"travel_guide/config"
	"travel_guide/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// 初始化数据库
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 初始化oss
	if err := config.InitOSS(); err != nil {
		log.Fatal("Failed to initialize OSS client:", err)
	}

	// 初始化路由
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r, db)

	//启动server
	addr := fmt.Sprintf(":%d", config.AppConfig.ServerConfig.Port)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}