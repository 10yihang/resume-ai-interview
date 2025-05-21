// OCR包提供光学字符识别功能，用于从图像中提取文本
package ocr

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TesseractOCR 使用Tesseract OCR进行文字识别
type TesseractOCR struct {
	tesseractPath string // Tesseract可执行文件路径
	tempDir       string // 临时文件目录
}

// NewTesseractOCR 创建一个新的TesseractOCR实例
func NewTesseractOCR(tesseractPath string) *TesseractOCR {
	// 如果未指定路径，尝试使用默认路径
	if tesseractPath == "" {
		tesseractPath = "tesseract"
	}

	return &TesseractOCR{
		tesseractPath: tesseractPath,
		tempDir:       os.TempDir(),
	}
}

// ExtractTextFromPDF 从PDF文件中提取文本
func (t *TesseractOCR) ExtractTextFromPDF(pdfPath string) (string, error) {
	// 第1步：确认Tesseract是否已安装
	err := t.checkTesseractInstallation()
	if err != nil {
		return "", fmt.Errorf("Tesseract OCR未正确安装: %w", err)
	}

	// 第2步：将PDF转换为图像（需要使用额外的工具如pdftoppm或Ghostscript）
	// 这里使用了临时目录
	timestamp := time.Now().UnixNano()
	tempImagePrefix := filepath.Join(t.tempDir, fmt.Sprintf("pdf_image_%d", timestamp))

	// 调用PDF转图像工具（默认使用pdftoppm）
	// pdftoppm是一个常见工具，通常安装了poppler-utils就会有
	cmd := exec.Command("pdftoppm", "-png", pdfPath, tempImagePrefix)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("将PDF转换为图像失败: %w", err)
	}

	// 第3步：使用Tesseract OCR处理所有生成的图像文件
	var allText strings.Builder

	// 查找所有生成的PNG文件
	imagePattern := tempImagePrefix + "*.png"
	imageFiles, err := filepath.Glob(imagePattern)
	if err != nil {
		return "", fmt.Errorf("查找生成的图像文件失败: %w", err)
	}

	// 处理每个图像文件
	for _, imgFile := range imageFiles {
		// 为每个图像创建一个临时输出文件
		outputBase := imgFile + "_ocr"
		outputFile := outputBase + ".txt"

		// 执行Tesseract OCR
		cmd = exec.Command(t.tesseractPath, imgFile, outputBase)
		err = cmd.Run()
		if err != nil {
			return "", fmt.Errorf("在图像上执行OCR失败: %w", err)
		}

		// 读取OCR结果
		textBytes, err := os.ReadFile(outputFile)
		if err != nil {
			return "", fmt.Errorf("读取OCR结果失败: %w", err)
		}

		// 追加到总文本
		allText.Write(textBytes)
		allText.WriteString("\n")

		// 清理临时文件
		os.Remove(outputFile)
		os.Remove(imgFile)
	}

	return allText.String(), nil
}

// ExtractTextFromImage 从图像文件中提取文本
func (t *TesseractOCR) ExtractTextFromImage(imagePath string) (string, error) {
	// 确认Tesseract是否已安装
	err := t.checkTesseractInstallation()
	if err != nil {
		return "", fmt.Errorf("Tesseract OCR未正确安装: %w", err)
	}

	// 创建临时输出文件
	timestamp := time.Now().UnixNano()
	outputBase := filepath.Join(t.tempDir, fmt.Sprintf("img_ocr_%d", timestamp))
	outputFile := outputBase + ".txt"

	// 执行Tesseract OCR
	cmd := exec.Command(t.tesseractPath, imagePath, outputBase)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("执行OCR失败: %w", err)
	}

	// 读取OCR结果
	textBytes, err := os.ReadFile(outputFile)
	if err != nil {
		return "", fmt.Errorf("读取OCR结果失败: %w", err)
	}

	// 清理临时文件
	os.Remove(outputFile)

	return string(textBytes), nil
}

// checkTesseractInstallation 检查Tesseract OCR是否已安装
func (t *TesseractOCR) checkTesseractInstallation() error {
	cmd := exec.Command(t.tesseractPath, "--version")
	return cmd.Run()
}
