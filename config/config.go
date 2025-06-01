package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Key string

const (
	Language Key = "language"
	// 未來還可擴充更多 key
)

type settings struct {
	Language string `json:"language"`
	// 可擴充其他設定欄位
}

var (
	current settings
	cfgPath string
)

func defaultSettings() settings {
	return settings{
		Language: "zh-TW",
	}
}

func init() {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("無法取得使用者設定目錄: %v", err)
	}
	cfgPath = filepath.Join(dir, "Insyra_Insights", "config.json")
}

// Load 讀取設定檔
func Load() error {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		current = defaultSettings()
		err := Save()
		return err
	}
	err = json.Unmarshal(data, &current)
	if err != nil {
		return err
	}

	// 填補新欄位預設值（若使用者是舊版 config）
	var configLost bool
	if current.Language == "" {
		configLost = true
		current.Language = defaultSettings().Language
	}

	if configLost {
		if err := Save(); err != nil {
			return err
		}
	}
	return nil
}

// Save 儲存設定檔
func Save() error {
	_ = os.MkdirAll(filepath.Dir(cfgPath), 0755)
	data, _ := json.MarshalIndent(current, "", "  ")
	return os.WriteFile(cfgPath, data, 0644)
}

func Get(k Key) string {
	switch k {
	case Language:
		return current.Language
	default:
		return ""
	}
}

func Set(k Key, v string) {
	switch k {
	case Language:
		current.Language = v
	}
}

// Path 傳回設定檔路徑
func Path() string {
	return cfgPath
}
