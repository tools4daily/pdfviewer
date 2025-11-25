package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ViewerUI PDF 阅读器界面
type ViewerUI struct {
	window       fyne.Window
	tabContainer *container.AppTabs
	tabs         []*PDFTab
	statusLabel  *widget.Label
	pageEntry    *widget.Entry
	zoomLabel    *widget.Button // 显示缩放比例的按钮
	currentLang  Language       // 当前语言
	tr           *Translations  // 翻译文本
}

// PDFTab 表示单个 PDF 标签页
type PDFTab struct {
	controller    *Controller
	imageCanvas   *canvas.Image
	scrollView    *container.Scroll
	loadingLabel  *widget.Label
	canvasWrapper *scrollableCanvas
	tabItem       *container.TabItem
}

// NewViewerUI 创建界面实例
func NewViewerUI(app fyne.App, _ *Controller) *ViewerUI {
	ui := &ViewerUI{
		tabs:        []*PDFTab{},
		currentLang: LangEnglish, // 默认英文
		tr:          GetTranslations(LangEnglish),
	}

	window := app.NewWindow(ui.tr.WindowTitle)
	window.Resize(fyne.NewSize(900, 700))

	// 设置窗口图标
	window.SetIcon(getAppIcon())

	ui.window = window

	ui.buildUI()
	ui.setupKeyBindings()

	// 创建初始标签页（空页面）
	ui.addNewTab("")

	return ui
}

// NewPDFTab 创建新的 PDF 标签页
func NewPDFTab(ui *ViewerUI, filePath string) *PDFTab {
	tab := &PDFTab{
		controller: NewController(),
	}

	// 创建标签页内容
	content := tab.createContent(ui)

	// 创建标签页项
	tabName := "新建标签"
	if filePath != "" {
		tabName = getFileName(filePath)
	}

	tab.tabItem = container.NewTabItem(tabName, content)

	// 如果指定了文件路径，加载文件
	if filePath != "" {
		go func() {
			tab.loadPDF(filePath, ui)
		}()
	}

	return tab
}

// createContent 创建标签页内容
func (tab *PDFTab) createContent(ui *ViewerUI) fyne.CanvasObject {
	// 中央显示区
	tab.imageCanvas = canvas.NewImageFromImage(nil)
	tab.imageCanvas.FillMode = canvas.ImageFillOriginal // 使用原始大小，支持缩放

	tab.loadingLabel = widget.NewLabel(ui.tr.MsgDoubleClickOpen)
	tab.loadingLabel.Alignment = fyne.TextAlignCenter

	centerContent := container.NewStack(
		tab.imageCanvas,
		container.NewCenter(tab.loadingLabel),
	)

	// 创建支持滚轮和双击的 canvas wrapper
	tab.canvasWrapper = newScrollableCanvas(
		centerContent,
		func(ev *fyne.ScrollEvent) { tab.onScrollWheel(ev, ui) },
		func() { tab.onDoubleTap(ui) },
	)

	tab.scrollView = container.NewScroll(tab.canvasWrapper)

	return tab.scrollView
}

// getFileName 从完整路径提取文件名
func getFileName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		if len(name) > 20 {
			return name[:17] + "..."
		}
		return name
	}
	return "PDF 文档"
}

// buildUI 构建界面
func (ui *ViewerUI) buildUI() {
	// 菜单栏
	ui.window.SetMainMenu(ui.createMenuBar())

	// 顶部工具栏
	toolbar := ui.createToolbar()

	// 创建标签页容器
	ui.tabContainer = container.NewAppTabs()
	ui.tabContainer.SetTabLocation(container.TabLocationTop)

	// 监听标签页切换事件
	ui.tabContainer.OnSelected = func(tab *container.TabItem) {
		ui.updateStatusBar()
		ui.updateZoomLabel()
	}

	// 底部状态栏
	statusBar := ui.createStatusBar()

	// 组合布局
	content := container.NewBorder(
		toolbar,
		statusBar,
		nil, nil,
		ui.tabContainer,
	)

	ui.window.SetContent(content)
}

// addNewTab 添加新标签页
func (ui *ViewerUI) addNewTab(filePath string) {
	tab := NewPDFTab(ui, filePath)
	ui.tabs = append(ui.tabs, tab)
	ui.tabContainer.Append(tab.tabItem)
	ui.tabContainer.Select(tab.tabItem)
	ui.updateStatusBar()
}

// getCurrentTab 获取当前激活的标签页
func (ui *ViewerUI) getCurrentTab() *PDFTab {
	if len(ui.tabs) == 0 {
		return nil
	}

	selectedItem := ui.tabContainer.Selected()
	for _, tab := range ui.tabs {
		if tab.tabItem == selectedItem {
			return tab
		}
	}

	return nil
}

// closeCurrentTab 关闭当前标签页
func (ui *ViewerUI) closeCurrentTab() {
	currentTab := ui.getCurrentTab()
	if currentTab == nil {
		return
	}

	// 找到并移除标签页
	for i, tab := range ui.tabs {
		if tab == currentTab {
			// 关闭 PDF 引擎
			if tab.controller.engine != nil {
				tab.controller.engine.Close()
			}

			// 从列表中移除
			ui.tabs = append(ui.tabs[:i], ui.tabs[i+1:]...)
			ui.tabContainer.Remove(tab.tabItem)
			break
		}
	}

	// 如果没有标签页了，创建一个新的空标签页
	if len(ui.tabs) == 0 {
		ui.addNewTab("")
	}

	ui.updateStatusBar()
}

// updateStatusBar 更新状态栏
func (ui *ViewerUI) updateStatusBar() {
	currentTab := ui.getCurrentTab()
	if currentTab == nil || !currentTab.controller.HasDocument() {
		ui.statusLabel.SetText(ui.tr.StatusNoDocument)
		return
	}

	ui.statusLabel.SetText(currentTab.controller.GetStatusText(ui.tr))
}

// updateZoomLabel 更新缩放标签
func (ui *ViewerUI) updateZoomLabel() {
	currentTab := ui.getCurrentTab()
	if currentTab == nil || !currentTab.controller.HasDocument() {
		ui.zoomLabel.SetText("100%")
		return
	}

	zoomPercent := int(currentTab.controller.zoomLevel * 100)
	ui.zoomLabel.SetText(fmt.Sprintf("%d%%", zoomPercent))
}

// createMenuBar 创建菜单栏
func (ui *ViewerUI) createMenuBar() *fyne.MainMenu {
	// 文件菜单
	fileMenu := fyne.NewMenu(ui.tr.MenuFile,
		fyne.NewMenuItem(ui.tr.MenuOpen, ui.onOpenFile),
		fyne.NewMenuItem(ui.tr.MenuNewTab, func() {
			ui.addNewTab("")
		}),
		fyne.NewMenuItem(ui.tr.MenuSaveAs, ui.onSaveAs),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem(ui.tr.MenuCloseTab, func() {
			ui.closeCurrentTab()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem(ui.tr.MenuExit, func() {
			ui.window.Close()
		}),
	)

	// 查看菜单
	viewMenu := fyne.NewMenu(ui.tr.MenuView,
		fyne.NewMenuItem(ui.tr.MenuFirstPage, ui.onFirstPage),
		fyne.NewMenuItem(ui.tr.MenuPrevPage, ui.onPrevPage),
		fyne.NewMenuItem(ui.tr.MenuNextPage, ui.onNextPage),
		fyne.NewMenuItem(ui.tr.MenuLastPage, ui.onLastPage),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem(ui.tr.MenuZoomIn, ui.onZoomIn),
		fyne.NewMenuItem(ui.tr.MenuZoomOut, ui.onZoomOut),
		fyne.NewMenuItem(ui.tr.MenuActualSize, ui.onZoomReset),
	)

	// 语言菜单
	langMenu := fyne.NewMenu(ui.tr.MenuLanguage,
		fyne.NewMenuItem(ui.tr.MenuEnglish, func() {
			ui.switchLanguage(LangEnglish)
		}),
		fyne.NewMenuItem(ui.tr.MenuChinese, func() {
			ui.switchLanguage(LangChinese)
		}),
	)

	// 帮助菜单
	helpMenu := fyne.NewMenu(ui.tr.MenuHelp,
		fyne.NewMenuItem(ui.tr.MenuShortcuts, ui.onShowShortcuts),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem(ui.tr.MenuAbout, ui.onShowAbout),
	)

	return fyne.NewMainMenu(fileMenu, viewMenu, langMenu, helpMenu)
}

// createToolbar 创建工具栏
func (ui *ViewerUI) createToolbar() fyne.CanvasObject {
	// 文件按钮
	openBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), ui.onOpenFile)
	openBtn.Importance = widget.HighImportance

	// 另存为按钮
	saveAsBtn := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), ui.onSaveAs)

	// 关闭标签页按钮
	closeTabBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		ui.closeCurrentTab()
	})
	closeTabBtn.Importance = widget.DangerImportance

	// 导航按钮
	firstBtn := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), ui.onFirstPage)
	prevBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), ui.onPrevPage)
	nextBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), ui.onNextPage)
	lastBtn := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), ui.onLastPage)

	// 页码跳转
	ui.pageEntry = widget.NewEntry()
	ui.pageEntry.SetPlaceHolder(ui.tr.HintPageEntry)
	ui.pageEntry.OnSubmitted = ui.onPageJump

	// 缩放按钮
	zoomOutBtn := widget.NewButtonWithIcon("", theme.ZoomOutIcon(), ui.onZoomOut)
	ui.zoomLabel = widget.NewButton("100%", ui.onZoomReset)
	zoomInBtn := widget.NewButtonWithIcon("", theme.ZoomInIcon(), ui.onZoomIn)

	// 组合工具栏
	toolbar := container.NewHBox(
		openBtn,
		saveAsBtn,
		closeTabBtn,
		widget.NewSeparator(),
		firstBtn,
		prevBtn,
		ui.pageEntry,
		nextBtn,
		lastBtn,
		widget.NewSeparator(),
		zoomOutBtn,
		ui.zoomLabel,
		zoomInBtn,
	)

	return toolbar
}

// createStatusBar 创建状态栏
func (ui *ViewerUI) createStatusBar() fyne.CanvasObject {
	ui.statusLabel = widget.NewLabel(ui.tr.StatusNoDocument)
	return container.NewPadded(ui.statusLabel)
}

// setupKeyBindings 设置键盘快捷键
func (ui *ViewerUI) setupKeyBindings() {
	// 添加 Ctrl+W 快捷键关闭当前标签
	ui.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyW,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		ui.closeCurrentTab()
	})

	ui.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		currentTab := ui.getCurrentTab()
		if currentTab == nil {
			return
		}

		switch key.Name {
		case fyne.KeyLeft, fyne.KeyPageUp:
			currentTab.onPrevPage(ui)
		case fyne.KeyRight, fyne.KeyPageDown, fyne.KeySpace:
			currentTab.onNextPage(ui)
		case fyne.KeyHome:
			currentTab.onFirstPage(ui)
		case fyne.KeyEnd:
			currentTab.onLastPage(ui)
		}
	})
}

// PDFTab 的滚轮事件处理
func (tab *PDFTab) onScrollWheel(ev *fyne.ScrollEvent, ui *ViewerUI) {
	if !tab.controller.HasDocument() {
		return
	}

	if ev.Scrolled.DY < 0 {
		// 向下滚动 → 下一页
		if tab.controller.NextPage() {
			tab.renderPage(ui)
		}
	} else if ev.Scrolled.DY > 0 {
		// 向上滚动 → 上一页
		if tab.controller.PrevPage() {
			tab.renderPage(ui)
		}
	}
}

// PDFTab 的双击事件处理
func (tab *PDFTab) onDoubleTap(ui *ViewerUI) {
	// 仅在未打开文档时响应双击
	if !tab.controller.HasDocument() {
		ui.onOpenFile()
	}
}

// openFileInCurrentTab 在当前标签页打开文件
func (ui *ViewerUI) openFileInCurrentTab(filePath string) {
	currentTab := ui.getCurrentTab()
	if currentTab != nil {
		go func() {
			currentTab.loadPDF(filePath, ui)
		}()
	}
}

// onSaveAs 另存为
func (ui *ViewerUI) onSaveAs() {
	currentTab := ui.getCurrentTab()
	if currentTab == nil || !currentTab.controller.HasDocument() {
		dialog.ShowInformation(ui.tr.MenuHelp, ui.tr.MsgNoDocumentToSave, ui.window)
		return
	}

	// 获取当前文件路径
	srcPath := currentTab.controller.engine.GetFilePath()

	// 创建保存对话框
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		// 复制文件
		dstPath := writer.URI().Path()
		err = copyFile(srcPath, dstPath)
		if err != nil {
			dialog.ShowError(fmt.Errorf(ui.tr.MsgSaveFailed, err), ui.window)
			return
		}

		dialog.ShowInformation(ui.tr.MenuHelp, ui.tr.MsgSaveSuccess, ui.window)
	}, ui.window)

	// 设置默认文件名
	fileName := currentTab.controller.engine.GetFileName()
	saveDialog.SetFileName(fileName)
	saveDialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
	saveDialog.Show()
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// onOpenFile 打开文件对话框
func (ui *ViewerUI) onOpenFile() {
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		filePath := reader.URI().Path()

		// 检查当前标签页是否为空
		currentTab := ui.getCurrentTab()
		if currentTab != nil && !currentTab.controller.HasDocument() {
			// 当前标签页为空，直接在当前标签页打开
			go func() {
				currentTab.loadPDF(filePath, ui)
			}()
		} else {
			// 当前标签页已有文档，创建新标签页
			ui.addNewTab(filePath)
		}
	}, ui.window)

	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
	fileDialog.Show()
}

// loadPDF 加载 PDF 文件（PDFTab 方法）
func (tab *PDFTab) loadPDF(filePath string, ui *ViewerUI) {
	tab.showLoading(ui.tr.MsgLoading)

	err := tab.controller.OpenPDF(filePath)
	if err != nil {
		tab.showError(fmt.Sprintf(ui.tr.MsgLoadFailed, err))
		return
	}

	// 更新标签页标题
	tab.tabItem.Text = getFileName(filePath)
	ui.tabContainer.Refresh()

	tab.renderPage(ui)
	ui.updateStatusBar()
	ui.updateZoomLabel()
}

// onFirstPage 跳转到首页
func (ui *ViewerUI) onFirstPage() {
	currentTab := ui.getCurrentTab()
	if currentTab != nil {
		currentTab.onFirstPage(ui)
	}
}

func (tab *PDFTab) onFirstPage(ui *ViewerUI) {
	if tab.controller.FirstPage() {
		tab.renderPage(ui)
		ui.updateStatusBar()
	}
}

// onPrevPage 上一页
func (ui *ViewerUI) onPrevPage() {
	currentTab := ui.getCurrentTab()
	if currentTab != nil {
		currentTab.onPrevPage(ui)
	}
}

func (tab *PDFTab) onPrevPage(ui *ViewerUI) {
	if tab.controller.PrevPage() {
		tab.renderPage(ui)
		ui.updateStatusBar()
	}
}

// onNextPage 下一页
func (ui *ViewerUI) onNextPage() {
	currentTab := ui.getCurrentTab()
	if currentTab != nil {
		currentTab.onNextPage(ui)
	}
}

func (tab *PDFTab) onNextPage(ui *ViewerUI) {
	if tab.controller.NextPage() {
		tab.renderPage(ui)
		ui.updateStatusBar()
	}
}

// onLastPage 跳转到末页
func (ui *ViewerUI) onLastPage() {
	currentTab := ui.getCurrentTab()
	if currentTab != nil {
		currentTab.onLastPage(ui)
	}
}

func (tab *PDFTab) onLastPage(ui *ViewerUI) {
	if tab.controller.LastPage() {
		tab.renderPage(ui)
		ui.updateStatusBar()
	}
}

// onPageJump 页码跳转
func (ui *ViewerUI) onPageJump(text string) {
	currentTab := ui.getCurrentTab()
	if currentTab == nil {
		return
	}

	pageNum, err := strconv.Atoi(text)
	if err != nil {
		dialog.ShowError(fmt.Errorf(ui.tr.MsgInvalidPage), ui.window)
		return
	}

	err = currentTab.controller.GoToPage(pageNum)
	if err != nil {
		dialog.ShowError(err, ui.window)
		return
	}

	currentTab.renderPage(ui)
	ui.updateStatusBar()
}

// onZoomIn 放大
func (ui *ViewerUI) onZoomIn() {
	currentTab := ui.getCurrentTab()
	if currentTab != nil {
		currentTab.onZoomIn(ui)
	}
}

func (tab *PDFTab) onZoomIn(ui *ViewerUI) {
	tab.controller.ZoomIn()
	tab.renderPage(ui)
	ui.updateStatusBar()
	ui.updateZoomLabel()
}

// onZoomOut 缩小
func (ui *ViewerUI) onZoomOut() {
	currentTab := ui.getCurrentTab()
	if currentTab != nil {
		currentTab.onZoomOut(ui)
	}
}

func (tab *PDFTab) onZoomOut(ui *ViewerUI) {
	tab.controller.ZoomOut()
	tab.renderPage(ui)
	ui.updateStatusBar()
	ui.updateZoomLabel()
}

// onZoomReset 重置缩放
func (ui *ViewerUI) onZoomReset() {
	currentTab := ui.getCurrentTab()
	if currentTab != nil {
		currentTab.onZoomReset(ui)
	}
}

func (tab *PDFTab) onZoomReset(ui *ViewerUI) {
	tab.controller.ResetZoom()
	tab.renderPage(ui)
	ui.updateStatusBar()
	ui.updateZoomLabel()
}

// renderPage 渲染当前页面（PDFTab 方法）
func (tab *PDFTab) renderPage(ui *ViewerUI) {
	if !tab.controller.HasDocument() {
		return
	}

	// 不再显示"正在渲染"提示，直接渲染
	go func() {
		img, err := tab.controller.RenderCurrentPage()
		if err != nil {
			tab.showError(fmt.Sprintf(ui.tr.MsgRenderFailed, err))
			return
		}

		// 更新界面
		tab.imageCanvas.Image = img
		tab.imageCanvas.Refresh()
		tab.hideLoading() // 隐藏加载提示
	}()
}

// showLoading 显示加载提示（PDFTab 方法）
func (tab *PDFTab) showLoading(message string) {
	tab.loadingLabel.SetText(message)
	tab.loadingLabel.Show()
}

// hideLoading 隐藏加载提示（PDFTab 方法）
func (tab *PDFTab) hideLoading() {
	tab.loadingLabel.Hide()
}

// showError 显示错误信息（PDFTab 方法）
func (tab *PDFTab) showError(message string) {
	tab.loadingLabel.SetText(message)
}

// Show 显示窗口
func (ui *ViewerUI) Show() {
	ui.window.ShowAndRun()
}

// onShowShortcuts 显示快捷键说明
func (ui *ViewerUI) onShowShortcuts() {
	dialog.ShowInformation(ui.tr.DialogShortcutsTitle, ui.tr.DialogShortcutsText, ui.window)
}

// onShowAbout 显示关于信息
func (ui *ViewerUI) onShowAbout() {
	dialog.ShowInformation(ui.tr.DialogAboutTitle, ui.tr.DialogAboutText, ui.window)
}

// switchLanguage 切换语言
func (ui *ViewerUI) switchLanguage(lang Language) {
	if ui.currentLang == lang {
		return // 已经是当前语言，无需切换
	}

	// 更新语言设置
	ui.currentLang = lang
	ui.tr = GetTranslations(lang)

	// 更新窗口标题
	ui.window.SetTitle(ui.tr.WindowTitle)

	// 重建菜单栏
	ui.window.SetMainMenu(ui.createMenuBar())

	// 更新页码输入框提示
	ui.pageEntry.SetPlaceHolder(ui.tr.HintPageEntry)

	// 更新状态栏
	ui.updateStatusBar()

	// 更新缩放标签
	ui.updateZoomLabel()

	// 刷新窗口
	ui.window.Canvas().Refresh(ui.window.Content())
}

// scrollableCanvas 支持滚轮翻页的自定义 widget
type scrollableCanvas struct {
	widget.BaseWidget
	content      fyne.CanvasObject
	onScroll     func(scrolled *fyne.ScrollEvent)
	onDoubleTap  func()
}

func newScrollableCanvas(content fyne.CanvasObject, onScroll func(*fyne.ScrollEvent), onDoubleTap func()) *scrollableCanvas {
	sc := &scrollableCanvas{
		content:     content,
		onScroll:    onScroll,
		onDoubleTap: onDoubleTap,
	}
	sc.ExtendBaseWidget(sc)
	return sc
}

func (sc *scrollableCanvas) CreateRenderer() fyne.WidgetRenderer {
	return &scrollableCanvasRenderer{
		canvas:  sc,
		content: sc.content,
	}
}

func (sc *scrollableCanvas) Scrolled(ev *fyne.ScrollEvent) {
	if sc.onScroll != nil {
		sc.onScroll(ev)
	}
}

func (sc *scrollableCanvas) DoubleTapped(ev *fyne.PointEvent) {
	if sc.onDoubleTap != nil {
		sc.onDoubleTap()
	}
}

func (sc *scrollableCanvas) Tapped(ev *fyne.PointEvent) {
	// 单击事件不处理，仅响应双击
}

type scrollableCanvasRenderer struct {
	canvas  *scrollableCanvas
	content fyne.CanvasObject
}

func (r *scrollableCanvasRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
}

func (r *scrollableCanvasRenderer) MinSize() fyne.Size {
	return r.content.MinSize()
}

func (r *scrollableCanvasRenderer) Refresh() {
	canvas.Refresh(r.content)
}

func (r *scrollableCanvasRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

func (r *scrollableCanvasRenderer) Destroy() {}
