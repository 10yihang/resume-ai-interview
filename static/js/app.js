// 全局变量
let resumeId = null;
let jdId = null;
let questionSetId = null;
let currentQuestions = [];
let currentQuestionId = null;

// DOM加载完成后执行
document.addEventListener('DOMContentLoaded', () => {
    // 上传简历表单处理
    const resumeForm = document.getElementById('resumeForm');
    resumeForm.addEventListener('submit', handleResumeUpload);

    // 上传JD表单处理
    const jdForm = document.getElementById('jdForm');
    jdForm.addEventListener('submit', handleJDUpload);

    // 生成问题按钮
    const generateBtn = document.getElementById('generateBtn');
    generateBtn.addEventListener('click', handleGenerateQuestions);

    // 提交回答按钮
    const submitAnswerBtn = document.getElementById('submitAnswerBtn');
    submitAnswerBtn.addEventListener('click', handleSubmitAnswer);

    // 检查按钮状态
    checkGenerateButtonStatus();
});

// 处理简历上传
async function handleResumeUpload(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const resumeInfoEl = document.getElementById('resumeInfo');
    const uploadBtn = document.getElementById('uploadResumeBtn');
    
    // 更改按钮状态
    uploadBtn.disabled = true;
    uploadBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> 上传中...';

    try {
        const response = await fetch('/upload/resume', {
            method: 'POST',
            body: formData
        });

        const data = await response.json();

        if (response.ok) {
            resumeId = data.resumeId;
            resumeInfoEl.classList.remove('d-none', 'alert-danger');
            resumeInfoEl.classList.add('alert-success');
            resumeInfoEl.textContent = `简历上传成功: ${data.resume.name || '未识别姓名'}`;
            checkGenerateButtonStatus();
        } else {
            throw new Error(data.error || '上传失败');
        }
    } catch (error) {
        resumeInfoEl.classList.remove('d-none', 'alert-success');
        resumeInfoEl.classList.add('alert-danger');
        resumeInfoEl.textContent = `错误: ${error.message}`;
    } finally {
        uploadBtn.disabled = false;
        uploadBtn.textContent = '上传简历';
    }
}

// 处理JD上传
async function handleJDUpload(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const jdInfoEl = document.getElementById('jdInfo');
    const uploadBtn = document.getElementById('uploadJDBtn');
    
    // 更改按钮状态
    uploadBtn.disabled = true;
    uploadBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> 上传中...';

    try {
        const response = await fetch('/upload/jd', {
            method: 'POST',
            body: formData
        });

        const data = await response.json();

        if (response.ok) {
            jdId = data.jdId;
            jdInfoEl.classList.remove('d-none', 'alert-danger');
            jdInfoEl.classList.add('alert-success');
            jdInfoEl.textContent = `JD上传成功: ${data.jd.title || '未识别职位名称'}`;
            checkGenerateButtonStatus();
        } else {
            throw new Error(data.error || '上传失败');
        }
    } catch (error) {
        jdInfoEl.classList.remove('d-none', 'alert-success');
        jdInfoEl.classList.add('alert-danger');
        jdInfoEl.textContent = `错误: ${error.message}`;
    } finally {
        uploadBtn.disabled = false;
        uploadBtn.textContent = '上传职位描述';
    }
}

// 检查生成按钮状态
function checkGenerateButtonStatus() {
    const generateBtn = document.getElementById('generateBtn');
    generateBtn.disabled = !(resumeId && jdId);
}

// 处理生成问题
async function handleGenerateQuestions() {
    if (!resumeId || !jdId) {
        alert('请先上传简历和职位描述');
        return;
    }

    const generateBtn = document.getElementById('generateBtn');
    generateBtn.disabled = true;
    generateBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> 生成中...';

    try {
        const response = await fetch('/generate/questions', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                resumeId: resumeId,
                jdId: jdId
            })
        });

        const data = await response.json();

        if (response.ok) {
            // 显示问题列表
            questionSetId = `${resumeId}_${jdId}`;
            currentQuestions = data.questions.questions;
            displayQuestions(currentQuestions);
        } else {
            throw new Error(data.error || '生成问题失败');
        }
    } catch (error) {
        alert(`错误: ${error.message}`);
    } finally {
        generateBtn.disabled = false;
        generateBtn.textContent = '生成面试问题';
    }
}

// 显示问题列表
function displayQuestions(questions) {
    const questionsContainer = document.getElementById('questionsContainer');
    const questionsList = document.getElementById('questionsList');
    
    // 清空现有问题
    questionsList.innerHTML = '';
    
    // 添加问题到列表
    questions.forEach(question => {
        const item = document.createElement('a');
        item.href = '#';
        item.className = 'list-group-item list-group-item-action question-item';
        item.dataset.id = question.id;
        item.innerHTML = `
            <div class="d-flex w-100 justify-content-between">
                <h6 class="mb-1">${question.content}</h6>
                <small>${question.category}</small>
            </div>
        `;
        item.addEventListener('click', () => selectQuestion(question));
        questionsList.appendChild(item);
    });
    
    // 显示问题容器
    questionsContainer.classList.remove('d-none');
}

// 选择问题
function selectQuestion(question) {
    // 高亮选中的问题
    document.querySelectorAll('.question-item').forEach(item => {
        if (parseInt(item.dataset.id) === question.id) {
            item.classList.add('active');
        } else {
            item.classList.remove('active');
        }
    });
    
    // 显示当前问题
    const currentQuestionContainer = document.getElementById('currentQuestion');
    const currentQuestionText = document.getElementById('currentQuestionText');
    const answerText = document.getElementById('answerText');
    
    currentQuestionId = question.id;
    currentQuestionText.textContent = question.content;
    answerText.value = '';
    
    currentQuestionContainer.classList.remove('d-none');
}

// 提交回答
async function handleSubmitAnswer() {
    if (!currentQuestionId || !questionSetId) {
        alert('请先选择问题');
        return;
    }

    const answerText = document.getElementById('answerText');
    if (!answerText.value.trim()) {
        alert('请输入回答');
        return;
    }

    const submitBtn = document.getElementById('submitAnswerBtn');
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> 评估中...';

    try {
        const response = await fetch('/evaluate/answer', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                questionSetId: questionSetId,
                questionId: currentQuestionId,
                answer: {
                    questionId: currentQuestionId,
                    content: answerText.value
                }
            })
        });

        const data = await response.json();

        if (response.ok) {
            displayEvaluation(data.evaluation);
        } else {
            throw new Error(data.error || '评估失败');
        }
    } catch (error) {
        alert(`错误: ${error.message}`);
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = '提交回答';
    }
}

// 显示评估结果
function displayEvaluation(evaluation) {
    const evaluationContainer = document.getElementById('evaluationContainer');
    const evaluationResult = document.getElementById('evaluationResult');
    
    // 创建评估卡片
    evaluationResult.innerHTML = `
        <div class="card evaluation-card">
            <div class="card-body">
                <div class="row">
                    <div class="col-md-3">
                        <div class="score-display">${evaluation.score}/10</div>
                        <p class="text-center mt-2">
                            ${getScoreDescription(evaluation.score)}
                        </p>
                    </div>
                    <div class="col-md-9">
                        <h5>评价</h5>
                        <p>${evaluation.feedback}</p>
                        <div class="feedback-section">
                            <h5>改进建议</h5>
                            <p>${evaluation.suggestions}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    // 显示评估容器
    evaluationContainer.classList.remove('d-none');
    
    // 滚动到评估结果
    evaluationResult.scrollIntoView({ behavior: 'smooth' });
}

// 根据分数获取评价描述
function getScoreDescription(score) {
    if (score >= 9) return '优秀';
    if (score >= 7) return '良好';
    if (score >= 5) return '一般';
    if (score >= 3) return '需改进';
    return '不足';
}
