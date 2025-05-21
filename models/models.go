package models

// Resume 表示解析后的简历
type Resume struct {
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Education  []string `json:"education"`
	Experience []string `json:"experience"`
	Skills     []string `json:"skills"`
	RawText    string   `json:"rawText"`
	FilePath   string   `json:"filePath"`
}

// JobDescription 表示岗位JD
type JobDescription struct {
	Title        string   `json:"title"`
	Company      string   `json:"company"`
	Description  string   `json:"description"`
	Requirements []string `json:"requirements"`
	RawText      string   `json:"rawText"`
	FilePath     string   `json:"filePath"`
}

// Question 表示面试问题
type Question struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

// QuestionSet 表示一组面试问题
type QuestionSet struct {
	ResumeID  string     `json:"resumeId"`
	JDID      string     `json:"jdId"`
	Questions []Question `json:"questions"`
}

// Answer 表示面试回答
type Answer struct {
	QuestionID int    `json:"questionId"`
	Content    string `json:"content"`
}

// Evaluation 表示面试评估
type Evaluation struct {
	AnswerID    int    `json:"answerId"`
	Score       int    `json:"score"`
	Feedback    string `json:"feedback"`
	Suggestions string `json:"suggestions"`
}
