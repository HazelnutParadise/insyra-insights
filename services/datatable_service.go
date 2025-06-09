package services

import (
	"context"
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
	var cellValue any
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
	dt.AddColUsingCCL(columnName, formula)
	// todo: 處理錯誤
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
	for i := range colCount {
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

	// 使用 AddColUsingCCL 方法來執行 CCL 公式並新增欄位
	dt.AddColUsingCCL(columnName, formula)
	// todo: 處理錯誤

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

// ===== 專案檔案操作 =====

// 專案狀態管理
type ProjectState struct {
	currentFilePath   string
	hasUnsavedChanges bool
}

var projectState = &ProjectState{
	currentFilePath:   "",
	hasUnsavedChanges: false,
}

// SaveProject 儲存整個專案（所有標籤頁）
func (s *DataTableService) SaveProject(filePath string) bool {
	// TODO: 實現專案檔案儲存
	// 應該保存所有 dataTables 到一個專案檔案中
	// 格式可以是 JSON 或自定義格式

	// 暫時實現：將所有表格保存到一個 JSON 檔案
	projectData := map[string]any{
		"version": "1.0",
		"tables":  make([]map[string]any, 0),
	}
	for i, dt := range s.dataTables {
		if dt != nil {
			rowCount, colCount := dt.Size()
			tableData := map[string]any{
				"id":       i,
				"name":     dt.GetName(),
				"rowCount": rowCount,
				"colCount": colCount,
				// TODO: 添加實際的資料匯出
			}
			projectData["tables"] = append(projectData["tables"].([]map[string]any), tableData)
		}
	}

	// 這裡應該實際寫入檔案
	// 暫時返回 true
	projectState.currentFilePath = filePath
	projectState.hasUnsavedChanges = false
	return true
}

// LoadProject 載入專案檔案
func (s *DataTableService) LoadProject(filePath string) bool {
	// TODO: 實現專案檔案載入
	// 清空現有資料表，載入專案檔案中的所有表格

	projectState.currentFilePath = filePath
	projectState.hasUnsavedChanges = false
	return true
}

// ===== 資料表匯出方法 =====

// ExportTableAsCSV 將指定資料表匯出為 CSV
func (s *DataTableService) ExportTableAsCSV(tableID int, filePath string) bool {
	dt := s.getTableByID(tableID)
	if dt == nil {
		return false
	}

	// TODO: 實現 CSV 匯出
	// 使用 insyra 的 CSV 輸出功能
	return true
}

// ExportTableAsJSON 將指定資料表匯出為 JSON
func (s *DataTableService) ExportTableAsJSON(tableID int, filePath string) bool {
	dt := s.getTableByID(tableID)
	if dt == nil {
		return false
	}

	// TODO: 實現 JSON 匯出
	return true
}

// ExportTableAsExcel 將指定資料表匯出為 Excel
func (s *DataTableService) ExportTableAsExcel(tableID int, filePath string) bool {
	dt := s.getTableByID(tableID)
	if dt == nil {
		return false
	}

	// TODO: 實現 Excel 匯出
	return true
}

// ===== 專案狀態管理 =====

// HasUnsavedChanges 檢查是否有未儲存的變更
func (s *DataTableService) HasUnsavedChanges() bool {
	return projectState.hasUnsavedChanges
}

// MarkAsSaved 標記專案為已儲存狀態
func (s *DataTableService) MarkAsSaved() {
	projectState.hasUnsavedChanges = false
}

// GetCurrentProjectPath 獲取當前專案檔案路徑
func (s *DataTableService) GetCurrentProjectPath() string {
	return projectState.currentFilePath
}

// MarkAsModified 標記專案有變更（在修改資料時調用）
func (s *DataTableService) MarkAsModified() {
	projectState.hasUnsavedChanges = true
}

// ===== 檔案開啟功能 =====

// OpenCSVFile 開啟CSV檔案並創建新的資料表
func (s *DataTableService) OpenCSVFile(filePath string) int {
	// TODO: 實現CSV檔案載入功能
	// 1. 讀取CSV檔案
	// 2. 解析CSV格式
	// 3. 創建新的資料表
	// 4. 設定表格名稱
	// 5. 返回新表格的ID
	return -1
}

// OpenJSONFile 開啟JSON檔案並創建新的資料表
func (s *DataTableService) OpenJSONFile(filePath string) int {
	// TODO: 實現JSON檔案載入功能
	// 1. 讀取JSON檔案
	// 2. 解析JSON格式（支援陣列格式和物件格式）
	// 3. 創建新的資料表
	// 4. 設定表格名稱和欄位
	// 5. 返回新表格的ID
	return -1
}

// OpenSQLiteFile 開啟SQLite檔案中的指定表格
func (s *DataTableService) OpenSQLiteFile(filePath string, tableName string) int {
	// TODO: 實現SQLite檔案載入功能
	// 1. 連接SQLite資料庫
	// 2. 查詢指定表格的結構和資料
	// 3. 創建新的資料表
	// 4. 設定表格名稱和欄位
	// 5. 載入所有資料行
	// 6. 返回新表格的ID
	return -1
}

// GetSQLiteTables 取得SQLite檔案中的表格列表
func (s *DataTableService) GetSQLiteTables(filePath string) []string {
	// TODO: 實現SQLite表格列表功能
	// 1. 連接SQLite資料庫
	// 2. 查詢sqlite_master表格
	// 3. 取得所有用戶定義的表格名稱
	// 4. 返回表格名稱列表
	return []string{}
}

// OpenFileDialog 開啟檔案選擇對話框
func (s *DataTableService) OpenFileDialog(ctx context.Context, filters string) string {
	// TODO: 實現檔案選擇對話框功能
	// 1. 使用Wails runtime.OpenFileDialog
	// 2. 根據filters參數設定檔案過濾器
	// 3. 顯示原生檔案選擇對話框
	// 4. 返回使用者選擇的檔案路徑
	return ""
}
