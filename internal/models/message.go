package models

import (
	"fmt"
	"time"
)

// Role 表示消息角色
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

// Message 表示一条聊天消息
type Message struct {
	ID        string    `json:"id"`
	Role      Role      `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// NewMessage 创建新消息
func NewMessage(role Role, content string) *Message {
	return &Message{
		ID:        generateID(),
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}
}

// generateID 生成消息 ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}
