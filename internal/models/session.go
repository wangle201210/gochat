package models

import "time"

// Session 表示一个会话
type Session struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewSession 创建新会话
func NewSession() *Session {
	now := time.Now()
	return &Session{
		ID:        generateID(),
		Title:     "新会话",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
