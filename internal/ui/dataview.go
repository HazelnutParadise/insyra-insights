package ui

import (
	"image/color"
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/HazelnutParadise/insyra"
)

// Tab 代表一個標籤頁
type Tab int

const (
	// TabData 代表數據表格標籤頁
	TabData Tab = iota
	// TabStats 代表統計分析標籤頁
	TabStats
	// TabCharts 代表圖表標籤頁
	TabCharts
)

// DataView 數據視圖
type DataView struct {
	// 標籤頁控制
	tabs       *widget.Enum
	currentTab Tab

	// insyra DataTable 組件
	dataTable *DataTable

	// 統計數據
	statsData map[string]string

	// 圖表相關
	chartButton widget.Clickable

	// 篩選和搜索
	filterEditor widget.Editor
	searchButton widget.Clickable

	// 返回按鈕
	backButton widget.Clickable

	// 導出按鈕
	exportButton widget.Clickable

	// 樣本數據
	sampleData *SampleData

	// 最後事件處理
	lastEvent interface{}

	// 視圖控制器參考
	viewController *ViewController

	// 搜索結果
	searchResults []int
	currentSearchIndex int
}

// NewDataView 創建一個新的數據視圖
func NewDataView() *DataView {
	// 初始化標籤頁控制
	tabs := new(widget.Enum)
	tabs.Value = "data" // 默認顯示數據表格

	// 初始化編輯器
	filterEditor := widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	filterEditor.SetText("") // 默認為空

	// 創建 DataTable 組件
	dataTable := NewDataTable()

	// 初始化數據視圖
	view := &DataView{
		tabs:         tabs,
		dataTable:    dataTable,
		filterEditor: filterEditor,
		statsData:    make(map[string]string),
	}

	// 生成樣本數據
	view.loadSampleData()

	return view
}

// loadSampleData 載入樣本數據
func (v *DataView) loadSampleData() {
	// 在實際應用中，這裡應該是從檔案或資料庫中載入的真實數據
	// 這裡只是為了示範
	v.sampleData = &SampleData{
		Headers: []string{"ID", "姓名", "年齡", "城市", "職業", "收入"},
		Rows: [][]string{
			{"1", "張小明", "28", "台北", "工程師", "85000"},
			{"2", "李美玲", "32", "台中", "設計師", "78000"},
			{"3", "王大偉", "45", "高雄", "經理", "120000"},
			{"4", "陳小華", "24", "新竹", "研究員", "76000"},
			{"5", "林志明", "36", "台北", "醫生", "160000"},
			{"6", "黃雅琪", "29", "台南", "老師", "72000"},
			{"7", "吳建宏", "41", "高雄", "建築師", "95000"},
			{"8", "趙小惠", "27", "台中", "護士", "68000"},
			{"9", "劉大為", "52", "台北", "律師", "130000"},
			{"10", "鄭美華", "38", "新北", "會計師", "92000"},
		},
	}

	// 設置 DataTable 數據
	v.dataTable.SetData(v.sampleData.Rows, v.sampleData.Headers)

	// 計算統計數據
	v.computeStatistics()
}

// computeStatistics 計算統計數據
func (v *DataView) computeStatistics() {
	insyraTable := v.dataTable.GetInsyraDataTable()
	rowCount, _ := insyraTable.Size()
	
	if rowCount == 0 {
		return
	}

	// 計算平均年齡 (假設年齡在第3列，索引為2)
	totalAge := 0
	ageCount := 0
	for i := 0; i < rowCount; i++ {
		ageValue := insyraTable.GetElementByNumberIndex(i, 2) // 年齡列
		if ageStr, ok := ageValue.(string); ok {
			if age, err := strconv.Atoi(ageStr); err == nil {
				totalAge += age
				ageCount++
			}
		}
	}
	if ageCount > 0 {
		avgAge := float64(totalAge) / float64(ageCount)
		v.statsData["平均年齡"] = strconv.FormatFloat(avgAge, 'f', 1, 64)
	}

	// 計算平均收入 (假設收入在第6列，索引為5)
	totalIncome := 0
	incomeCount := 0
	for i := 0; i < rowCount; i++ {
		incomeValue := insyraTable.GetElementByNumberIndex(i, 5) // 收入列
		if incomeStr, ok := incomeValue.(string); ok {
			if income, err := strconv.Atoi(incomeStr); err == nil {
				totalIncome += income
				incomeCount++
			}
		}
	}
	if incomeCount > 0 {
		avgIncome := float64(totalIncome) / float64(incomeCount)
		v.statsData["平均收入"] = strconv.FormatFloat(avgIncome, 'f', 0, 64)
	}

	// 計算城市分布 (假設城市在第4列，索引為3)
	cityCount := make(map[string]int)
	for i := 0; i < rowCount; i++ {
		cityValue := insyraTable.GetElementByNumberIndex(i, 3) // 城市列
		if city, ok := cityValue.(string); ok {
			cityCount[city]++
		}
	}

	cities := []string{}
	for city, count := range cityCount {
		cities = append(cities, city+": "+strconv.Itoa(count)+"人")
	}
	v.statsData["城市分布"] = strings.Join(cities, ", ")

	// 計算職業分布 (假設職業在第5列，索引為4)
	jobCount := make(map[string]int)
	for i := 0; i < rowCount; i++ {
		jobValue := insyraTable.GetElementByNumberIndex(i, 4) // 職業列
		if job, ok := jobValue.(string); ok {
			jobCount[job]++
		}
	}

	jobs := []string{}
	for job, count := range jobCount {
		jobs = append(jobs, job+": "+strconv.Itoa(count)+"人")
	}
	v.statsData["職業分布"] = strings.Join(jobs, ", ")
}

// Layout 實現視圖布局
func (v *DataView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 處理按鈕點擊事件
	v.processEvents(gtx)

	// 獲取當前選擇的標籤頁
	switch v.tabs.Value {
	case "data":
		v.currentTab = TabData
	case "stats":
		v.currentTab = TabStats
	case "charts":
		v.currentTab = TabCharts
	}

	// 主布局
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Start,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// 頂部工具欄
			return v.layoutToolbar(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// 標籤頁選擇器
			return v.layoutTabs(gtx, th)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// 基於選擇的標籤頁顯示相應內容
			switch v.currentTab {
			case TabData:
				return v.layoutDataTable(gtx, th)
			case TabStats:
				return v.layoutStats(gtx, th)
			case TabCharts:
				return v.layoutCharts(gtx, th)
			default:
				return layout.Dimensions{}
			}
		}),
	)
}

// layoutToolbar 顯示工具欄
func (v *DataView) layoutToolbar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 返回按鈕
	backBtn := material.Button(th, &v.backButton, "返回")
	// 搜索輸入框和按鈕
	searchEditor := material.Editor(th, &v.filterEditor, "搜尋...")
	searchBtn := material.Button(th, &v.searchButton, "搜尋")
	// 匯出按鈕
	exportBtn := material.Button(th, &v.exportButton, "匯出")

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, backBtn.Layout)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, searchEditor.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, searchBtn.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, exportBtn.Layout)
		}),
	)
}

// layoutTabs 顯示標籤頁
func (v *DataView) layoutTabs(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 標籤頁按鈕
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			dataTab := material.RadioButton(th, v.tabs, "data", "資料表格")
			if v.tabs.Value == "data" {
				dataTab.Color = color.NRGBA{R: 63, G: 81, B: 181, A: 255} // 靛藍色
			}
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, dataTab.Layout)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			statsTab := material.RadioButton(th, v.tabs, "stats", "統計分析")
			if v.tabs.Value == "stats" {
				statsTab.Color = color.NRGBA{R: 63, G: 81, B: 181, A: 255} // 靛藍色
			}
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, statsTab.Layout)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			chartsTab := material.RadioButton(th, v.tabs, "charts", "圖表視覺化")
			if v.tabs.Value == "charts" {
				chartsTab.Color = color.NRGBA{R: 63, G: 81, B: 181, A: 255} // 靛藍色
			}
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, chartsTab.Layout)
		}),
	)
}

// layoutDataTable 顯示數據表格
func (v *DataView) layoutDataTable(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return v.dataTable.Layout(gtx, th)
}

// layoutStats 顯示統計分析
func (v *DataView) layoutStats(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H5(th, "基本統計數據")
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, title.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.layoutStatItem(gtx, th, "平均年齡", v.statsData["平均年齡"]+" 歲")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.layoutStatItem(gtx, th, "平均收入", v.statsData["平均收入"]+" 元")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.layoutStatItem(gtx, th, "城市分布", v.statsData["城市分布"])
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.layoutStatItem(gtx, th, "職業分布", v.statsData["職業分布"])
		}),
	)
}

// layoutStatItem 顯示單個統計項目
func (v *DataView) layoutStatItem(gtx layout.Context, th *material.Theme, label, value string) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(8),
		Bottom: unit.Dp(8),
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		labelText := material.Body1(th, label+":")
		labelText.Font.Weight = font.Bold

		valueText := material.Body1(th, value)

		return layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(labelText.Layout),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.E.Layout(gtx, valueText.Layout)
			}),
		)
	})
}

// layoutCharts 顯示圖表視覺化
func (v *DataView) layoutCharts(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			msg := material.H5(th, "圖表功能開發中")
			return layout.Center.Layout(gtx, msg.Layout)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			desc := material.Body1(th, "此功能將在未來版本中提供。")
			return layout.Center.Layout(gtx, desc.Layout)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(40)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &v.chartButton, "生成範例圖表")
			return layout.Center.Layout(gtx, btn.Layout)
		}),
	)
}

// processEvents 處理按鈕點擊事件
func (v *DataView) processEvents(gtx layout.Context) {
	// 返回按鈕點擊
	if v.backButton.Clicked(gtx) {
		// TODO: 實現返回功能
		if v.viewController != nil {
			v.viewController.SwitchView("home")
		}
	}

	// 搜索按鈕點擊
	if v.searchButton.Clicked(gtx) {
		query := v.filterEditor.Text()
		if query != "" {
			// 使用 DataTable 的搜索功能
			v.searchResults = v.dataTable.SearchInTable(query)
			v.currentSearchIndex = 0
			
			// 如果有搜索結果，高亮第一個
			if len(v.searchResults) > 0 {
				v.dataTable.SetSelectedRow(v.searchResults[0])
			}
		}
	}

	// 匯出按鈕點擊
	if v.exportButton.Clicked(gtx) {
		// 使用 DataTable 的匯出功能
		err := v.dataTable.ExportToCSV("export_data.csv")
		if err != nil {
			// TODO: 顯示錯誤訊息
		}
	}

	// 圖表按鈕點擊
	if v.chartButton.Clicked(gtx) {
		// TODO: 實現圖表生成功能
	}
}

// SetViewController 設置視圖控制器參考
func (v *DataView) SetViewController(controller *ViewController) {
	v.viewController = controller
}

// AddDataFromInsyraTable 從 insyra DataTable 添加數據
func (v *DataView) AddDataFromInsyraTable(table *insyra.DataTable) {
	v.dataTable.LoadFromInsyraDataTable(table)
	v.computeStatistics()
}

// GetDataTable 獲取 DataTable 組件
func (v *DataView) GetDataTable() *DataTable {
	return v.dataTable
}

// Update 實現視圖更新
func (v *DataView) Update(e interface{}) {
	v.lastEvent = e
}

// Event 實現事件處理
func (v *DataView) Event(e interface{}) {
	// 事件處理將在 Layout 中完成
}
