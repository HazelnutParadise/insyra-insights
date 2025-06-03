package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/HazelnutParadise/insyra-insights/i18n"
)

// WelcomeView 歡迎頁面視圖
type WelcomeView struct {
	// 匯入按鈕
	importCSVButton    widget.Clickable
	importExcelButton  widget.Clickable
	importJSONButton   widget.Clickable
	importSQLiteButton widget.Clickable

	// 其他控件
	startAnalysisButton widget.Clickable

	// 狀態訊息
	message      string
	messageColor color.NRGBA

	// 視圖控制器參考
	viewController *ViewController
}

// NewWelcomeView 創建一個新的歡迎頁面視圖
func NewWelcomeView() *WelcomeView {
	return &WelcomeView{
		messageColor: color.NRGBA{R: 0, G: 128, B: 0, A: 255}, // 綠色
	}
}

// Layout 實現視圖布局
func (v *WelcomeView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 定義顏色
	primaryColor := color.NRGBA{R: 63, G: 81, B: 181, A: 255} // 靛藍色

	// 處理按鈕點擊事件
	v.handleEvents(gtx)
	// 創建標題和副標題，使用 i18n 翻譯
	title := material.H1(th, i18n.T("welcome.title"))
	title.Alignment = text.Middle
	title.Color = primaryColor
	title.TextSize = unit.Sp(32) // 設置標題字體大小

	subtitle := material.H4(th, i18n.T("welcome.subtitle"))
	subtitle.Alignment = text.Middle
	// 獲取當前日期 - 此處顯示指定的測試日期
	// 實際應用中，你可以改回使用 time.Now()
	currentDate := "2025年6月2日" // 硬編碼為當前日期，用於演示
	dateText := material.Body2(th, i18n.T("welcome.today_is")+" "+currentDate)
	dateText.Alignment = text.Middle

	description := material.Body1(th, i18n.T("welcome.description"))
	description.Alignment = text.Middle

	// 佈局所有元素
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: unit.Dp(40)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return title.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return subtitle.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return dateText.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(40)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return description.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.importButtonsLayout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// 開始分析按鈕
			btn := material.Button(th, &v.startAnalysisButton, i18n.T("welcome.start_analysis"))
			btn.Background = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // 綠色
			btn.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}    // 白色
			return layout.Center.Layout(gtx, btn.Layout)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if v.message != "" {
				msg := material.Body1(th, v.message)
				msg.Color = v.messageColor
				msg.Alignment = text.Middle
				return msg.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}

// importButtonsLayout 匯入按鈕區域佈局
func (v *WelcomeView) importButtonsLayout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 匯入按鈕佈局
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
				Spacing:   layout.SpaceEvenly,
			}.Layout(gtx,
				layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &v.importCSVButton, i18n.T("welcome.import_csv"))
					btn.Background = color.NRGBA{R: 33, G: 150, B: 243, A: 255} // 藍色
					return layout.UniformInset(unit.Dp(5)).Layout(gtx, btn.Layout)
				}),
				layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &v.importExcelButton, i18n.T("welcome.import_excel"))
					btn.Background = color.NRGBA{R: 0, G: 150, B: 136, A: 255} // 青綠色
					return layout.UniformInset(unit.Dp(5)).Layout(gtx, btn.Layout)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
				Spacing:   layout.SpaceEvenly,
			}.Layout(gtx, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &v.importJSONButton, i18n.T("welcome.import_json"))
				btn.Background = color.NRGBA{R: 156, G: 39, B: 176, A: 255} // 紫色
				return layout.UniformInset(unit.Dp(5)).Layout(gtx, btn.Layout)
			}),
				layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &v.importSQLiteButton, i18n.T("welcome.import_sqlite"))
					btn.Background = color.NRGBA{R: 233, G: 30, B: 99, A: 255} // 粉紅色
					return layout.UniformInset(unit.Dp(5)).Layout(gtx, btn.Layout)
				}),
			)
		}),
	)
}

// handleEvents 處理按鈕點擊事件
func (v *WelcomeView) handleEvents(gtx layout.Context) {
	if v.importCSVButton.Clicked(gtx) {
		v.message = i18n.T("messages.select_file")
		v.messageColor = color.NRGBA{R: 33, G: 150, B: 243, A: 255} // 藍色

		// 開啟檔案選擇對話框
		go func() {
			result := OpenFileDialog([]string{".csv"})
			if result.Err != nil {
				v.message = i18n.T("messages.dialog_error") + result.Err.Error()
				v.messageColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // 紅色
				return
			}

			// 匯入資料
			importResult := ImportData(result.Path, ImportCSV)
			if importResult.Success {
				v.message = importResult.Message
				v.messageColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // 綠色
			} else {
				v.message = i18n.T("messages.import_fail") + importResult.Message
				v.messageColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // 紅色
			}
		}()
	}

	if v.importExcelButton.Clicked(gtx) {
		v.message = "選擇 Excel 檔案..."
		v.messageColor = color.NRGBA{R: 0, G: 150, B: 136, A: 255} // 青綠色

		// 開啟檔案選擇對話框
		go func() {
			result := OpenFileDialog([]string{".xlsx", ".xls"})
			if result.Err != nil {
				v.message = "無法開啟檔案選擇對話框: " + result.Err.Error()
				v.messageColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // 紅色
				return
			}

			// 匯入資料
			importResult := ImportData(result.Path, ImportExcel)
			if importResult.Success {
				v.message = importResult.Message
				v.messageColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // 綠色
			} else {
				v.message = "匯入失敗: " + importResult.Message
				v.messageColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // 紅色
			}
		}()
	}

	if v.importJSONButton.Clicked(gtx) {
		v.message = "選擇 JSON 檔案..."
		v.messageColor = color.NRGBA{R: 156, G: 39, B: 176, A: 255} // 紫色

		// 開啟檔案選擇對話框
		go func() {
			result := OpenFileDialog([]string{".json"})
			if result.Err != nil {
				v.message = "無法開啟檔案選擇對話框: " + result.Err.Error()
				v.messageColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // 紅色
				return
			}

			// 匯入資料
			importResult := ImportData(result.Path, ImportJSON)
			if importResult.Success {
				v.message = importResult.Message
				v.messageColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // 綠色
			} else {
				v.message = "匯入失敗: " + importResult.Message
				v.messageColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // 紅色
			}
		}()
	}

	if v.importSQLiteButton.Clicked(gtx) {
		v.message = "選擇 SQLite 資料庫..."
		v.messageColor = color.NRGBA{R: 233, G: 30, B: 99, A: 255} // 粉紅色

		// 開啟檔案選擇對話框
		go func() {
			result := OpenFileDialog([]string{".db", ".sqlite"})
			if result.Err != nil {
				v.message = "無法開啟檔案選擇對話框: " + result.Err.Error()
				v.messageColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // 紅色
				return
			} // 匯入資料
			importResult := ImportData(result.Path, ImportSQLite)
			if importResult.Success {
				v.message = importResult.Message
				v.messageColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // 綠色
			} else {
				v.message = "匯入失敗: " + importResult.Message
				v.messageColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // 紅色
			}
		}()
	}

	if v.startAnalysisButton.Clicked(gtx) {
		// 直接切換到資料頁面，不需要等待
		if v.viewController != nil {
			v.viewController.SwitchView("data")
		}
	}
}

// Update 實現視圖更新
func (v *WelcomeView) Update(e interface{}) {
	// 這裡可以處理更新邏輯
}

// Event 實現事件處理
func (v *WelcomeView) Event(e interface{}) {
	// 處理其他類型的事件
}
