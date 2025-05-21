/*
 * @author: yihang_01
 * @Date: 2025-05-21 19:13:37
 * @LastEditTime: 2025-05-21 19:26:55
 * QwQ 加油加油
 */
package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/10yihang/resume-ai-interview/internal/ocr"
)

// FileParser 定义了文件解析器的接口
type FileParser interface {
	// ParseFile 从文件中解析内容
	ParseFile(filePath string) (string, error)
}

// ResumeFileParser 使用OCR和传统方法解析简历文件
type ResumeFileParser struct {
	ocrProcessor ocr.OCRProcessor
	useOCR       bool
}

// NewResumeFileParser 创建一个新的简历文件解析器
func NewResumeFileParser(ocrProcessor ocr.OCRProcessor, useOCR bool) *ResumeFileParser {
	return &ResumeFileParser{
		ocrProcessor: ocrProcessor,
		useOCR:       useOCR,
	}
}

// ParseFile 解析简历文件
func (p *ResumeFileParser) ParseFile(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	// 根据文件扩展名选择处理方式
	switch ext {
	case ".pdf":
		// 优先使用OCR处理PDF以解决token限制问题
		if p.useOCR && p.ocrProcessor != nil {
			// 使用增强的OCR处理功能
			text, err := ocr.ProcessFile(p.ocrProcessor, filePath)
			if err == nil {
				return text, nil
			}
			// OCR失败时记录错误并使用传统方法
			fmt.Printf("OCR处理PDF失败: %v，尝试使用传统解析方法\n", err)
		}
		// 当OCR未启用或失败时，使用传统方法
		return extractTextFromPDF(filePath)
	case ".txt":
		// 文本文件直接读取
		return extractTextFromTXT(filePath)
	case ".png", ".jpg", ".jpeg":
		// 图像文件使用OCR
		if p.ocrProcessor != nil {
			return ocr.ProcessFile(p.ocrProcessor, filePath)
		}
		return "", fmt.Errorf("无法处理图像文件：OCR处理器未初始化")
	default:
		return "", fmt.Errorf("不支持的文件格式: %s", ext)
	}
}
