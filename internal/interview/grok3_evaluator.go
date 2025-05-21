/*
 * @author: yihang_01
 * @Date: 2025-05-21 16:39:39
 * @LastEditTime: 2025-05-21 16:59:41
 * QwQ 加油加油
 */
package interview

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/10yihang/resume-ai-interview/internal/ai"
	"github.com/10yihang/resume-ai-interview/models"
)

// Grok3AnswerEvaluator 用于使用Grok 3评估面试答案
type Grok3AnswerEvaluator struct {
	client *ai.Grok3Client
}

// NewGrok3AnswerEvaluator 创建使用Grok 3的答案评估器
func NewGrok3AnswerEvaluator(apiKey string) *Grok3AnswerEvaluator {
	if apiKey == "" {
		log.Fatal("Grok 3 API密钥未配置")
	}

	client := ai.NewGrok3Client(apiKey)
	return &Grok3AnswerEvaluator{
		client: client,
	}
}

// EvaluateAnswer 评估面试回答
func (e *Grok3AnswerEvaluator) EvaluateAnswer(question models.Question, answer models.Answer, jd *models.JobDescription) (*models.Evaluation, error) {
	// 构建提示词
	prompt := buildEvaluationPrompt(question, answer, jd)

	// 调用Grok 3 API
	resp, err := e.client.CreateChatCompletion(
		context.Background(),
		ai.Grok3ChatRequest{
			Model: "grok-3", // 使用Grok 3模型
			Messages: []ai.Grok3Message{
				{
					Role:    "system",
					Content: "你是一位专业的HR面试官，需要评估候选人的面试回答。请基于面试问题、候选人的回答以及职位要求，评估回答质量，给出分数（1-10）、反馈和改进建议。",
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
			MaxTokens:   1024,
			Temperature: 0.5,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("调用Grok 3接口评估回答失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("Grok 3返回了空的回复")
	}
	// 解析评估
	evaluation := parseGrokEvaluation(answer, resp.Choices[0].Message.Content)
	return evaluation, nil
}

// parseGrokEvaluation 解析Grok 3返回的评估结果
func parseGrokEvaluation(answer models.Answer, content string) *models.Evaluation {
	// 提取JSON部分
	jsonContent := extractJSONFromEvalContent(content)

	var result struct {
		Score       int    `json:"score"`
		Feedback    string `json:"feedback"`
		Suggestions string `json:"suggestions"`
	}

	err := json.Unmarshal([]byte(jsonContent), &result)
	if err != nil {
		fmt.Printf("解析Grok 3返回的JSON失败：%v\n", err)
		// 尝试修复JSON
		fixedJson := tryFixEvaluationJSON(jsonContent)
		err = json.Unmarshal([]byte(fixedJson), &result)

		if err != nil {
			// 如果仍然失败，通过文本分析提取评估信息
			return extractEvaluationFromText(answer, content)
		}
	}

	// 验证分数范围
	if result.Score < 1 {
		result.Score = 1
	} else if result.Score > 10 {
		result.Score = 10
	}

	// 确保有反馈和建议
	if result.Feedback == "" {
		result.Feedback = "你的回答需要进一步改进。"
	}
	if result.Suggestions == "" {
		result.Suggestions = "尝试提供更多具体的例子和细节，使用STAR方法（情境、任务、行动、结果）来结构化回答。"
	}

	return &models.Evaluation{
		AnswerID:    answer.QuestionID,
		Score:       result.Score,
		Feedback:    result.Feedback,
		Suggestions: result.Suggestions,
	}
}

// extractJSONFromEvalContent 从评估内容中提取JSON
func extractJSONFromEvalContent(content string) string {
	// 查找JSON开始位置
	start := strings.Index(content, "{")
	if start < 0 {
		// 尝试查找Markdown代码块
		codeStart := strings.Index(content, "```json")
		if codeStart >= 0 {
			// 从代码块开始查找JSON
			codeContentStart := codeStart + 7 // "```json" 长度为7
			jsonStart := strings.Index(content[codeContentStart:], "{")
			if jsonStart >= 0 {
				start = codeContentStart + jsonStart
			}
		} else {
			// 尝试查找普通代码块
			codeStart := strings.Index(content, "```")
			if codeStart >= 0 {
				codeContentStart := codeStart + 3 // "```" 长度为3
				jsonStart := strings.Index(content[codeContentStart:], "{")
				if jsonStart >= 0 {
					start = codeContentStart + jsonStart
				}
			}
		}
	}

	// 如果找不到JSON开始，返回原始内容
	if start < 0 {
		return content
	}

	// 查找JSON结束位置
	end := -1
	braceCount := 0
	for i := start; i < len(content); i++ {
		if content[i] == '{' {
			braceCount++
		} else if content[i] == '}' {
			braceCount--
			if braceCount == 0 {
				end = i + 1
				break
			}
		}
	}

	// 如果找不到匹配的结束大括号，返回从开始到结尾的内容
	if end < 0 {
		end = len(content)
	}

	// 提取JSON内容
	jsonStr := content[start:end]

	// 清理可能的代码块标记
	jsonStr = strings.TrimSuffix(jsonStr, "```")

	return jsonStr
}

// tryFixEvaluationJSON 尝试修复评估JSON
func tryFixEvaluationJSON(jsonStr string) string {
	// 移除可能的非法字符
	r := strings.NewReplacer("\n", " ", "\r", " ", "\t", " ")
	jsonStr = r.Replace(jsonStr)

	// 确保有引号的键
	jsonStr = strings.Replace(jsonStr, "score:", "\"score\":", -1)
	jsonStr = strings.Replace(jsonStr, "feedback:", "\"feedback\":", -1)
	jsonStr = strings.Replace(jsonStr, "suggestions:", "\"suggestions\":", -1)

	// 确保值有引号
	// 这里简化处理，实际可能需要更复杂的正则处理

	return jsonStr
}

// extractEvaluationFromText 从文本中提取评估信息
func extractEvaluationFromText(answer models.Answer, content string) *models.Evaluation {
	// 默认评估
	eval := &models.Evaluation{
		AnswerID:    answer.QuestionID,
		Score:       6,
		Feedback:    "回答基本符合要求，但可以提供更多具体的例子和细节。",
		Suggestions: "考虑使用STAR方法（情境、任务、行动、结果）来结构化你的回答，使其更有条理。",
	}

	// 尝试从文本中提取分数
	scoreIdx := strings.Index(strings.ToLower(content), "score")
	if scoreIdx >= 0 {
		// 在"score"后寻找数字
		for i := scoreIdx + 5; i < len(content); i++ {
			if content[i] >= '0' && content[i] <= '9' {
				score := int(content[i] - '0')
				if score > 0 && score <= 10 {
					eval.Score = score
				}
				break
			}
		}
	}

	// 尝试提取反馈
	feedbackIdx := strings.Index(strings.ToLower(content), "feedback")
	if feedbackIdx >= 0 {
		feedbackEnd := strings.Index(content[feedbackIdx:], "\n\n")
		if feedbackEnd > 0 {
			feedback := content[feedbackIdx+8 : feedbackIdx+feedbackEnd]
			feedback = strings.TrimSpace(feedback)
			if len(feedback) > 0 {
				eval.Feedback = feedback
			}
		}
	}

	// 尝试提取建议
	suggestionsIdx := strings.Index(strings.ToLower(content), "suggestion")
	if suggestionsIdx >= 0 {
		suggestionsEnd := strings.Index(content[suggestionsIdx:], "\n\n")
		if suggestionsEnd > 0 {
			suggestions := content[suggestionsIdx+11 : suggestionsIdx+suggestionsEnd]
			suggestions = strings.TrimSpace(suggestions)
			if len(suggestions) > 0 {
				eval.Suggestions = suggestions
			}
		}
	}

	return eval
}
