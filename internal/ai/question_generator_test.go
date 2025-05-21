package ai

import (
	"fmt"
	"os"
	"testing"

	"github.com/10yihang/resume-ai-interview/config"
	"github.com/10yihang/resume-ai-interview/models"
)

func TestQuestionGeneration(t *testing.T) {
	// 从环境变量或配置文件加载API密钥
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		// 如果没有设置环境变量，可以从配置加载
		cfg := config.NewConfig()
		if cfg.APIKey != "" {
			apiKey = cfg.APIKey
		} else {
			t.Skip("没有设置API密钥，跳过测试")
		}
	}

	// 创建测试数据
	resume, jd := createTestResumeAndJD()

	// 测试Grok问题生成器
	grokGenerator := NewGrok3QuestionGenerator(apiKey)
	grokQuestions, err := grokGenerator.GenerateQuestions(resume, jd)
	if err != nil {
		t.Logf("Grok问题生成失败: %v", err)
	} else {
		fmt.Printf("Grok生成的问题:\n")
		for i, q := range grokQuestions.Questions {
			fmt.Printf("%d. [%s] %s\n", i+1, q.Category, q.Content)
		}
		fmt.Println()
	}

	// 测试模拟问题生成器
	mockGenerator := NewMockQuestionGenerator()
	mockQuestions, err := mockGenerator.GenerateQuestions(resume, jd)
	if err != nil {
		t.Errorf("模拟问题生成失败: %v", err)
	} else {
		fmt.Printf("模拟生成的问题:\n")
		for i, q := range mockQuestions.Questions {
			fmt.Printf("%d. [%s] %s\n", i+1, q.Category, q.Content)
		}
	}
}

// 创建一个简历和JD供测试使用
func createTestResumeAndJD() (*models.Resume, *models.JobDescription) {
	resume := &models.Resume{
		Name:       "张三",
		Email:      "zhangsan@example.com",
		Phone:      "13800138000",
		Education:  []string{"北京大学，计算机科学与技术，学士学位，2015-2019", "清华大学，软件工程，硕士学位，2019-2022"},
		Experience: []string{"云智科技，高级开发工程师，2022至今", "ABC公司，初级开发工程师，2019-2022"},
		Skills:     []string{"Go", "Python", "JavaScript", "React", "Docker", "Kubernetes", "MySQL", "MongoDB"},
		RawText:    "这是一份简历的原始文本",
		FilePath:   "测试简历.pdf",
	}

	jd := &models.JobDescription{
		Title:        "高级后端工程师",
		Company:      "未来科技有限公司",
		Description:  "我们正在寻找一位经验丰富的高级后端工程师加入我们的团队，帮助构建和扩展我们的微服务架构。",
		Requirements: []string{"5年以上Go语言开发经验", "熟悉微服务架构和相关技术", "具有大规模分布式系统开发经验"},
		RawText:      "这是一份职位描述的原始文本",
		FilePath:     "测试JD.txt",
	}

	return resume, jd
}
