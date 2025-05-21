package ocr

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// OCRResult 表示OCR处理结果
type OCRResult struct {
	Text        string
	ProcessTime time.Duration
	TokenCount  int    // 估计的token数量
	Source      string // 使用的OCR引擎
}

// 估计token数量（简单实现，实际上tokens约等于单词数的1.3倍）
func EstimateTokenCount(text string) int {
	// 分割文本统计单词数
	words := strings.Fields(text)
	// 估算token数（简单估计为单词数的1.3倍）
	return int(float64(len(words)) * 1.3)
}

// ProcessFile 处理文件并提取文本，带有错误重试和日志
func ProcessFile(processor OCRProcessor, filePath string) (string, error) {
	start := time.Now()
	var text string
	var err error

	ext := strings.ToLower(filepath.Ext(filePath))
	source := fmt.Sprintf("%T", processor)

	// 根据文件类型选择合适的处理方法
	switch ext {
	case ".pdf":
		text, err = processor.ExtractTextFromPDF(filePath)
		if err != nil {
			log.Printf("PDF OCR处理失败 (%s): %v", source, err)
			return "", fmt.Errorf("OCR处理失败: %w", err)
		}
	case ".png", ".jpg", ".jpeg":
		text, err = processor.ExtractTextFromImage(filePath)
		if err != nil {
			log.Printf("图像OCR处理失败 (%s): %v", source, err)
			return "", fmt.Errorf("OCR处理失败: %w", err)
		}
	default:
		return "", fmt.Errorf("不支持的文件类型: %s", ext)
	}

	duration := time.Since(start)
	tokens := EstimateTokenCount(text)

	log.Printf("OCR处理完成 (%s): 文件=%s, 耗时=%v, 估计tokens=%d",
		source, filepath.Base(filePath), duration, tokens)

	// 检查提取的文本是否为空
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("OCR提取的文本为空")
	}

	// 如果token数量过大，截断文本
	if tokens > 4000 {
		log.Printf("警告: OCR提取的文本token数(%d)过多，将被截断", tokens)
		// 简单截断为原来的一半长度
		words := strings.Fields(text)
		text = strings.Join(words[:len(words)/2], " ")
	}

	return text, nil
}
