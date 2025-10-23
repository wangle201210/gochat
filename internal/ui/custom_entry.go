package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

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
