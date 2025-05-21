package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/10yihang/resume-ai-interview/models"
	"github.com/sashabaranov/go-openai"
)

// QuestionGenerator 用于生成面试问题
type QuestionGenerator struct {
	client *openai.Client
}

// NewQuestionGenerator 创建问题生成器
func NewQuestionGenerator(apiKey string) *QuestionGenerator {
	if apiKey == "" {
		log.Fatal("OpenAI API密钥未配置")
	}

	client := openai.NewClient(apiKey)
	return &QuestionGenerator{
		client: client,
	}
}

// GenerateQuestions 根据简历和JD生成面试问题
func (g *QuestionGenerator) GenerateQuestions(resume *models.Resume, jd *models.JobDescription) (*models.QuestionSet, error) {
	// 构建提示词
	prompt := buildQuestionPrompt(resume, jd)

	// 调用OpenAI API
	resp, err := g.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是一位经验丰富的HR面试官，需要根据简历和职位描述生成有针对性的面试问题。请生成10个问题，包括技术能力、项目经验、职业规划、团队协作等方面。问题要有针对性，能够考察候选人是否符合岗位需求。",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 2048,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("调用AI接口生成问题失败: %w", err)
	}

	// 解析问题
	questionSet := parseQuestions(resume, jd, resp.Choices[0].Message.Content)
	return questionSet, nil
}

// 构建问题生成的提示词
func buildQuestionPrompt(resume *models.Resume, jd *models.JobDescription) string {
	return fmt.Sprintf(`
请根据以下简历和职位描述生成10个针对性的面试问题：

==== 简历信息 ====
姓名: %s
技能: %s
教育经历: %s
工作经验: %s
简历内容: %s

==== 职位描述 ====
职位: %s
公司: %s
描述: %s
要求: %s

请确保问题涵盖以下几个方面：
1. 专业技能核实（3个问题）
2. 工作经验相关（3个问题）
3. 个人能力和团队协作（2个问题）
4. 职业规划（2个问题）

请以JSON格式输出，格式如下：
{
  "questions": [
    {
      "id": 1,
      "content": "问题内容",
      "category": "问题类别"
    }
  ]
}
`,
		resume.Name,
		strings.Join(resume.Skills, ", "),
		strings.Join(resume.Education, ", "),
		strings.Join(resume.Experience, ", "),
		resume.RawText,
		jd.Title,
		jd.Company,
		jd.Description,
		strings.Join(jd.Requirements, ", "),
	)
}

// 解析AI返回的问题
func parseQuestions(resume *models.Resume, jd *models.JobDescription, content string) *models.QuestionSet {
	// 提取JSON部分
	jsonStr := extractJSONFromContent(content)

	// 解析JSON
	var result struct {
		Questions []struct {
			ID       int    `json:"id"`
			Content  string `json:"content"`
			Category string `json:"category"`
		} `json:"questions"`
	}

	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		fmt.Printf("解析问题JSON失败: %v\n", err)
		// 解析失败时返回默认问题集
		return getDefaultQuestions(resume, jd)
	}

	// 转换为模型
	questions := make([]models.Question, 0, len(result.Questions))
	for _, q := range result.Questions {
		question := models.Question{
			ID:       q.ID,
			Content:  q.Content,
			Category: q.Category,
		}
		questions = append(questions, question)
	}

	// 如果没有解析到问题，返回默认问题集
	if len(questions) == 0 {
		return getDefaultQuestions(resume, jd)
	}

	return &models.QuestionSet{
		ResumeID:  filepath.Base(resume.FilePath),
		JDID:      filepath.Base(jd.FilePath),
		Questions: questions,
	}
}

// 提取JSON内容
func extractJSONFromContent(content string) string {
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

// 获取默认问题集
func getDefaultQuestions(resume *models.Resume, jd *models.JobDescription) *models.QuestionSet {
	questions := []models.Question{
		{ID: 1, Content: "请介绍一下你的技术背景和专长?", Category: "专业技能"},
		{ID: 2, Content: "你在简历中提到了" + getSkillsPrompt(resume) + "，能详细讲讲你在这方面的经验吗?", Category: "专业技能"},
		{ID: 3, Content: "请描述一个你最有挑战性的项目经历，以及你如何解决遇到的问题?", Category: "工作经验"},
		{ID: 4, Content: "你对" + getJobTitlePrompt(jd) + "这个职位的理解是什么?", Category: "职业规划"},
		{ID: 5, Content: "你认为自己适合这个职位的优势是什么?", Category: "个人能力"},
		{ID: 6, Content: "你如何与团队成员协作完成项目?", Category: "团队协作"},
		{ID: 7, Content: "你过去的工作经历中，有哪些经验可以应用到这个职位?", Category: "工作经验"},
		{ID: 8, Content: "你对技术发展趋势的看法是什么?", Category: "专业技能"},
		{ID: 9, Content: "你未来五年的职业规划是什么?", Category: "职业规划"},
		{ID: 10, Content: "你如何处理工作中的压力和挑战?", Category: "个人能力"},
	}

	return &models.QuestionSet{
		ResumeID:  filepath.Base(resume.FilePath),
		JDID:      filepath.Base(jd.FilePath),
		Questions: questions,
	}
}

// 从简历中提取技能关键词
func getSkillsPrompt(resume *models.Resume) string {
	if len(resume.Skills) > 0 {
		if len(resume.Skills) > 2 {
			return strings.Join(resume.Skills[:2], "、")
		}
		return strings.Join(resume.Skills, "、")
	}
	return "相关技能"
}

// 获取职位名称
func getJobTitlePrompt(jd *models.JobDescription) string {
	if jd.Title != "" {
		return jd.Title
	}
	return "这个"
}
