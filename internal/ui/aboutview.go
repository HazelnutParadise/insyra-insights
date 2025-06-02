package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
)

// AboutView 關於頁視圖
type AboutView struct {
	// 這裡可以添加視圖特有的狀態和屬性
}

// NewAboutView 創建一個新的關於頁視圖
func NewAboutView() *AboutView {
	return &AboutView{}
}

// Layout 實現視圖布局
func (v *AboutView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 創建一個大標籤
	title := material.H1(th, "關於我們")

	// 更改標籤的顏色
	title.Color = color.NRGBA{R: 0, G: 0, B: 127, A: 255} // 藍色

	// 更改標籤的位置
	title.Alignment = text.Middle

	// 繪製標籤
	return title.Layout(gtx)
}

// Update 實現視圖更新
func (v *AboutView) Update(e interface{}) {
	// 處理視圖狀態更新
}

// Event 實現事件處理
func (v *AboutView) Event(e interface{}) {
	// 處理視圖事件
}
