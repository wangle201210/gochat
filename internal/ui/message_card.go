package ui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/wangle201210/gochat/internal/models"
)

// normalizeEmoji 规范化 emoji，替换带变体选择器的 emoji 为兼容版本
func normalizeEmoji(content string) string {
	// 替换数字 emoji (0️⃣-9️⃣) 为普通数字加圆圈
	replacements := map[string]string{
		"0️⃣": "⓪", "1️⃣": "①", "2️⃣": "②", "3️⃣": "③", "4️⃣": "④",
		"5️⃣": "⑤", "6️⃣": "⑥", "7️⃣": "⑦", "8️⃣": "⑧", "9️⃣": "⑨",
		"🔟": "⑩",
	}

	result := content
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

// addMessage 添加消息到列表，返回消息内容的 RichText 引用
func (cw *ChatWindow) addMessage(msg *models.Message) *widget.RichText {
	cw.messages = append(cw.messages, msg)

	// 规范化消息内容中的 emoji
	displayContent := normalizeEmoji(msg.Content)

	var richText *widget.RichText
	var messageCard fyne.CanvasObject

	switch msg.Role {
	case models.RoleUser:
		// 用户消息
		roleLabel := widget.NewLabel("※ 我")
		roleLabel.TextStyle = fyne.TextStyle{Bold: true}

		// 用户消息使用 Label 保留换行符
		contentLabel := widget.NewLabel(displayContent)
		contentLabel.Wrapping = fyne.TextWrapWord

		// 创建内容容器
		contentBox := container.NewVBox(
			roleLabel,
			contentLabel,
		)

		// 创建带柔和边距的背景
		bg := canvas.NewRectangle(userMessageBg)

		// 使用适度的内边距
		cardContent := container.NewPadded(contentBox)
		messageCard = container.NewStack(bg, cardContent)

	case models.RoleAssistant:
		// AI 消息
		roleLabel := widget.NewLabel("✨ 助手")
		roleLabel.TextStyle = fyne.TextStyle{Bold: true}

		// AI 消息使用 RichText 渲染 Markdown
		richText = widget.NewRichTextFromMarkdown(displayContent)
		richText.Wrapping = fyne.TextWrapWord

		// 创建内容容器
		contentBox := container.NewVBox(
			roleLabel,
			richText,
		)

		// 创建带柔和边距的背景
		bg := canvas.NewRectangle(assistantBg)

		// 使用适度的内边距
		cardContent := container.NewPadded(contentBox)
		messageCard = container.NewStack(bg, cardContent)

	case models.RoleSystem:
		// 系统消息 - 简单样式
		roleLabel := widget.NewLabel("⚙️ 系统")
		roleLabel.TextStyle = fyne.TextStyle{Bold: true, Italic: true}

		contentLabel := widget.NewLabel(displayContent)
		contentLabel.Wrapping = fyne.TextWrapWord

		contentBox := container.NewVBox(
			roleLabel,
			contentLabel,
		)
		messageCard = container.NewPadded(contentBox)
	}

	// 添加更大的间距，营造清爽感
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(1, 12)) // 12 像素间距

	spacedCard := container.NewVBox(
		messageCard,
		spacer,
	)

	// 添加到消息容器，左右添加边距
	cw.messageContainer.Add(
		container.NewPadded(spacedCard),
	)
	cw.scrollToBottom()

	return richText
}

// scrollToBottom 滚动到底部
func (cw *ChatWindow) scrollToBottom() {
	cw.scrollContainer.ScrollToBottom()
}
