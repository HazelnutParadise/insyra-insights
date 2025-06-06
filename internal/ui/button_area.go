// GenericDataTable 相關按鈕區域的實現
package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// drawButtonArea 繪製表格底部的按鈕操作區域
func (dt *GenericDataTable) drawButtonArea(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// 為按鈕區域創建一個固定高度的區域，帶有卡片陰影效果
	return layout.Stack{}.Layout(gtx,
		// 背景和陰影效果層
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Max
			buttonAreaHeight := gtx.Dp(unit.Dp(48)) // 固定的按鈕區域高度

			// 繪製卡片背景
			roundedRect := clip.RRect{
				Rect: image.Rectangle{
					Max: image.Point{X: size.X, Y: buttonAreaHeight},
				},
				NE: 6, SE: 6, SW: 6, NW: 6, // 稍大的圓角
			}

			// 使用微妙的背景色
			bgColor := color.NRGBA{R: 248, G: 248, B: 252, A: 255}
			paint.FillShape(gtx.Ops, bgColor, roundedRect.Op(gtx.Ops))

			// 頂部陰影效果（顯示在按鈕區域上方）
			shadowHeight := 6
			for i := 0; i < shadowHeight; i++ {
				y := i

				// 使用指數衰減為陰影創建更自然的效果
				alpha := uint8(40 - i*7)
				if alpha < 5 {
					alpha = 5
				}

				// 繪製陰影線
				paint.FillShape(gtx.Ops, color.NRGBA{0, 0, 0, alpha}, clip.Rect{
					Min: image.Pt(4, y),
					Max: image.Pt(size.X-4, y+1),
				}.Op())
			}

			// 頂部光澤效果
			paint.FillShape(gtx.Ops, color.NRGBA{255, 255, 255, 120}, clip.Rect{
				Min: image.Pt(4, buttonAreaHeight-3),
				Max: image.Pt(size.X-4, buttonAreaHeight-1),
			}.Op())

			return layout.Dimensions{Size: image.Point{X: size.X, Y: buttonAreaHeight}}
		}),

		// 按鈕內容層 - 在這裡可以添加實際的按鈕
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// 使用水平 Flex 佈局來排列按鈕
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceBetween, // 均勻分布
				}.Layout(gtx,
					// 這裡可以添加實際按鈕，目前僅作為佔位區域
					// 左側按鈕區（如果需要）
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// 預留空間，可稍後添加實際按鈕
						return layout.Dimensions{Size: image.Point{X: gtx.Dp(unit.Dp(100)), Y: gtx.Dp(unit.Dp(32))}}
					}),

					// 中間彈性空間
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: 0}}
					}),

					// 右側按鈕區（如果需要）
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// 預留空間，可稍後添加實際按鈕
						return layout.Dimensions{Size: image.Point{X: gtx.Dp(unit.Dp(100)), Y: gtx.Dp(unit.Dp(32))}}
					}),
				)
			})
		}),
	)
}
