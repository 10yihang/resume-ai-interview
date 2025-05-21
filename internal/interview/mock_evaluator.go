package interview

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/10yihang/resume-ai-interview/models"
)

// MockAnswerEvaluator 模拟答案评估器，用于测试
type MockAnswerEvaluator struct{}

// NewMockAnswerEvaluator 创建模拟答案评估器
func NewMockAnswerEvaluator() *MockAnswerEvaluator {
	return &MockAnswerEvaluator{}
}

// EvaluateAnswer 评估面试回答
func (e *MockAnswerEvaluator) EvaluateAnswer(question models.Question, answer models.Answer, jd *models.JobDescription) (*models.Evaluation, error) {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 根据回答长度和内容生成模拟分数
	score := 5 // 默认中等分数

	// 答案长度影响分数
	if len(answer.Content) > 200 {
		score += 2 // 答案较长，加分
	} else if len(answer.Content) < 50 {
		score -= 2 // 答案太短，减分
	}

	// 根据关键词影响分数
	keywords := getKeywords(question.Category)
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(answer.Content), strings.ToLower(keyword)) {
			score += 1 // 包含关键词，加分
			if score > 10 {
				score = 10 // 最高分10分
			}
		}
	}

	// 随机调整分数，增加一些变化
	score += rand.Intn(3) - 1 // -1到+1的随机调整

	// 确保分数在1-10范围内
	if score < 1 {
		score = 1
	} else if score > 10 {
		score = 10
	}

	// 生成反馈和建议
	feedback := generateFeedback(score, answer.Content)
	suggestions := generateSuggestions(score, question.Category)

	return &models.Evaluation{
		AnswerID:    answer.QuestionID,
		Score:       score,
		Feedback:    feedback,
		Suggestions: suggestions,
	}, nil
}

// 根据问题类别获取关键词
func getKeywords(category string) []string {
	switch category {
	case "专业技能":
		return []string{"经验", "技术", "项目", "解决方案", "工具", "框架", "语言", "开发", "实现"}
	case "职业规划":
		return []string{"目标", "规划", "发展", "成长", "学习", "提升", "career", "职业", "未来"}
	case "工作经验":
		return []string{"负责", "参与", "项目", "团队", "完成", "实施", "贡献", "经验", "挑战"}
	case "团队协作":
		return []string{"合作", "沟通", "协调", "团队", "配合", "分享", "协作", "团队精神", "共同"}
	default:
		return []string{"经验", "技能", "项目", "思考", "学习", "工作", "解决", "理解"}
	}
}

// 根据分数生成反馈
func generateFeedback(score int, answerContent string) string {
	if score >= 9 {
		return "回答非常出色！展示了深入的专业知识和丰富的实践经验。表达清晰、有条理，并且有具体的例子支持观点。"
	} else if score >= 7 {
		return "回答良好，展示了扎实的知识基础和一定的实践经验。思路清晰，但可以提供更多具体的例子来支持观点。"
	} else if score >= 5 {
		return "回答基本达到要求，表达了对问题的理解，但缺乏深度和细节。可以更有条理地组织回答，并提供更多实际例子。"
	} else if score >= 3 {
		return "回答较为简单，缺乏必要的细节和深度。建议更全面地思考问题，并结合自身经验提供具体例子。"
	} else {
		return "回答不够充分，对问题的理解有限。需要更深入地学习相关知识，并思考如何将其应用到实际工作中。"
	}
}

// 根据分数和问题类别生成建议
func generateSuggestions(score int, category string) string {
	if score >= 8 {
		switch category {
		case "专业技能":
			return "可以尝试分享更多关于如何解决复杂技术挑战的细节，以及你在项目中的独特贡献和创新点。"
		case "职业规划":
			return "建议在回答中加入一些具体的短期和长期职业目标，以及你计划如何达到这些目标的行动步骤。"
		case "工作经验":
			return "可以更详细地描述你在项目中解决的挑战，以及从中获得的经验教训，展示你的成长轨迹。"
		case "团队协作":
			return "可以分享更多关于你在团队中如何处理分歧和冲突的实际例子，展示你的沟通和协调能力。"
		default:
			return "整体回答很好，可以尝试提供更多具体的例子和数据来支持你的观点，让回答更有说服力。"
		}
	} else if score >= 5 {
		switch category {
		case "专业技能":
			return "建议深入描述你如何应用这些技能解决实际问题，提供具体的项目案例和技术细节。"
		case "职业规划":
			return "考虑将你的职业目标与公司的发展方向和岗位需求更紧密地结合起来，展示你对公司的了解和价值。"
		case "工作经验":
			return "在描述你的工作经历时，尝试使用STAR方法（情境、任务、行动、结果）来结构化你的回答，使其更有条理。"
		case "团队协作":
			return "可以分享一些具体的例子，说明你如何有效地与不同类型的团队成员合作，以及你在团队中扮演的角色。"
		default:
			return "建议使用更多具体的例子来支持你的观点，确保回答更有针对性和relevance。"
		}
	} else {
		switch category {
		case "专业技能":
			return "建议系统学习相关技术知识，并通过实践项目来巩固技能。准备具体例子来展示你如何应用这些技能。"
		case "职业规划":
			return "花些时间思考你的职业目标和发展路径，了解行业趋势和公司需求，确保你的规划既有抱负又现实可行。"
		case "工作经验":
			return "即使经验有限，也可以分享学校项目、实习或志愿者经历。关注你如何克服挑战，以及你学到了什么。"
		case "团队协作":
			return "反思你的团队合作经历，思考你的优势和不足。准备具体例子来说明你如何有效地与他人合作。"
		default:
			return "建议在回答前仔细思考问题的核心，确保你的回答直接回应了问题。使用具体的例子来支持你的观点。"
		}
	}
}

// 解析模拟评估结果
func (e *MockAnswerEvaluator) parseMockEvaluation(answer models.Answer, content string) *models.Evaluation {
	var result struct {
		Score       int    `json:"score"`
		Feedback    string `json:"feedback"`
		Suggestions string `json:"suggestions"`
	}

	err := json.Unmarshal([]byte(content), &result)
	if err != nil {
		fmt.Printf("解析JSON失败：%v\n", err)
		// 返回一个默认评估
		return &models.Evaluation{
			AnswerID:    answer.QuestionID,
			Score:       6,
			Feedback:    "回答基本符合要求，但可以提供更多具体的例子和细节。",
			Suggestions: "考虑使用STAR方法（情境、任务、行动、结果）来结构化你的回答，使其更有条理。",
		}
	}

	return &models.Evaluation{
		AnswerID:    answer.QuestionID,
		Score:       result.Score,
		Feedback:    result.Feedback,
		Suggestions: result.Suggestions,
	}
}
