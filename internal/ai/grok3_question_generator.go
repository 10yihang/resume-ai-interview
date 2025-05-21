package ai

import (
	"context"
	"fmt"
	"log"

	"github.com/10yihang/resume-ai-interview/models"
)

// Grok3QuestionGenerator 用于使用Grok 3生成面试问题
type Grok3QuestionGenerator struct {
	client *Grok3Client
}

// NewGrok3QuestionGenerator 创建使用Grok 3的问题生成器
func NewGrok3QuestionGenerator(apiKey string) *Grok3QuestionGenerator {
	if apiKey == "" {
		log.Fatal("Grok 3 API密钥未配置")
	}

	client := NewGrok3Client(apiKey)
	return &Grok3QuestionGenerator{
		client: client,
	}
}

// GenerateQuestions 根据简历和JD生成面试问题
func (g *Grok3QuestionGenerator) GenerateQuestions(resume *models.Resume, jd *models.JobDescription) (*models.QuestionSet, error) {
	// 构建提示词
	prompt := buildQuestionPrompt(resume, jd)

	// 调用Grok 3 API
	resp, err := g.client.CreateChatCompletion(
		context.Background(),
		Grok3ChatRequest{
			Model: "grok-3", // 使用Grok 3模型
			Messages: []Grok3Message{
				{
					Role:    "system",
					Content: "你是一位经验丰富的HR面试官，需要根据简历和职位描述生成有针对性的面试问题。请生成10个问题，包括技术能力、项目经验、职业规划、团队协作等方面。问题要有针对性，能够考察候选人是否符合岗位需求。",
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
			MaxTokens:   2048,
			Temperature: 0.7,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("调用Grok 3接口生成问题失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("Grok 3返回了空的回复")
	}

	// 解析问题
	questionSet := parseQuestions(resume, jd, resp.Choices[0].Message.Content)
	return questionSet, nil
}
