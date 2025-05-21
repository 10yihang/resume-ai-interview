package ocr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOCR(t *testing.T) {
	// 获取OCR API密钥
	apiKey := os.Getenv("OCR_SPACE_API_KEY")

	// 创建测试目录
	tempDir, err := os.MkdirTemp("", "ocr-test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试文本文件
	testText := "这是一个OCR测试\nThis is an OCR test\n123456"
	textFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(textFile, []byte(testText), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 测试OCRSpaceAPI (仅当提供了API密钥时)
	if apiKey != "" {
		t.Run("TestOCRSpaceAPI", func(t *testing.T) {
			ocrProcessor := NewOCRSpaceAPI(apiKey)

			// 对于真正的OCR测试，需要准备一个测试PDF或图像文件
			// 这里仅检查接口实现
			if ocrProcessor == nil {
				t.Fatal("创建OCRSpaceAPI处理器失败")
			}

			t.Logf("成功创建OCRSpaceAPI处理器")
		})
	} else {
		t.Log("未提供OCR_SPACE_API_KEY，跳过OCRSpaceAPI测试")
	}

	// 测试TesseractOCR (不检查实际的OCR功能，只测试接口)
	t.Run("TestTesseractOCR", func(t *testing.T) {
		ocrProcessor := NewTesseractOCR("")

		if ocrProcessor == nil {
			t.Fatal("创建TesseractOCR处理器失败")
		}

		t.Logf("成功创建TesseractOCR处理器")
	})

	// 测试GetOCRProcessor函数
	t.Run("TestGetOCRProcessor", func(t *testing.T) {
		// 测试基于API密钥的选择
		processor := GetOCRProcessor(apiKey, "")
		if processor == nil {
			t.Fatal("GetOCRProcessor返回nil")
		}

		// 测试无API密钥的情况
		tesseractProcessor := GetOCRProcessor("", "tesseract")
		if tesseractProcessor == nil {
			t.Fatal("GetOCRProcessor(无API密钥)返回nil")
		}

		t.Logf("成功测试GetOCRProcessor")
	})
}
