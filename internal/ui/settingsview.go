package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// SettingsView 設置頁視圖
type SettingsView struct {
	// 表單元素
	nameInput     widget.Editor
	emailInput    widget.Editor
	saveButton    widget.Clickable
	resetButton   widget.Clickable
	darkModeCheck widget.Bool

	// 表單數據
	formData FormData

	// 狀態
	message string
}

// FormData 表單數據結構
type FormData struct {
	Name     string
	Email    string
	DarkMode bool
}

// NewSettingsView 創建一個新的設置頁視圖
func NewSettingsView() *SettingsView {
	view := &SettingsView{}

	// 設置文本輸入框屬性
	view.nameInput.SingleLine = true
	view.nameInput.SetText("使用者")
	view.emailInput.SingleLine = true
	view.emailInput.SetText("user@example.com")

	// 設置初始表單數據
	view.formData = FormData{
		Name:     "使用者",
		Email:    "user@example.com",
		DarkMode: false,
	}

	return view
}

// Layout 實現視圖布局
func (v *SettingsView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 創建標題
	title := material.H2(th, "設定")
	title.Alignment = text.Middle
	title.Color = color.NRGBA{R: 0, G: 100, B: 100, A: 255} // 藍綠色

	// 處理按鈕點擊
	v.handleEvents(gtx)

	// 外層佈局
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return title.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.formLayout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.buttonsLayout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if v.message != "" {
				msg := material.Body1(th, v.message)
				msg.Color = color.NRGBA{R: 0, G: 150, B: 0, A: 255} // 綠色
				return msg.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}

// formLayout 佈局表單區塊
func (v *SettingsView) formLayout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// 姓名輸入框
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(th, "姓名:")
					return label.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, &v.nameInput, "請輸入姓名")
					return editor.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		// 電子郵件輸入框
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(th, "電子郵件:")
					return label.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, &v.emailInput, "請輸入電子郵件")
					return editor.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		// 暗黑模式選項
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.CheckBox(th, &v.darkModeCheck, "啟用暗黑模式").Layout(gtx)
		}),
	)
}

// buttonsLayout 佈局按鈕區塊
func (v *SettingsView) buttonsLayout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceBetween,
	}.Layout(gtx,
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &v.saveButton, "儲存設定")
			btn.Background = color.NRGBA{R: 0, G: 150, B: 100, A: 255}
			return btn.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &v.resetButton, "重置")
			btn.Background = color.NRGBA{R: 150, G: 100, B: 0, A: 255}
			return btn.Layout(gtx)
		}),
	)
}

// handleEvents 處理事件
func (v *SettingsView) handleEvents(gtx layout.Context) {
	if v.saveButton.Clicked(gtx) {
		// 儲存表單數據
		v.formData.Name = v.nameInput.Text()
		v.formData.Email = v.emailInput.Text()
		v.formData.DarkMode = v.darkModeCheck.Value
		v.message = "設定已儲存!"
	}

	if v.resetButton.Clicked(gtx) {
		// 重置表單
		v.nameInput.SetText("使用者")
		v.emailInput.SetText("user@example.com")
		v.darkModeCheck.Value = false
		v.message = "設定已重置!"
	}
}

// Update 實現視圖更新
func (v *SettingsView) Update(e interface{}) {
	// 更新視圖狀態
}

// Event 實現事件處理
func (v *SettingsView) Event(e interface{}) {
	// 處理視圖事件
}
