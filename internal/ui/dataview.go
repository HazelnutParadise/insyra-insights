package ui

import (
	"fmt"
	"image"
	"image/color"
	"log" // Re-added import
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/HazelnutParadise/insyra"
)

// DataTabInfo 每個標籤頁的資訊
type DataTabInfo struct {
	Name      string
	DataTable *GenericDataTable
	StatsData map[string]string
}

// DataView 支援多標籤頁的數據視圖
type DataView struct {
	// 多標籤頁支援
	tabs             []*DataTabInfo     // 標籤頁列表
	tabButtons       []widget.Clickable // 標籤按鈕
	addTabButton     widget.Clickable   // 新增標籤按鈕
	currentTabIndex  int                // 當前選中的標籤索引	// 功能列按鈕
	addColButton     widget.Clickable   // 新增欄按鈕
	addRowButton     widget.Clickable   // 新增列按鈕
	addCalcColButton widget.Clickable   // 新增計算欄按鈕

	// 計算欄功能
	showColumnInput     bool             // 是否顯示計算欄輸入條
	columnFormulaEditor widget.Editor    // CCL 表達式輸入編輯器
	columnNameEditor    widget.Editor    // 新欄位名稱輸入編輯器
	addColumnConfirmBtn widget.Clickable // 確認添加按鈕
	addColumnCancelBtn  widget.Clickable // 取消按鈕
	errorMessage        string           // 錯誤訊息
	showError           bool             // 是否顯示錯誤訊息

	// 底部工具列按鈕
	openButton     widget.Clickable // 開啟檔案
	saveButton     widget.Clickable // 存檔
	exportButton   widget.Clickable // 匯出
	settingsButton widget.Clickable // 設定

	// 視圖控制器參考
	viewController *ViewController
}

// NewDataView 創建一個新的多標籤頁數據視圖
func NewDataView() *DataView {
	// 創建第一個標籤頁
	firstTab := &DataTabInfo{
		Name:      "Tab 1",
		DataTable: NewGenericDataTable(insyra.NewDataTable()),
		StatsData: make(map[string]string),
	}
	view := &DataView{
		tabs:            []*DataTabInfo{firstTab},
		tabButtons:      make([]widget.Clickable, 1),
		currentTabIndex: 0,
		showColumnInput: false, // 初始不顯示輸入條
		showError:       false, // 初始不顯示錯誤訊息
		errorMessage:    "",    // 初始化錯誤訊息為空
	}

	// 設定公式編輯器
	view.columnFormulaEditor.SingleLine = true
	view.columnFormulaEditor.Submit = true

	// 設定名稱編輯器
	view.columnNameEditor.SingleLine = true
	view.columnNameEditor.Submit = true

	return view
}

// Layout 實現多標籤頁界面
func (v *DataView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 處理新增標籤頁按鈕
	if v.addTabButton.Clicked(gtx) {
		v.addNewTab()
	}

	// 處理標籤按鈕點擊
	for i := range v.tabButtons {
		if v.tabButtons[i].Clicked(gtx) {
			v.currentTabIndex = i
		}
	}

	// 處理功能按鈕
	v.handleFunctionButtons(gtx)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// 標籤列
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.layoutTabBar(gtx, th)
		}),
		// 功能列（新增欄/列）
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.layoutFunctionBar(gtx, th)
		}),
		// 主要內容區域（表格區 + 資訊區）
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return v.layoutMainContent(gtx, th)
		}),
		// 底部工具列
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.layoutBottomToolbar(gtx, th)
		}),
	)
}

// addNewTab 新增標籤頁
func (v *DataView) addNewTab() {
	newTabName := fmt.Sprintf("Tab %d", len(v.tabs)+1)
	newTab := &DataTabInfo{
		Name:      newTabName,
		DataTable: NewGenericDataTable(insyra.NewDataTable()),
		StatsData: make(map[string]string),
	}

	v.tabs = append(v.tabs, newTab)
	v.tabButtons = append(v.tabButtons, widget.Clickable{})
	v.currentTabIndex = len(v.tabs) - 1 // 切換到新建的標籤頁
}

// handleFunctionButtons 處理功能按鈕
func (v *DataView) handleFunctionButtons(gtx layout.Context) {
	if v.addColButton.Clicked(gtx) {
		v.addColumn()
	}
	if v.addRowButton.Clicked(gtx) {
		v.addRow()
	}
	if v.addCalcColButton.Clicked(gtx) {
		v.showColumnInput = true
	}

	// 處理計算欄按鈕
	if v.addColumnConfirmBtn.Clicked(gtx) {
		v.addCalculatedColumn()
	}
	if v.addColumnCancelBtn.Clicked(gtx) {
		v.cancelColumnInput()
	}
}

// layoutTabBar 繪製標籤列，並支援自動換行
func (v *DataView) layoutTabBar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 1. 定義小工具佈局函式
	widgetLayoutFuncs := make([]func(gtx layout.Context) layout.Dimensions, 0, len(v.tabs)+1)

	for i, tabInfo := range v.tabs {
		capturedIndex := i
		capturedTabInfo := tabInfo
		widgetLayoutFuncs = append(widgetLayoutFuncs, func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &v.tabButtons[capturedIndex], capturedTabInfo.Name)
			if capturedIndex == v.currentTabIndex {
				// 選中標籤使用與表格選中列相同的淡綠色背景
				btn.Background = color.NRGBA{R: 235, G: 250, B: 235, A: 255} // 淡綠色
				btn.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}            // 黑色文字
			} else {
				// 未選中標籤使用與表格標題相同的淡藍色背景
				btn.Background = color.NRGBA{R: 225, G: 235, B: 250, A: 255} // 淡藍色
				btn.Color = color.NRGBA{R: 0, G: 90, B: 180, A: 255}         // 藍色文字
			}
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
		})
	}

	// 新增標籤按鈕
	widgetLayoutFuncs = append(widgetLayoutFuncs, func(gtx layout.Context) layout.Dimensions {
		btn := material.Button(th, &v.addTabButton, "+")
		// 新增按鈕使用與計算欄按鈕相同的藍色樣式
		btn.Background = color.NRGBA{R: 225, G: 245, B: 254, A: 255} // 淡藍色背景
		btn.Color = color.NRGBA{R: 33, G: 150, B: 243, A: 255}       // 藍色文字
		return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
	})

	// 2. 實作換行邏輯
	var linesFlexChildren []layout.FlexChild // 儲存每一行 (本身是一個 FlexChild)

	currentLineWidgets := []func(gtx layout.Context) layout.Dimensions{}
	currentLineWidthPixels := 0
	maxLineWidthPixels := gtx.Constraints.Max.X
	itemSpacingDp := unit.Dp(4)                // 按鈕間的間距 (可調整)
	itemSpacingPixels := gtx.Dp(itemSpacingDp) // Convert Dp to pixels for calculations
	lineSpacing := unit.Dp(4)                  // 行間距 (可調整)

	for _, widgetFunc := range widgetLayoutFuncs {
		// 測量 widgetFunc 的寬度
		macro := op.Record(gtx.Ops)
		mgtx := gtx
		mgtx.Constraints.Min = image.Point{} // 允許小工具自行決定最小尺寸
		// Max.X 保持 gtx.Constraints.Max.X，以便小工具知道可用空間，但它應該回報其偏好寬度
		widgetDims := widgetFunc(mgtx)
		_ = macro.Stop() // 停止錄製，我們只關心尺寸

		widgetWidthPixels := widgetDims.Size.X

		// 檢查是否需要換行
		// 如果目前行非空，且加入此小工具 (包含間距) 會超出最大寬度
		if len(currentLineWidgets) > 0 && currentLineWidthPixels+itemSpacingPixels+widgetWidthPixels > maxLineWidthPixels {
			// 目前行已滿，將其作為一個 FlexChild 加入到 linesFlexChildren
			// 複製 currentLineWidgets 以避免閉包問題
			lineToLayout := make([]func(gtx layout.Context) layout.Dimensions, len(currentLineWidgets))
			copy(lineToLayout, currentLineWidgets)

			linesFlexChildren = append(linesFlexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				rowChildren := []layout.FlexChild{}
				for i, f := range lineToLayout {
					f := f     // 捕獲
					if i > 0 { // 在項目之間加入間距
						rowChildren = append(rowChildren, layout.Rigid(layout.Spacer{Width: itemSpacingDp}.Layout))
					}
					rowChildren = append(rowChildren, layout.Rigid(f))
				}
				// 使用 Flex 佈局目前行的所有小工具
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx, rowChildren...)
			}))

			// 開始新的一行
			currentLineWidgets = []func(gtx layout.Context) layout.Dimensions{}
			currentLineWidthPixels = 0
		}

		// 將目前小工具加入到目前行
		if len(currentLineWidgets) > 0 {
			currentLineWidthPixels += itemSpacingPixels // 加上項目間距
		}
		currentLineWidgets = append(currentLineWidgets, widgetFunc)
		currentLineWidthPixels += widgetWidthPixels
	}

	// 加入最後一行 (如果有的話)
	if len(currentLineWidgets) > 0 {
		lineToLayout := make([]func(gtx layout.Context) layout.Dimensions, len(currentLineWidgets))
		copy(lineToLayout, currentLineWidgets)

		linesFlexChildren = append(linesFlexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			rowChildren := []layout.FlexChild{}
			for i, f := range lineToLayout {
				f := f     // 捕獲
				if i > 0 { // 在項目之間加入間距
					rowChildren = append(rowChildren, layout.Rigid(layout.Spacer{Width: itemSpacingDp}.Layout))
				}
				rowChildren = append(rowChildren, layout.Rigid(f))
			}
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx, rowChildren...)
		}))
	}

	// 3. 使用垂直 Flex 佈局所有行
	if len(linesFlexChildren) > 1 { // 如果有多行，才加入行間距
		spacedLines := make([]layout.FlexChild, 0, len(linesFlexChildren)*2-1)
		for i, lineChild := range linesFlexChildren {
			spacedLines = append(spacedLines, lineChild)
			if i < len(linesFlexChildren)-1 {
				spacedLines = append(spacedLines, layout.Rigid(layout.Spacer{Height: lineSpacing}.Layout))
			}
		}
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx, spacedLines...)
	}
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx, linesFlexChildren...)
}

// layoutFunctionBar 繪製功能列
func (v *DataView) layoutFunctionBar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 使用堆疊佈局，增加功能列背景色和陰影效果
	return layout.Stack{}.Layout(gtx,
		// 背景與陰影層
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Max
			// 使用藍紫色調背景，提升現代感
			bgColor := color.NRGBA{R: 63, G: 81, B: 181, A: 255} // Material Design 的 Indigo 500
			paint.FillShape(gtx.Ops, bgColor, clip.Rect{
				Max: image.Pt(size.X, gtx.Dp(unit.Dp(52))),
			}.Op())

			// 底部陰影效果
			shadowHeight := 6
			for i := range shadowHeight {
				y := gtx.Dp(unit.Dp(52)) + i
				alpha := uint8(40 - i*7)
				if alpha < 3 {
					alpha = 3
				}

				paint.FillShape(gtx.Ops, color.NRGBA{0, 0, 0, alpha}, clip.Rect{
					Min: image.Pt(0, y),
					Max: image.Pt(size.X, y+1),
				}.Op())
			}

			return layout.Dimensions{Size: image.Pt(size.X, gtx.Dp(unit.Dp(52)))}
		}),

		// 內容層
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				// 按鈕行
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &v.addColButton, "新增欄")
							// 更改按鈕顏色以搭配新背景
							btn.Background = color.NRGBA{R: 255, G: 255, B: 255, A: 40} // 半透明白色背景
							btn.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}     // 白色文字
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, btn.Layout)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &v.addRowButton, "新增列")
							// 更改按鈕顏色以搭配新背景
							btn.Background = color.NRGBA{R: 255, G: 255, B: 255, A: 40} // 半透明白色背景
							btn.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}     // 白色文字
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, btn.Layout)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &v.addCalcColButton, "新增計算欄")
							// 更改按鈕顏色以搭配新背景
							btn.Background = color.NRGBA{R: 255, G: 255, B: 255, A: 40} // 半透明白色背景
							btn.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}     // 白色文字
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, btn.Layout)
						}),
					)
				}),
				// 計算欄輸入區域
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if !v.showColumnInput {
						return layout.Dimensions{}
					}
					return v.layoutColumnInput(gtx, th)
				}),
			)
		}),
	)
}

// layoutMainContent 繪製主要內容區域
func (v *DataView) layoutMainContent(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if len(v.tabs) == 0 {
		return material.Body1(th, "沒有打開的標籤頁").Layout(gtx)
	}

	currentTab := v.tabs[v.currentTabIndex]

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		// 左側：表格區域
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			return layout.NW.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return v.layoutTableArea(gtx, th, currentTab)
			})
		}),
		// 右側：資訊區域
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return v.layoutInfoArea(gtx, th, currentTab)
		}),
	)
}

// layoutTableArea 繪製表格區域
func (v *DataView) layoutTableArea(gtx layout.Context, th *material.Theme, tab *DataTabInfo) layout.Dimensions {
	if tab.DataTable == nil || tab.DataTable.Table == nil {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, "資料表為空，請新增資料").Layout(gtx)
		})
	}

	rowCount, colCount := tab.DataTable.Table.Size()
	if rowCount == 0 && colCount == 0 {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, "資料表為空，請新增資料").Layout(gtx)
		})
	}

	return tab.DataTable.Layout(gtx, th)
}

// layoutInfoArea 繪製資訊區域
func (v *DataView) layoutInfoArea(gtx layout.Context, th *material.Theme, tab *DataTabInfo) layout.Dimensions {
	// 設定資訊區背景顏色
	bgColor := color.NRGBA{R: 245, G: 248, B: 255, A: 255} // 淡藍灰色
	return layout.Stack{}.Layout(gtx,
		// 背景層
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Max
			paint.FillShape(gtx.Ops, bgColor, clip.Rect{Max: size}.Op())
			return layout.Dimensions{Size: size}
		}),
		// 內容層
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := material.H6(th, "資訊區")
					title.Font.Weight = font.Bold
					title.Color = color.NRGBA{R: 0, G: 90, B: 180, A: 255} // 藍色
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, title.Layout)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return v.layoutStats(gtx, th, tab)
				}),
			)
		}),
	)
}

// layoutStats 顯示統計資訊
func (v *DataView) layoutStats(gtx layout.Context, th *material.Theme, tab *DataTabInfo) layout.Dimensions {
	children := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.Body1(th, "基本統計")
			title.Font.Weight = font.Bold
			// 使用與表格標題相同的藍色文字
			title.Color = color.NRGBA{R: 0, G: 90, B: 180, A: 255} // 藍色
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, title.Layout)
		}),
	}

	// 顯示統計數據
	// todo: 改用struct切片，維持順序
	for key, value := range tab.StatsData {
		key, value := key, value // 捕獲循環變數
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.layoutStatItem(gtx, th, key, value)
		}))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

// layoutStatItem 顯示單個統計項目
func (v *DataView) layoutStatItem(gtx layout.Context, th *material.Theme, label, value string) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(4),
		Bottom: unit.Dp(4),
		Left:   unit.Dp(8),
		Right:  unit.Dp(8),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				labelText := material.Caption(th, label+":")
				return labelText.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				valueText := material.Body2(th, value)
				return valueText.Layout(gtx)
			}),
		)
	})
}

// layoutBottomToolbar 繪製底部工具列
func (v *DataView) layoutBottomToolbar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 使用清爽的淺色背景，營造年輕感
	bgColor := color.NRGBA{R: 180, G: 220, B: 255, A: 255} // 更藍的淺藍紫色背景
	height := gtx.Dp(unit.Dp(52))
	return layout.Stack{}.Layout(gtx,
		// 背景層
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Constrain(image.Pt(gtx.Constraints.Max.X, height))
			// 主背景
			paint.FillShape(gtx.Ops, bgColor, clip.Rect{Max: size}.Op())

			// 頂部分隔線
			paint.FillShape(gtx.Ops, color.NRGBA{R: 220, G: 225, B: 230, A: 255}, clip.Rect{
				Max: image.Pt(size.X, 1),
			}.Op())

			return layout.Dimensions{Size: size}
		}),
		// 按鈕層
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			btns := []*widget.Clickable{&v.openButton, &v.saveButton, &v.exportButton, &v.settingsButton}
			labels := []string{"開啟", "存檔", "匯出", "設定"}
			// 不同按鈕使用不同淡色背景
			backgroundColors := []color.NRGBA{
				{R: 240, G: 253, B: 244, A: 255}, // 開啟 - 淡綠色背景
				{R: 239, G: 246, B: 255, A: 255}, // 存檔 - 淡藍色背景
				{R: 255, G: 251, B: 235, A: 255}, // 匯出 - 淡橙色背景
				{R: 249, G: 250, B: 251, A: 255}, // 設定 - 淡灰色背景
			}
			// 文字顏色
			textColors := []color.NRGBA{
				{R: 34, G: 197, B: 94, A: 255},   // 開啟 - 綠色文字
				{R: 59, G: 130, B: 246, A: 255},  // 存檔 - 藍色文字
				{R: 245, G: 158, B: 11, A: 255},  // 匯出 - 橙色文字
				{R: 107, G: 114, B: 128, A: 255}, // 設定 - 灰色文字
			}

			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, btns[0], labels[0])
					btn.Background = backgroundColors[0] // 淡綠色背景
					btn.Color = textColors[0]            // 綠色文字
					btn.CornerRadius = unit.Dp(0)
					return layout.UniformInset(unit.Dp(6)).Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, btns[1], labels[1])
					btn.Background = backgroundColors[1] // 淡藍色背景
					btn.Color = textColors[1]            // 藍色文字
					btn.CornerRadius = unit.Dp(0)
					return layout.UniformInset(unit.Dp(6)).Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, btns[2], labels[2])
					btn.Background = backgroundColors[2] // 淡橙色背景
					btn.Color = textColors[2]            // 橙色文字
					btn.CornerRadius = unit.Dp(0)
					return layout.UniformInset(unit.Dp(6)).Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, btns[3], labels[3])
					btn.Background = backgroundColors[3] // 淡灰色背景
					btn.Color = textColors[3]            // 灰色文字
					btn.CornerRadius = unit.Dp(0)
					return layout.UniformInset(unit.Dp(6)).Layout(gtx, btn.Layout)
				}),
			)
		}),
	)
}

// addColumn 新增欄位到當前標籤頁
func (v *DataView) addColumn() {
	if len(v.tabs) == 0 {
		return
	}
	tab := v.tabs[v.currentTabIndex]
	if tab.DataTable == nil || tab.DataTable.Table == nil {
		log.Println("addColumn: DataTable or Table is nil")
		return
	}

	_, colCount := tab.DataTable.Table.Size()
	columnName := "var" + strconv.Itoa(colCount+1)

	// 獲取當前表格的行數，以確定新欄位的長度
	rowCount, _ := tab.DataTable.Table.Size()
	newColumn := make([]interface{}, rowCount)
	// 如果沒有行數，至少添加一個預設值，以便欄位存在
	if rowCount == 0 {
		newColumn = append(newColumn, nil)
	}

	dl := insyra.NewDataList().SetName(columnName)
	for _, item := range newColumn {
		dl.Append(item) // Use Append to add items
	}
	tab.DataTable.Table.AppendCols(dl)

	// 重新計算統計數據
	v.computeStatistics(tab)
}

// addRow 新增列到當前標籤頁
func (v *DataView) addRow() {
	if len(v.tabs) == 0 {
		return
	}

	currentTab := v.tabs[v.currentTabIndex]
	if currentTab.DataTable == nil || currentTab.DataTable.Table == nil {
		return
	}

	// 獲取目前欄位數量
	_, colCount := currentTab.DataTable.Table.Size()

	// 如果沒有欄位，先創建一個預設欄位
	if colCount == 0 {
		defaultCol := insyra.NewDataList().SetName("var1")
		defaultCol.Append(nil)
		currentTab.DataTable.Table.AppendCols(defaultCol)
		v.computeStatistics(currentTab)
		return
	}

	var rowDL = insyra.NewDataList()
	// 為每個現有欄位新增一個值
	for i := 0; i < colCount; i++ {
		rowDL.Append(nil)
	}

	// 將新行添加到 DataTable
	currentTab.DataTable.Table.AppendRowsFromDataList(rowDL)
	// 重新計算統計數據
	v.computeStatistics(currentTab)
}

// computeStatistics 計算統計數據
func (v *DataView) computeStatistics(tab *DataTabInfo) {
	if tab.DataTable == nil || tab.DataTable.Table == nil {
		return
	}

	insyraTable := tab.DataTable.Table
	rowCount, colCount := insyraTable.Size()

	// 創建有序的統計數據 map（使用固定順序確保顯示一致性）
	tab.StatsData = make(map[string]string)

	if rowCount == 0 || colCount == 0 {
		tab.StatsData["總行數"] = "0"
		tab.StatsData["總欄數"] = "0"
		return
	}

	// 基本統計（按固定順序）
	tab.StatsData["總行數"] = strconv.Itoa(rowCount)
	tab.StatsData["總欄數"] = strconv.Itoa(colCount)
	// 計算數值欄位的額外統計
	var numericCols int
	for i := 0; i < colCount; i++ {
		colData := insyraTable.GetColByNumber(i)
		if colData != nil {
			// 檢查是否為數值欄位
			hasNumeric := false
			for j := 0; j < colData.Len() && j < 10; j++ { // 檢查前10個值
				if val := colData.Get(j); val != nil {
					if _, err := strconv.ParseFloat(fmt.Sprint(val), 64); err == nil {
						hasNumeric = true
						break
					}
				}
			}
			if hasNumeric {
				numericCols++
			}
		}
	}

	if numericCols > 0 {
		tab.StatsData["數值欄數"] = strconv.Itoa(numericCols)
	}
}

// LoadSampleData 載入樣本數據到當前標籤頁
func (v *DataView) LoadSampleData() {
	if len(v.tabs) == 0 {
		return
	}

	currentTab := v.tabs[v.currentTabIndex]
	v.loadSampleDataToTab(currentTab)
}

// loadSampleDataToTab 載入樣本數據到指定標籤頁
func (v *DataView) loadSampleDataToTab(tab *DataTabInfo) {
	// 創建新的 DataTable
	tab.DataTable.Table = insyra.NewDataTable()
	tab.DataTable.ResetEditors()

	// 創建列數據
	idCol := insyra.NewDataList("c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8", "c9", "c10").SetName("ID")
	nameCol := insyra.NewDataList("張小明", "李美玲", "王大偉", "陳小華", "林志明", "黃雅琪", "吳建宏", "趙小惠", "劉大為", "鄭美華").SetName("姓名")
	ageCol := insyra.NewDataList("28", "32", "45", "24", "36", "29", "41", "27", "52", "38").SetName("年齡")
	cityCol := insyra.NewDataList("台北", "台中", "高雄", "新竹", "台北", "台南", "高雄", "台中", "台北", "新北").SetName("城市")
	jobCol := insyra.NewDataList("工程師", "設計師", "經理", "研究員", "醫生", "老師", "建築師", "護士", "律師", "會計師").SetName("職業")
	salaryCol := insyra.NewDataList("85000", "78000", "120000", "76000", "160000", "72000", "95000", "68000", "130000", "92000").SetName("收入")

	// 將列添加到 DataTable
	tab.DataTable.Table.AppendCols(idCol, nameCol, ageCol, cityCol, jobCol, salaryCol)
	tab.DataTable.Table.SetColToRowNames("A")

	// 計算統計數據
	v.computeStatistics(tab)
}

// SetViewController 設置視圖控制器參考
func (v *DataView) SetViewController(controller *ViewController) {
	v.viewController = controller
}

// GetDataTable 獲取當前標籤頁的 DataTable 組件
func (v *DataView) GetDataTable() *GenericDataTable {
	if len(v.tabs) == 0 {
		return nil
	}
	return v.tabs[v.currentTabIndex].DataTable
}

// AddDataFromInsyraTable 從 insyra DataTable 添加數據到當前標籤頁
func (v *DataView) AddDataFromInsyraTable(table *insyra.DataTable) {
	if len(v.tabs) == 0 {
		return
	}

	currentTab := v.tabs[v.currentTabIndex]
	currentTab.DataTable.Table = table
	currentTab.DataTable.ResetEditors()

	// 重新計算統計數據
	v.computeStatistics(currentTab)
}

// Update 實現視圖更新
func (v *DataView) Update(e interface{}) {
	// 事件處理
}

// Event 實現事件處理
func (v *DataView) Event(e interface{}) {
	// 事件處理將在 Layout 中完成
}

// layoutColumnInput 繪製計算欄輸入區域
func (v *DataView) layoutColumnInput(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// 主要輸入行
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(4), Bottom: unit.Dp(4),
				Left: unit.Dp(8), Right: unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					// fx 標籤
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						fxLabel := material.Body2(th, "fx")
						fxLabel.Color = color.NRGBA{R: 90, G: 90, B: 90, A: 255}
						return fxLabel.Layout(gtx)
					}),

					// 名稱輸入框
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Max.X = gtx.Dp(100)
							editor := material.Editor(th, &v.columnNameEditor, "名稱")
							editor.TextSize = unit.Sp(14)
							return editor.Layout(gtx)
						})
					}),

					// 等號
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						equalLabel := material.Body1(th, "=")
						return equalLabel.Layout(gtx)
					}),

					// CCL 表達式輸入框
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							editor := material.Editor(th, &v.columnFormulaEditor, "CCL 表達式")
							editor.TextSize = unit.Sp(14)
							return editor.Layout(gtx)
						})
					}),

					// 確認按鈕
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &v.addColumnConfirmBtn, "✓")
						btn.Background = color.NRGBA{R: 0, G: 150, B: 0, A: 255}
						btn.TextSize = unit.Sp(12)
						return btn.Layout(gtx)
					}),

					// 取消按鈕
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &v.addColumnCancelBtn, "✕")
							btn.Background = color.NRGBA{R: 150, G: 0, B: 0, A: 255}
							btn.TextSize = unit.Sp(12)
							return btn.Layout(gtx)
						})
					}),
				)
			})
		}),
		// 錯誤訊息
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !v.showError || v.errorMessage == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{
				Top: unit.Dp(2), Left: unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				errorLabel := material.Caption(th, v.errorMessage)
				errorLabel.Color = color.NRGBA{R: 200, G: 30, B: 30, A: 255}
				return errorLabel.Layout(gtx)
			})
		}),
	)
}

// addCalculatedColumn 添加計算欄
func (v *DataView) addCalculatedColumn() {
	if len(v.tabs) == 0 {
		return
	}

	currentTab := v.tabs[v.currentTabIndex]
	if currentTab.DataTable == nil || currentTab.DataTable.Table == nil {
		return
	}

	// 獲取輸入的 CCL 表達式和欄位名稱
	formula := v.columnFormulaEditor.Text()
	colName := v.columnNameEditor.Text()

	// 如果輸入有效，添加新的計算欄
	if formula != "" && colName != "" {
		// 使用 AddColUsingCCL 方法添加計算欄
		currentTab.DataTable.Table.AddColUsingCCL(colName, formula)

		// 檢查是否有錯誤發生
		_, errMsg := insyra.PopError(insyra.ErrPoppingModeFIFO)
		if errMsg != "" {
			// 如果有錯誤，顯示錯誤訊息但不關閉輸入面板
			v.errorMessage = "計算欄錯誤: " + errMsg
			v.showError = true
		} else {
			// 如果沒有錯誤，更新表格顯示並清空輸入
			currentTab.DataTable.Table.Show()
			v.columnFormulaEditor.SetText("")
			v.columnNameEditor.SetText("")
			v.showColumnInput = false
			v.showError = false

			// 重新計算統計數據
			v.computeStatistics(currentTab)
		}
	} else {
		v.errorMessage = "請輸入欄位名稱與 CCL 表達式"
		v.showError = true
	}
}

// cancelColumnInput 取消計算欄輸入
func (v *DataView) cancelColumnInput() {
	// 清空輸入並隱藏輸入面板
	v.columnFormulaEditor.SetText("")
	v.columnNameEditor.SetText("")
	v.showColumnInput = false
	v.showError = false // 重設錯誤訊息狀態
	v.errorMessage = ""
}
