package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/HazelnutParadise/Go-Utils/conv"
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
	tabs            []*DataTabInfo     // 標籤頁列表
	tabButtons      []widget.Clickable // 標籤按鈕
	addTabButton    widget.Clickable   // 新增標籤按鈕
	currentTabIndex int                // 當前選中的標籤索引

	// 功能列按鈕
	addColButton widget.Clickable // 新增欄按鈕
	addRowButton widget.Clickable // 新增列按鈕

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
	}

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
}

// layoutTabBar 繪製標籤列
func (v *DataView) layoutTabBar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	children := make([]layout.FlexChild, len(v.tabs)+1)

	// 現有標籤
	for i, tab := range v.tabs {
		idx := i
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// 如果是當前選中的標籤，改變樣式
			btn := material.Button(th, &v.tabButtons[idx], tab.Name)
			if idx == v.currentTabIndex {
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
	children[len(v.tabs)] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		btn := material.Button(th, &v.addTabButton, "+")
		// 新增按鈕使用與計算欄按鈕相同的藍色樣式
		btn.Background = color.NRGBA{R: 225, G: 245, B: 254, A: 255} // 淡藍色背景
		btn.Color = color.NRGBA{R: 33, G: 150, B: 243, A: 255}       // 藍色文字
		return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
	})

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

// layoutFunctionBar 繪製功能列
func (v *DataView) layoutFunctionBar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &v.addColButton, "新增欄")
			// 使用與表格計算欄按鈕相同的藍色樣式
			btn.Background = color.NRGBA{R: 225, G: 245, B: 254, A: 255} // 淡藍色背景
			btn.Color = color.NRGBA{R: 33, G: 150, B: 243, A: 255}       // 藍色文字
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &v.addRowButton, "新增列")
			// 使用與表格計算欄按鈕相同的藍色樣式
			btn.Background = color.NRGBA{R: 225, G: 245, B: 254, A: 255} // 淡藍色背景
			btn.Color = color.NRGBA{R: 33, G: 150, B: 243, A: 255}       // 藍色文字
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
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
			return v.layoutTableArea(gtx, th, currentTab)
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
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H6(th, "資訊區")
			title.Font.Weight = font.Bold
			// 使用與表格標題相同的藍色文字
			title.Color = color.NRGBA{R: 0, G: 90, B: 180, A: 255} // 藍色
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, title.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.layoutStats(gtx, th, tab)
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
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &v.openButton, "開啟")
			// 使用與表格計算欄按鈕相同的藍色樣式
			btn.Background = color.NRGBA{R: 225, G: 245, B: 254, A: 255} // 淡藍色背景
			btn.Color = color.NRGBA{R: 33, G: 150, B: 243, A: 255}       // 藍色文字
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &v.saveButton, "存檔")
			// 使用與表格計算欄按鈕相同的藍色樣式
			btn.Background = color.NRGBA{R: 225, G: 245, B: 254, A: 255} // 淡藍色背景
			btn.Color = color.NRGBA{R: 33, G: 150, B: 243, A: 255}       // 藍色文字
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &v.exportButton, "匯出")
			// 使用與表格計算欄按鈕相同的藍色樣式
			btn.Background = color.NRGBA{R: 225, G: 245, B: 254, A: 255} // 淡藍色背景
			btn.Color = color.NRGBA{R: 33, G: 150, B: 243, A: 255}       // 藍色文字
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &v.settingsButton, "設定")
			// 使用與表格計算欄按鈕相同的藍色樣式
			btn.Background = color.NRGBA{R: 225, G: 245, B: 254, A: 255} // 淡藍色背景
			btn.Color = color.NRGBA{R: 33, G: 150, B: 243, A: 255}       // 藍色文字
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, btn.Layout)
		}),
	)
}

// addColumn 新增欄位到當前標籤頁
func (v *DataView) addColumn() {
	if len(v.tabs) == 0 {
		return
	}

	currentTab := v.tabs[v.currentTabIndex]
	if currentTab.DataTable == nil || currentTab.DataTable.Table == nil {
		return
	}

	// 獲取目前欄位數量和行數
	rowCount, colCount := currentTab.DataTable.Table.Size()

	// 新增一個欄位
	columnName := "var" + conv.ToString(colCount+1)
	newCol := insyra.NewDataList().SetName(columnName)

	// 為新欄位填入預設值
	for i := 0; i < rowCount; i++ {
		newCol.Append("")
	}

	// 如果沒有行數，至少添加一個預設值
	if rowCount == 0 {
		newCol.Append("")
	}

	currentTab.DataTable.Table.AppendCols(newCol)

	// 重新計算統計數據
	v.computeStatistics(currentTab)
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
		defaultCol.Append("")
		currentTab.DataTable.Table.AppendCols(defaultCol)
		v.computeStatistics(currentTab)
		return
	}

	var rowDL = insyra.NewDataList()
	// 為每個現有欄位新增一個值
	for i := 0; i < colCount; i++ {
		rowDL.Append("")
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

	// 清空統計數據
	tab.StatsData = make(map[string]string)

	if rowCount == 0 || colCount == 0 {
		tab.StatsData["總行數"] = "0"
		tab.StatsData["總欄數"] = "0"
		return
	}

	// 基本統計
	tab.StatsData["總行數"] = strconv.Itoa(rowCount)
	tab.StatsData["總欄數"] = strconv.Itoa(colCount)
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
