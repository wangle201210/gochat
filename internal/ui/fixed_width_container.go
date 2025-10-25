package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// fixedWidthContainer 固定宽度的容器
type fixedWidthContainer struct {
	widget.BaseWidget
	content fyne.CanvasObject
	width   float32
}

// newFixedWidthContainer 创建固定宽度容器
func newFixedWidthContainer(width float32, content fyne.CanvasObject) *fixedWidthContainer {
	f := &fixedWidthContainer{
		content: content,
		width:   width,
	}
	f.ExtendBaseWidget(f)
	return f
}

// CreateRenderer 创建渲染器
func (f *fixedWidthContainer) CreateRenderer() fyne.WidgetRenderer {
	return &fixedWidthRenderer{
		container: f,
		content:   f.content,
	}
}

// fixedWidthRenderer 固定宽度渲染器
type fixedWidthRenderer struct {
	container *fixedWidthContainer
	content   fyne.CanvasObject
}

func (r *fixedWidthRenderer) Layout(size fyne.Size) {
	// 固定宽度，高度使用传入的高度
	r.content.Resize(fyne.NewSize(r.container.width, size.Height))
	r.content.Move(fyne.NewPos(0, 0))
}

func (r *fixedWidthRenderer) MinSize() fyne.Size {
	// 返回固定宽度和内容的最小高度
	minSize := r.content.MinSize()
	return fyne.NewSize(r.container.width, minSize.Height)
}

func (r *fixedWidthRenderer) Refresh() {
	r.content.Refresh()
}

func (r *fixedWidthRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

func (r *fixedWidthRenderer) Destroy() {}
