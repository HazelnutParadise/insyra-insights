package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
)

// HomeView 首頁視圖
type HomeView struct {
	// 這裡可以添加視圖特有的狀態和屬性
}

// NewHomeView 創建一個新的首頁視圖
func NewHomeView() *HomeView {
	return &HomeView{}
}

// Layout 實現視圖布局
func (v *HomeView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 創建一個大標籤
	title := material.H1(th, "首頁")

	// 更改標籤的顏色
	title.Color = color.NRGBA{R: 0, G: 127, B: 0, A: 255} // 綠色

	// 更改標籤的位置
	title.Alignment = text.Middle

	// 繪製標籤
	return title.Layout(gtx)
}

// Update 實現視圖更新
func (v *HomeView) Update(e interface{}) {
	// 處理視圖狀態更新
}

// Event 實現事件處理
func (v *HomeView) Event(e interface{}) {
	// 處理視圖事件
}
