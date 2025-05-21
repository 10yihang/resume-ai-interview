package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/10yihang/resume-ai-interview/config"
	"github.com/10yihang/resume-ai-interview/internal/ocr"
	"github.com/10yihang/resume-ai-interview/models"
)

func TestAIParser(t *testing.T) {
	// 从环境变量或配置文件加载API密钥
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		// // 如果没有设置环境变量，可以从配置加载
		cfg := config.NewConfig()
		if cfg.APIKey != "" {
			apiKey = cfg.APIKey
		} else {
			t.Skip("没有设置API密钥，跳过测试")
		}
	}
	// 创建AI解析器，优先使用Grok，暂时传入nil作为fileParser
	parser := NewAITextParser(apiKey, true, nil)

	// 测试简历文本
	resumeText := `
姓名: 张三
邮箱: zhangsan@example.com
电话: 13800138000

教育经历:
北京大学，计算机科学与技术，学士学位，2015-2019
清华大学，软件工程，硕士学位，2019-2022

工作经验:
云智科技，高级开发工程师，2022至今
- 负责企业级SaaS应用的后端开发
- 使用Go语言和微服务架构构建高性能API
- 优化数据库查询，提高系统响应速度50%

ABC公司，初级开发工程师，2019-2022
- 开发和维护Web应用
- 参与前端React组件开发
- 协助设计RESTful API

技能:
Go, Python, JavaScript, React, Docker, Kubernetes, MySQL, MongoDB
`

	// 测试JD文本
	jdText := `
职位: 高级后端工程师
公司: 未来科技有限公司
职位描述:
我们正在寻找一位经验丰富的高级后端工程师加入我们的团队，帮助构建和扩展我们的微服务架构。

职位要求:
- 5年以上Go语言开发经验
- 熟悉微服务架构和相关技术
- 具有大规模分布式系统开发经验
- 精通SQL和NoSQL数据库优化
- 熟悉容器技术和Kubernetes
- 良好的沟通能力和团队协作精神
`

	// 解析简历
	resume, err := parser.ParseResumeText(resumeText)
	if err != nil {
		t.Fatalf("解析简历失败: %v", err)
	}

	// 验证简历解析结果
	fmt.Printf("简历解析结果:\n")
	fmt.Printf("姓名: %s\n", resume.Name)
	fmt.Printf("邮箱: %s\n", resume.Email)
	fmt.Printf("电话: %s\n", resume.Phone)
	fmt.Printf("教育经历: %v\n", resume.Education)
	fmt.Printf("工作经验: %v\n", resume.Experience)
	fmt.Printf("技能: %v\n\n", resume.Skills)

	// 解析JD
	jd, err := parser.ParseJDText(jdText)
	if err != nil {
		t.Fatalf("解析JD失败: %v", err)
	}

	// 验证JD解析结果
	fmt.Printf("JD解析结果:\n")
	fmt.Printf("职位: %s\n", jd.Title)
	fmt.Printf("公司: %s\n", jd.Company)
	fmt.Printf("描述: %s\n", jd.Description)
	fmt.Printf("要求: %v\n", jd.Requirements)

	// 如果简历和JD都成功解析，测试通过
	if resume.Name != "" && jd.Title != "" {
		t.Log("AI解析测试通过")
	} else {
		t.Error("AI解析测试失败")
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
	}

	jd := &models.JobDescription{
		Title:        "高级后端工程师",
		Company:      "未来科技有限公司",
		Description:  "我们正在寻找一位经验丰富的高级后端工程师加入我们的团队，帮助构建和扩展我们的微服务架构。",
		Requirements: []string{"5年以上Go语言开发经验", "熟悉微服务架构和相关技术", "具有大规模分布式系统开发经验"},
		RawText:      "这是一份职位描述的原始文本",
	}

	return resume, jd
}

func TestFileParserWithOCR(t *testing.T) {
	// 获取OCR API密钥
	ocrAPIKey := os.Getenv("OCR_SPACE_API_KEY")

	// 创建OCR处理器
	var ocrProcessor ocr.OCRProcessor
	if ocrAPIKey != "" {
		ocrProcessor = ocr.GetOCRProcessor(ocrAPIKey, "")
	} else {
		// 使用Tesseract OCR
		ocrProcessor = ocr.GetOCRProcessor("", "tesseract")
	}

	// 创建文件解析器
	fileParser := NewResumeFileParser(ocrProcessor, true)

	// 测试文本文件解析（不依赖于OCR）
	t.Run("TestTextFileParser", func(t *testing.T) {
		// 创建测试文件
		tempDir, err := os.MkdirTemp("", "parser-test")
		if err != nil {
			t.Fatalf("创建临时目录失败: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// 创建测试文本
		testText := "姓名: 张三\n电话: 13800138000\n教育经历: 北京大学\n技能: Go, Python"
		textFile := filepath.Join(tempDir, "resume.txt")
		err = os.WriteFile(textFile, []byte(testText), 0644)
		if err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}

		// 解析文件
		text, err := fileParser.ParseFile(textFile)
		if err != nil {
			t.Fatalf("解析文本文件失败: %v", err)
		}

		// 验证结果
		if text != testText {
			t.Errorf("解析结果不匹配\n期望: %s\n实际: %s", testText, text)
		} else {
			t.Logf("成功解析文本文件")
		}
	})

	// 测试AI解析器与文件解析器的集成
	t.Run("TestAIParserWithFileParser", func(t *testing.T) {
		// 从环境变量获取API密钥
		apiKey := os.Getenv("GROK_API_KEY")
		if apiKey == "" {
			t.Skip("未设置GROK_API_KEY，跳过测试")
		}

		// 创建AI解析器
		aiParser := NewAITextParser(apiKey, true, fileParser)

		// 创建测试文件
		tempDir, err := os.MkdirTemp("", "parser-test")
		if err != nil {
			t.Fatalf("创建临时目录失败: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// 创建测试文本
		testText := `姓名: 李四
电子邮件: lisi@example.com
电话: 13900139000
教育经历: 
  - 清华大学, 计算机科学, 硕士, 2018-2021
  - 北京大学, 软件工程, 学士, 2014-2018
工作经验:
  - ABC科技有限公司, 高级开发工程师, 2021至今
  - XYZ公司, 软件开发实习生, 2019-2020
技能: Java, Python, Docker, Kubernetes, MySQL, Redis`

		textFile := filepath.Join(tempDir, "resume.txt")
		err = os.WriteFile(textFile, []byte(testText), 0644)
		if err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}

		// 解析文件
		resume, err := aiParser.ParseResumeFile(textFile)
		if err != nil {
			t.Fatalf("AI解析简历文件失败: %v", err)
		}

		// 验证结果
		if resume == nil {
			t.Fatal("解析结果为nil")
		}

		t.Logf("AI解析简历结果:\n")
		t.Logf("姓名: %s", resume.Name)
		t.Logf("邮箱: %s", resume.Email)
		t.Logf("电话: %s", resume.Phone)
		t.Logf("教育经历: %v", resume.Education)
		t.Logf("工作经验: %v", resume.Experience)
		t.Logf("技能: %v", resume.Skills)

		// 简单验证
		if resume.Name == "" || resume.Email == "" || resume.Phone == "" {
			t.Error("基本信息解析失败")
		}

		if len(resume.Education) == 0 || len(resume.Experience) == 0 || len(resume.Skills) == 0 {
			t.Error("数组字段解析失败")
		}
	})
}
