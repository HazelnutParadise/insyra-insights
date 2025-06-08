package main

import (
	"embed"

	"github.com/HazelnutParadise/insyra"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 使用命令行參數處理 UUID 和文件路徑
	// uuid := flag.String("uuid", "", "UUID from loader")
	// dataFilePath := flag.String("filepath", "", "Path to data file")
	// flag.Parse()

	// // 在建置時，可能沒有 UUID，這是正常的
	// if *uuid == "" {
	// 	fmt.Println("⚠️ 沒有收到 UUID，可能是建置模式")
	// } else {
	// 	fmt.Printf("收到 UUID: %s\n", *uuid)
	// 	// 移除鎖定文件（僅在有 UUID 時）
	// 	lockPath := filepath.Join(os.TempDir(), "insyra_starting_"+*uuid+".lock")
	// 	time.Sleep(2 * time.Second)
	// 	err := os.Remove(lockPath)
	// 	if err != nil {
	// 		log.Printf("警告：無法刪除鎖定檔案: %v", err)
	// 	}
	// }

	insyra.Config.SetDontPanic(true)

	// Create an instance of the app structure
	app := NewApp()

	// 如果有指定資料檔案路徑，可以在這裡進行預載
	// if *dataFilePath != "" {
	// 	fmt.Printf("收到資料檔案路徑: %s\n", *dataFilePath)
	// 	// 這裡可以預載資料，或者傳遞給前端使用
	// }

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Insyra Insights",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 245, G: 245, B: 245, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
