package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
)

func init() {
	LoadLanguage("en-US") // fallback 預設
	LoadLanguage("zh-TW") // 預設語系
}

//go:embed lang/*.json
var langFiles embed.FS

var (
	strings      map[string]string
	fallback     map[string]string
	currentLang  = "zh-TW"
	fallbackLang = "en-US"
	loadedLangs  = map[string]map[string]string{}
)

func LoadLanguage(code string) error {
	filePath := fmt.Sprintf("lang/%s.json", code)
	data, err := langFiles.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s failed: %w", filePath, err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("json parse %s failed: %w", filePath, err)
	}
	result := make(map[string]string)
	flatten(raw, "", result)
	loadedLangs[code] = result

	if code == currentLang {
		strings = result
	}
	if code == fallbackLang {
		fallback = result
	}
	return nil
}

func SetLanguage(code string) {
	currentLang = code
	if val, ok := loadedLangs[code]; ok {
		strings = val
	} else {
		log.Printf("⚠️ 語言 %s 尚未載入", code)
	}
}

func T(key string) string {
	if val, ok := strings[key]; ok {
		return val
	}
	if val, ok := fallback[key]; ok {
		return val
	}
	log.Printf("⚠️ Missing key: %s", key)
	return "??" + key + "??"
}

func flatten(input map[string]any, prefix string, output map[string]string) {
	for k, v := range input {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		switch val := v.(type) {
		case string:
			output[fullKey] = val
		case map[string]any:
			flatten(val, fullKey, output)
		}
	}
}
