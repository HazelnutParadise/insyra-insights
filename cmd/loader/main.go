package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/HazelnutParadise/insyra-insights/config"
	"github.com/HazelnutParadise/insyra-insights/i18n"
	"github.com/HazelnutParadise/insyra-insights/version"
)

func main() {
	config.Load()

	// config.Set(config.Language, "en-US")

	i18n.SetLanguage(config.Get(config.Language))
	go func() {
		// 建立固定大小的視窗，寬度800，高度600，可移動
		window := new(app.Window)
		window.Option(
			app.Title("Insyra Insights"),
			app.Size(unit.Dp(600), unit.Dp(400)),
			app.Decorated(false), // 允許視窗被移動
		)

		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	// 讀取 logo 圖片
	logoFile, err := os.Open("assets/logo_transparent.png")
	if err != nil {
		return err
	}
	defer logoFile.Close()

	logoImg, _, err := image.Decode(logoFile)
	if err != nil {
		return err
	}
	logoOp := paint.NewImageOp(logoImg)

	// 創建主題，使用我們自定義的字體
	theme := material.NewTheme()

	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// 主佈局
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceAround,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// Logo 部分
						// 設置 logo 的大小
						imgSize := logoOp.Size()

						// 縮放到適當大小 - 增加 logo 大小
						maxWidth := gtx.Dp(unit.Dp(400))                        // 從原本的 300 增加到 400
						scale := min(float32(maxWidth)/float32(imgSize.X), 1.2) // 增加最大縮放比例到 1.2

						// 計算縮放後的尺寸
						scaledWidth := int(float32(imgSize.X) * scale)
						scaledHeight := int(float32(imgSize.Y) * scale)

						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = scaledWidth
							gtx.Constraints.Min.Y = scaledHeight

							// 應用縮放
							defer op.Affine(f32.Affine2D{}.Scale(f32.Point{}, f32.Point{X: scale, Y: scale})).Push(gtx.Ops).Pop()

							// 繪製圖片
							logoOp.Add(gtx.Ops)
							paint.PaintOp{}.Add(gtx.Ops)

							return layout.Dimensions{
								Size: image.Point{X: scaledWidth, Y: scaledHeight},
							}
						})
					}), // 標題文字
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// 定義標題
						title := material.H1(theme, "Insights")

						// 設置更現代化的漸層藍色
						elegantBlue := color.NRGBA{R: 41, G: 98, B: 255, A: 255}
						title.Color = elegantBlue // 設置文字對齊
						title.Alignment = text.Middle

						// 使用斜體
						title.Font.Style = font.Italic

						// 設定字體粗細為粗體
						title.Font.Weight = font.SemiBold

						// 增加字體大小
						title.TextSize = unit.Sp(52)

						// 顯示標題並添加一些上下間距
						return layout.Inset{Top: unit.Dp(-20), Bottom: unit.Dp(45), Left: unit.Dp(250)}.Layout(gtx, title.Layout)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						slogan := material.Body1(theme, i18n.T("starting.slogan"))
						// 使用更柔和的深灰色
						slogan.Color = color.NRGBA{R: 80, G: 80, B: 90, A: 255}
						slogan.Alignment = text.Middle
						slogan.TextSize = unit.Sp(22)
						slogan.Font.Weight = font.Medium
						return layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(30)}.Layout(gtx, slogan.Layout)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						loadingText := material.Body1(theme, i18n.T("starting.loading"))
						// 使用淡藍色調，與標題呼應
						loadingText.Color = color.NRGBA{R: 100, G: 120, B: 180, A: 255}
						loadingText.Alignment = text.Middle
						loadingText.TextSize = unit.Sp(24)
						loadingText.Font.Style = font.Italic
						return layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(20)}.Layout(gtx, loadingText.Layout)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// 設定載入動畫的大小
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(50))
						gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(50))

						// 建立載入動畫
						loader := material.Loader(theme)
						return layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(20)}.Layout(gtx, loader.Layout)

					}),
					// 移除了版本號的 layout.Rigid
				)
			})

			// 在主佈局之後，單獨處理版本號 - 將其放在右下角
			versionText := material.Caption(theme, i18n.T("starting.version")+": "+version.Version)
			versionText.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 200}
			versionText.Alignment = text.End
			versionText.TextSize = unit.Sp(14)

			// 使用 Stack 和 SE (South East) 將版本號放在右下角
			layout.Stack{Alignment: layout.SW}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					// 透明背景，充滿整個視窗
					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					// 添加右下角邊距
					return layout.Inset{
						Left:   unit.Dp(10),
						Right:  unit.Dp(10),
						Bottom: unit.Dp(10),
					}.Layout(gtx, versionText.Layout)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
