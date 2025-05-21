/*
 * @author: yihang_01
 * @Date: 2025-05-21 16:25:10
 * @LastEditTime: 2025-05-21 19:50:34
 * QwQ 加油加油
 */
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/10yihang/resume-ai-interview/api/handlers"
	"github.com/10yihang/resume-ai-interview/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	Version = "1.0.0"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// 初始化配置
	cfg := config.NewConfig()

	// 打印版本和AI提供商信息
	fmt.Printf("AI简历面试助手 v%s\n", Version)
	if cfg.UseGrok {
		fmt.Println("AI提供商: Grok 3")
	} else if cfg.APIKey != "" {
		fmt.Println("AI提供商: OpenAI")
	} else {
		fmt.Println("AI提供商: 模拟模式 (未配置API密钥)")
	}
	// 创建Gin引擎
	r := gin.Default()

	// 加载静态文件
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// 初始化handlers，传入配置
	handlers.InitHandlers(cfg)

	// 设置路由
	r.GET("/", handlers.IndexHandler)
	r.POST("/upload/resume", handlers.UploadResumeHandler)
	r.POST("/upload/jd", handlers.UploadJDHandler)
	r.POST("/generate/questions", handlers.GenerateQuestionsHandler)
	r.POST("/evaluate/answer", handlers.EvaluateAnswerHandler)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
