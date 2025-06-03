// GenericDataTable 是一個可繪製 Insyra DataTable 的 UI 表格組件，支援雙向捲動
package ui

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/HazelnutParadise/insyra"
)

type GenericDataTable struct {
	Table         *insyra.DataTable
	CellWidth     unit.Dp
	CellHeight    unit.Dp
	HeaderBgColor color.NRGBA
	BorderColor   color.NRGBA

	// 編輯功能
	cellEditors    map[string]*widget.Editor    // 儲存格編輯器 (key: "row:col")
	cellClickers   map[string]*widget.Clickable // 儲存格點擊器 (key: "row:col")
	editingCell    string                       // 當前編輯的儲存格	// 捲動控制
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
}

func (dt *GenericDataTable) drawColumnIndexRow(gtx layout.Context, th *material.Theme, cols int) layout.Dimensions {
	var children []layout.FlexChild
	// 只使用一個儲存格作為行標頭
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return dt.cell(gtx, th, "")
	}))
	for i := 0; i < cols; i++ {
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

	for i := 0; i < cols; i++ {
		var text string

		// 使用列名作為鍵來獲取資料
		colName := dt.Table.GetColByNumber(i).GetName()
		if colData, exists := data[colName]; exists && row < len(colData) {
			text = fmt.Sprint(colData[row])
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
	// 檢查是否有編輯器，如果沒有則創建
	if _, exists := dt.cellEditors[cellKey]; !exists {
		editor := &widget.Editor{
			SingleLine: true,
			Submit:     true,
		}
		editor.SetText(text)
		dt.cellEditors[cellKey] = editor
	}

	// 檢查是否有點擊器，如果沒有則創建
	if _, exists := dt.cellClickers[cellKey]; !exists {
		dt.cellClickers[cellKey] = &widget.Clickable{}
	}

	editor := dt.cellEditors[cellKey]
	clicker := dt.cellClickers[cellKey]

	// 檢查是否有提交事件
	if editor.Submit {
		newValue := editor.Text()
		dt.updateCellValue(row, col, newValue)
		dt.editingCell = ""
		editor.Submit = false
	}
	// 檢查點擊事件來進入編輯模式
	if clicker.Clicked(gtx) {
		dt.editingCell = cellKey
		editor.SetText(text)
	}
	// 檢查點擊事件來進入編輯模式
	if dt.editingCell != cellKey {
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				// 繪製儲存格背景
				paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.Rect{
					Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
				}.Op())

				// 繪製邊框
				paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
					Min: image.Pt(gtx.Dp(dt.CellWidth)-1, 0),
					Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
				}.Op())

				return layout.Dimensions{Size: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight))}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					// 普通模式：顯示文字，可點擊進入編輯
					return clicker.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Body2(th, text).Layout(gtx)
					})
				})
			}),
		)
	} else {
		// 編輯模式
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				// 繪製編輯模式的背景（淺藍色）
				paint.FillShape(gtx.Ops, color.NRGBA{R: 240, G: 248, B: 255, A: 255}, clip.Rect{
					Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
				}.Op())

				// 繪製加粗的編輯模式邊框（藍色，3像素）
				borderColor := color.NRGBA{R: 0, G: 123, B: 255, A: 255} // 藍色邊框

				// 上邊框
				paint.FillShape(gtx.Ops, borderColor, clip.Rect{
					Min: image.Pt(0, 0),
					Max: image.Pt(gtx.Dp(dt.CellWidth), 3),
				}.Op())

				// 下邊框
				paint.FillShape(gtx.Ops, borderColor, clip.Rect{
					Min: image.Pt(0, gtx.Dp(dt.CellHeight)-3),
					Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
				}.Op())

				// 左邊框
				paint.FillShape(gtx.Ops, borderColor, clip.Rect{
					Min: image.Pt(0, 0),
					Max: image.Pt(3, gtx.Dp(dt.CellHeight)),
				}.Op())

				// 右邊框
				paint.FillShape(gtx.Ops, borderColor, clip.Rect{
					Min: image.Pt(gtx.Dp(dt.CellWidth)-3, 0),
					Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
				}.Op())

				return layout.Dimensions{Size: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight))}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(6)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					// 編輯模式：顯示輸入框
					editorWidget := material.Editor(th, editor, "")
					editorWidget.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255} // 黑色文字
					return editorWidget.Layout(gtx)
				})
			}),
		)
	}
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
