package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/wangle201210/gochat/internal/config"
	"github.com/wangle201210/gochat/internal/models"
	"github.com/wangle201210/gochat/internal/service/ai"
)

// ChatWindow 聊天窗口
type ChatWindow struct {
	window           fyne.Window
	app              fyne.App
	aiService        *ai.Service
	uiConfig         *config.UIConfig
	messageContainer *fyne.Container
	scrollContainer  *container.Scroll
	inputEntry       *customEntry
	sendButton       *widget.Button
	clearButton      *widget.Button
	messages         []*models.Message
}

// NewChatWindow 创建聊天窗口
func NewChatWindow(app fyne.App, aiService *ai.Service, uiConfig *config.UIConfig) *ChatWindow {
	window := app.NewWindow("GoChat - AI 对话助手")

	cw := &ChatWindow{
		window:    window,
		app:       app,
		aiService: aiService,
		uiConfig:  uiConfig,
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

	// 使用配置中的窗口尺寸
	windowWidth := cw.uiConfig.WindowWidth
	windowHeight := cw.uiConfig.WindowHeight
	if windowWidth <= 0 {
		windowWidth = 900 // 默认宽度
	}
	if windowHeight <= 0 {
		windowHeight = 700 // 默认高度
	}
	cw.window.Resize(fyne.NewSize(float32(windowWidth), float32(windowHeight)))
}

// Show 显示窗口
func (cw *ChatWindow) Show() {
	cw.window.ShowAndRun()
}
