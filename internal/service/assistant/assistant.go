package assistant

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/wangle201210/gochat/internal/config"
	"github.com/wangle201210/gochat/internal/models"
)

// Service 助手服务（用于生成标题等辅助任务）
type Service struct {
	chatModel model.ChatModel
	config    *config.AssistantConfig
}

// NewService 创建助手服务
func NewService(cfg *config.AssistantConfig) (*Service, error) {
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
		return nil, fmt.Errorf("不支持的助手模型 provider: %s", cfg.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("初始化助手模型失败: %w", err)
	}

	return &Service{
		chatModel: chatModel,
		config:    cfg,
	}, nil
}

// GenerateTitle 根据最近的消息生成会话标题
// messages: 最近的消息列表（建议传入最近4组对话）
func (s *Service) GenerateTitle(ctx context.Context, messages []*models.Message) (string, error) {
	if len(messages) == 0 {
		return "新会话", nil
	}

	// 构建 prompt
	systemPrompt := "你是一个会话标题生成助手。根据用户与AI的对话内容，生成一个简洁、准确的会话标题。标题应该在10个字以内，能够概括对话的主题。只输出标题，不要有其他内容。"

	// 构建对话上下文
	conversationText := ""
	for _, msg := range messages {
		if msg.Role == models.RoleUser {
			conversationText += fmt.Sprintf("用户: %s\n", msg.Content)
		} else if msg.Role == models.RoleAssistant {
			conversationText += fmt.Sprintf("助手: %s\n", msg.Content)
		}
	}

	userPrompt := fmt.Sprintf("请根据以下对话内容生成一个简洁的标题:\n\n%s", conversationText)

	// 构建消息列表
	schemaMessages := []*schema.Message{
		{
			Role:    schema.System,
			Content: systemPrompt,
		},
		{
			Role:    schema.User,
			Content: userPrompt,
		},
	}

	// 调用模型
	resp, err := s.chatModel.Generate(ctx, schemaMessages)
	if err != nil {
		return "", fmt.Errorf("生成标题失败: %w", err)
	}

	title := resp.Content
	if title == "" {
		title = "新会话"
	}

	return title, nil
}
