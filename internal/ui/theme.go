package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	// 颜色方案
	userMessageBg   = color.NRGBA{R: 240, G: 248, B: 255, A: 255} // 淡蓝白
	assistantBg     = color.NRGBA{R: 255, G: 253, B: 245, A: 255} // 温暖米白
	backgroundColor = color.NRGBA{R: 250, G: 252, B: 252, A: 255} // 清新白
)

// customTheme 自定义主题
type customTheme struct {
	fyne.Theme
}

func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameSeparator {
		return color.NRGBA{R: 230, G: 230, B: 230, A: 255}
	}
	return t.Theme.Color(name, variant)
}

func (t *customTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameSeparatorThickness {
		return 1
	}
	return t.Theme.Size(name)
}

// newCustomTheme 创建自定义主题
func newCustomTheme() fyne.Theme {
	return &customTheme{Theme: theme.DefaultTheme()}
}
