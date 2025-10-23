package ui

import (
	"context"
	"fmt"
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

	// 立即清空输入框（不阻塞）
	cw.inputEntry.SetText("")

	// 禁用发送按钮，防止重复发送
	cw.sendButton.Disable()

	// 立即添加用户消息到界面（不阻塞）
	userMsg := models.NewMessage(models.RoleUser, userInput)
	cw.addMessage(userMsg)

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
			}

			// 完成后滚动到底部并重新启用发送按钮
			cw.scrollToBottom()
			cw.sendButton.Enable()
		})
	}()
}

// handleClear 处理清空历史
func (cw *ChatWindow) handleClear() {
	dialog.ShowConfirm("确认清空", "确定要清空所有聊天历史吗？", func(ok bool) {
		if ok {
			cw.messages = make([]*models.Message, 0)
			cw.aiService.ClearHistory()
			cw.messageContainer.Objects = []fyne.CanvasObject{}
			cw.messageContainer.Refresh()
		}
	}, cw.window)
}
