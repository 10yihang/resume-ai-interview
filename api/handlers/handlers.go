package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/10yihang/resume-ai-interview/config"
	"github.com/10yihang/resume-ai-interview/internal/ai"
	"github.com/10yihang/resume-ai-interview/internal/interview"
	"github.com/10yihang/resume-ai-interview/internal/ocr"
	"github.com/10yihang/resume-ai-interview/internal/parser"
	"github.com/10yihang/resume-ai-interview/models"
	"github.com/gin-gonic/gin"
)

// 存储上传的文件和处理过的数据
var (
	resumeStore = make(map[string]*models.Resume)
	jdStore     = make(map[string]*models.JobDescription)
	questions   = make(map[string]*models.QuestionSet)
	// 加载配置
	cfg *config.Config
)

// InitHandlers 初始化处理器
func InitHandlers(_config *config.Config) {
	cfg = _config
	// 如果配置为空，创建默认配置
	if cfg == nil {
		cfg = config.NewConfig()
	}
}

// IndexHandler 处理首页请求
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "AI简历面试助手",
	})
}

// UploadResumeHandler 处理简历上传
func UploadResumeHandler(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取上传文件: " + err.Error()})
		return
	}
	defer file.Close()

	// 创建上传目录
	uploadDir := "./uploads/resumes"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败: " + err.Error()})
		return
	}

	// 保存文件
	filename := filepath.Join(uploadDir, header.Filename)
	out, err := os.Create(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "复制文件失败: " + err.Error()})
		return
	}

	// 初始化OCR处理器
	var ocrProcessor ocr.OCRProcessor
	if cfg.UseOCR {
		ocrProcessor = ocr.GetOCRProcessor(cfg.OCRAPIKey, cfg.TesseractPath)
	}

	// 创建文件解析器
	fileParser := parser.NewResumeFileParser(ocrProcessor, cfg.UseOCR)
	// 使用AI解析简历文件
	aiParser := parser.NewAITextParser(cfg.APIKey, cfg.UseGrok, fileParser)
	resume, err := aiParser.ParseResumeFile(filename)
	if err != nil {
		// 如果OCR失败，尝试使用传统方法解析
		if cfg.UseOCR && err.Error() == "文件解析失败: OCR处理失败" {
			// 创建不使用OCR的文件解析器
			fileParser := parser.NewResumeFileParser(nil, false)
			aiParser := parser.NewAITextParser(cfg.APIKey, cfg.UseGrok, fileParser)
			resume, err = aiParser.ParseResumeFile(filename)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "简历解析失败: " + err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "简历解析失败: " + err.Error()})
			return
		}
	}

	// 保存解析后的简历
	resumeID := header.Filename
	resumeStore[resumeID] = resume

	c.JSON(http.StatusOK, gin.H{
		"message":  "简历上传成功",
		"resumeId": resumeID,
		"resume":   resume,
	})
}

// UploadJDHandler 处理JD上传
func UploadJDHandler(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("jd")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取上传文件: " + err.Error()})
		return
	}
	defer file.Close()

	// 创建上传目录
	uploadDir := "./uploads/jds"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败: " + err.Error()})
		return
	}

	// 保存文件
	filename := filepath.Join(uploadDir, header.Filename)
	out, err := os.Create(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "复制文件失败: " + err.Error()})
		return
	}

	// 初始化OCR处理器
	var ocrProcessor ocr.OCRProcessor
	if cfg.UseOCR {
		ocrProcessor = ocr.GetOCRProcessor(cfg.OCRAPIKey, cfg.TesseractPath)
	}

	// 创建文件解析器
	fileParser := parser.NewResumeFileParser(ocrProcessor, cfg.UseOCR)
	// 使用AI解析JD文件
	aiParser := parser.NewAITextParser(cfg.APIKey, cfg.UseGrok, fileParser)
	jd, err := aiParser.ParseJDFile(filename)
	if err != nil {
		// 如果OCR失败，尝试使用传统方法解析
		if cfg.UseOCR && err.Error() == "文件解析失败: OCR处理失败" {
			// 创建不使用OCR的文件解析器
			fileParser := parser.NewResumeFileParser(nil, false)
			aiParser := parser.NewAITextParser(cfg.APIKey, cfg.UseGrok, fileParser)
			jd, err = aiParser.ParseJDFile(filename)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "JD解析失败: " + err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JD解析失败: " + err.Error()})
			return
		}
	}

	// 保存解析后的JD
	jdID := header.Filename
	jdStore[jdID] = jd

	c.JSON(http.StatusOK, gin.H{
		"message": "JD上传成功",
		"jdId":    jdID,
		"jd":      jd,
	})
}

// GenerateQuestionsHandler 生成面试问题
func GenerateQuestionsHandler(c *gin.Context) {
	var request struct {
		ResumeID string `json:"resumeId" binding:"required"`
		JDID     string `json:"jdId" binding:"required"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	// 获取简历和JD
	resume, ok := resumeStore[request.ResumeID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "简历不存在"})
		return
	}
	jd, ok := jdStore[request.JDID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "JD不存在"})
		return
	}

	// 生成问题
	generator := ai.GetQuestionGenerator(cfg.APIKey, cfg.UseGrok)
	questionSet, err := generator.GenerateQuestions(resume, jd)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成问题失败: " + err.Error()})
		return
	}

	// 保存生成的问题
	questionID := request.ResumeID + "_" + request.JDID
	questions[questionID] = questionSet

	c.JSON(http.StatusOK, gin.H{
		"message":   "问题生成成功",
		"questions": questionSet,
	})
}

// EvaluateAnswerHandler 评估面试回答
func EvaluateAnswerHandler(c *gin.Context) {
	var request struct {
		QuestionSetID string        `json:"questionSetId" binding:"required"`
		QuestionID    int           `json:"questionId" binding:"required"`
		Answer        models.Answer `json:"answer" binding:"required"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	// 获取问题集
	questionSet, ok := questions[request.QuestionSetID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "问题集不存在"})
		return
	}

	// 获取问题
	var question models.Question
	found := false
	for _, q := range questionSet.Questions {
		if q.ID == request.QuestionID {
			question = q
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "问题不存在"})
		return
	}

	// 获取JD
	jd, ok := jdStore[questionSet.JDID]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JD数据不存在"})
		return
	}
	// 评估回答
	// 评估回答
	evaluator := interview.GetAnswerEvaluator(cfg.APIKey, cfg.UseGrok)
	evaluation, err := evaluator.EvaluateAnswer(question, request.Answer, jd)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评估回答失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "回答评估成功",
		"evaluation": evaluation,
	})
}
