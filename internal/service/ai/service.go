package ai

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/wangle201210/gochat/internal/config"
	"github.com/wangle201210/gochat/internal/models"
)

// Service AI 服务
type Service struct {
	chatModel model.ChatModel
	config    *config.AIConfig
	history   []*models.Message
}

// NewService 创建 AI 服务
func NewService(cfg *config.AIConfig) (*Service, error) {
	var chatModel model.ChatModel
	var err error

	switch cfg.Provider {
	case "openai":
		chatModel, err = openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
			BaseURL: cfg.BaseURL,
			Model:   cfg.Model,
			APIKey:  cfg.APIKey,
		})
	default:
		return nil, fmt.Errorf("不支持的 AI provider: %s", cfg.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("初始化 AI 模型失败: %w", err)
	}

	return &Service{
		chatModel: chatModel,
		config:    cfg,
		history:   make([]*models.Message, 0),
	}, nil
}

// Chat 发送消息并获取回复
func (s *Service) Chat(ctx context.Context, userMessage string) (string, error) {
	// 添加用户消息到历史
	userMsg := models.NewMessage(models.RoleUser, userMessage)
	s.history = append(s.history, userMsg)

	// 转换消息历史为 Eino 格式
	messages := s.convertMessages()

	// 调用 AI 模型
	resp, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("AI 生成失败: %w", err)
	}

	// 获取回复内容
	assistantContent := resp.Content

	// 添加助手消息到历史
	assistantMsg := models.NewMessage(models.RoleAssistant, assistantContent)
	s.history = append(s.history, assistantMsg)

	return assistantContent, nil
}

// StreamChat 流式发送消息并获取回复
func (s *Service) StreamChat(ctx context.Context, userMessage string, callback func(string) error) error {
	// 添加用户消息到历史
	userMsg := models.NewMessage(models.RoleUser, userMessage)
	s.history = append(s.history, userMsg)

	// 转换消息历史为 Eino 格式
	messages := s.convertMessages()

	// 调用流式 AI 模型
	streamReader, err := s.chatModel.Stream(ctx, messages)
	if err != nil {
		return fmt.Errorf("AI 流式生成失败: %w", err)
	}

	var fullContent string

	// 读取流式响应
	for {
		chunk, err := streamReader.Recv()
		if err != nil {
			// 流结束
			break
		}

		content := chunk.Content
		fullContent += content

		// 回调处理每个流式块
		if callback != nil {
			if err := callback(content); err != nil {
				return err
			}
		}
	}

	// 添加完整的助手消息到历史
	assistantMsg := models.NewMessage(models.RoleAssistant, fullContent)
	s.history = append(s.history, assistantMsg)

	return nil
}

// GetHistory 获取消息历史
func (s *Service) GetHistory() []*models.Message {
	return s.history
}

// ClearHistory 清空消息历史
func (s *Service) ClearHistory() {
	s.history = make([]*models.Message, 0)
}

// convertMessages 将内部消息格式转换为 Eino 格式
func (s *Service) convertMessages() []*schema.Message {
	messages := make([]*schema.Message, 0, len(s.history))

	for _, msg := range s.history {
		var role schema.RoleType
		switch msg.Role {
		case models.RoleUser:
			role = schema.User
		case models.RoleAssistant:
			role = schema.Assistant
		case models.RoleSystem:
			role = schema.System
		default:
			role = schema.User
		}

		messages = append(messages, &schema.Message{
			Role:    role,
			Content: msg.Content,
		})
	}

	return messages
}
