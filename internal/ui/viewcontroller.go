package ui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

// ViewController 負責管理和切換不同的視圖
type ViewController struct {
	// 當前視圖
	currentView View
	// 註冊的視圖集合
	views map[string]View
	// 主題
	theme *material.Theme
}

// NewViewController 創建一個新的視圖控制器
func NewViewController(theme *material.Theme) *ViewController {
	return &ViewController{
		views: make(map[string]View),
		theme: theme,
	}
}

// RegisterView 註冊一個新視圖
func (vc *ViewController) RegisterView(name string, view View) {
	vc.views[name] = view

	// 設置 viewController 參考
	if welcomeView, ok := view.(*WelcomeView); ok {
		welcomeView.viewController = vc
	}
	if dataView, ok := view.(*DataView); ok {
		dataView.SetViewController(vc)
	}
}

// SwitchView 切換到指定的視圖
func (vc *ViewController) SwitchView(name string) bool {
	view, exists := vc.views[name]
	if !exists {
		return false
	}
	vc.currentView = view

	// 嘗試設置 viewController 參考
	if welcomeView, ok := view.(*WelcomeView); ok {
		welcomeView.viewController = vc
	}
	if dataView, ok := view.(*DataView); ok {
		dataView.SetViewController(vc)
	}

	return true
}

// GetCurrentView 獲取當前視圖
func (vc *ViewController) GetCurrentView() View {
	return vc.currentView
}

// Layout 處理當前視圖的佈局和繪製
func (vc *ViewController) Layout(gtx layout.Context) layout.Dimensions {
	if vc.currentView == nil {
		return layout.Dimensions{}
	}
	return vc.currentView.Layout(gtx, vc.theme)
}

// Update 更新當前視圖狀態
func (vc *ViewController) Update(e interface{}) {
	if vc.currentView != nil {
		vc.currentView.Update(e)
	}
}

// Event 處理事件並轉發給當前視圖
func (vc *ViewController) Event(e interface{}) {
	if vc.currentView != nil {
		vc.currentView.Event(e)
	}
}
