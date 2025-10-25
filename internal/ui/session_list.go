package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/wangle201210/gochat/internal/models"
)

// sessionListItem 会话列表项
type sessionListItem struct {
	widget.BaseWidget
	label      *widget.Label
	deleteBtn  *widget.Button
	background *canvas.Rectangle
	content    *fyne.Container
	container  *fyne.Container
	onTapped   func()
	onDelete   func()
}

func newSessionListItem(text string, onTapped func(), onDelete func()) *sessionListItem {
	item := &sessionListItem{
		label:    widget.NewLabel(text),
		onTapped: onTapped,
		onDelete: onDelete,
	}

	// 创建删除按钮，使用低优先级样式让它不那么显眼
	item.deleteBtn = widget.NewButton("✕", func() {
		if item.onDelete != nil {
			item.onDelete()
		}
	})
	item.deleteBtn.Importance = widget.LowImportance

	// 创建背景矩形（默认透明）
	item.background = canvas.NewRectangle(color.Transparent)

	// 创建内容容器
	item.content = container.NewBorder(nil, nil, nil, item.deleteBtn, item.label)

	// 使用 Stack 将背景和内容叠加
	item.container = container.NewStack(item.background, container.NewPadded(item.content))
	item.ExtendBaseWidget(item)
	return item
}

func (i *sessionListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(i.container)
}

func (i *sessionListItem) Tapped(_ *fyne.PointEvent) {
	if i.onTapped != nil {
		i.onTapped()
	}
}

func (i *sessionListItem) SetText(text string) {
	i.label.SetText(text)
}

func (i *sessionListItem) SetBold(bold bool) {
	if bold {
		i.label.TextStyle = fyne.TextStyle{Bold: true}
	} else {
		i.label.TextStyle = fyne.TextStyle{}
	}
	i.label.Refresh()
}

func (i *sessionListItem) SetHighlight(highlight bool) {
	if highlight {
		// 高亮背景色 - 淡蓝色
		i.background.FillColor = color.NRGBA{R: 230, G: 240, B: 255, A: 255}
		i.label.TextStyle = fyne.TextStyle{Bold: true}
	} else {
		// 透明背景
		i.background.FillColor = color.Transparent
		i.label.TextStyle = fyne.TextStyle{}
	}
	i.background.Refresh()
	i.label.Refresh()
}

// SessionList 会话列表组件
type SessionList struct {
	widget.BaseWidget
	sessions        []*models.Session
	currentSession  *models.Session
	onSessionSelect func(*models.Session)
	onNewSession    func()
	onDeleteSession func(*models.Session)
	list            *widget.List
}

// NewSessionList 创建会话列表
func NewSessionList(onSessionSelect func(*models.Session), onNewSession func(), onDeleteSession func(*models.Session)) *SessionList {
	sl := &SessionList{
		sessions:        make([]*models.Session, 0),
		onSessionSelect: onSessionSelect,
		onNewSession:    onNewSession,
		onDeleteSession: onDeleteSession,
	}
	sl.ExtendBaseWidget(sl)
	return sl
}

// CreateRenderer 实现 Widget 接口
func (sl *SessionList) CreateRenderer() fyne.WidgetRenderer {
	// 创建新会话按钮
	newSessionBtn := widget.NewButton("开启新会话", func() {
		if sl.onNewSession != nil {
			sl.onNewSession()
		}
	})
	newSessionBtn.Importance = widget.HighImportance

	// 创建会话列表
	sl.list = widget.NewList(
		func() int {
			return len(sl.sessions)
		},
		func() fyne.CanvasObject {
			return newSessionListItem("会话标题", nil, nil)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < 0 || id >= len(sl.sessions) {
				return
			}

			session := sl.sessions[id]
			listItem := item.(*sessionListItem)

			// 设置标题
			listItem.SetText(session.Title)

			// 高亮当前会话 - 使用背景色和粗体
			isCurrentSession := sl.currentSession != nil && session.ID == sl.currentSession.ID
			listItem.SetHighlight(isCurrentSession)

			// 设置回调
			listItem.onTapped = func() {
				if sl.onSessionSelect != nil {
					sl.onSessionSelect(session)
				}
			}
			listItem.onDelete = func() {
				if sl.onDeleteSession != nil {
					sl.onDeleteSession(session)
				}
			}
		},
	)

	content := container.NewBorder(
		container.NewVBox(
			newSessionBtn,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		sl.list,
	)

	return widget.NewSimpleRenderer(content)
}

// SetSessions 设置会话列表
func (sl *SessionList) SetSessions(sessions []*models.Session) {
	sl.sessions = sessions
	if sl.list != nil {
		sl.list.Refresh()
	}
}

// SetCurrentSession 设置当前会话
func (sl *SessionList) SetCurrentSession(session *models.Session) {
	sl.currentSession = session
	if sl.list != nil {
		sl.list.Refresh()
	}
}

// GetCurrentSession 获取当前会话
func (sl *SessionList) GetCurrentSession() *models.Session {
	return sl.currentSession
}

// ShowRenameDialog 显示重命名对话框
func ShowRenameDialog(window fyne.Window, session *models.Session, onRename func(string)) {
	entry := widget.NewEntry()
	entry.SetText(session.Title)

	dialog.ShowForm("重命名会话", "确定", "取消", []*widget.FormItem{
		widget.NewFormItem("会话标题", entry),
	}, func(ok bool) {
		if ok && entry.Text != "" {
			onRename(entry.Text)
		}
	}, window)
}
