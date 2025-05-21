package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// OCRSpaceAPI 使用免费的OCR.space API进行文字识别
// https://ocr.space/OCRAPI
type OCRSpaceAPI struct {
	apiKey string
}

// OCRSpaceResponse API响应结构
type OCRSpaceResponse struct {
	ParsedResults []struct {
		ParsedText string `json:"ParsedText"`
	} `json:"ParsedResults"`
	IsErroredOnProcessing bool   `json:"IsErroredOnProcessing"`
	ErrorMessage          string `json:"ErrorMessage"`
}

// NewOCRSpaceAPI 创建一个新的OCRSpaceAPI实例
func NewOCRSpaceAPI(apiKey string) *OCRSpaceAPI {
	return &OCRSpaceAPI{
		apiKey: apiKey,
	}
}

// ExtractTextFromPDF 从PDF文件中提取文本
func (o *OCRSpaceAPI) ExtractTextFromPDF(pdfPath string) (string, error) {
	return o.extractTextFromFile(pdfPath)
}

// ExtractTextFromImage 从图像文件中提取文本
func (o *OCRSpaceAPI) ExtractTextFromImage(imagePath string) (string, error) {
	return o.extractTextFromFile(imagePath)
}

// extractTextFromFile 从文件中提取文本（支持PDF、PNG、JPG等）
func (o *OCRSpaceAPI) extractTextFromFile(filePath string) (string, error) {
	// 准备请求
	url := "https://api.ocr.space/parse/image"
	method := "POST"

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加API密钥
	_ = writer.WriteField("apikey", o.apiKey)

	// 添加其他参数
	_ = writer.WriteField("language", "chs") // 简体中文和英文
	_ = writer.WriteField("isOverlayRequired", "false")
	_ = writer.WriteField("filetype", filepath.Ext(filePath)[1:]) // 去掉前面的点

	// 添加文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("创建表单文件失败: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("复制文件内容失败: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("关闭表单写入器失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON响应
	var ocrResp OCRSpaceResponse
	err = json.Unmarshal(respBody, &ocrResp)
	if err != nil {
		return "", fmt.Errorf("解析JSON响应失败: %w", err)
	}

	// 检查是否有错误
	if ocrResp.IsErroredOnProcessing {
		return "", fmt.Errorf("OCR处理错误: %s", ocrResp.ErrorMessage)
	}

	// 提取文本
	var allText string
	for _, result := range ocrResp.ParsedResults {
		allText += result.ParsedText
	}

	return allText, nil
}
