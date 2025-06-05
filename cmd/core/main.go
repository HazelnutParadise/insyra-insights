package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/HazelnutParadise/insyra"
	"github.com/HazelnutParadise/insyra-insights/internal/ui"
)

func main() {
	uuid := flag.String("uuid", "", "UUID from loader")
	flag.Parse()

	if *uuid == "" {
		fmt.Println("⚠️ 沒有收到 UUID，直接退出")
		os.Exit(1)
	}
	insyra.Config.SetDontPanic(true)

	go func() {
		window := new(app.Window)
		window.Option(
			app.Title("Insyra Insights"),
		)
		lockPath := filepath.Join(os.TempDir(), "insyra_starting_"+*uuid+".lock")
		time.Sleep(2 * time.Second)
		err := os.Remove(lockPath)
		if err != nil {
			log.Fatalf("無法刪除鎖定檔案: %v", err)
		}
		err = run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func run(window *app.Window) error {
	ui.GlobalWindow = window
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
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// 創建圖形上下文
			gtx := app.NewContext(&ops, e)

			// 佈局界面 - 移除底部導覽按鈕，只顯示視圖區域
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				// 視圖區域佔滿整個空間
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return viewController.Layout(gtx)
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
