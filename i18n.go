package main

// Language 语言类型
type Language string

const (
	LangEnglish Language = "en"
	LangChinese Language = "zh"
)

// Translations 翻译文本集合
type Translations struct {
	// Window
	WindowTitle string

	// Menu - File
	MenuFile          string
	MenuOpen          string
	MenuNewTab        string
	MenuSaveAs        string
	MenuCloseTab      string
	MenuExit          string

	// Menu - View
	MenuView          string
	MenuFirstPage     string
	MenuPrevPage      string
	MenuNextPage      string
	MenuLastPage      string
	MenuZoomIn        string
	MenuZoomOut       string
	MenuActualSize    string

	// Menu - Help
	MenuHelp          string
	MenuShortcuts     string
	MenuAbout         string

	// Menu - Language
	MenuLanguage      string
	MenuEnglish       string
	MenuChinese       string

	// Status
	StatusNoDocument  string
	StatusPage        string
	StatusZoom        string
	StatusSize        string

	// Messages
	MsgDoubleClickOpen    string
	MsgLoading            string
	MsgLoadFailed         string
	MsgRenderFailed       string
	MsgNoDocumentToSave   string
	MsgSaveSuccess        string
	MsgSaveFailed         string
	MsgInvalidPage        string

	// Dialogs
	DialogShortcutsTitle  string
	DialogShortcutsText   string
	DialogAboutTitle      string
	DialogAboutText       string

	// Toolbar hints
	HintOpen              string
	HintSaveAs            string
	HintCloseTab          string
	HintFirstPage         string
	HintPrevPage          string
	HintNextPage          string
	HintLastPage          string
	HintPageEntry         string
	HintZoomOut           string
	HintZoomIn            string
}

// GetTranslations 获取指定语言的翻译
func GetTranslations(lang Language) *Translations {
	switch lang {
	case LangChinese:
		return getChineseTranslations()
	default:
		return getEnglishTranslations()
	}
}

// getEnglishTranslations 英文翻译
func getEnglishTranslations() *Translations {
	return &Translations{
		WindowTitle:       "PDF Reader",

		MenuFile:          "File",
		MenuOpen:          "Open...",
		MenuNewTab:        "New Tab",
		MenuSaveAs:        "Save As...",
		MenuCloseTab:      "Close Tab",
		MenuExit:          "Exit",

		MenuView:          "View",
		MenuFirstPage:     "First Page",
		MenuPrevPage:      "Previous Page",
		MenuNextPage:      "Next Page",
		MenuLastPage:      "Last Page",
		MenuZoomIn:        "Zoom In",
		MenuZoomOut:       "Zoom Out",
		MenuActualSize:    "Actual Size",

		MenuHelp:          "Help",
		MenuShortcuts:     "Shortcuts",
		MenuAbout:         "About",

		MenuLanguage:      "Language",
		MenuEnglish:       "English",
		MenuChinese:       "中文",

		StatusNoDocument:  "No document open",
		StatusPage:        "Page",
		StatusZoom:        "Zoom",
		StatusSize:        "Size",

		MsgDoubleClickOpen:    "Double-click to open PDF file",
		MsgLoading:            "Loading...",
		MsgLoadFailed:         "Load failed: %v",
		MsgRenderFailed:       "Render failed: %v",
		MsgNoDocumentToSave:   "No document to save",
		MsgSaveSuccess:        "File saved successfully",
		MsgSaveFailed:         "Save failed: %v",
		MsgInvalidPage:        "Invalid page number",

		DialogShortcutsTitle: "Shortcuts",
		DialogShortcutsText: `Keyboard Shortcuts:

Navigation:
  Left / PageUp      - Previous page
  Right / PageDown   - Next page
  Space              - Next page
  Home               - First page
  End                - Last page

Other:
  Ctrl+W             - Close current tab
`,

		DialogAboutTitle: "About",
		DialogAboutText: `PDF Reader v1.2.4

Built with Fyne + go-fitz
A simple and efficient PDF reading tool

Open Source Licenses:
- Fyne: BSD-3-Clause
- go-fitz: AGPL-3.0
`,

		HintOpen:      "Open PDF file",
		HintSaveAs:    "Save as",
		HintCloseTab:  "Close current tab",
		HintFirstPage: "First page",
		HintPrevPage:  "Previous page",
		HintNextPage:  "Next page",
		HintLastPage:  "Last page",
		HintPageEntry: "Page number",
		HintZoomOut:   "Zoom out",
		HintZoomIn:    "Zoom in",
	}
}

// getChineseTranslations 中文翻译
func getChineseTranslations() *Translations {
	return &Translations{
		WindowTitle:       "PDF 阅读器",

		MenuFile:          "文件",
		MenuOpen:          "打开...",
		MenuNewTab:        "新建标签页",
		MenuSaveAs:        "另存为...",
		MenuCloseTab:      "关闭标签页",
		MenuExit:          "退出",

		MenuView:          "查看",
		MenuFirstPage:     "首页",
		MenuPrevPage:      "上一页",
		MenuNextPage:      "下一页",
		MenuLastPage:      "末页",
		MenuZoomIn:        "放大",
		MenuZoomOut:       "缩小",
		MenuActualSize:    "实际大小",

		MenuHelp:          "帮助",
		MenuShortcuts:     "快捷键",
		MenuAbout:         "关于",

		MenuLanguage:      "语言",
		MenuEnglish:       "English",
		MenuChinese:       "中文",

		StatusNoDocument:  "未打开文档",
		StatusPage:        "第",
		StatusZoom:        "缩放",
		StatusSize:        "大小",

		MsgDoubleClickOpen:    "双击打开 PDF 文件",
		MsgLoading:            "正在加载...",
		MsgLoadFailed:         "加载失败: %v",
		MsgRenderFailed:       "渲染失败: %v",
		MsgNoDocumentToSave:   "没有打开的文档可保存",
		MsgSaveSuccess:        "文件已保存",
		MsgSaveFailed:         "保存失败: %v",
		MsgInvalidPage:        "无效的页码",

		DialogShortcutsTitle: "快捷键",
		DialogShortcutsText: `快捷键列表:

导航:
  左箭头 / PageUp    - 上一页
  右箭头 / PageDown  - 下一页
  空格键             - 下一页
  Home              - 首页
  End               - 末页

其他:
  Ctrl+W            - 关闭当前标签页
`,

		DialogAboutTitle: "关于",
		DialogAboutText: `PDF 阅读器 v1.2.4

基于 Fyne + go-fitz 开发
简洁、高效的 PDF 阅读工具

开源许可:
- Fyne: BSD-3-Clause
- go-fitz: AGPL-3.0
`,

		HintOpen:      "打开 PDF 文件",
		HintSaveAs:    "另存为",
		HintCloseTab:  "关闭当前标签页",
		HintFirstPage: "首页",
		HintPrevPage:  "上一页",
		HintNextPage:  "下一页",
		HintLastPage:  "末页",
		HintPageEntry: "页码",
		HintZoomOut:   "缩小",
		HintZoomIn:    "放大",
	}
}
