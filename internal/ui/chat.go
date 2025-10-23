package ui

import (
	"context"
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/wangle201210/gochat/internal/models"
	"github.com/wangle201210/gochat/internal/service/ai"
)

var (
	// 小清新颜色方案
	userMessageBg   = color.NRGBA{R: 240, G: 248, B: 255, A: 255} // 淡蓝白 - 柔和不明显
	assistantBg     = color.NRGBA{R: 255, G: 253, B: 245, A: 255} // 温暖米白
	backgroundColor = color.NRGBA{R: 250, G: 252, B: 252, A: 255} // 清新白
	inputAreaBg     = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // 纯白
	separatorColor  = color.NRGBA{R: 230, G: 240, B: 235, A: 255} // 淡绿分隔线
)

// ChatWindow 聊天窗口
type ChatWindow struct {
	window           fyne.Window
	app              fyne.App
	aiService        *ai.Service
	messageContainer *fyne.Container
	scrollContainer  *container.Scroll
	inputEntry       *customEntry
	sendButton       *widget.Button
	clearButton      *widget.Button
	messages         []*models.Message
}

// customEntry 自定义输入框，支持 Enter 发送
type customEntry struct {
	widget.Entry
	onEnter func()
}

// newCustomEntry 创建自定义输入框
func newCustomEntry(onEnter func()) *customEntry {
	entry := &customEntry{onEnter: onEnter}
	entry.MultiLine = true
	entry.Wrapping = fyne.TextWrapWord
	entry.ExtendBaseWidget(entry)
	return entry
}

// TypedKey 处理键盘按键
func (e *customEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn, fyne.KeyEnter:
		// Enter 键发送消息
		if e.onEnter != nil {
			e.onEnter()
		}
	default:
		// 其他键使用默认处理
		e.Entry.TypedKey(key)
	}
}

// TypedShortcut 处理快捷键
func (e *customEntry) TypedShortcut(shortcut fyne.Shortcut) {
	// Shift+Enter 插入换行
	if _, ok := shortcut.(*desktop.CustomShortcut); ok {
		e.TypedRune('\n')
		return
	}
	e.Entry.TypedShortcut(shortcut)
}

// NewChatWindow 创建聊天窗口
func NewChatWindow(app fyne.App, aiService *ai.Service) *ChatWindow {
	window := app.NewWindow("GoChat - AI 对话助手")

	cw := &ChatWindow{
		window:    window,
		app:       app,
		aiService: aiService,
		messages:  make([]*models.Message, 0),
	}

	cw.setupUI()
	return cw
}

// setupUI 设置 UI 组件
func (cw *ChatWindow) setupUI() {
	// 消息容器 - 使用 VBox 允许动态高度
	cw.messageContainer = container.NewVBox()

	// 创建带背景的消息区域
	messageAreaBg := canvas.NewRectangle(backgroundColor)
	messagesWithBg := container.NewStack(messageAreaBg, cw.messageContainer)

	// 滚动容器
	cw.scrollContainer = container.NewScroll(messagesWithBg)

	// 创建自定义输入框（Enter 发送）
	cw.inputEntry = newCustomEntry(cw.handleSend)
	cw.inputEntry.SetPlaceHolder("输入消息... (Enter 发送, Shift+Enter 换行)")
	cw.inputEntry.SetMinRowsVisible(3)

	// 发送按钮 - 使用重要样式
	cw.sendButton = widget.NewButton("发送消息", cw.handleSend)
	cw.sendButton.Importance = widget.HighImportance

	// 清空按钮
	cw.clearButton = widget.NewButton("清空对话", cw.handleClear)

	// 底部按钮栏 - 添加间距
	buttonBar := container.NewHBox(
		cw.clearButton,
		layout.NewSpacer(),
		cw.sendButton,
	)

	// 输入区域容器 - 添加内边距
	inputCard := container.NewVBox(
		widget.NewSeparator(),
		container.NewPadded(cw.inputEntry),
		container.NewPadded(buttonBar),
	)

	// 主布局
	mainContent := container.NewBorder(
		nil,
		inputCard,
		nil,
		nil,
		cw.scrollContainer,
	)

	cw.window.SetContent(mainContent)
	cw.window.Resize(fyne.NewSize(900, 700))
}

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

// addMessage 添加消息到列表，返回消息内容的 RichText 引用
func (cw *ChatWindow) addMessage(msg *models.Message) *widget.RichText {
	cw.messages = append(cw.messages, msg)

	var richText *widget.RichText
	var messageCard *fyne.Container

	switch msg.Role {
	case models.RoleUser:
		// 用户消息 - 薄荷绿卡片（小清新风格）
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

// Show 显示窗口
func (cw *ChatWindow) Show() {
	cw.window.ShowAndRun()
}
