package ai

import (
	"encoding/json"
	"fmt"

	"github.com/10yihang/resume-ai-interview/models"
)

// MockQuestionGenerator 模拟问题生成器，用于测试
type MockQuestionGenerator struct{}

// NewMockQuestionGenerator 创建模拟问题生成器
func NewMockQuestionGenerator() *MockQuestionGenerator {
	return &MockQuestionGenerator{}
}

// GenerateQuestions 生成模拟面试问题
func (g *MockQuestionGenerator) GenerateQuestions(resume *models.Resume, jd *models.JobDescription) (*models.QuestionSet, error) {
	// 创建一些模拟问题
	questions := []models.Question{
		{ID: 1, Content: "请介绍一下你的技术背景和专长？", Category: "专业技能"},
		{ID: 2, Content: "你对这个职位的理解是什么？", Category: "职业规划"},
		{ID: 3, Content: "请描述一个你曾经解决的技术难题及解决方案？", Category: "专业技能"},
		{ID: 4, Content: "你是如何处理团队合作中的冲突的？", Category: "团队协作"},
		{ID: 5, Content: "你最近学习了哪些新技术？为什么选择学习它们？", Category: "专业技能"},
		{ID: 6, Content: "请分享一个你在工作中犯过的错误，以及从中学到了什么？", Category: "工作经验"},
		{ID: 7, Content: "你期望的职业发展路径是什么？", Category: "职业规划"},
		{ID: 8, Content: "你如何确保你的代码质量？", Category: "专业技能"},
		{ID: 9, Content: "你如何应对工作中的压力和截止日期？", Category: "工作经验"},
		{ID: 10, Content: "你如何保持自己的技术更新？", Category: "专业技能"},
	}

	return &models.QuestionSet{
		ResumeID:  "resume_id",
		JDID:      "jd_id",
		Questions: questions,
	}, nil
}

// 解析模拟问题响应
func (g *MockQuestionGenerator) parseMockQuestions(resume *models.Resume, jd *models.JobDescription, content string) *models.QuestionSet {
	// 解析JSON内容
	var result struct {
		Questions []models.Question `json:"questions"`
	}

	err := json.Unmarshal([]byte(content), &result)
	if err != nil {
		fmt.Printf("解析JSON失败：%v\n", err)
		// 返回一些默认问题
		questions := []models.Question{
			{ID: 1, Content: "请介绍一下你的技术背景和专长？", Category: "专业技能"},
			{ID: 2, Content: "你对这个职位的理解是什么？", Category: "职业规划"},
			// 更多默认问题...
		}
		return &models.QuestionSet{
			ResumeID:  "resume_id",
			JDID:      "jd_id",
			Questions: questions,
		}
	}

	return &models.QuestionSet{
		ResumeID:  "resume_id",
		JDID:      "jd_id",
		Questions: result.Questions,
	}
}
