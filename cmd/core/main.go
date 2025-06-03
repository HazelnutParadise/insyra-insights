package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/HazelnutParadise/insyra-insights/internal/ui"
)

func main() {
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops
	// 創建視圖控制器
	viewController := ui.NewViewController(theme)

	// 創建視圖實例
	dataView := ui.NewDataView()

	// 載入樣本數據到 DataView
	dataView.LoadSampleData()

	// 設置視圖控制器引用
	dataView.SetViewController(viewController)

	// 註冊視圖
	viewController.RegisterView("welcome", ui.NewWelcomeView())
	viewController.RegisterView("about", ui.NewAboutView())
	viewController.RegisterView("settings", ui.NewSettingsView())
	viewController.RegisterView("data", dataView)

	// 設置初始視圖
	viewController.SwitchView("welcome")
	// 創建切換按鈕
	var welcomeButton widget.Clickable
	var aboutButton widget.Clickable
	var settingsButton widget.Clickable
	var dataButton widget.Clickable

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// 創建圖形上下文
			gtx := app.NewContext(&ops, e)
			// 處理按鈕點擊
			if welcomeButton.Clicked(gtx) {
				viewController.SwitchView("welcome")
			}
			if aboutButton.Clicked(gtx) {
				viewController.SwitchView("about")
			}
			if settingsButton.Clicked(gtx) {
				viewController.SwitchView("settings")
			}
			if dataButton.Clicked(gtx) {
				viewController.SwitchView("data")
			}

			// 佈局界面
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				// 視圖區域
				layout.Flexed(0.9, func(gtx layout.Context) layout.Dimensions {
					return viewController.Layout(gtx)
				}),
				// 按鈕區域
				layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
							return material.Button(theme, &welcomeButton, "歡迎").Layout(gtx)
						}),
						layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
							return material.Button(theme, &aboutButton, "關於").Layout(gtx)
						}),
						layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
							return material.Button(theme, &settingsButton, "設定").Layout(gtx)
						}),
						layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
							return material.Button(theme, &dataButton, "資料").Layout(gtx)
						}),
					)
				}),
			)
			// 更新視圖控制器中的當前視圖
			viewController.Update(e)

			// 轉發事件給視圖控制器
			viewController.Event(e)

			// 傳送繪製操作到 GPU
			e.Frame(gtx.Ops)
		}
	}
}
