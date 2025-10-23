package main

import (
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"github.com/wangle201210/gochat/internal/config"
	"github.com/wangle201210/gochat/internal/service/ai"
	"github.com/wangle201210/gochat/internal/ui"
)

func main() {
	// 加载配置
	configPath := config.GetConfigPath()
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 检查 API Key
	if cfg.AI.APIKey == "" {
		fmt.Println("警告: 未配置 API Key")
		fmt.Printf("请编辑配置文件: %s\n", configPath)
		fmt.Println("添加您的 API Key 后重新运行程序")

		// 保存默认配置
		if err := cfg.Save(configPath); err != nil {
			log.Printf("保存配置失败: %v", err)
		} else {
			fmt.Printf("已创建默认配置文件: %s\n", configPath)
		}

		os.Exit(1)
	}

	// 初始化 AI 服务
	aiService, err := ai.NewService(&cfg.AI)
	if err != nil {
		log.Fatalf("初始化 AI 服务失败: %v", err)
	}

	// 创建 Fyne 应用
	fyneApp := app.New()

	// 创建聊天窗口，传入 UI 配置
	chatWindow := ui.NewChatWindow(fyneApp, aiService, &cfg.UI)

	// 显示窗口并运行应用
	chatWindow.Show()
}
