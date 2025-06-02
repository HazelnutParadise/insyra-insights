package ui

// SampleData 代表樣本數據
type SampleData struct {
	Headers []string
	Rows    [][]string
}

// GetSampleData 獲取樣本數據
func GetSampleData() *SampleData {
	return &SampleData{
		Headers: []string{"日期", "收入", "支出", "淨利", "備註"},
		Rows: [][]string{
			{"2025-01-15", "15,000", "8,500", "6,500", "季度開始"},
			{"2025-02-10", "18,200", "9,100", "9,100", "促銷活動"},
			{"2025-03-05", "16,800", "8,300", "8,500", "普通月份"},
			{"2025-04-20", "19,500", "10,200", "9,300", "業績提升"},
			{"2025-05-12", "17,300", "9,600", "7,700", "庫存調整"},
			{"2025-06-01", "20,100", "10,800", "9,300", "季度結束"},
		},
	}
}

// GetSampleChartData 獲取樣本圖表數據
func GetSampleChartData() map[string][]float64 {
	return map[string][]float64{
		"收入": {15000, 18200, 16800, 19500, 17300, 20100},
		"支出": {8500, 9100, 8300, 10200, 9600, 10800},
		"淨利": {6500, 9100, 8500, 9300, 7700, 9300},
	}
}

// GetStatisticsData 獲取統計數據
func GetStatisticsData() map[string]float64 {
	return map[string]float64{
		"平均收入":  17816.67,
		"最高收入":  20100.00,
		"最低收入":  15000.00,
		"收入標準差": 1712.35,
		"平均淨利":  8400.00,
		"淨利率":   47.14,
	}
}
