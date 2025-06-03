// GenericDataTable 是一個可繪製 Insyra DataTable 的 UI 表格組件，支援雙向捲動
package ui

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/HazelnutParadise/insyra"
	"github.com/mattn/go-runewidth"
)

type GenericDataTable struct {
	Table         *insyra.DataTable
	CellWidth     unit.Dp
	CellHeight    unit.Dp
	HeaderBgColor color.NRGBA
	BorderColor   color.NRGBA

	// 編輯功能
	cellEditors     map[string]*widget.Editor    // 儲存格編輯器 (key: "row:col")
	cellClickers    map[string]*widget.Clickable // 儲存格點擊器 (key: "row:col")
	editingCell     string                       // 當前編輯的儲存格
	selectedContent string                       // 已選中格子的完整內容
	selectedCellKey string                       // 已選中格子的索引

	// 捲動控制
	verticalList   widget.List
	horizontalList widget.List
}

func NewGenericDataTable(tbl *insyra.DataTable) *GenericDataTable {
	return &GenericDataTable{
		Table:         tbl,
		CellWidth:     unit.Dp(80),
		CellHeight:    unit.Dp(32),
		HeaderBgColor: color.NRGBA{R: 240, G: 240, B: 240, A: 255},
		BorderColor:   color.NRGBA{R: 180, G: 180, B: 180, A: 255},
		cellEditors:   make(map[string]*widget.Editor),
		cellClickers:  make(map[string]*widget.Clickable),
		editingCell:   "",
	}
}

func (dt *GenericDataTable) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if dt.Table == nil {
		return layout.Dimensions{}
	}
	rows, cols := dt.Table.Size()

	// 安全檢查：如果表格為空，返回空佈局
	if rows == 0 || cols == 0 {
		return layout.Dimensions{}
	}

	// 設置垂直捲動
	dt.verticalList.Axis = layout.Vertical
	// 設置水平捲動
	dt.horizontalList.Axis = layout.Horizontal

	// 計算表格總寬度
	totalWidth := int(dt.CellWidth) * (cols + 2) // +2 for row index and name

	// 使用垂直 Flex 佈局來組合選中內容顯示和表格
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// 選中內容區域顯示
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if dt.selectedContent == "" {
				return layout.Dimensions{}
			}

			// 顯示選中儲存格的信息
			var cellInfo string
			if dt.selectedCellKey != "" {
				parts := strings.Split(dt.selectedCellKey, ":")
				if len(parts) == 2 {
					row, _ := strconv.Atoi(parts[0])
					col, _ := strconv.Atoi(parts[1])
					colLetter := indexToLetters(col)
					cellInfo = fmt.Sprintf("已選中 %s%d: ", colLetter, row+1) // 加1讓行號從1開始計數，更直觀
				}
			}

			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// 使用水平佈局分開顯示位置和內容
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Start,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// 儲存格位置標籤
						infoLabel := material.Body1(th, cellInfo)
						infoLabel.Color = color.NRGBA{0, 0, 128, 255} // 藍色
						return infoLabel.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						// 儲存格內容標籤
						contentLabel := material.Body1(th, dt.selectedContent)
						contentLabel.Color = color.NRGBA{0, 0, 0, 255} // 黑色						// 設置最小寬度為0，允許文本擴展
						contentGtx := gtx
						contentGtx.Constraints.Min.X = 0
						return contentLabel.Layout(contentGtx)
					}),
				)
			})
		}),

		// 表格區域
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// 使用嵌套的 List 來實現雙向滾動
			return material.List(th, &dt.verticalList).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
				// 內層水平滾動
				return material.List(th, &dt.horizontalList).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
					// 設置完整的表格尺寸
					gtx.Constraints.Max.X = gtx.Dp(unit.Dp(totalWidth))

					// 渲染完整的表格
					return dt.layoutFullTable(gtx, th, rows, cols)
				})
			})
		}),
	)
}

func (dt *GenericDataTable) drawColumnIndexRow(gtx layout.Context, th *material.Theme, cols int) layout.Dimensions {
	var children []layout.FlexChild
	// 只使用一個儲存格作為行標頭
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return dt.cell(gtx, th, "")
	}))
	for i := range cols {
		label := indexToLetters(i)
		currentLabel := label // 捕獲迴圈變數
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return dt.headerCell(gtx, th, currentLabel)
		}))
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func (dt *GenericDataTable) drawColumnNameRow(gtx layout.Context, th *material.Theme, cols int) layout.Dimensions {
	var children []layout.FlexChild
	// 只使用一個儲存格作為行標頭
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return dt.headerCell(gtx, th, "行/欄")
	}))
	for i := 0; i < cols; i++ {
		name := dt.Table.GetColByNumber(i).GetName()
		currentName := name // 捕獲迴圈變數
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return dt.headerCell(gtx, th, currentName)
		}))
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func (dt *GenericDataTable) drawDataRow(gtx layout.Context, th *material.Theme, row, cols int, rowName string) layout.Dimensions {
	var children []layout.FlexChild
	// 將行索引和行名稱合併在同一格內
	combinedText := fmt.Sprintf("%s: %s", strconv.Itoa(row), rowName)
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return dt.headerCell(gtx, th, combinedText)
	}))
	// 使用 Data() 方法獲取所有資料
	data := dt.Table.Data()

	for i := range cols {
		var text string

		// 使用列名作為鍵來獲取資料
		colName := dt.Table.GetColByNumber(i).GetName()
		if colData, exists := data[colName]; exists && row < len(colData) {
			el := colData[row]
			if el != nil {
				text = fmt.Sprint(colData[row])
			} else {
				text = "."
			}
		} else {
			text = "N/A"
		}

		// 為每個儲存格添加編輯功能
		cellKey := fmt.Sprintf("%d:%d", row, i)

		// 捕獲迴圈變數
		currentText := text
		currentCellKey := cellKey
		currentRow := row
		currentCol := i

		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return dt.editableCell(gtx, th, currentText, currentCellKey, currentRow, currentCol)
		}))
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func (dt *GenericDataTable) cell(gtx layout.Context, th *material.Theme, text string) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.Rect{
				Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
			}.Op())
			paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
				Min: image.Pt(gtx.Dp(dt.CellWidth)-1, 0),
				Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
			}.Op())
			return layout.Dimensions{Size: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight))}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Body2(th, text).Layout(gtx)
			})
		}),
	)
}

func (dt *GenericDataTable) headerCell(gtx layout.Context, th *material.Theme, text string) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			paint.FillShape(gtx.Ops, dt.HeaderBgColor, clip.Rect{
				Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
			}.Op())
			paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
				Min: image.Pt(gtx.Dp(dt.CellWidth)-1, 0),
				Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
			}.Op())
			return layout.Dimensions{Size: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight))}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(th, text)
				return lbl.Layout(gtx)
			})
		}),
	)
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

// editableCell 創建可編輯的儲存格
func (dt *GenericDataTable) editableCell(gtx layout.Context, th *material.Theme, text, cellKey string, row, col int) layout.Dimensions {
	// 建立 editor（如尚未存在）
	if _, exists := dt.cellEditors[cellKey]; !exists {
		editor := &widget.Editor{}
		editor.SetText(text)
		dt.cellEditors[cellKey] = editor
	}

	// 建立 clickable（如尚未存在）
	if _, exists := dt.cellClickers[cellKey]; !exists {
		dt.cellClickers[cellKey] = &widget.Clickable{}
	}

	editor := dt.cellEditors[cellKey]
	clicker := dt.cellClickers[cellKey]
	// 若使用者點擊儲存格，進入編輯模式
	if clicker.Clicked(gtx) {
		dt.editingCell = cellKey
		dt.selectedCellKey = cellKey
		dt.selectedContent = text
		editor.SetText(text)
	}

	// 若不在編輯模式，呈現可點擊文字模式
	if dt.editingCell != cellKey {
		return clicker.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					paint.FillShape(gtx.Ops, color.NRGBA{255, 255, 255, 255}, clip.Rect{
						Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
					}.Op())
					paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
						Min: image.Pt(gtx.Dp(dt.CellWidth)-1, 0),
						Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
					}.Op())
					return layout.Dimensions{Size: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight))}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// 計算可用的字符數，中文字符寬度為2，英文為1
						// 這裡我們根據儲存格寬度估算可顯示的字符數
						// 調整為更保守的值以確保不會溢出
						maxChars := int(dt.CellWidth) / 10
						if maxChars < 4 {
							maxChars = 4 // 至少顯示幾個字符
						}

						// 使用改進的 truncateText 函數截斷文字
						displayText := truncateText(text, maxChars)

						// 當被點擊時，記錄選中的儲存格內容
						if clicker.Clicked(gtx) {
							dt.selectedCellKey = cellKey
							dt.selectedContent = text
						}

						return material.Body2(th, displayText).Layout(gtx)
					})
				}),
			)
		})
	}
	// 若在編輯模式，持續檢查是否按下 Enter（Text 包含 \n）
	enteredText := editor.Text()
	if len(enteredText) > 0 && strings.Contains(enteredText, "\n") {
		trimmed := strings.ReplaceAll(enteredText, "\n", "")
		dt.updateCellValue(row, col, trimmed)
		dt.Table.Show()
		editor.SetText(trimmed)

		// 更新選中的內容
		dt.selectedContent = trimmed

		dt.editingCell = ""
	} // 編輯模式介面
	// 先保存原始約束條件
	origConstraints := gtx.Constraints

	// 設置固定的儲存格大小，防止編輯模式影響佈局
	cellWidth := gtx.Dp(dt.CellWidth)
	cellHeight := gtx.Dp(dt.CellHeight)
	gtx.Constraints.Max = image.Point{X: cellWidth, Y: cellHeight}
	gtx.Constraints.Min = image.Point{X: cellWidth, Y: cellHeight}

	// 當使用者點擊儲存格時，讓編輯器保持焦點，實現直接輸入
	if clicker.Clicked(gtx) {
		// 我們確保編輯器已經在編輯模式，不需要做更多操作
		// 因為已經是編輯模式，所以點擊後編輯器仍會保持焦點
	}

	// 使用 Stack 布局
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// 繪製背景和邊框，但不包裹點擊器
			paint.FillShape(gtx.Ops, color.NRGBA{240, 248, 255, 255}, clip.Rect{
				Max: image.Pt(cellWidth, cellHeight),
			}.Op())

			// 恢復原始約束條件以正確繪製內容
			gtx.Constraints = origConstraints
			gtx.Constraints.Max = image.Point{X: cellWidth, Y: cellHeight}

			borderColor := color.NRGBA{0, 123, 255, 255}

			// 畫四邊藍色邊框
			paint.FillShape(gtx.Ops, borderColor, clip.Rect{Min: image.Pt(0, 0), Max: image.Pt(cellWidth, 3)}.Op())
			paint.FillShape(gtx.Ops, borderColor, clip.Rect{Min: image.Pt(0, cellHeight-3), Max: image.Pt(cellWidth, cellHeight)}.Op())
			paint.FillShape(gtx.Ops, borderColor, clip.Rect{Min: image.Pt(0, 0), Max: image.Pt(3, cellHeight)}.Op())
			paint.FillShape(gtx.Ops, borderColor, clip.Rect{Min: image.Pt(cellWidth-3, 0), Max: image.Pt(cellWidth, cellHeight)}.Op())

			return layout.Dimensions{Size: image.Point{X: cellWidth, Y: cellHeight}}
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// 使用 Expanded 讓編輯器能夠占滿整個儲存格，並能接收點擊
			// 使用最小的內邊距，讓輸入框盡可能佔滿整個儲存格
			return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions { // 設置編輯器屬性，確保單行顯示和按 Enter 提交

				// 建立編輯器小工具
				editorWidget := material.Editor(th, editor, "")
				editorWidget.Color = color.NRGBA{0, 0, 0, 255}

				// 留出更多邊距確保文字不會貼邊
				maxWidth := cellWidth - 8
				maxHeight := cellHeight - 4

				// 設置最大和最小約束，確保編輯器大小固定
				gtx.Constraints.Min = image.Point{X: maxWidth, Y: maxHeight}
				gtx.Constraints.Max = image.Point{X: maxWidth, Y: maxHeight}

				// 確保此格子在選中時，顯示完整內容在頂部區域
				if dt.editingCell == cellKey {
					dt.selectedCellKey = cellKey
					dt.selectedContent = editor.Text()
				}

				return editorWidget.Layout(gtx)
			})
		}),
	)
}

// updateCellValue 更新儲存格的值
func (dt *GenericDataTable) updateCellValue(row, col int, newValue string) {
	if dt.Table == nil {
		return
	}

	// 獲取對應的列
	if column := dt.Table.GetColByNumber(col); column != nil {
		// 使用 Insyra 的更新方法
		iCol := indexToLetters(col)
		if newValue == "" || newValue == "." {
			dt.Table.UpdateElement(row, iCol, nil)
			return
		}
		dt.Table.UpdateElement(row, iCol, newValue)
	}
}

func (dt *GenericDataTable) layoutFullTable(gtx layout.Context, th *material.Theme, rows, cols int) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// 列索引行 (A, B, C, ...)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return dt.drawColumnIndexRow(gtx, th, cols)
		}),
		// 列名稱行
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return dt.drawColumnNameRow(gtx, th, cols)
		}),
		// 數據行
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var children []layout.FlexChild
			for i := 0; i < rows; i++ {
				rowName := dt.Table.GetRowNameByIndex(i)
				rowIndex := i // 捕獲迴圈變數
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return dt.drawDataRow(gtx, th, rowIndex, cols, rowName)
				}))
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		}),
	)
}

// truncateText 根據最大寬度截斷文字，支援中文等寬字符
func truncateText(text string, maxWidth int) string {
	// 處理空字符串的情況
	if text == "" || maxWidth <= 0 {
		return ""
	}

	// 判斷文字是否超過限制寬度
	textWidth := runewidth.StringWidth(text)
	if textWidth <= maxWidth {
		return text
	}

	// 若字符數很少但寬度超出，可能是表情符號或其他特殊Unicode字符
	// 在這種情況下，簡單地限制字符數
	if len(text) <= 3 && textWidth > maxWidth {
		return "..."
	}

	// 保守處理：確保省略號有足夠空間
	if maxWidth <= 3 {
		return "..."
	}

	// 使用 runewidth 的 Truncate 方法，這會正確處理各種 Unicode 字符
	return runewidth.Truncate(text, maxWidth-1, "…") // 使用單個省略號字符節省空間
}
