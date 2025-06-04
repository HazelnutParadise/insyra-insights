package ui

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// editableColumnName 創建可編輯的欄名儲存格
func (dt *GenericDataTable) editableColumnName(gtx layout.Context, th *material.Theme, colIndex int, colName string) layout.Dimensions {
	// 建立 editor（如尚未存在）
	if _, exists := dt.colNameEditors[colIndex]; !exists {
		editor := &widget.Editor{}
		editor.SetText(colName)
		editor.SingleLine = true
		editor.Submit = true
		dt.colNameEditors[colIndex] = editor
	}

	// 建立 clickable（如尚未存在）
	if _, exists := dt.colNameClickers[colIndex]; !exists {
		dt.colNameClickers[colIndex] = &widget.Clickable{}
	}

	editor := dt.colNameEditors[colIndex]
	clicker := dt.colNameClickers[colIndex]

	// 若使用者點擊儲存格，進入編輯模式
	if clicker.Clicked(gtx) {
		dt.editingColName = colIndex
		dt.selectedCol = colIndex
		editor.SetText(colName)
		// 讓選中的儲存格內容顯示在頂部
		dt.selectedCellKey = fmt.Sprintf("col:%d", colIndex)
		dt.selectedContent = colName
	}

	// 若不在編輯模式，呈現可點擊文字模式
	if dt.editingColName != colIndex {
		return clicker.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// 使用與一般標題儲存格相同的背景處理
			var bgColor color.NRGBA
			if colIndex == dt.selectedCol {
				// 對選中列的標題使用更深的綠色高亮
				bgColor = color.NRGBA{R: 200, G: 240, B: 200, A: 255}
			} else {
				// 欄名一律使用淡藍色背景
				bgColor = color.NRGBA{R: 225, G: 235, B: 250, A: 255}
			}

			// 繪製背景矩形
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					// 繪製背景和邊框
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
					// 欄名文本
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body2(th, colName)
						lbl.Font.Weight = font.SemiBold
						lbl.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
						return lbl.Layout(gtx)
					})
				}),
			)
		})
	}

	// 若在編輯模式，檢查是否按下 Enter
	enteredText := editor.Text()
	if len(enteredText) > 0 && strings.Contains(enteredText, "\n") {
		trimmed := strings.ReplaceAll(enteredText, "\n", "")
		// 更新欄名
		dt.updateColumnName(colIndex, trimmed)
		editor.SetText(trimmed)
		dt.editingColName = -1 // 退出編輯模式
		dt.selectedContent = trimmed

		// 刷新表格
		dt.Table.Show()
	}

	// 編輯模式介面
	// 先保存原始約束條件
	origConstraints := gtx.Constraints

	// 設置固定的儲存格大小
	cellWidth := gtx.Dp(dt.CellWidth)
	cellHeight := gtx.Dp(dt.CellHeight)
	gtx.Constraints.Max = image.Point{X: cellWidth, Y: cellHeight}
	gtx.Constraints.Min = image.Point{X: cellWidth, Y: cellHeight}

	// 使用 Stack 布局
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// 繪製編輯模式背景
			roundedRect := clip.RRect{
				Rect: image.Rectangle{
					Max: image.Point{X: cellWidth, Y: cellHeight},
				},
				NE: 4, // 所有角都有小圓角
				SE: 4,
				SW: 4,
				NW: 4,
			}
			bgColor := color.NRGBA{R: 240, G: 250, B: 255, A: 255} // 淡藍色背景，與普通儲存格區分
			paint.FillShape(gtx.Ops, bgColor, roundedRect.Op(gtx.Ops))

			// 恢復原始約束條件以正確繪製內容
			gtx.Constraints = origConstraints
			gtx.Constraints.Max = image.Point{X: cellWidth, Y: cellHeight}

			// 使用藍色邊框
			borderColor := color.NRGBA{R: 100, G: 150, B: 240, A: 255} // 藍色邊框

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

			return layout.Dimensions{Size: image.Point{X: cellWidth, Y: cellHeight}}
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// 使用 Expanded 讓編輯器能夠占滿整個儲存格
			return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// 建立編輯器小工具
				editorWidget := material.Editor(th, editor, "")
				editorWidget.Color = color.NRGBA{0, 0, 0, 255}

				// 留出更多邊距確保文字不會貼邊
				maxWidth := cellWidth - 8
				maxHeight := cellHeight - 4

				// 設置最大和最小約束，確保編輯器大小固定
				gtx.Constraints.Min = image.Point{X: maxWidth, Y: maxHeight}
				gtx.Constraints.Max = image.Point{X: maxWidth, Y: maxHeight}

				return editorWidget.Layout(gtx)
			})
		}),
	)
}

// editableRowName 創建可編輯的行名儲存格
func (dt *GenericDataTable) editableRowName(gtx layout.Context, th *material.Theme, rowIndex int, rowName string) layout.Dimensions {
	// 建立 editor（如尚未存在）
	if _, exists := dt.rowNameEditors[rowIndex]; !exists {
		editor := &widget.Editor{}
		editor.SetText(rowName)
		editor.SingleLine = true
		editor.Submit = true
		dt.rowNameEditors[rowIndex] = editor
	}

	// 建立 clickable（如尚未存在）
	if _, exists := dt.rowNameClickers[rowIndex]; !exists {
		dt.rowNameClickers[rowIndex] = &widget.Clickable{}
	}

	editor := dt.rowNameEditors[rowIndex]
	clicker := dt.rowNameClickers[rowIndex]

	// 若使用者點擊儲存格，進入編輯模式
	if clicker.Clicked(gtx) {
		dt.editingRowName = rowIndex
		dt.selectedRow = rowIndex
		editor.SetText(rowName)
		// 讓選中的儲存格內容顯示在頂部
		dt.selectedCellKey = fmt.Sprintf("row:%d", rowIndex)
		dt.selectedContent = rowName
	}

	// 若不在編輯模式，呈現可點擊文字模式
	if dt.editingRowName != rowIndex {
		return clicker.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// 使用與一般標題儲存格相同的背景處理
			var bgColor color.NRGBA
			if rowIndex == dt.selectedRow {
				// 對選中行的標題使用更深的綠色高亮
				bgColor = color.NRGBA{R: 200, G: 240, B: 200, A: 255}
			} else {
				// 列索引使用淡紫色背景
				bgColor = color.NRGBA{R: 230, G: 220, B: 255, A: 255}
			}

			// 繪製背景矩形
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					// 繪製背景和邊框
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
					// 顯示行索引和行名，但將行名設成可編輯
					// 將行索引顯示為 "索引:"，行名則可以編輯
					indexPart := fmt.Sprintf("%d: ", rowIndex)

					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// 使用 Flex 排列索引和名稱
						return layout.Flex{
							Axis:    layout.Horizontal,
							Spacing: layout.SpaceStart,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								// 索引部分 - 不可編輯
								indexLbl := material.Body2(th, indexPart)
								indexLbl.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255} // 灰色
								indexLbl.Font.Weight = font.SemiBold
								return indexLbl.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								// 名稱部分 - 可點擊編輯
								nameLbl := material.Body2(th, rowName)
								nameLbl.Font.Weight = font.SemiBold
								nameLbl.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255} // 黑色
								return nameLbl.Layout(gtx)
							}),
						)
					})
				}),
			)
		})
	}

	// 若在編輯模式，檢查是否按下 Enter
	enteredText := editor.Text()
	if len(enteredText) > 0 && strings.Contains(enteredText, "\n") {
		trimmed := strings.ReplaceAll(enteredText, "\n", "")
		// 更新行名
		dt.updateRowName(rowIndex, trimmed)
		editor.SetText(trimmed)
		dt.editingRowName = -1 // 退出編輯模式
		dt.selectedContent = trimmed

		// 刷新表格
		dt.Table.Show()
	}

	// 編輯模式介面
	// 先保存原始約束條件
	origConstraints := gtx.Constraints

	// 設置固定的儲存格大小
	cellWidth := gtx.Dp(dt.CellWidth)
	cellHeight := gtx.Dp(dt.CellHeight)
	gtx.Constraints.Max = image.Point{X: cellWidth, Y: cellHeight}
	gtx.Constraints.Min = image.Point{X: cellWidth, Y: cellHeight}

	// 使用 Stack 布局
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// 繪製編輯模式背景
			roundedRect := clip.RRect{
				Rect: image.Rectangle{
					Max: image.Point{X: cellWidth, Y: cellHeight},
				},
				NE: 4, // 所有角都有小圓角
				SE: 4,
				SW: 4,
				NW: 4,
			}
			bgColor := color.NRGBA{R: 250, G: 245, B: 255, A: 255} // 淡紫色背景，與普通儲存格區分
			paint.FillShape(gtx.Ops, bgColor, roundedRect.Op(gtx.Ops))

			// 恢復原始約束條件以正確繪製內容
			gtx.Constraints = origConstraints
			gtx.Constraints.Max = image.Point{X: cellWidth, Y: cellHeight}

			// 使用紫色邊框
			borderColor := color.NRGBA{R: 180, G: 140, B: 220, A: 255} // 紫色邊框

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

			return layout.Dimensions{Size: image.Point{X: cellWidth, Y: cellHeight}}
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// 在編輯模式下，先繪製行索引（不可編輯），再繪製編輯器
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// 顯示行索引，不可編輯
					indexPart := fmt.Sprintf("%d: ", rowIndex)
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						indexLbl := material.Body2(th, indexPart)
						indexLbl.Font.Weight = font.SemiBold
						indexLbl.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255} // 灰色
						return indexLbl.Layout(gtx)
					})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					// 編輯行名部分
					return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// 建立編輯器小工具
						editorWidget := material.Editor(th, editor, "")
						editorWidget.Color = color.NRGBA{0, 0, 0, 255}

						// 計算索引文字寬度後調整編輯器寬度
						indexWidth := gtx.Dp(20) // 估計索引寬度
						maxWidth := cellWidth - indexWidth - 8
						maxHeight := cellHeight - 4

						// 設置最大和最小約束，確保編輯器大小固定
						gtx.Constraints.Min = image.Point{X: maxWidth, Y: maxHeight}
						gtx.Constraints.Max = image.Point{X: maxWidth, Y: maxHeight}

						return editorWidget.Layout(gtx)
					})
				}),
			)
		}),
	)
}

// updateColumnName 更新欄名
func (dt *GenericDataTable) updateColumnName(colIndex int, newName string) {
	if dt.Table == nil {
		return
	}

	// 獲取舊欄名
	oldName := dt.Table.GetColByNumber(colIndex).GetName()
	if oldName == newName {
		return // 沒有變更
	}

	// 修改欄名
	dt.Table.SetColNameByNumber(colIndex, newName)
}

// updateRowName 更新行名
func (dt *GenericDataTable) updateRowName(rowIndex int, newName string) {
	if dt.Table == nil {
		return
	}

	// 獲取舊行名
	oldName := dt.Table.GetRowNameByIndex(rowIndex)
	if oldName == newName {
		return // 沒有變更
	}

	// 修改行名
	dt.Table.SetRowNameByIndex(rowIndex, newName)
}
