// GenericDataTable 是一個可繪製 Insyra DataTable 的 UI 表格組件，支援雙向捲動
package ui

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"

	"gioui.org/font"
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
	selectedRow     int                          // 已選中格子的行索引
	selectedCol     int                          // 已選中格子的列索引

	// 欄名編輯功能
	colNameEditors  map[int]*widget.Editor    // 欄名編輯器 (key: col索引)
	colNameClickers map[int]*widget.Clickable // 欄名點擊器 (key: col索引)
	editingColName  int                       // 當前正在編輯的欄名索引，-1為無
	// 顏色設定
	selectedRowColor  color.NRGBA // 選中行的背景色
	selectedColColor  color.NRGBA // 選中列的背景色
	selectedCellColor color.NRGBA // 選中單元格的背景色

	// 捲動控制
	verticalList   widget.List
	horizontalList widget.List
}

func NewGenericDataTable(tbl *insyra.DataTable) *GenericDataTable {
	dt := &GenericDataTable{
		Table:         tbl,
		CellWidth:     unit.Dp(80),
		CellHeight:    unit.Dp(32),
		HeaderBgColor: color.NRGBA{R: 245, G: 246, B: 250, A: 255}, // 更柔和的標題背景色
		BorderColor:   color.NRGBA{R: 225, G: 228, B: 232, A: 255}, // 更柔和的邊框色
		cellEditors:   make(map[string]*widget.Editor),
		cellClickers:  make(map[string]*widget.Clickable),
		editingCell:   "",
		selectedRow:   -1, // 初始化為 -1 表示未選中任何行
		selectedCol:   -1, // 初始化為 -1 表示未選中任何列

		// 初始化欄名編輯功能
		colNameEditors:    make(map[int]*widget.Editor),
		colNameClickers:   make(map[int]*widget.Clickable),
		editingColName:    -1,                                          // -1 表示尚未編輯任何欄名
		selectedRowColor:  color.NRGBA{R: 235, G: 250, B: 235, A: 255}, // 淡綠色 (選中行背景)
		selectedColColor:  color.NRGBA{R: 235, G: 250, B: 235, A: 255}, // 淡綠色 (選中列背景)
		selectedCellColor: color.NRGBA{R: 220, G: 200, B: 250, A: 255}, // 中紫色 (選中單元格)
	}

	return dt
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

// 舊的 drawColumnIndexRow 和 drawColumnNameRow 函數已被 drawColumnHeader 函數替代

func (dt *GenericDataTable) drawDataRow(gtx layout.Context, th *material.Theme, row, cols int) layout.Dimensions {
	var children []layout.FlexChild
	// row index cell: 對齊資料格、選中行變色
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		cellWidth := gtx.Dp(dt.CellWidth)
		cellHeight := gtx.Dp(dt.CellHeight)
		var bgColor color.NRGBA
		if row == dt.selectedRow {
			bgColor = color.NRGBA{R: 200, G: 240, B: 200, A: 255} // 選中行高亮
		} else {
			bgColor = color.NRGBA{R: 230, G: 220, B: 255, A: 255} // 普通行
		}
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				paint.FillShape(gtx.Ops, bgColor, clip.Rect{
					Max: image.Pt(cellWidth, cellHeight),
				}.Op())
				// 右側格線
				paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
					Min: image.Pt(cellWidth-1, 0),
					Max: image.Pt(cellWidth, cellHeight),
				}.Op())
				// 底部格線
				paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
					Min: image.Pt(0, cellHeight-1),
					Max: image.Pt(cellWidth, cellHeight),
				}.Op())
				return layout.Dimensions{Size: image.Pt(cellWidth, cellHeight)}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				indexLbl := material.Body2(th, fmt.Sprintf("%d", row+1))
				indexLbl.Font.Weight = font.SemiBold
				indexLbl.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, indexLbl.Layout)
			}),
		)
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

// 已移除未使用的 cell 函數

func (dt *GenericDataTable) headerCell(gtx layout.Context, th *material.Theme, text string) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// 獲取行號
			rowNum := -1
			if strings.Contains(text, ": ") {
				parts := strings.Split(text, ": ")
				if len(parts) > 0 {
					rowNum, _ = strconv.Atoi(parts[0])
				}
			} // 根據是否是選中的行來決定背景色
			var bgColor color.NRGBA
			if rowNum == dt.selectedRow {
				// 對選中行的標題使用更深的綠色高亮
				bgColor = color.NRGBA{R: 200, G: 240, B: 200, A: 255}
			} else if text == "列/欄" {
				// 左上角指示格使用灰色
				bgColor = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
			} else if strings.Contains(text, ": ") {
				// 列索引使用淡紫色背景
				bgColor = color.NRGBA{R: 230, G: 220, B: 255, A: 255}
			} else {
				// 欄名一律使用淡藍色背景
				bgColor = color.NRGBA{R: 225, G: 235, B: 250, A: 255}
			}

			// 使用微妙的斜向陰影效果增強擬物感
			paint.FillShape(gtx.Ops, bgColor, clip.Rect{
				Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
			}.Op())

			// 頂部亮色增強立體感
			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 100}, clip.Rect{
				Min: image.Pt(1, 1),
				Max: image.Pt(gtx.Dp(dt.CellWidth)-1, 2),
			}.Op())

			// 右側垂直格線
			paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
				Min: image.Pt(gtx.Dp(dt.CellWidth)-1, 0),
				Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
			}.Op())
			// 底部水平格線
			paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
				Min: image.Pt(0, gtx.Dp(dt.CellHeight)-1),
				Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
			}.Op())

			return layout.Dimensions{Size: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight))}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(th, text)
				lbl.Font.Weight = font.SemiBold
				// 使用藍色文字
				lbl.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
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
	clicker := dt.cellClickers[cellKey] // 若使用者點擊儲存格，進入編輯模式
	if clicker.Clicked(gtx) {
		dt.editingCell = cellKey
		dt.selectedCellKey = cellKey
		dt.selectedContent = text
		dt.selectedRow = row
		dt.selectedCol = col
		editor.SetText(text)
	}
	// 若不在編輯模式，呈現可點擊文字模式
	if dt.editingCell != cellKey {
		return clicker.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions { // 根據選中狀態決定背景色
					var bgColor color.NRGBA
					var isSelected bool = false

					// 如果是選中的儲存格
					if row == dt.selectedRow && col == dt.selectedCol {
						bgColor = dt.selectedCellColor
						isSelected = true
					} else if row == dt.selectedRow { // 如果在選中的行
						bgColor = dt.selectedRowColor
						isSelected = true
					} else if col == dt.selectedCol { // 如果在選中的列
						bgColor = dt.selectedColColor
						isSelected = true
					} else { // 普通儲存格
						bgColor = color.NRGBA{255, 255, 255, 255}
					}

					// 繪製微妙圓角的矩形背景
					roundedRect := clip.RRect{
						Rect: image.Rectangle{
							Max: image.Point{X: gtx.Dp(dt.CellWidth), Y: gtx.Dp(dt.CellHeight)},
						},
						// 所有角都有非常微小的圓角，提升美感
						NE: 1,
						SE: 1,
						SW: 1,
						NW: 1,
					}
					paint.FillShape(gtx.Ops, bgColor, roundedRect.Op(gtx.Ops))

					// 如果是選中的儲存格，增加微妙的內陰影效果增強立體感
					if isSelected {
						// 左側和頂部細微陰影
						paint.FillShape(gtx.Ops, color.NRGBA{0, 0, 0, 20}, clip.Rect{
							Min: image.Pt(0, 0),
							Max: image.Pt(1, gtx.Dp(dt.CellHeight)),
						}.Op())
						paint.FillShape(gtx.Ops, color.NRGBA{0, 0, 0, 20}, clip.Rect{
							Min: image.Pt(0, 0),
							Max: image.Pt(gtx.Dp(dt.CellWidth), 1),
						}.Op())
					} else {
						// 非選中格子增加微妙光澤效果，頂部有細微高亮
						paint.FillShape(gtx.Ops, color.NRGBA{255, 255, 255, 120}, clip.Rect{
							Min: image.Pt(0, 0),
							Max: image.Pt(gtx.Dp(dt.CellWidth), 2),
						}.Op())
					}
					// 右側垂直格線
					paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
						Min: image.Pt(gtx.Dp(dt.CellWidth)-1, 0),
						Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
					}.Op())
					// 底部水平格線
					paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
						Min: image.Pt(0, gtx.Dp(dt.CellHeight)-1),
						Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
					}.Op())
					return layout.Dimensions{Size: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight))}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// 計算可用的字符數，中文字符寬度為2，英文為1
						// 這裡我們根據儲存格寬度估算可顯示的字符數
						// 調整為更保守的值以確保不會溢出
						maxChars := max(int(dt.CellWidth)/10-1, 4)

						// 使用改進的 truncateText 函數截斷文字
						displayText := truncateText(text, maxChars) // 當被點擊時，記錄選中的儲存格內容和行列
						if clicker.Clicked(gtx) {
							dt.selectedCellKey = cellKey
							dt.selectedContent = text
							dt.selectedRow = row
							dt.selectedCol = col
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
		// 更新選中的內容(樂觀更新)
		dt.selectedContent = trimmed
		editor.SetText(trimmed)
		dt.editingCell = ""
		dt.updateCellValue(row, col, trimmed)

		dt.Table.Show()
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
		layout.Stacked(func(gtx layout.Context) layout.Dimensions { // 繪製稍帶圓角的光亮背景 - 編輯模式更突出
			roundedRect := clip.RRect{
				Rect: image.Rectangle{
					Max: image.Point{X: cellWidth, Y: cellHeight},
				},
				NE: 4, // 所有角都有小圓角
				SE: 4,
				SW: 4,
				NW: 4,
			} // 用淡紫色作為編輯模式背景
			bgColor := color.NRGBA{R: 245, G: 240, B: 255, A: 255} // 淡紫色背景
			paint.FillShape(gtx.Ops, bgColor, roundedRect.Op(gtx.Ops))

			// 恢復原始約束條件以正確繪製內容
			gtx.Constraints = origConstraints
			gtx.Constraints.Max = image.Point{X: cellWidth, Y: cellHeight}

			// 使用更柔和的紫色邊框
			borderColor := color.NRGBA{R: 150, G: 120, B: 200, A: 255} // 淡紫色邊框

			// 畫圓角邊框 - 使用多層繪製實現邊框效果
			outerRect := clip.RRect{
				Rect: image.Rectangle{
					Max: image.Point{X: cellWidth, Y: cellHeight},
				},
				NE: 6,
				SE: 6,
				SW: 6,
				NW: 6,
			}
			innerRect := clip.RRect{
				Rect: image.Rectangle{
					Min: image.Point{X: 2, Y: 2},
					Max: image.Point{X: cellWidth - 2, Y: cellHeight - 2},
				},
				NE: 4,
				SE: 4,
				SW: 4,
				NW: 4,
			}

			// 先畫外框
			paint.FillShape(gtx.Ops, borderColor, outerRect.Op(gtx.Ops))
			// 再用背景色填充內部，形成邊框效果
			paint.FillShape(gtx.Ops, bgColor, innerRect.Op(gtx.Ops))

			// 添加陰影效果增強立體感
			for i := 0; i < 3; i++ {
				shadowColor := color.NRGBA{R: 0, G: 0, B: 0, A: uint8(10 - i*3)}
				shadowRect := clip.RRect{
					Rect: image.Rectangle{
						Min: image.Point{X: 2 + i, Y: cellHeight + i},
						Max: image.Point{X: cellWidth + i, Y: cellHeight + i + 1},
					},
					NE: 0,
					SE: 0,
					SW: 0,
					NW: 0,
				}
				paint.FillShape(gtx.Ops, shadowColor, shadowRect.Op(gtx.Ops))
			}

			// 頂部光澤效果
			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 100}, clip.Rect{
				Min: image.Pt(3, 3),
				Max: image.Pt(cellWidth-3, 5),
			}.Op())

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
		// 合併列索引行和列名稱行，無格線
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return dt.drawColumnHeader(gtx, th, cols)
		}),
		// 數據行
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var children []layout.FlexChild
			for i := 0; i < rows; i++ {
				rowIndex := i // 捕獲迴圈變數
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return dt.drawDataRow(gtx, th, rowIndex, cols)
				}))
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		}),
	)
}

// drawColumnHeader 繪製欄位標頭，包含欄索引和欄名，中間沒有分隔線
func (dt *GenericDataTable) drawColumnHeader(gtx layout.Context, th *material.Theme, cols int) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// 欄索引行 (A, B, C, ...)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var children []layout.FlexChild
			// 左上角空白格
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// 不需要底部格線
				return dt.headerCellNoBorder(gtx, th, "", false)
			}))
			for i := range cols {
				label := indexToLetters(i)
				currentLabel := label // 捕獲迴圈變數
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// 不需要底部格線，確保使用紫色背景
					return dt.headerCellNoBorder(gtx, th, currentLabel, false)
				}))
			}
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
		}),
		// 欄名稱行
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var children []layout.FlexChild
			// 左上角指示格 (行/欄)
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return dt.headerCell(gtx, th, "列/欄")
			}))
			for i := 0; i < cols; i++ {
				name := dt.Table.GetColByNumber(i).GetName()
				currentName := name // 捕獲迴圈變數
				currentIndex := i   // 捕獲迴圈變數
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// 使用可編輯的欄名單元格
					return dt.editableColumnName(gtx, th, currentIndex, currentName)
				}))
			}
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
		}),
	)
}

// headerCellNoBorder 繪製沒有底部邊框的標題儲存格
func (dt *GenericDataTable) headerCellNoBorder(gtx layout.Context, th *material.Theme, text string, showBottomBorder bool) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// 檢查是否是選中列的欄標題
			colNum := -1
			if text != "" {
				// 將欄標題轉換為欄號，支持複合索引如 AA, AB 等
				for i := 0; i < 100; i++ { // 設置合理的上限以避免無限循環
					if indexToLetters(i) == text {
						colNum = i
						break
					}
				}
			}

			var bgColor color.NRGBA
			if text == "" {
				// 左上角空白格使用灰色背景，與「列/欄」格子保持一致
				bgColor = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
			} else if colNum == dt.selectedCol {
				// 對選中列的標題使用更深的綠色高亮
				bgColor = color.NRGBA{R: 200, G: 240, B: 200, A: 255}
			} else {
				// 使用淡紫色作為欄索引背景
				bgColor = color.NRGBA{R: 230, G: 220, B: 255, A: 255}
			}

			// 繪製稍帶圓角的背景 (但僅適用於上半部分)
			roundedRect := clip.RRect{
				Rect: image.Rectangle{
					Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
				},
				NE: 4, // 右上稍微圓角
				SE: 0,
				SW: 0,
				NW: 4, // 左上稍微圓角
			}
			paint.FillShape(gtx.Ops, bgColor, roundedRect.Op(gtx.Ops))

			// 繪製微妙的頂部光澤效果 - 使用淺色條紋增強立體感
			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 100}, clip.Rect{
				Min: image.Pt(1, 1),
				Max: image.Pt(gtx.Dp(dt.CellWidth)-1, 3),
			}.Op())

			// 右側垂直格線
			paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
				Min: image.Pt(gtx.Dp(dt.CellWidth)-1, 0),
				Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
			}.Op())

			// 只在需要時繪製底部水平格線
			if showBottomBorder {
				paint.FillShape(gtx.Ops, dt.BorderColor, clip.Rect{
					Min: image.Pt(0, gtx.Dp(dt.CellHeight)-1),
					Max: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight)),
				}.Op())
			}

			return layout.Dimensions{Size: image.Pt(gtx.Dp(dt.CellWidth), gtx.Dp(dt.CellHeight))}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(th, text)
				// 使用粗體字型
				lbl.Font.Weight = font.Bold
				// 使用藍色文字
				lbl.Color = color.NRGBA{R: 0, G: 90, B: 180, A: 255}
				return lbl.Layout(gtx)
			})
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
	return runewidth.Truncate(text, maxWidth, "…") // 使用單個省略號字符節省空間
}

// ResetEditors 清空所有 editor/clicker pool，避免 pool 汙染
func (dt *GenericDataTable) ResetEditors() {
	dt.cellEditors = make(map[string]*widget.Editor)
	dt.colNameEditors = make(map[int]*widget.Editor)
	dt.cellClickers = make(map[string]*widget.Clickable)
	dt.colNameClickers = make(map[int]*widget.Clickable)
}
