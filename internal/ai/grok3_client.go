package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// Grok3Client 是Grok 3 API的客户端
type Grok3Client struct {
	client *openai.Client
}

// Grok3Message 表示Grok 3 API的消息格式
type Grok3Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Grok3ChatRequest 表示Grok 3 聊天请求
type Grok3ChatRequest struct {
	Model       string         `json:"model"`
	Messages    []Grok3Message `json:"messages"`
	MaxTokens   int            `json:"max_tokens,omitempty"`
	Temperature float32        `json:"temperature,omitempty"`
}

// Grok3ChatResponse 表示Grok 3 聊天回复，实际使用OpenAI的结构
type Grok3ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int          `json:"index"`
		Message      Grok3Message `json:"message"`
		FinishReason string       `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewGrok3Client 创建一个新的Grok 3 客户端
func NewGrok3Client(apiKey string) *Grok3Client {
	// 从环境变量获取API基础URL，如果没有设置则使用默认值
	baseURL := os.Getenv("GROK3_API_URL")
	if baseURL == "" {
		// X.AI的API与OpenAI兼容，但URL不同
		baseURL = "https://api.x.ai/v1"
	}

	// 配置OpenAI客户端
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	client := openai.NewClientWithConfig(config)

	// 打印API信息
	fmt.Printf("初始化Grok API客户端，基础URL: %s\n", baseURL)

	return &Grok3Client{
		client: client,
	}
}

// CreateChatCompletion 发送聊天请求到Grok 3 API
func (c *Grok3Client) CreateChatCompletion(ctx context.Context, request Grok3ChatRequest) (Grok3ChatResponse, error) {
	var response Grok3ChatResponse

	// 转换Grok3的请求格式为OpenAI的请求格式
	messages := make([]openai.ChatCompletionMessage, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 创建OpenAI的请求
	openaiRequest := openai.ChatCompletionRequest{
		Model:       request.Model,
		Messages:    messages,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
	}
	// fmt.Printf("Grok3 API请求：%v\n", openaiRequest)

	// 发送请求
	openaiResp, err := c.client.CreateChatCompletion(ctx, openaiRequest)
	if err != nil {
		return response, fmt.Errorf("错误：发送请求到Grok API失败：%w", err)
	}

	// 转换OpenAI的响应为Grok3的响应格式
	response.ID = openaiResp.ID
	response.Object = openaiResp.Object
	response.Created = openaiResp.Created
	response.Model = openaiResp.Model

	// 转换选择
	response.Choices = make([]struct {
		Index        int          `json:"index"`
		Message      Grok3Message `json:"message"`
		FinishReason string       `json:"finish_reason"`
	}, len(openaiResp.Choices))

	for i, choice := range openaiResp.Choices {
		response.Choices[i] = struct {
			Index        int          `json:"index"`
			Message      Grok3Message `json:"message"`
			FinishReason string       `json:"finish_reason"`
		}{
			Index: choice.Index,
			Message: Grok3Message{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
			FinishReason: string(choice.FinishReason),
		}
	}

	// 转换使用情况
	response.Usage.PromptTokens = openaiResp.Usage.PromptTokens
	response.Usage.CompletionTokens = openaiResp.Usage.CompletionTokens
	response.Usage.TotalTokens = openaiResp.Usage.TotalTokens

	return response, nil
}
