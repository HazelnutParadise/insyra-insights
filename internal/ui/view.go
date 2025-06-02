package ui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

// View 介面定義了一個畫面應該具備的功能
type View interface {
	// Layout 處理畫面的佈局和繪製
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
	// Update 更新畫面狀態，例如處理使用者輸入
	Update(e interface{})
	// Event 處理事件
	Event(e interface{})
}
