package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/wangle201210/gochat/internal/config"
	"github.com/wangle201210/gochat/internal/models"
	"github.com/wangle201210/gochat/internal/service/ai"
	"github.com/wangle201210/gochat/internal/service/assistant"
	"github.com/wangle201210/gochat/internal/storage"
)

// ChatWindow 聊天窗口
type ChatWindow struct {
	window               fyne.Window
	app                  fyne.App
	aiService            *ai.Service
	assistantService     *assistant.Service
	uiConfig             *config.UIConfig
	db                   *storage.Database
	messageContainer     *fyne.Container
	scrollContainer      *container.Scroll
	inputEntry           *customEntry
	sendButton           *widget.Button
	messages             []*models.Message
	currentSession       *models.Session
	sessionList          *SessionList
	sessionListContainer *fyne.Container
	toggleButton         *widget.Button
	mainContent          *fyne.Container
	sessionListVisible   bool
}

// NewChatWindow 创建聊天窗口
func NewChatWindow(app fyne.App, aiService *ai.Service, assistantService *assistant.Service, uiConfig *config.UIConfig, db *storage.Database) *ChatWindow {
	window := app.NewWindow("GoChat - AI 对话助手")

	// 应用自定义主题
	app.Settings().SetTheme(newCustomTheme())

	cw := &ChatWindow{
		window:             window,
		app:                app,
		aiService:          aiService,
		assistantService:   assistantService,
		uiConfig:           uiConfig,
		db:                 db,
		messages:           make([]*models.Message, 0),
		sessionListVisible: true, // 默认显示会话列表
	}

	cw.setupUI()
	cw.initializeSession()
	return cw
}

// initializeSession 初始化会话
func (cw *ChatWindow) initializeSession() {
	// 尝试加载最近的会话
	sessions, err := cw.db.ListSessions()
	if err != nil {
		log.Printf("加载会话列表失败: %v", err)
	}

	if len(sessions) > 0 {
		// 加载最近的会话
		cw.loadSession(sessions[0])
	} else {
		// 创建新会话
		cw.createNewSession()
	}

	// 刷新会话列表
	cw.refreshSessionList()
}

// setupUI 设置 UI 组件
func (cw *ChatWindow) setupUI() {
	// 消息容器
	cw.messageContainer = container.NewVBox()

	// 创建带背景的消息区域
	messageAreaBg := canvas.NewRectangle(backgroundColor)
	messagesWithBg := container.NewStack(messageAreaBg, cw.messageContainer)

	// 滚动容器
	cw.scrollContainer = container.NewScroll(messagesWithBg)

	// 创建自定义输入框
	cw.inputEntry = newCustomEntry(cw.handleSend)
	cw.inputEntry.SetPlaceHolder("输入消息... (Enter 发送, Shift+Enter 换行)")
	cw.inputEntry.SetMinRowsVisible(3)

	// 发送按钮
	cw.sendButton = widget.NewButton("发送消息", cw.handleSend)
	cw.sendButton.Importance = widget.HighImportance

	// 创建切换按钮
	cw.toggleButton = widget.NewButton("☰", cw.toggleSessionList)
	cw.toggleButton.Importance = widget.LowImportance

	// 底部按钮栏
	buttonBar := container.NewHBox(
		cw.toggleButton,
		layout.NewSpacer(),
		cw.sendButton,
	)

	// 输入区域容器
	inputCard := container.NewVBox(
		widget.NewSeparator(),
		container.NewPadded(cw.inputEntry),
		container.NewPadded(buttonBar),
	)

	// 创建会话列表
	cw.sessionList = NewSessionList(
		cw.onSessionSelect,
		cw.onNewSession,
		cw.onDeleteSession,
	)

	// 会话列表区域
	cw.sessionListContainer = container.NewBorder(
		nil, nil, nil, nil,
		cw.sessionList,
	)

	// 使用固定宽度容器包装会话列表
	sessionListFixed := newFixedWidthContainer(200, cw.sessionListContainer)

	// 主聊天区域
	chatArea := container.NewBorder(
		nil,
		inputCard,
		nil,
		nil,
		cw.scrollContainer,
	)

	// 主布局 - 左侧会话列表，右侧聊天区域
	cw.mainContent = container.NewBorder(
		nil, nil,
		sessionListFixed,
		nil,
		chatArea,
	)

	cw.window.SetContent(cw.mainContent)

	// 使用配置中的窗口尺寸
	windowWidth := cw.uiConfig.WindowWidth
	windowHeight := cw.uiConfig.WindowHeight
	if windowWidth <= 0 {
		windowWidth = 1000
	}
	if windowHeight <= 0 {
		windowHeight = 700
	}
	cw.window.Resize(fyne.NewSize(float32(windowWidth), float32(windowHeight)))
}

// toggleSessionList 切换会话列表的显示/隐藏
func (cw *ChatWindow) toggleSessionList() {
	cw.sessionListVisible = !cw.sessionListVisible

	// 创建聊天区域
	chatArea := container.NewBorder(
		nil,
		container.NewVBox(
			widget.NewSeparator(),
			container.NewPadded(cw.inputEntry),
			container.NewPadded(container.NewHBox(
				cw.toggleButton,
				layout.NewSpacer(),
				cw.sendButton,
			)),
		),
		nil,
		nil,
		cw.scrollContainer,
	)

	if cw.sessionListVisible {
		// 显示会话列表
		cw.toggleButton.SetText("☰")
		sessionListFixed := newFixedWidthContainer(200, cw.sessionListContainer)
		cw.mainContent = container.NewBorder(nil, nil, sessionListFixed, nil, chatArea)
	} else {
		// 隐藏会话列表
		cw.toggleButton.SetText("→")
		cw.mainContent = chatArea
	}

	cw.window.SetContent(cw.mainContent)
	cw.window.Canvas().Refresh(cw.mainContent)
}

// Show 显示窗口
func (cw *ChatWindow) Show() {
	cw.window.ShowAndRun()
}

// createNewSession 创建新会话
func (cw *ChatWindow) createNewSession() {
	// 保存当前会话的消息
	if cw.currentSession != nil {
		cw.saveCurrentMessages()
	}

	// 创建新会话
	newSession := models.NewSession()
	if err := cw.db.SaveSession(newSession); err != nil {
		log.Printf("保存新会话失败: %v", err)
		dialog.ShowError(err, cw.window)
		return
	}

	// 清空当前消息和 AI 历史
	cw.messages = make([]*models.Message, 0)
	cw.aiService.ClearHistory()
	cw.messageContainer.Objects = []fyne.CanvasObject{}
	cw.messageContainer.Refresh()

	// 设置当前会话
	cw.currentSession = newSession
	cw.sessionList.SetCurrentSession(newSession)
	cw.refreshSessionList()
}

// loadSession 加载会话
func (cw *ChatWindow) loadSession(session *models.Session) {
	if session == nil {
		return
	}

	// 保存当前会话的消息
	if cw.currentSession != nil && cw.currentSession.ID != session.ID {
		cw.saveCurrentMessages()
	}

	// 加载会话消息
	messages, err := cw.db.GetMessages(session.ID)
	if err != nil {
		log.Printf("加载会话消息失败: %v", err)
		dialog.ShowError(err, cw.window)
		return
	}

	// 清空当前界面
	cw.messages = make([]*models.Message, 0)
	cw.messageContainer.Objects = []fyne.CanvasObject{}

	// 重新加载所有消息到界面
	for _, msg := range messages {
		cw.addMessage(msg)
	}

	// 恢复 AI 服务的历史记录
	cw.aiService.SetHistory(messages)

	cw.currentSession = session
	cw.sessionList.SetCurrentSession(session)
	cw.scrollToBottom()
}

// saveCurrentMessages 保存当前会话的消息
func (cw *ChatWindow) saveCurrentMessages() {
	if cw.currentSession == nil {
		return
	}

	// 获取数据库中已有的消息
	existingMessages, err := cw.db.GetMessages(cw.currentSession.ID)
	if err != nil {
		log.Printf("获取已有消息失败: %v", err)
		return
	}

	// 创建已存在消息的 ID 映射
	existingIDs := make(map[string]bool)
	for _, msg := range existingMessages {
		existingIDs[msg.ID] = true
	}

	// 只保存新消息
	for _, msg := range cw.messages {
		if !existingIDs[msg.ID] {
			if err := cw.db.SaveMessage(cw.currentSession.ID, msg); err != nil {
				log.Printf("保存消息失败: %v", err)
			}
		}
	}
}

// refreshSessionList 刷新会话列表
func (cw *ChatWindow) refreshSessionList() {
	sessions, err := cw.db.ListSessions()
	if err != nil {
		log.Printf("刷新会话列表失败: %v", err)
		return
	}
	cw.sessionList.SetSessions(sessions)
}

// onNewSession 新建会话回调
func (cw *ChatWindow) onNewSession() {
	cw.createNewSession()
}

// onSessionSelect 选择会话回调
func (cw *ChatWindow) onSessionSelect(session *models.Session) {
	if session != nil && (cw.currentSession == nil || session.ID != cw.currentSession.ID) {
		cw.loadSession(session)
	}
}

// onDeleteSession 删除会话回调
func (cw *ChatWindow) onDeleteSession(session *models.Session) {
	dialog.ShowConfirm("确认删除", "确定要删除这个会话吗？所有消息将被删除。", func(ok bool) {
		if ok {
			if err := cw.db.DeleteSession(session.ID); err != nil {
				log.Printf("删除会话失败: %v", err)
				dialog.ShowError(err, cw.window)
				return
			}

			// 如果删除的是当前会话
			if cw.currentSession != nil && cw.currentSession.ID == session.ID {
				// 获取剩余会话列表
				sessions, err := cw.db.ListSessions()
				if err != nil {
					log.Printf("获取会话列表失败: %v", err)
				}

				if len(sessions) > 0 {
					// 加载第一个会话
					cw.loadSession(sessions[0])
					cw.refreshSessionList()
				} else {
					// 没有其他会话，创建新会话
					cw.createNewSession()
				}
			} else {
				// 删除的不是当前会话，只刷新列表
				cw.refreshSessionList()
			}
		}
	}, cw.window)
}
