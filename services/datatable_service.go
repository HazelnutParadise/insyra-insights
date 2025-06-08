package services

import (
	"fmt"

	"slices"

	"github.com/HazelnutParadise/insyra"
)

// DataTableService 提供資料表的核心操作
type DataTableService struct {
	dataTables []*insyra.DataTable
}

// NewDataTableService 創建一個新的 DataTableService 實例
func NewDataTableService() *DataTableService {
	return &DataTableService{
		dataTables: make([]*insyra.DataTable, 0),
	}
}

// findTableByName 根據表格名稱查找表格
func (s *DataTableService) findTableByName(tableName string) *insyra.DataTable {
	for _, dt := range s.dataTables {
		if dt.GetName() == tableName {
			return dt
		}
	}
	return nil
}

// getTableByID 根據slice索引獲取表格
func (s *DataTableService) getTableByID(tableID int) *insyra.DataTable {
	if tableID < 0 || tableID >= len(s.dataTables) {
		return nil
	}
	return s.dataTables[tableID]
}

// findTableIndex 根據表格名稱查找索引
func (s *DataTableService) findTableIndex(tableName string) int {
	for i, dt := range s.dataTables {
		if dt.GetName() == tableName {
			return i
		}
	}
	return -1
}

// LoadTable 加載資料表
func (s *DataTableService) LoadTable(tableName string, filePath string) bool {
	dt := insyra.NewDataTable()

	// 使用 insyra 的 LoadFromJSON 方法載入資料
	err := dt.LoadFromJSON(filePath)
	if err != nil {
		return false
	}

	// 設定表格名稱
	dt.SetName(tableName)
	s.dataTables = append(s.dataTables, dt)
	return true
}

// CreateEmptyTable 創建一個空白資料表
func (s *DataTableService) CreateEmptyTable(tableName string) bool {
	dt := insyra.NewDataTable()
	dt.SetName(tableName)
	s.dataTables = append(s.dataTables, dt)
	return true
}

// GetTableData 獲取資料表的完整資料
func (s *DataTableService) GetTableData(tableName string) map[string]any {
	dt := s.findTableByName(tableName)
	if dt == nil {
		return nil
	}

	result := make(map[string]any)

	// 獲取表格大小
	rowCount, colCount := dt.Size()

	// 獲取所有欄名
	var columns []map[string]any
	for i := 0; i < colCount; i++ {
		col := dt.GetColByNumber(i)
		columns = append(columns, map[string]any{
			"id":   i,
			"name": col.GetName(),
		})
	}
	result["columns"] = columns
	// 獲取所有行資料
	var rows []map[string]any
	for i := range rowCount {
		row := make(map[string]any)
		row["id"] = i

		cells := make(map[string]any)
		for j := range colCount {
			// 使用 GetElementByNumberIndex 獲取每個單元格的值
			cellValue := dt.GetElementByNumberIndex(i, j)
			cells[columns[j]["name"].(string)] = cellValue
		}
		row["cells"] = cells
		rows = append(rows, row)
	}
	result["rows"] = rows

	return result
}

// UpdateCellValue 更新儲存格的值
func (s *DataTableService) UpdateCellValue(tableName string, rowIndex int, colIndex int, value string) bool {
	dt := s.findTableByName(tableName)
	if dt == nil {
		return false
	}

	// 處理特殊值：空字符串轉換為 nil
	var cellValue interface{}
	if value == "" {
		cellValue = nil
	} else {
		cellValue = value
	}

	// 使用 UpdateElement 設置單元格值
	// 需要將列索引轉換為列字母標識符 (A, B, C, ...)
	colLetter := indexToLetters(colIndex)
	dt.UpdateElement(rowIndex, colLetter, cellValue)
	return true
}

// UpdateColumnName 更新欄名
func (s *DataTableService) UpdateColumnName(tableName string, colIndex int, newName string) bool {
	dt := s.findTableByName(tableName)
	if dt == nil {
		return false
	}

	oldName := dt.GetColByNumber(colIndex).GetName()
	if oldName == newName {
		return false // 沒有變更
	}

	dt.SetColNameByNumber(colIndex, newName)
	return true
}

// SaveTable 保存資料表
func (s *DataTableService) SaveTable(tableName string, filePath string) bool {
	dt := s.findTableByName(tableName)
	if dt == nil {
		return false
	}
	// 使用 insyra 的 ToJSON 方法保存為 JSON
	err := dt.ToJSON(filePath, true) // useColNames = true
	return err == nil
}

// AddColumn 新增欄
func (s *DataTableService) AddColumn(tableName string, columnName string) bool {
	fmt.Printf("AddColumn 被調用: tableName=%s, columnName=%s\n", tableName, columnName)
	dt := s.findTableByName(tableName)
	if dt == nil {
		fmt.Printf("錯誤: 找不到 tableName=%s 的資料表\n", tableName)
		return false
	}

	fmt.Printf("資料表存在，正在新增欄位: %s\n", columnName)
	newCol := insyra.NewDataList(nil).SetName(columnName)
	dt.AppendCols(newCol)
	fmt.Printf("成功新增欄位: %s\n", columnName)
	return true
}

// AddRow 新增列
func (s *DataTableService) AddRow(tableName string) bool {
	fmt.Printf("AddRow 被調用: tableName=%s\n", tableName)
	dt := s.findTableByName(tableName)
	if dt == nil {
		fmt.Printf("錯誤: 找不到 tableName=%s 的資料表\n", tableName)
		return false
	}

	fmt.Printf("資料表存在，正在新增行\n")

	// 獲取當前資料表的欄位數量
	_, colCount := dt.Size()
	fmt.Printf("當前資料表有 %d 個欄位\n", colCount)

	// 如果沒有欄位，先創建一個預設欄位
	if colCount == 0 {
		fmt.Printf("資料表沒有欄位，先創建預設欄位\n")
		defaultCol := insyra.NewDataList(nil).SetName("Column1")
		dt.AppendCols(defaultCol)
		colCount = 1
	}

	// 創建新行資料，為每個欄位設置空值
	newRowData := make([]interface{}, colCount)
	for i := range newRowData {
		newRowData[i] = nil
	}

	newRow := insyra.NewDataList(newRowData...)
	dt.AppendRowsFromDataList(newRow)
	fmt.Printf("成功新增行\n")
	return true
}

// AddCalculatedColumn 新增計算欄位
func (s *DataTableService) AddCalculatedColumn(tableName string, columnName string, formula string) bool {
	dt := s.findTableByName(tableName)
	if dt == nil {
		return false
	}

	// TODO: 實現 CCL 表達式解析和計算
	// 目前先創建一個空的計算欄位
	rowCount, _ := dt.Size()

	// 創建計算結果欄位，暫時填充公式字符串
	newData := make([]any, rowCount)
	for i := range newData {
		newData[i] = formula // 暫時顯示公式，實際應該計算結果
	}

	newCol := insyra.NewDataList(newData...).SetName(columnName)
	dt.AppendCols(newCol)
	return true
}

// GetTableNames 獲取所有表格名稱
func (s *DataTableService) GetTableNames() []string {
	names := make([]string, len(s.dataTables))
	for i, dt := range s.dataTables {
		names[i] = dt.GetName()
	}
	return names
}

// RemoveTable 移除指定名稱的表格
func (s *DataTableService) RemoveTable(tableName string) bool {
	for i, dt := range s.dataTables {
		if dt.GetName() == tableName {
			// 從切片中移除
			s.dataTables = append(s.dataTables[:i], s.dataTables[i+1:]...)
			return true
		}
	}
	return false
}

// indexToLetters 將數字索引轉換為字母索引 (A, B, C, ..., AA, AB, ...)
func indexToLetters(index int) string {
	if index < 0 {
		return "A"
	}

	result := ""
	for index >= 0 {
		result = string(rune('A'+(index%26))) + result
		index = index/26 - 1
		if index < 0 {
			break
		}
	}
	return result
}

// ===== 基於 ID 的操作方法 =====

// LoadTableByID 加載資料表到指定位置 (如果ID超出範圍則添加到末尾)
func (s *DataTableService) LoadTableByID(tableID int, tableName string, filePath string) int {
	dt := insyra.NewDataTable()

	// 使用 insyra 的 LoadFromJSON 方法載入資料
	err := dt.LoadFromJSON(filePath)
	if err != nil {
		return -1
	}

	// 設定表格名稱
	dt.SetName(tableName)

	// 如果 tableID 有效，插入到指定位置
	if tableID >= 0 && tableID < len(s.dataTables) {
		// 插入到指定位置
		s.dataTables = append(s.dataTables[:tableID+1], s.dataTables[tableID:]...)
		s.dataTables[tableID] = dt
		return tableID
	} else {
		// 添加到末尾
		s.dataTables = append(s.dataTables, dt)
		return len(s.dataTables) - 1
	}
}

// CreateEmptyTableByID 在指定位置創建空白資料表
func (s *DataTableService) CreateEmptyTableByID(tableID int, tableName string) int {
	dt := insyra.NewDataTable()
	dt.SetName(tableName)

	// 創建一個預設的欄位以確保表格有基本結構
	defaultCol := insyra.NewDataList(nil).SetName("Column1")
	dt.AppendCols(defaultCol)

	// 如果 tableID 有效，插入到指定位置
	if tableID >= 0 && tableID < len(s.dataTables) {
		// 插入到指定位置
		s.dataTables = append(s.dataTables[:tableID+1], s.dataTables[tableID:]...)
		s.dataTables[tableID] = dt
		return tableID
	} else {
		// 添加到末尾
		s.dataTables = append(s.dataTables, dt)
		return len(s.dataTables) - 1
	}
}

// GetTableDataByID 根據ID獲取資料表的完整資料
func (s *DataTableService) GetTableDataByID(tableID int) map[string]any {
	dt := s.getTableByID(tableID)
	if dt == nil {
		return nil
	}

	result := make(map[string]any)

	// 獲取表格大小
	rowCount, colCount := dt.Size()

	// 獲取所有欄名
	var columns []map[string]any
	for i := 0; i < colCount; i++ {
		col := dt.GetColByNumber(i)
		columns = append(columns, map[string]any{
			"id":   i,
			"name": col.GetName(),
		})
	}
	result["columns"] = columns

	// 獲取所有行資料
	var rows []map[string]any
	for i := range rowCount {
		row := make(map[string]any)
		row["id"] = i

		cells := make(map[string]any)
		for j := range colCount {
			// 使用 GetElementByNumberIndex 獲取每個單元格的值
			cellValue := dt.GetElementByNumberIndex(i, j)
			cells[columns[j]["name"].(string)] = cellValue
		}
		row["cells"] = cells
		rows = append(rows, row)
	}
	result["rows"] = rows

	return result
}

// UpdateCellValueByID 根據ID更新儲存格的值
func (s *DataTableService) UpdateCellValueByID(tableID int, rowIndex int, colIndex int, value string) bool {
	dt := s.getTableByID(tableID)
	if dt == nil {
		return false
	}
	// 處理特殊值：點和空字串都轉換為 nil
	var cellValue any
	if value == "." || value == "" {
		cellValue = nil
	} else {
		cellValue = value
	}

	// 使用 UpdateElement 設置單元格值
	colLetter := indexToLetters(colIndex)
	dt.UpdateElement(rowIndex, colLetter, cellValue)
	return true
}

// UpdateColumnNameByID 根據ID更新欄名
func (s *DataTableService) UpdateColumnNameByID(tableID int, colIndex int, newName string) bool {
	dt := s.getTableByID(tableID)
	if dt == nil {
		return false
	}

	oldName := dt.GetColByNumber(colIndex).GetName()
	if oldName == newName {
		return false // 沒有變更
	}

	dt.SetColNameByNumber(colIndex, newName)
	return true
}

// SaveTableByID 根據ID保存資料表
func (s *DataTableService) SaveTableByID(tableID int, filePath string) bool {
	dt := s.getTableByID(tableID)
	if dt == nil {
		return false
	}
	// 使用 insyra 的 ToJSON 方法保存為 JSON
	err := dt.ToJSON(filePath, true) // useColNames = true
	return err == nil
}

// AddColumnByID 根據ID新增欄
func (s *DataTableService) AddColumnByID(tableID int, columnName string) bool {
	fmt.Printf("AddColumnByID 被調用: tableID=%d, columnName=%s\n", tableID, columnName)
	dt := s.getTableByID(tableID)
	if dt == nil {
		fmt.Printf("錯誤: 找不到 tableID=%d 的資料表\n", tableID)
		return false
	}

	fmt.Printf("資料表存在，正在新增欄位: %s\n", columnName)
	newCol := insyra.NewDataList(nil).SetName(columnName)
	dt.AppendCols(newCol)
	fmt.Printf("成功新增欄位: %s\n", columnName)
	return true
}

// AddRowByID 根據ID新增列
func (s *DataTableService) AddRowByID(tableID int) bool {
	fmt.Printf("AddRowByID 被調用: tableID=%d\n", tableID)
	dt := s.getTableByID(tableID)
	if dt == nil {
		fmt.Printf("錯誤: 找不到 tableID=%d 的資料表\n", tableID)
		return false
	}

	fmt.Printf("資料表存在，正在新增行\n")

	// 獲取當前資料表的欄位數量
	_, colCount := dt.Size()
	fmt.Printf("當前資料表有 %d 個欄位\n", colCount)

	// 如果沒有欄位，先創建一個預設欄位
	if colCount == 0 {
		fmt.Printf("資料表沒有欄位，先創建預設欄位\n")
		defaultCol := insyra.NewDataList(nil).SetName("Column1")
		dt.AppendCols(defaultCol)
		colCount = 1
	}

	// 創建新行資料，為每個欄位設置空值
	newRowData := make([]any, colCount)
	for i := range newRowData {
		newRowData[i] = nil
	}

	newRow := insyra.NewDataList(newRowData...)
	dt.AppendRowsFromDataList(newRow)
	fmt.Printf("成功新增行\n")
	return true
}

// AddCalculatedColumnByID 根據ID新增計算欄位
func (s *DataTableService) AddCalculatedColumnByID(tableID int, columnName string, formula string) bool {
	dt := s.getTableByID(tableID)
	if dt == nil {
		return false
	}

	// TODO: 實現 CCL 表達式解析和計算
	// 目前先創建一個空的計算欄位
	rowCount, _ := dt.Size()

	// 創建計算結果欄位，暫時填充公式字符串
	newData := make([]any, rowCount)
	for i := range newData {
		newData[i] = formula // 暫時顯示公式，實際應該計算結果
	}

	newCol := insyra.NewDataList(newData...).SetName(columnName)
	dt.AppendCols(newCol)
	return true
}

// GetTableCount 獲取表格總數
func (s *DataTableService) GetTableCount() int {
	return len(s.dataTables)
}

// GetTableInfo 獲取指定ID表格的基本信息
func (s *DataTableService) GetTableInfo(tableID int) map[string]any {
	dt := s.getTableByID(tableID)
	if dt == nil {
		return nil
	}

	rowCount, colCount := dt.Size()
	return map[string]any{
		"id":       tableID,
		"name":     dt.GetName(),
		"rowCount": rowCount,
		"colCount": colCount,
	}
}

// RemoveTableByID 根據ID移除表格
func (s *DataTableService) RemoveTableByID(tableID int) bool {
	if tableID < 0 || tableID >= len(s.dataTables) {
		return false
	}
	// 從切片中移除
	s.dataTables = slices.Delete(s.dataTables, tableID, tableID+1)
	return true
}
