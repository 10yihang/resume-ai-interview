/*
 * @author: yihang_01
 * @Date: 2025-05-21 17:20:03
 * @LastEditTime: 2025-05-21 17:40:38
 * QwQ 加油加油
 */
package interview

import (
	"fmt"
	"os"
	"testing"

	"github.com/10yihang/resume-ai-interview/config"
	"github.com/10yihang/resume-ai-interview/models"
)

func TestAnswerEvaluation(t *testing.T) {
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
	question, answer, jd := createTestData()

	// 测试Grok答案评估器
	grokEvaluator := NewGrok3AnswerEvaluator(apiKey)
	grokEvaluation, err := grokEvaluator.EvaluateAnswer(question, answer, jd)
	if err != nil {
		t.Logf("Grok评估失败: %v", err)
	} else {
		fmt.Printf("Grok评估结果:\n")
		fmt.Printf("分数: %d\n", grokEvaluation.Score)
		fmt.Printf("反馈: %s\n", grokEvaluation.Feedback)
		fmt.Printf("建议: %s\n\n", grokEvaluation.Suggestions)
	}

	// 测试模拟答案评估器
	mockEvaluator := NewMockAnswerEvaluator()
	mockEvaluation, err := mockEvaluator.EvaluateAnswer(question, answer, jd)
	if err != nil {
		t.Errorf("模拟评估失败: %v", err)
	} else {
		fmt.Printf("模拟评估结果:\n")
		fmt.Printf("分数: %d\n", mockEvaluation.Score)
		fmt.Printf("反馈: %s\n", mockEvaluation.Feedback)
		fmt.Printf("建议: %s\n", mockEvaluation.Suggestions)
	}
}

// 创建测试数据
func createTestData() (models.Question, models.Answer, *models.JobDescription) {
	question := models.Question{
		ID:       1,
		Content:  "请介绍一下你在Go语言方面的经验，特别是在微服务架构中的应用?",
		Category: "专业技能",
	}

	answer := models.Answer{
		QuestionID: 1,
		Content: `我有5年的Go语言开发经验，主要应用在微服务架构中。
在上一家公司，我负责将一个单体应用拆分成微服务架构。我们使用Go语言重写了核心服务，包括用户认证、订单处理和支付系统。
我们采用了gRPC作为服务间通信协议，使用Kubernetes进行容器编排，同时实现了服务发现、负载均衡和熔断等功能。
这个重构项目使系统性能提升了约40%，并且大大提高了开发团队的开发效率。我们能够更快地交付新功能，并且能更容易地进行单元测试和集成测试。`,
	}

	jd := &models.JobDescription{
		Title:        "高级后端工程师",
		Company:      "未来科技有限公司",
		Description:  "我们正在寻找一位经验丰富的高级后端工程师加入我们的团队，帮助构建和扩展我们的微服务架构。",
		Requirements: []string{"5年以上Go语言开发经验", "熟悉微服务架构和相关技术", "具有大规模分布式系统开发经验"},
		RawText:      "这是一份职位描述的原始文本",
	}

	return question, answer, jd
}
