package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/HazelnutParadise/insyra"
)

// DataTable 是簡化的數據表格組件
type DataTable struct {
	// 數據存儲在這裡
	headers []string
	rows    [][]string
	
	// insyra DataTable 用於進階功能
	insyraTable *insyra.DataTable

	// UI 組件
	list *widget.List

	// 選中狀態
	selectedRow int

	// 樣式設定
	cellPadding   unit.Dp
	rowHeight     unit.Dp
	headerHeight  unit.Dp
	borderColor   color.NRGBA
	headerBgColor color.NRGBA
	selectedColor color.NRGBA
	alternateColor color.NRGBA
}

// NewDataTable 創建一個新的 DataTable 組件
func NewDataTable() *DataTable {
	// 創建垂直列表用於行
	list := &widget.List{}
	list.Axis = layout.Vertical

	// 創建水平列表用於表頭
	headerList := &widget.List{}
	headerList.Axis = layout.Horizontal

	// 創建滾動條
	scrollBar := &widget.Scrollbar{}

	return &DataTable{
		table:          insyra.NewDataTable(),
		list:           list,
		headerList:     headerList,
		scrollBar:      scrollBar,
		selectedRow:    -1,
		selectedColumn: "",
		cellPadding:    unit.Dp(8),
		headerHeight:   unit.Dp(40),
		rowHeight:      unit.Dp(32),
		borderColor:    color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		headerBgColor:  color.NRGBA{R: 240, G: 240, B: 240, A: 255},
		selectedColor:  color.NRGBA{R: 63, G: 81, B: 181, A: 50},
		alternateColor: color.NRGBA{R: 248, G: 248, B: 248, A: 255},
		pageSize:       50,
		currentPage:    0,
		showRowNames:   true,
	}
}

// SetData 設置表格數據
func (dt *DataTable) SetData(data [][]string, headers []string) {
	// 清空現有數據
	dt.table = insyra.NewDataTable()
	// 如果有標題，先添加列
	if len(headers) > 0 {
		for _, header := range headers {
			col := insyra.NewDataList()
			col.SetName(header)
			// 為每列添加空數據以匹配行數
			for range data {
				col.Append(nil)
			}
			dt.table.AppendCols(col)
		}
	}

	// 添加數據行
	for _, row := range data {
		rowData := make(map[string]any)
		for i, cell := range row {
			if i < len(headers) {
				rowData[headers[i]] = cell
			} else {
				// 如果沒有足夠的標題，使用默認列名
				colIndex := generateColumnIndex(i)
				rowData[colIndex] = cell
			}
		}
		dt.table.AppendRowsByColName(rowData)
	}
}

// LoadFromInsyraDataTable 從 insyra DataTable 載入數據
func (dt *DataTable) LoadFromInsyraDataTable(srcTable *insyra.DataTable) {
	dt.table = srcTable
}

// GetInsyraDataTable 獲取底層的 insyra DataTable
func (dt *DataTable) GetInsyraDataTable() *insyra.DataTable {
	return dt.table
}

// Layout 實現表格的佈局渲染
func (dt *DataTable) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if dt.table == nil {
		return layout.Dimensions{}
	}

	rowCount, colCount := dt.table.Size()
	if rowCount == 0 || colCount == 0 {
		return dt.layoutEmpty(gtx, th)
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// 表頭
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return dt.layoutHeader(gtx, th)
		}),
		// 表格內容
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return dt.layoutTable(gtx, th)
		}),
	)
}

// layoutEmpty 顯示空表格狀態
func (dt *DataTable) layoutEmpty(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.Body1(th, "沒有數據")
		label.Color = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
		return label.Layout(gtx)
	})
}

// layoutHeader 顯示表格標題行
func (dt *DataTable) layoutHeader(gtx layout.Context, th *material.Theme) layout.Dimensions {
	_, colCount := dt.table.Size()
	
	// 繪製標題背景
	paint.ColorOp{Color: dt.headerBgColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.UniformInset(dt.cellPadding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx, dt.createHeaderItems(gtx, th, colCount)...)
	})
}

// createHeaderItems 創建標題項目
func (dt *DataTable) createHeaderItems(gtx layout.Context, th *material.Theme, colCount int) []layout.FlexChild {
	items := make([]layout.FlexChild, 0, colCount)

	// 如果顯示行名，添加行名列標題
	if dt.showRowNames {
		items = append(items, layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(th, "行")
			label.Font.Weight = font.Bold
			return label.Layout(gtx)
		}))
	}

	// 添加數據列標題
	for i := 0; i < colCount; i++ {
		colIndex := i
		col := dt.table.GetColByNumber(colIndex)
		headerText := col.GetName()
		if headerText == "" {
			headerText = generateColumnIndex(colIndex)
		}

		items = append(items, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(th, headerText)
			label.Font.Weight = font.Bold
			return label.Layout(gtx)
		}))
	}

	return items
}

// layoutTable 顯示表格內容
func (dt *DataTable) layoutTable(gtx layout.Context, th *material.Theme) layout.Dimensions {
	rowCount, colCount := dt.table.Size()

	return material.List(th, dt.list).Layout(gtx, rowCount, func(gtx layout.Context, rowIndex int) layout.Dimensions {
		return dt.layoutRow(gtx, th, rowIndex, colCount)
	})
}

// layoutRow 顯示單行數據
func (dt *DataTable) layoutRow(gtx layout.Context, th *material.Theme, rowIndex, colCount int) layout.Dimensions {
	// 交替行顏色
	var bgColor color.NRGBA
	if rowIndex%2 == 1 {
		bgColor = dt.alternateColor
	} else {
		bgColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	}

	// 選中行高亮
	if rowIndex == dt.selectedRow {
		bgColor = dt.selectedColor
	}

	// 繪製行背景
	paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.UniformInset(dt.cellPadding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx, dt.createRowItems(gtx, th, rowIndex, colCount)...)
	})
}

// createRowItems 創建行項目
func (dt *DataTable) createRowItems(gtx layout.Context, th *material.Theme, rowIndex, colCount int) []layout.FlexChild {
	items := make([]layout.FlexChild, 0, colCount+1)

	// 如果顯示行名，添加行名
	if dt.showRowNames {
		rowName := dt.table.GetRowNameByIndex(rowIndex)
		if rowName == "" {
			rowName = strconv.Itoa(rowIndex)
		}

		items = append(items, layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
			label := material.Body2(th, rowName)
			label.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
			return label.Layout(gtx)
		}))
	}

	// 添加數據列
	for colIndex := 0; colIndex < colCount; colIndex++ {
		value := dt.table.GetElementByNumberIndex(rowIndex, colIndex)
		cellText := formatCellValue(value)

		items = append(items, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			label := material.Body2(th, cellText)
			return label.Layout(gtx)
		}))
	}

	return items
}

// SetSelectedRow 設置選中的行
func (dt *DataTable) SetSelectedRow(rowIndex int) {
	dt.selectedRow = rowIndex
}

// GetSelectedRow 獲取選中的行
func (dt *DataTable) GetSelectedRow() int {
	return dt.selectedRow
}

// SetSelectedColumn 設置選中的列
func (dt *DataTable) SetSelectedColumn(columnIndex string) {
	dt.selectedColumn = columnIndex
}

// GetSelectedColumn 獲取選中的列
func (dt *DataTable) GetSelectedColumn() string {
	return dt.selectedColumn
}

// SetShowRowNames 設置是否顯示行名
func (dt *DataTable) SetShowRowNames(show bool) {
	dt.showRowNames = show
}

// GetRowCount 獲取行數
func (dt *DataTable) GetRowCount() int {
	rowCount, _ := dt.table.Size()
	return rowCount
}

// GetColumnCount 獲取列數
func (dt *DataTable) GetColumnCount() int {
	_, colCount := dt.table.Size()
	return colCount
}

// FilterTable 過濾表格數據
func (dt *DataTable) FilterTable(filterFunc func(rowIndex int, columnIndex string, value any) bool) *DataTable {
	filtered := dt.table.Filter(func(rowIndex int, columnIndex string, value any) bool {
		return filterFunc(rowIndex, columnIndex, value)
	})

	newDt := NewDataTable()
	newDt.LoadFromInsyraDataTable(filtered)
	return newDt
}

// SearchInTable 在表格中搜索
func (dt *DataTable) SearchInTable(searchText string) []int {
	rowCount, colCount := dt.table.Size()
	matchingRows := make([]int, 0)

	for i := 0; i < rowCount; i++ {
		for j := 0; j < colCount; j++ {
			value := dt.table.GetElementByNumberIndex(i, j)
			cellText := formatCellValue(value)
			
			if containsIgnoreCase(cellText, searchText) {
				matchingRows = append(matchingRows, i)
				break // 找到匹配項就跳出內層循環
			}
		}
	}

	return matchingRows
}

// AddCalculatedColumn 使用 CCL (Column Calculation Language) 添加計算列
func (dt *DataTable) AddCalculatedColumn(columnName, formula string) error {
	if dt.table == nil {
		return fmt.Errorf("數據表未初始化")
	}
	
	// 使用 insyra 的 CCL 功能添加計算列
	err := dt.table.AddColUsingCCL(columnName, formula)
	if err != nil {
		return fmt.Errorf("添加計算列失敗: %v", err)
	}
	
	return nil
}

// GetColumnStatistics 獲取指定列的統計信息
func (dt *DataTable) GetColumnStatistics(columnName string) (map[string]interface{}, error) {
	if dt.table == nil {
		return nil, fmt.Errorf("數據表未初始化")
	}
	
	col := dt.table.GetColByName(columnName)
	if col == nil {
		return nil, fmt.Errorf("列 '%s' 不存在", columnName)
	}
	
	stats := make(map[string]interface{})
	
	// 基本統計
	stats["count"] = col.Count()
	stats["mean"] = col.Mean()
	stats["median"] = col.Median()
	stats["mode"] = col.Mode()
	stats["std"] = col.Std()
	stats["variance"] = col.Variance()
	stats["min"] = col.Min()
	stats["max"] = col.Max()
	stats["sum"] = col.Sum()
	
	// 分位數
	stats["q1"] = col.Quantile(0.25)
	stats["q3"] = col.Quantile(0.75)
	stats["iqr"] = col.IQR()
	
	return stats, nil
}

// ExportToCSV 匯出為 CSV 格式
func (dt *DataTable) ExportToCSV(filename string) error {
	if dt.table == nil {
		return fmt.Errorf("數據表未初始化")
	}
	
	return dt.table.SaveAsCSV(filename)
}

// ExportToJSON 匯出為 JSON 格式
func (dt *DataTable) ExportToJSON(filename string) error {
	if dt.table == nil {
		return fmt.Errorf("數據表未初始化")
	}
	
	return dt.table.SaveAsJSON(filename)
}

// FilterByCondition 使用條件過濾數據
func (dt *DataTable) FilterByCondition(condition string) (*DataTable, error) {
	if dt.table == nil {
		return nil, fmt.Errorf("數據表未初始化")
	}
	
	// 使用 insyra 的過濾功能
	filteredTable := dt.table.FilterRowsWhere(condition)
	
	// 創建新的 DataTable 包裝器
	newDT := NewDataTable()
	newDT.table = filteredTable
	
	// 複製樣式設定
	newDT.cellPadding = dt.cellPadding
	newDT.headerHeight = dt.headerHeight
	newDT.rowHeight = dt.rowHeight
	newDT.borderColor = dt.borderColor
	newDT.headerBgColor = dt.headerBgColor
	newDT.selectedColor = dt.selectedColor
	newDT.alternateColor = dt.alternateColor
	newDT.pageSize = dt.pageSize
	newDT.showRowNames = dt.showRowNames
	
	return newDT, nil
}

// GetUniqueValues 獲取列的唯一值
func (dt *DataTable) GetUniqueValues(columnName string) ([]interface{}, error) {
	if dt.table == nil {
		return nil, fmt.Errorf("數據表未初始化")
	}
	
	col := dt.table.GetColByName(columnName)
	if col == nil {
		return nil, fmt.Errorf("列 '%s' 不存在", columnName)
	}
	
	return col.Unique(), nil
}

// SortByColumn 按指定列排序
func (dt *DataTable) SortByColumn(columnName string, ascending bool) error {
	if dt.table == nil {
		return fmt.Errorf("數據表未初始化")
	}
	
	if ascending {
		dt.table.SortByColAsc(columnName)
	} else {
		dt.table.SortByColDesc(columnName)
	}
	
	return nil
}

// 輔助函數

// generateColumnIndex 生成列索引（A, B, C, ...）
func generateColumnIndex(index int) string {
	if index < 26 {
		return string(rune('A' + index))
	}
	return fmt.Sprintf("Col%d", index+1)
}

// formatCellValue 格式化儲存格值
func formatCellValue(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// containsIgnoreCase 不區分大小寫的字符串包含檢查
func containsIgnoreCase(str, substr string) bool {
	// 簡單的不區分大小寫搜索實現
	if len(substr) == 0 {
		return true
	}
	if len(str) < len(substr) {
		return false
	}

	strLower := toLower(str)
	substrLower := toLower(substr)

	for i := 0; i <= len(strLower)-len(substrLower); i++ {
		if strLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

// toLower 簡單的轉小寫函數
func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}
