package main

import (
	"os"

	"fyne.io/fyne/v2/app"
)

func main() {
	// 创建 Fyne 应用
	myApp := app.New()
	myApp.Settings().SetTheme(&customTheme{})

	// 创建界面（不再需要传递 controller）
	ui := NewViewerUI(myApp, nil)

	// 如果有命令行参数，在当前标签页打开文件
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		ui.openFileInCurrentTab(filePath)
	}

	// 显示窗口并运行
	ui.Show()
}
