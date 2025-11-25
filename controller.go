package main

import (
	"fmt"
	"image"
)

// Controller 管理 PDF 阅读器的状态和逻辑
type Controller struct {
	engine      *PDFEngine
	currentPage int
	zoomLevel   float64
	baseDPI     int
}

// NewController 创建控制器实例
func NewController() *Controller {
	return &Controller{
		currentPage: 1,
		zoomLevel:   1.0,
		baseDPI:     150, // 默认 150 DPI
	}
}

// OpenPDF 打开 PDF 文件
func (c *Controller) OpenPDF(filePath string) error {
	engine, err := NewPDFEngine(filePath)
	if err != nil {
		return err
	}

	c.engine = engine
	c.currentPage = 1
	c.zoomLevel = 1.0

	return nil
}

// HasDocument 检查是否已加载文档
func (c *Controller) HasDocument() bool {
	return c.engine != nil
}

// GetCurrentPage 获取当前页码
func (c *Controller) GetCurrentPage() int {
	return c.currentPage
}

// GetPageCount 获取总页数
func (c *Controller) GetPageCount() int {
	if c.engine == nil {
		return 0
	}
	return c.engine.GetPageCount()
}

// GetZoomLevel 获取缩放级别
func (c *Controller) GetZoomLevel() float64 {
	return c.zoomLevel
}

// NextPage 下一页
func (c *Controller) NextPage() bool {
	if c.engine == nil {
		return false
	}

	if c.currentPage < c.engine.GetPageCount() {
		c.currentPage++
		return true
	}
	return false
}

// PrevPage 上一页
func (c *Controller) PrevPage() bool {
	if c.engine == nil {
		return false
	}

	if c.currentPage > 1 {
		c.currentPage--
		return true
	}
	return false
}

// FirstPage 跳转到第一页
func (c *Controller) FirstPage() bool {
	if c.engine == nil {
		return false
	}

	if c.currentPage != 1 {
		c.currentPage = 1
		return true
	}
	return false
}

// LastPage 跳转到最后一页
func (c *Controller) LastPage() bool {
	if c.engine == nil {
		return false
	}

	lastPage := c.engine.GetPageCount()
	if c.currentPage != lastPage {
		c.currentPage = lastPage
		return true
	}
	return false
}

// GoToPage 跳转到指定页
func (c *Controller) GoToPage(pageNum int) error {
	if c.engine == nil {
		return fmt.Errorf("未打开文档")
	}

	if pageNum < 1 || pageNum > c.engine.GetPageCount() {
		return fmt.Errorf("页码超出范围: %d (1-%d)", pageNum, c.engine.GetPageCount())
	}

	c.currentPage = pageNum
	return nil
}

// SetZoom 设置缩放级别
func (c *Controller) SetZoom(level float64) {
	if level < 0.5 {
		level = 0.5
	} else if level > 3.0 {
		level = 3.0
	}
	c.zoomLevel = level
}

// ZoomIn 放大
func (c *Controller) ZoomIn() {
	newZoom := c.zoomLevel * 1.25
	if newZoom > 3.0 {
		newZoom = 3.0
	}
	c.zoomLevel = newZoom
}

// ZoomOut 缩小
func (c *Controller) ZoomOut() {
	newZoom := c.zoomLevel / 1.25
	if newZoom < 0.5 {
		newZoom = 0.5
	}
	c.zoomLevel = newZoom
}

// ResetZoom 重置缩放
func (c *Controller) ResetZoom() {
	c.zoomLevel = 1.0
}

// RenderCurrentPage 渲染当前页面
func (c *Controller) RenderCurrentPage() (image.Image, error) {
	if c.engine == nil {
		return nil, fmt.Errorf("未打开文档")
	}

	dpi := int(float64(c.baseDPI) * c.zoomLevel)
	return c.engine.RenderPage(c.currentPage, dpi)
}

// GetStatusText 获取状态栏文本
func (c *Controller) GetStatusText(tr *Translations) string {
	if c.engine == nil {
		return tr.StatusNoDocument
	}

	// 获取文件名
	fileName := c.engine.GetFileName()

	// 获取文件大小
	fileSize, err := c.engine.GetFileSize()
	fileSizeStr := ""
	if err == nil {
		if fileSize < 1024 {
			fileSizeStr = fmt.Sprintf("%d B", fileSize)
		} else if fileSize < 1024*1024 {
			fileSizeStr = fmt.Sprintf("%.1f KB", float64(fileSize)/1024)
		} else {
			fileSizeStr = fmt.Sprintf("%.1f MB", float64(fileSize)/(1024*1024))
		}
	}

	return fmt.Sprintf("%s  |  %s %d / %d %s  |  %s: %d%%  |  %s: %s",
		fileName,
		tr.StatusPage,
		c.currentPage,
		c.engine.GetPageCount(),
		"页",
		tr.StatusZoom,
		int(c.zoomLevel*100),
		tr.StatusSize,
		fileSizeStr)
}
