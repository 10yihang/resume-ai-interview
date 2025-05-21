package ai

import (
	"github.com/10yihang/resume-ai-interview/models"
)

// QuestionGeneratorInterface 定义了问题生成器的接口
type QuestionGeneratorInterface interface {
	GenerateQuestions(resume *models.Resume, jd *models.JobDescription) (*models.QuestionSet, error)
}

// GetQuestionGenerator 根据配置返回适当的问题生成器
func GetQuestionGenerator(apiKey string, useGrok bool) QuestionGeneratorInterface {
	if apiKey == "" {
		// 如果没有API密钥，使用模拟生成器
		return NewMockQuestionGenerator()
	}

	if useGrok {
		// 使用Grok 3
		return NewGrok3QuestionGenerator(apiKey)
	}

	// 默认使用OpenAI
	return NewQuestionGenerator(apiKey)
}
