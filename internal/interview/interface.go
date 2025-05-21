/*
 * @author: yihang_01
 * @Date: 2025-05-21 16:44:20
 * @LastEditTime: 2025-05-21 16:55:09
 * QwQ 加油加油
 */
package interview

import (
	"github.com/10yihang/resume-ai-interview/models"
)

// AnswerEvaluatorInterface 定义了面试答案评估器的接口
type AnswerEvaluatorInterface interface {
	EvaluateAnswer(question models.Question, answer models.Answer, jd *models.JobDescription) (*models.Evaluation, error)
}

// GetAnswerEvaluator 根据配置返回适当的答案评估器
func GetAnswerEvaluator(apiKey string, useGrok bool) AnswerEvaluatorInterface {
	if apiKey == "" {
		// 如果没有API密钥，使用模拟评估器
		return NewMockAnswerEvaluator()
	}

	if useGrok {
		// 使用Grok 3
		return NewGrok3AnswerEvaluator(apiKey)
	}
	// 默认使用OpenAI
	return NewAnswerEvaluator(apiKey)
}
