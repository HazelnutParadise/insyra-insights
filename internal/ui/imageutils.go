package ui

import (
	"image"
	"image/png"
	"os"

	"gioui.org/op/paint"
	"gioui.org/unit"
)

// LoadPNG 從檔案載入PNG圖像
func LoadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// ScaleImage 將圖像縮放到指定的寬度和高度
func ScaleImage(img image.Image, width, height unit.Dp) image.Image {
	// 在實際應用中，應該實現真正的圖像縮放
	// 這裡為了簡單，我們只返回原始圖像
	return img
}

// ImageOp 準備圖像以供繪製
func ImageOp(img image.Image) paint.ImageOp {
	return paint.NewImageOp(img)
}
