package ocr

// OCRProcessor 定义了OCR文本提取处理器的接口
type OCRProcessor interface {
	// ExtractTextFromPDF 从PDF文件中提取文本
	ExtractTextFromPDF(pdfPath string) (string, error)

	// ExtractTextFromImage 从图像文件中提取文本
	ExtractTextFromImage(imagePath string) (string, error)
}

// GetOCRProcessor 根据配置返回合适的OCR处理器
// 如果有OCR.space API密钥，优先使用云端OCR，否则使用本地Tesseract
func GetOCRProcessor(ocrAPIKey string, tesseractPath string) OCRProcessor {
	if ocrAPIKey != "" {
		return NewOCRSpaceAPI(ocrAPIKey)
	}
	return NewTesseractOCR(tesseractPath)
}
