package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIClient 是OpenAI API的客户端
type OpenAIClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// OpenAIMessage 表示OpenAI API的消息格式
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIChatRequest 表示OpenAI聊天请求
type OpenAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float32         `json:"temperature,omitempty"`
}

// OpenAIChatResponse 表示OpenAI聊天回复
type OpenAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewOpenAIClient 创建一个新的OpenAI客户端
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		httpClient: &http.Client{
			Timeout: time.Second * 60,
		},
	}
}

// CreateChatCompletion 发送聊天请求到OpenAI API
func (c *OpenAIClient) CreateChatCompletion(ctx context.Context, request OpenAIChatRequest) (OpenAIChatResponse, error) {
	var response OpenAIChatResponse

	jsonReq, err := json.Marshal(request)
	if err != nil {
		return response, fmt.Errorf("错误：无法序列化请求：%w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewBuffer(jsonReq))
	if err != nil {
		return response, fmt.Errorf("错误：创建HTTP请求失败：%w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("错误：发送HTTP请求失败：%w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("错误：读取响应失败：%w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("错误：API返回非200状态码：%d，响应：%s", resp.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("错误：解析响应失败：%w", err)
	}

	return response, nil
}
