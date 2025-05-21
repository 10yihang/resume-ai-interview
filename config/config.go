package config

import (
	"fmt"
	"os"
)

// Config 保存应用程序配置信息
type Config struct {
	APIKey        string
	UseGrok       bool
	MaxFileSize   int64
	DataDir       string
	OCRAPIKey     string
	TesseractPath string
	UseOCR        bool
}

// NewConfig 创建一个新的配置实例
func NewConfig() *Config {
	// 优先尝试使用Grok 3 API密钥
	apiKey := getEnvOrDefault("GROK3_API_KEY", "")
	useGrok := true

	// 如果没有配置Grok 3 API密钥，尝试使用OpenAI API密钥
	if apiKey == "" {
		apiKey = getEnvOrDefault("OPENAI_API_KEY", "")
		useGrok = false
	}

	// OCR配置
	ocrAPIKey := getEnvOrDefault("OCR_SPACE_API_KEY", "")
	tesseractPath := getEnvOrDefault("TESSERACT_PATH", "tesseract")
	useOCR := getEnvOrDefault("USE_OCR", "true") == "true"

	config := &Config{
		APIKey:        apiKey,
		UseGrok:       useGrok,
		MaxFileSize:   getEnvAsInt64OrDefault("MAX_FILE_SIZE", 10*1024*1024), // 默认10MB
		DataDir:       getEnvOrDefault("DATA_DIR", "./data"),
		OCRAPIKey:     ocrAPIKey,
		TesseractPath: tesseractPath,
		UseOCR:        useOCR,
	}

	// 打印配置信息
	if config.UseGrok {
		fmt.Printf("使用Grok 3 API，API密钥长度: %d\n", len(config.APIKey))
	} else if config.APIKey != "" {
		fmt.Printf("使用OpenAI API，API密钥长度: %d\n", len(config.APIKey))
	} else {
		fmt.Println("警告: 未配置API密钥，将使用模拟模式")
	}

	return config
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	config := NewConfig()

	if config.APIKey == "" {
		return config, fmt.Errorf("未找到API密钥配置，请设置GROK3_API_KEY或OPENAI_API_KEY环境变量")
	}

	return config, nil
}

// getEnvOrDefault 获取环境变量或返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt64OrDefault 获取环境变量并转换为int64或返回默认值
func getEnvAsInt64OrDefault(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if parsed, err := parseInt64(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// parseInt64 将字符串解析为int64
func parseInt64(value string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(value, "%d", &result)
	return result, err
}
