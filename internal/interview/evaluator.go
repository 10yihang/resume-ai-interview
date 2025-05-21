package interview

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/10yihang/resume-ai-interview/models"
	"github.com/sashabaranov/go-openai"
)

// AnswerEvaluator 用于评估面试答案
type AnswerEvaluator struct {
	client *openai.Client
}

// NewAnswerEvaluator 创建答案评估器
func NewAnswerEvaluator(apiKey string) *AnswerEvaluator {
	if apiKey == "" {
		log.Fatal("OpenAI API密钥未配置")
	}

	client := openai.NewClient(apiKey)
	return &AnswerEvaluator{
		client: client,
	}
}

// EvaluateAnswer 评估面试回答
func (e *AnswerEvaluator) EvaluateAnswer(question models.Question, answer models.Answer, jd *models.JobDescription) (*models.Evaluation, error) {
	// 构建提示词
	prompt := buildEvaluationPrompt(question, answer, jd)

	// 调用OpenAI API
	resp, err := e.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是一位专业的HR面试官，需要评估候选人的面试回答。请基于面试问题、候选人的回答以及职位要求，评估回答质量，给出分数（1-10）、反馈和改进建议。",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 1024,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("调用AI接口评估答案失败: %w", err)
	}

	// 解析评估结果
	evaluation := parseEvaluation(answer, resp.Choices[0].Message.Content)
	return evaluation, nil
}

// 构建评估提示词
func buildEvaluationPrompt(question models.Question, answer models.Answer, jd *models.JobDescription) string {
	return fmt.Sprintf(`
请评估以下面试回答：

==== 面试问题 ====
问题：%s
问题类别：%s

==== 候选人回答 ====
%s

==== 职位信息 ====
职位：%s
公司：%s
职位要求：%s

请评估回答质量，并以下面的JSON格式给出评分和反馈：
{
  "score": 7,
  "feedback": "你的评价内容...",
  "suggestions": "改进建议..."
}

评分标准：
1-3分：不满足基本要求，回答模糊或错误
4-6分：基本符合要求，但缺乏深度或细节
7-8分：良好的回答，体现了专业知识和经验
9-10分：优秀的回答，全面、深入且有洞察力
`,
		question.Content,
		question.Category,
		answer.Content,
		jd.Title,
		jd.Company,
		strings.Join(jd.Requirements, ", "),
	)
}

// 解析AI返回的评估结果
func parseEvaluation(answer models.Answer, content string) *models.Evaluation {
	// 提取JSON部分
	jsonStr := extractEvaluationJSON(content)

	var result struct {
		Score       int    `json:"score"`
		Feedback    string `json:"feedback"`
		Suggestions string `json:"suggestions"`
	}

	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		fmt.Printf("解析评估JSON失败: %v\n", err)
		// 尝试修复JSON格式
		fixedJSON := tryFixJSON(jsonStr)
		err = json.Unmarshal([]byte(fixedJSON), &result)

		if err != nil {
			// 如果仍然失败，返回默认评估
			return &models.Evaluation{
				AnswerID:    answer.QuestionID,
				Score:       7,
				Feedback:    "回答展示了对问题的理解，但可以提供更多具体例子。",
				Suggestions: "建议增加更多实际工作中的案例来支持你的观点。",
			}
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
		result.Feedback = "回答展示了对问题的理解，但可以提供更多具体例子。"
	}
	if result.Suggestions == "" {
		result.Suggestions = "建议增加更多实际工作中的案例来支持你的观点。"
	}

	return &models.Evaluation{
		AnswerID:    answer.QuestionID,
		Score:       result.Score,
		Feedback:    result.Feedback,
		Suggestions: result.Suggestions,
	}
}

// 从评估内容中提取JSON
func extractEvaluationJSON(content string) string {
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

// 尝试修复JSON格式
func tryFixJSON(jsonStr string) string {
	// 移除可能的非法字符
	r := strings.NewReplacer("\n", " ", "\r", " ", "\t", " ")
	jsonStr = r.Replace(jsonStr)

	// 确保有引号的键
	jsonStr = strings.Replace(jsonStr, "score:", "\"score\":", -1)
	jsonStr = strings.Replace(jsonStr, "feedback:", "\"feedback\":", -1)
	jsonStr = strings.Replace(jsonStr, "suggestions:", "\"suggestions\":", -1)

	// 确保值有引号（简化处理）

	return jsonStr
}
