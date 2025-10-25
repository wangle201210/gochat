package ui

import (
	"context"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/wangle201210/gochat/internal/models"
)

// handleSend 处理发送消息
func (cw *ChatWindow) handleSend() {
	userInput := strings.TrimSpace(cw.inputEntry.Text)
	if userInput == "" {
		return
	}

	// 确保有当前会话
	if cw.currentSession == nil {
		cw.createNewSession()
	}

	// 立即清空输入框（不阻塞）
	cw.inputEntry.SetText("")

	// 禁用发送按钮，防止重复发送
	cw.sendButton.Disable()

	// 立即添加用户消息到界面（不阻塞）
	userMsg := models.NewMessage(models.RoleUser, userInput)
	cw.addMessage(userMsg)

	// 保存用户消息到数据库
	if err := cw.db.SaveMessage(cw.currentSession.ID, userMsg); err != nil {
		dialog.ShowError(err, cw.window)
	}

	// 创建一个占位消息用于流式更新
	assistantMsg := models.NewMessage(models.RoleAssistant, "正在思考...")
	assistantRichText := cw.addMessage(assistantMsg)
	assistantIndex := len(cw.messages) - 1

	// 异步获取 AI 回复（不阻塞 UI）
	go func() {
		ctx := context.Background()
		var fullResponse strings.Builder

		err := cw.aiService.StreamChat(ctx, userInput, func(chunk string) error {
			fullResponse.WriteString(chunk)
			currentContent := fullResponse.String()

			// 在主线程中更新 UI - 使用 Fyne 提供的线程安全方法
			fyne.Do(func() {
				cw.messages[assistantIndex].Content = currentContent
				// 更新 RichText 的 Markdown 内容
				assistantRichText.ParseMarkdown(currentContent)
				cw.scrollToBottom()
			})

			return nil
		})

		// 在主线程中处理错误和完成操作
		fyne.Do(func() {
			if err != nil {
				errMsg := fmt.Sprintf("错误: %v", err)
				cw.messages[assistantIndex].Content = errMsg
				assistantRichText.ParseMarkdown(errMsg)
				dialog.ShowError(err, cw.window)
			} else {
				// 保存 AI 回复到数据库
				if err := cw.db.SaveMessage(cw.currentSession.ID, cw.messages[assistantIndex]); err != nil {
					dialog.ShowError(err, cw.window)
				}

				// 异步生成/更新会话标题
				go cw.generateSessionTitle()
			}

			// 完成后滚动到底部并重新启用发送按钮
			cw.scrollToBottom()
			cw.sendButton.Enable()

			// 刷新会话列表以更新时间戳
			cw.refreshSessionList()
		})
	}()
}

// generateSessionTitle 生成会话标题
func (cw *ChatWindow) generateSessionTitle() {
	if cw.currentSession == nil || cw.assistantService == nil {
		return
	}

	// 获取最近4组对话（最多8条消息）
	recentMessages := cw.messages
	if len(recentMessages) > 8 {
		recentMessages = recentMessages[len(recentMessages)-8:]
	}

	// 如果消息太少，不生成标题
	if len(recentMessages) < 2 {
		return
	}

	// 调用助手服务生成标题
	ctx := context.Background()
	title, err := cw.assistantService.GenerateTitle(ctx, recentMessages)
	if err != nil {
		log.Printf("生成会话标题失败: %v", err)
		return
	}

	// 更新数据库中的标题
	if err := cw.db.UpdateSessionTitle(cw.currentSession.ID, title); err != nil {
		log.Printf("更新会话标题失败: %v", err)
		return
	}

	// 在主线程更新界面
	fyne.Do(func() {
		cw.currentSession.Title = title
		cw.refreshSessionList()
	})
}
