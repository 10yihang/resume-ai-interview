package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/10yihang/resume-ai-interview/internal/ai"
	"github.com/10yihang/resume-ai-interview/models"
)

// AITextParser 使用AI解析文本
type AITextParser struct {
	client     interface{} // 可以是OpenAI或Grok3客户端
	useGrok    bool
	apiKey     string
	fileParser FileParser // 文件解析器
}

// NewAITextParser 创建一个新的AI文本解析器
func NewAITextParser(apiKey string, useGrok bool, fileParser FileParser) *AITextParser {
	return &AITextParser{
		apiKey:     apiKey,
		useGrok:    useGrok,
		fileParser: fileParser,
	}
}

// ParseResumeText 使用AI解析简历文本
func (p *AITextParser) ParseResumeText(text string) (*models.Resume, error) {
	if p.apiKey == "" {
		// 如果没有API密钥，仅返回原始文本
		return &models.Resume{
			RawText: text,
		}, nil
	}

	// 构建提示词
	prompt := buildResumeParsePrompt(text)

	var content string
	var err error

	// 根据配置选择使用Grok3或OpenAI
	if p.useGrok {
		content, err = p.callGrok3API(prompt)
	} else {
		content, err = p.callOpenAIAPI(prompt)
	}

	if err != nil {
		return nil, fmt.Errorf("AI解析简历失败: %w", err)
	}

	// 解析AI返回的JSON
	resume, err := parseResumeJSON(content, text)
	if err != nil {
		return nil, fmt.Errorf("解析AI返回的JSON失败: %w", err)
	}

	return resume, nil
}

// ParseResumeFile 使用AI解析简历文件
func (p *AITextParser) ParseResumeFile(filePath string) (*models.Resume, error) {
	if p.fileParser == nil {
		return nil, fmt.Errorf("文件解析器未初始化")
	}

	// 从文件中提取文本
	text, err := p.fileParser.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("文件解析失败: %w", err)
	}

	// 使用提取的文本解析简历
	resume, err := p.ParseResumeText(text)
	if err != nil {
		return nil, err
	}

	// 设置文件路径
	resume.FilePath = filePath
	return resume, nil
}

// ParseJDText 使用AI解析职位描述文本
func (p *AITextParser) ParseJDText(text string) (*models.JobDescription, error) {
	if p.apiKey == "" {
		// 如果没有API密钥，仅返回原始文本
		return &models.JobDescription{
			RawText: text,
		}, nil
	}

	// 构建提示词
	prompt := buildJDParsePrompt(text)

	var content string
	var err error

	// 根据配置选择使用Grok3或OpenAI
	if p.useGrok {
		content, err = p.callGrok3API(prompt)
	} else {
		content, err = p.callOpenAIAPI(prompt)
	}

	if err != nil {
		return nil, fmt.Errorf("AI解析职位描述失败: %w", err)
	}

	// 解析AI返回的JSON
	jd, err := parseJDJSON(content, text)
	if err != nil {
		return nil, fmt.Errorf("解析AI返回的JSON失败: %w", err)
	}

	return jd, nil
}

// ParseJDFile 使用AI解析JD文件
func (p *AITextParser) ParseJDFile(filePath string) (*models.JobDescription, error) {
	if p.fileParser == nil {
		return nil, fmt.Errorf("文件解析器未初始化")
	}

	// 从文件中提取文本
	text, err := p.fileParser.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("文件解析失败: %w", err)
	}

	// 使用提取的文本解析JD
	jd, err := p.ParseJDText(text)
	if err != nil {
		return nil, err
	}

	// 设置文件路径
	jd.FilePath = filePath
	return jd, nil
}

// callGrok3API 调用Grok3 API
func (p *AITextParser) callGrok3API(prompt string) (string, error) {
	// 延迟初始化客户端
	if p.client == nil {
		p.client = ai.NewGrok3Client(p.apiKey)
	}

	client := p.client.(*ai.Grok3Client)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		ai.Grok3ChatRequest{
			Model: "grok-3-latest",
			Messages: []ai.Grok3Message{
				{
					Role:    "system",
					Content: "你是一个专业的简历分析助手，擅长从文本中提取结构化信息。请尽可能准确地提取所有相关信息，并按照要求的格式输出JSON。",
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
			MaxTokens:   1024,
			Temperature: 0.2,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("Grok3返回了空的回复")
	}

	return resp.Choices[0].Message.Content, nil
}

// callOpenAIAPI 调用OpenAI API
func (p *AITextParser) callOpenAIAPI(prompt string) (string, error) {
	// 延迟初始化客户端
	if p.client == nil {
		p.client = ai.NewOpenAIClient(p.apiKey)
	}

	client := p.client.(*ai.OpenAIClient)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		ai.OpenAIChatRequest{
			Model: "gpt-4o",
			Messages: []ai.OpenAIMessage{
				{
					Role:    "system",
					Content: "你是一个专业的简历分析助手，擅长从文本中提取结构化信息。请尽可能准确地提取所有相关信息，并按照要求的格式输出JSON。",
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
			MaxTokens:   1024,
			Temperature: 0.2,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("OpenAI返回了空的回复")
	}

	return resp.Choices[0].Message.Content, nil
}

// 构建简历解析的提示词
func buildResumeParsePrompt(text string) string {
	return fmt.Sprintf(`
请从以下简历文本中提取关键信息，并以JSON格式返回：

==== 简历文本 ====
%s

请提取并返回以下字段（如果信息不可用，请返回空字符串或空数组）：
1. 姓名
2. 电子邮箱
3. 电话号码
4. 教育经历（包括学校、学位、专业、时间段）
5. 工作经验（包括公司、职位、时间段、职责描述）
6. 技能列表

请按照以下JSON格式返回：
{
  "name": "姓名",
  "email": "电子邮箱",
  "phone": "电话号码",
  "education": ["教育经历1", "教育经历2", ...],
  "experience": ["工作经验1", "工作经验2", ...],
  "skills": ["技能1", "技能2", ...]
}

只返回JSON，不要包含额外的解释或修饰文字。
`, text)
}

// 构建JD解析的提示词
func buildJDParsePrompt(text string) string {
	return fmt.Sprintf(`
请从以下职位描述文本中提取关键信息，并以JSON格式返回：

==== 职位描述文本 ====
%s

请提取并返回以下字段（如果信息不可用，请返回空字符串或空数组）：
1. 职位标题
2. 公司名称
3. 职位描述概述
4. 职位要求列表

请按照以下JSON格式返回：
{
  "title": "职位标题",
  "company": "公司名称",
  "description": "职位描述概述",
  "requirements": ["要求1", "要求2", ...]
}

只返回JSON，不要包含额外的解释或修饰文字。
`, text)
}

// 解析AI返回的简历JSON
func parseResumeJSON(content string, originalText string) (*models.Resume, error) {
	// 提取并清理JSON内容
	jsonContent := extractJSONFromText(content)

	// 解析JSON
	var result struct {
		Name       string   `json:"name"`
		Email      string   `json:"email"`
		Phone      string   `json:"phone"`
		Education  []string `json:"education"`
		Experience []string `json:"experience"`
		Skills     []string `json:"skills"`
	}

	err := json.Unmarshal([]byte(jsonContent), &result)
	if err != nil {
		// 如果解析失败，尝试修复常见的JSON格式错误
		fixedJSON := tryFixJSONFormat(jsonContent)
		err = json.Unmarshal([]byte(fixedJSON), &result)

		if err != nil {
			// 如果仍然失败，返回一个包含原始文本的简单Resume对象
			fmt.Printf("JSON解析错误: %v\nJSON内容: %s\n", err, jsonContent)
			return &models.Resume{
				RawText: originalText,
			}, nil
		}
	}

	// 验证并清理数据
	resume := &models.Resume{
		Name:       sanitizeField(result.Name),
		Email:      sanitizeField(result.Email),
		Phone:      sanitizeField(result.Phone),
		Education:  sanitizeStringArray(result.Education),
		Experience: sanitizeStringArray(result.Experience),
		Skills:     sanitizeStringArray(result.Skills),
		RawText:    originalText,
	}

	return resume, nil
}

// 解析AI返回的JD JSON
func parseJDJSON(content string, originalText string) (*models.JobDescription, error) {
	// 提取并清理JSON内容
	jsonContent := extractJSONFromText(content)

	// 解析JSON
	var result struct {
		Title        string   `json:"title"`
		Company      string   `json:"company"`
		Description  string   `json:"description"`
		Requirements []string `json:"requirements"`
	}

	err := json.Unmarshal([]byte(jsonContent), &result)
	if err != nil {
		// 如果解析失败，尝试修复常见的JSON格式错误
		fixedJSON := tryFixJSONFormat(jsonContent)
		err = json.Unmarshal([]byte(fixedJSON), &result)

		if err != nil {
			// 如果仍然失败，返回一个包含原始文本的简单JD对象
			fmt.Printf("JSON解析错误: %v\nJSON内容: %s\n", err, jsonContent)
			return &models.JobDescription{
				RawText: originalText,
			}, nil
		}
	}

	// 验证并清理数据
	jd := &models.JobDescription{
		Title:        sanitizeField(result.Title),
		Company:      sanitizeField(result.Company),
		Description:  sanitizeField(result.Description),
		Requirements: sanitizeStringArray(result.Requirements),
		RawText:      originalText,
	}

	return jd, nil
}

// 查找JSON开始位置
func findJSONStartIndex(text string) int {
	// 首先尝试找到标准的JSON开始位置
	jsonStart := strings.Index(text, "{")

	// 如果找不到，可能内容被包裹在代码块中
	if jsonStart < 0 {
		// 查找Markdown代码块标记
		codeBlockStart := strings.Index(text, "```json")
		if codeBlockStart >= 0 {
			// 找到代码块开始后的第一个{
			blockContentStart := codeBlockStart + 7 // "```json" 长度为7
			jsonStart = strings.Index(text[blockContentStart:], "{")
			if jsonStart >= 0 {
				jsonStart += blockContentStart
			}
		}

		// 如果还是没有找到，尝试查找不带语言标记的代码块
		if jsonStart < 0 {
			codeBlockStart = strings.Index(text, "```")
			if codeBlockStart >= 0 {
				// 找到代码块开始后的第一个{
				blockContentStart := codeBlockStart + 3 // "```" 长度为3
				jsonStart = strings.Index(text[blockContentStart:], "{")
				if jsonStart >= 0 {
					jsonStart += blockContentStart
				}
			}
		}
	}

	return jsonStart
}

// 查找JSON结束位置
func findJSONEndIndex(text string) int {
	// 首先尝试简单查找最后一个}位置
	jsonEnd := strings.LastIndex(text, "}")

	// 如果找到了}，检查它后面是否有代码块结束标记
	if jsonEnd >= 0 {
		codeBlockEnd := strings.Index(text[jsonEnd+1:], "```")
		if codeBlockEnd >= 0 {
			// 不需要调整jsonEnd，因为我们只关心}的位置
		}
	}

	return jsonEnd
}

// 从文本中提取JSON内容
func extractJSONFromText(content string) string {
	// 尝试找到JSON部分
	jsonStart := 0
	jsonEnd := len(content)

	// 如果AI返回了额外的文字，尝试定位JSON部分
	if startIdx := findJSONStartIndex(content); startIdx >= 0 {
		jsonStart = startIdx
	}
	if endIdx := findJSONEndIndex(content[jsonStart:]); endIdx >= 0 {
		jsonEnd = jsonStart + endIdx + 1 // 加1是为了包含最后一个花括号
	} else {
		jsonEnd = len(content)
	}

	// 提取JSON部分
	jsonContent := content[jsonStart:jsonEnd]

	// 清理JSON内容（移除可能的Markdown代码块标记等）
	jsonContent = cleanJSONContent(jsonContent)

	return jsonContent
}

// 清理JSON内容
func cleanJSONContent(content string) string {
	// 移除开头的```json或```标记
	content = strings.TrimPrefix(strings.TrimPrefix(content, "```json"), "```")

	// 移除结尾的```标记
	content = strings.TrimSuffix(content, "```")

	// 移除可能的前后空白
	content = strings.TrimSpace(content)

	return content
}
