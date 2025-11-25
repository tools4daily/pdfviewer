package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"fyne.io/fyne/v2"
)

// generateAppIcon 生成应用图标
// 设计：现代简约风格的 PDF 文档图标
// 配色：红色渐变背景 + 白色文档图形
func generateAppIcon() *image.RGBA {
	const size = 256
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// 背景渐变（从红色到深红色）
	for y := 0; y < size; y++ {
		ratio := float64(y) / float64(size)
		r := uint8(220 - int(40*ratio))
		g := uint8(50 - int(20*ratio))
		b := uint8(50 - int(20*ratio))
		bgColor := color.RGBA{r, g, b, 255}

		for x := 0; x < size; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// 绘制圆角矩形背景
	drawRoundRect(img, 20, 20, size-40, size-40, 20, color.RGBA{255, 255, 255, 30})

	// 绘制白色文档形状
	docColor := color.RGBA{255, 255, 255, 255}

	// 文档主体
	drawRect(img, 60, 40, 136, 180, docColor)

	// 文档折角
	drawTriangle(img, 196, 40, 196, 70, 166, 40, color.RGBA{200, 200, 200, 255})

	// PDF 文字
	drawPDFText(img, 80, 100)

	// 文档线条（表示文本内容）
	lineColor := color.RGBA{200, 50, 50, 255}
	drawRect(img, 80, 140, 96, 4, lineColor)
	drawRect(img, 80, 155, 116, 4, lineColor)
	drawRect(img, 80, 170, 86, 4, lineColor)

	return img
}

// drawRect 绘制矩形
func drawRect(img *image.RGBA, x, y, w, h int, col color.Color) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

// drawRoundRect 绘制圆角矩形
func drawRoundRect(img *image.RGBA, x, y, w, h, radius int, col color.Color) {
	// 绘制中心矩形
	drawRect(img, x+radius, y, w-2*radius, h, col)
	drawRect(img, x, y+radius, w, h-2*radius, col)

	// 绘制四个圆角
	drawCircleQuarter(img, x+radius, y+radius, radius, col, 2)         // 左上
	drawCircleQuarter(img, x+w-radius, y+radius, radius, col, 1)       // 右上
	drawCircleQuarter(img, x+radius, y+h-radius, radius, col, 3)       // 左下
	drawCircleQuarter(img, x+w-radius, y+h-radius, radius, col, 0)     // 右下
}

// drawCircleQuarter 绘制四分之一圆
func drawCircleQuarter(img *image.RGBA, cx, cy, radius int, col color.Color, quarter int) {
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radius*radius {
				var draw bool
				switch quarter {
				case 0: // 右下
					draw = dx >= 0 && dy >= 0
				case 1: // 右上
					draw = dx >= 0 && dy <= 0
				case 2: // 左上
					draw = dx <= 0 && dy <= 0
				case 3: // 左下
					draw = dx <= 0 && dy >= 0
				}
				if draw {
					img.Set(cx+dx, cy+dy, col)
				}
			}
		}
	}
}

// drawTriangle 绘制三角形
func drawTriangle(img *image.RGBA, x1, y1, x2, y2, x3, y3 int, col color.Color) {
	// 简单的扫描线填充三角形
	minY := min(y1, min(y2, y3))
	maxY := max(y1, max(y2, y3))

	for y := minY; y <= maxY; y++ {
		// 计算该行的 x 范围
		xCoords := []int{}

		if isIntersect(y, y1, y2) {
			x := interpolateX(y, x1, y1, x2, y2)
			xCoords = append(xCoords, x)
		}
		if isIntersect(y, y2, y3) {
			x := interpolateX(y, x2, y2, x3, y3)
			xCoords = append(xCoords, x)
		}
		if isIntersect(y, y3, y1) {
			x := interpolateX(y, x3, y3, x1, y1)
			xCoords = append(xCoords, x)
		}

		if len(xCoords) >= 2 {
			minX := min(xCoords[0], xCoords[1])
			maxX := max(xCoords[0], xCoords[1])
			for x := minX; x <= maxX; x++ {
				img.Set(x, y, col)
			}
		}
	}
}

func isIntersect(y, y1, y2 int) bool {
	return (y >= min(y1, y2) && y <= max(y1, y2))
}

func interpolateX(y, x1, y1, x2, y2 int) int {
	if y2 == y1 {
		return x1
	}
	return x1 + (x2-x1)*(y-y1)/(y2-y1)
}

// drawPDFText 绘制 "PDF" 文字
func drawPDFText(img *image.RGBA, x, y int) {
	textColor := color.RGBA{200, 50, 50, 255}

	// 简化的像素字体 "PDF"
	// P
	drawRect(img, x, y, 4, 24, textColor)
	drawRect(img, x, y, 16, 4, textColor)
	drawRect(img, x, y+10, 16, 4, textColor)
	drawRect(img, x+12, y, 4, 14, textColor)

	// D
	drawRect(img, x+25, y, 4, 24, textColor)
	drawRect(img, x+25, y, 12, 4, textColor)
	drawRect(img, x+25, y+20, 12, 4, textColor)
	drawRect(img, x+33, y+4, 4, 16, textColor)

	// F
	drawRect(img, x+48, y, 4, 24, textColor)
	drawRect(img, x+48, y, 16, 4, textColor)
	drawRect(img, x+48, y+10, 12, 4, textColor)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getAppIcon 获取应用图标资源
func getAppIcon() fyne.Resource {
	iconImg := generateAppIcon()
	return &fyne.StaticResource{
		StaticName:    "icon.png",
		StaticContent: imgToBytes(iconImg),
	}
}

// imgToBytes 将图像转换为 PNG 字节数据
func imgToBytes(img image.Image) []byte {
	// 使用标准库将图像编码为 PNG 格式
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		// 如果编码失败，返回空数据
		return []byte{}
	}
	return buf.Bytes()
}
