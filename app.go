package main

import (
	"context"
	"flag"
	"insyra-insights-wails/i18n"
	"insyra-insights-wails/services"
)

// App struct
type App struct {
	ctx         context.Context
	dataService *services.DataTableService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		dataService: services.NewDataTableService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// I18n 相關方法

// GetText 獲取翻譯文字
func (a *App) GetText(key string) string {
	return i18n.T(key)
}

// SetLanguage 設定語言
func (a *App) SetLanguage(lang string) {
	i18n.SetLanguage(lang)
}

// GetCurrentLanguage 獲取當前語言
func (a *App) GetCurrentLanguage() string {
	return "zh-TW" // 預設為繁體中文
}

// GetParamValue 獲取命令行參數值
func (a *App) GetParamValue(key string) string {
	// 從命令行獲取參數
	// 這個方法允許前端訪問命令行參數
	switch key {
	case "uuid":
		uuid := flag.Lookup("uuid")
		if uuid != nil {
			return uuid.Value.String()
		}
	case "filepath":
		filepath := flag.Lookup("filepath")
		if filepath != nil {
			return filepath.Value.String()
		}
	}
	return ""
}

// LoadTable 載入資料表
func (a *App) LoadTable(tableName string, filePath string) bool {
	return a.dataService.LoadTable(tableName, filePath)
}

// CreateEmptyTable 創建空白資料表
func (a *App) CreateEmptyTable(tableName string) bool {
	return a.dataService.CreateEmptyTable(tableName)
}

// GetTableData 獲取資料表資料
func (a *App) GetTableData(tableName string) map[string]interface{} {
	return a.dataService.GetTableData(tableName)
}

// UpdateCellValue 更新儲存格值
func (a *App) UpdateCellValue(tableName string, rowIndex int, colIndex int, value string) bool {
	return a.dataService.UpdateCellValue(tableName, rowIndex, colIndex, value)
}

// UpdateColumnName 更新欄名
func (a *App) UpdateColumnName(tableName string, colIndex int, newName string) bool {
	return a.dataService.UpdateColumnName(tableName, colIndex, newName)
}

// SaveTable 保存資料表
func (a *App) SaveTable(tableName string, filePath string) bool {
	return a.dataService.SaveTable(tableName, filePath)
}

// AddColumn 新增欄位
func (a *App) AddColumn(tableName string, columnName string) bool {
	return a.dataService.AddColumn(tableName, columnName)
}

// AddRow 新增行
func (a *App) AddRow(tableName string) bool {
	return a.dataService.AddRow(tableName)
}

// AddCalculatedColumn 新增計算欄位
func (a *App) AddCalculatedColumn(tableName string, columnName string, formula string) bool {
	return a.dataService.AddCalculatedColumn(tableName, columnName, formula)
}

// GetTableNames 獲取所有表格名稱
func (a *App) GetTableNames() []string {
	return a.dataService.GetTableNames()
}

// RemoveTable 移除指定名稱的表格
func (a *App) RemoveTable(tableName string) bool {
	return a.dataService.RemoveTable(tableName)
}

// ===== 基於 ID 的操作方法 =====

// LoadTableByID 加載資料表到指定ID位置
func (a *App) LoadTableByID(tableID int, tableName string, filePath string) int {
	return a.dataService.LoadTableByID(tableID, tableName, filePath)
}

// CreateEmptyTableByID 在指定ID位置創建空白資料表
func (a *App) CreateEmptyTableByID(tableID int, tableName string) int {
	return a.dataService.CreateEmptyTableByID(tableID, tableName)
}

// GetTableDataByID 根據ID獲取資料表資料
func (a *App) GetTableDataByID(tableID int) map[string]interface{} {
	return a.dataService.GetTableDataByID(tableID)
}

// UpdateCellValueByID 根據ID更新儲存格值
func (a *App) UpdateCellValueByID(tableID int, rowIndex int, colIndex int, value string) bool {
	return a.dataService.UpdateCellValueByID(tableID, rowIndex, colIndex, value)
}

// UpdateColumnNameByID 根據ID更新欄名
func (a *App) UpdateColumnNameByID(tableID int, colIndex int, newName string) bool {
	return a.dataService.UpdateColumnNameByID(tableID, colIndex, newName)
}

// SaveTableByID 根據ID保存資料表
func (a *App) SaveTableByID(tableID int, filePath string) bool {
	return a.dataService.SaveTableByID(tableID, filePath)
}

// AddColumnByID 根據ID新增欄位
func (a *App) AddColumnByID(tableID int, columnName string) bool {
	return a.dataService.AddColumnByID(tableID, columnName)
}

// AddRowByID 根據ID新增行
func (a *App) AddRowByID(tableID int) bool {
	return a.dataService.AddRowByID(tableID)
}

// AddCalculatedColumnByID 根據ID新增計算欄位
func (a *App) AddCalculatedColumnByID(tableID int, columnName string, formula string) bool {
	return a.dataService.AddCalculatedColumnByID(tableID, columnName, formula)
}

// GetTableCount 獲取表格總數
func (a *App) GetTableCount() int {
	return a.dataService.GetTableCount()
}

// GetTableInfo 獲取指定ID表格的基本信息
func (a *App) GetTableInfo(tableID int) map[string]interface{} {
	return a.dataService.GetTableInfo(tableID)
}

// RemoveTableByID 根據ID移除表格
func (a *App) RemoveTableByID(tableID int) bool {
	return a.dataService.RemoveTableByID(tableID)
}

// ===== 專案檔案操作方法 =====

// SaveProject 儲存整個專案（所有標籤頁）
func (a *App) SaveProject(filePath string) bool {
	return a.dataService.SaveProject(filePath)
}

// SaveProjectAs 另存新檔（所有標籤頁）
func (a *App) SaveProjectAs(filePath string) bool {
	return a.dataService.SaveProject(filePath)
}

// LoadProject 載入專案檔案
func (a *App) LoadProject(filePath string) bool {
	return a.dataService.LoadProject(filePath)
}

// ===== 資料表匯出方法 =====

// ExportTableAsCSV 將指定資料表匯出為 CSV
func (a *App) ExportTableAsCSV(tableID int, filePath string) bool {
	return a.dataService.ExportTableAsCSV(tableID, filePath)
}

// ExportTableAsJSON 將指定資料表匯出為 JSON
func (a *App) ExportTableAsJSON(tableID int, filePath string) bool {
	return a.dataService.ExportTableAsJSON(tableID, filePath)
}

// ExportTableAsExcel 將指定資料表匯出為 Excel
func (a *App) ExportTableAsExcel(tableID int, filePath string) bool {
	return a.dataService.ExportTableAsExcel(tableID, filePath)
}

// ===== 專案狀態管理 =====

// HasUnsavedChanges 檢查是否有未儲存的變更
func (a *App) HasUnsavedChanges() bool {
	return a.dataService.HasUnsavedChanges()
}

// MarkAsSaved 標記專案為已儲存狀態
func (a *App) MarkAsSaved() {
	a.dataService.MarkAsSaved()
}

// GetCurrentProjectPath 獲取當前專案檔案路徑
func (a *App) GetCurrentProjectPath() string {
	return a.dataService.GetCurrentProjectPath()
}
