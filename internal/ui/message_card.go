package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/wangle201210/gochat/internal/models"
)

// addMessage 添加消息到列表，返回消息内容的 RichText 引用
func (cw *ChatWindow) addMessage(msg *models.Message) *widget.RichText {
	cw.messages = append(cw.messages, msg)

	var richText *widget.RichText
	var messageCard *fyne.Container

	switch msg.Role {
	case models.RoleUser:
		// 用户消息 - 淡蓝白卡片（小清新风格）
		roleLabel := widget.NewLabel("※ 我")
		roleLabel.TextStyle = fyne.TextStyle{Bold: true}

		// 用户消息使用 Label 保留换行符
		contentLabel := widget.NewLabel(msg.Content)
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
		// AI 消息 - 温暖米白卡片（小清新风格）
		roleLabel := widget.NewLabel("✨ 助手")
		roleLabel.TextStyle = fyne.TextStyle{Bold: true}

		// AI 消息使用 RichText 渲染 Markdown
		richText = widget.NewRichTextFromMarkdown(msg.Content)
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

		contentLabel := widget.NewLabel(msg.Content)
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
