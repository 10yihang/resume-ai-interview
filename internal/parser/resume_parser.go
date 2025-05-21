package parser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/10yihang/resume-ai-interview/models"
	"github.com/ledongthuc/pdf"
)

// ResumeParser 提供简历解析功能
type ResumeParser struct{}

// NewResumeParser 创建新的简历解析器
func NewResumeParser() *ResumeParser {
	return &ResumeParser{}
}

// ParseFromFile 从文件解析简历
func (p *ResumeParser) ParseFromFile(filePath string) (*models.Resume, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	var text string
	var err error

	switch ext {
	case ".pdf":
		text, err = extractTextFromPDF(filePath)
	case ".txt":
		text, err = extractTextFromTXT(filePath)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return nil, err
	}

	// 简单解析简历文本
	resume := &models.Resume{
		RawText:  text,
		FilePath: filePath,
	}

	// TODO: 使用更复杂的解析逻辑提取姓名、邮箱等信息
	resume.Name = extractName(text)
	resume.Email = extractEmail(text)
	resume.Phone = extractPhone(text)
	resume.Education = extractEducation(text)
	resume.Experience = extractExperience(text)
	resume.Skills = extractSkills(text)

	return resume, nil
}

// 从PDF文件中提取文本
func extractTextFromPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	buf.ReadFrom(b)
	return buf.String(), nil
}

// 从TXT文件中提取文本
func extractTextFromTXT(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// 从文本中提取姓名（简单实现）
func extractName(text string) string {
	// TODO: 使用更复杂的算法提取姓名
	return ""
}

// 从文本中提取邮箱
func extractEmail(text string) string {
	// TODO: 使用正则表达式提取邮箱
	return ""
}

// 从文本中提取电话
func extractPhone(text string) string {
	// TODO: 使用正则表达式提取电话
	return ""
}

// 从文本中提取教育经历
func extractEducation(text string) []string {
	// TODO: 提取教育经历
	return []string{}
}

// 从文本中提取工作经验
func extractExperience(text string) []string {
	// TODO: 提取工作经验
	return []string{}
}

// 从文本中提取技能
func extractSkills(text string) []string {
	// TODO: 提取技能
	return []string{}
}
