package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
)

// PDFEngine 封装 PDF 处理功能
type PDFEngine struct {
	filePath  string
	document  *fitz.Document
	pageCount int
}

// NewPDFEngine 创建 PDF 引擎实例
func NewPDFEngine(filePath string) (*PDFEngine, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("文件不存在: %s", filePath)
	}

	// 打开 PDF 文件
	doc, err := fitz.New(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 PDF: %w", err)
	}

	return &PDFEngine{
		filePath:  filePath,
		document:  doc,
		pageCount: doc.NumPage(),
	}, nil
}

// RenderPage 渲染指定页面为图像
func (e *PDFEngine) RenderPage(pageNum int, dpi int) (image.Image, error) {
	if pageNum < 1 || pageNum > e.pageCount {
		return nil, fmt.Errorf("页码超出范围: %d (总页数: %d)", pageNum, e.pageCount)
	}

	// go-fitz 页码从 0 开始
	img, err := e.document.ImageDPI(pageNum-1, float64(dpi))
	if err != nil {
		return nil, fmt.Errorf("渲染失败: %w", err)
	}

	return img, nil
}

// GetPageCount 返回总页数
func (e *PDFEngine) GetPageCount() int {
	return e.pageCount
}

// GetFilePath 返回文件路径
func (e *PDFEngine) GetFilePath() string {
	return e.filePath
}

// GetFileName 返回文件名（不含路径）
func (e *PDFEngine) GetFileName() string {
	return filepath.Base(e.filePath)
}

// GetFileSize 返回文件大小（字节）
func (e *PDFEngine) GetFileSize() (int64, error) {
	info, err := os.Stat(e.filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// Close 关闭文档
func (e *PDFEngine) Close() error {
	if e.document != nil {
		return e.document.Close()
	}
	return nil
}
