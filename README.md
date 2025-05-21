# AI简历面试助手

AI简历面试助手是一个基于Go语言的Web应用程序，可以帮助求职者准备面试。通过上传简历和职位描述，系统会生成针对性的面试问题，并对回答进行评估和反馈。

## 功能特点

- 解析简历文件（支持PDF和TXT格式）
- 解析职位描述文件（支持PDF和TXT格式）
- 基于简历和职位描述生成针对性面试问题
- 评估面试回答质量并提供反馈
- 提供改进建议和评分
- 用户友好的Web界面
- 集成OCR功能，支持多种文件格式的文本提取

## OCR功能

为了更好地处理各种格式的简历和职位描述文件，本系统集成了高级OCR（光学字符识别）功能：

- **多种OCR选项**：支持通过OCR.space API（云端）或Tesseract（本地）进行文本提取
- **智能格式处理**：从PDF、PNG、JPG等格式中提取文本内容
- **Token优化**：自动处理提取文本，确保不超过AI处理的token限制
- **自动回退**：OCR处理失败时，自动回退到传统解析方法
- **多语言支持**：支持中英文简历和职位描述的OCR处理

这些功能确保系统能够有效地从各种文件格式中提取文本，为后续的AI分析提供高质量的输入。

## 技术栈

- 后端：Go (Gin Web Framework)
- 前端：HTML, CSS, JavaScript (Bootstrap 5)
- AI：支持Grok 3 API和OpenAI API
- 文件处理：PDF解析库
- OCR：支持Tesseract OCR（本地）和OCR.space API（云端）

## 快速开始

### 前提条件

- Go 1.16+
- Grok 3 API或OpenAI API密钥
- （可选）Tesseract OCR（用于本地OCR处理）
- （可选）OCR.space API密钥（用于云端OCR处理）

### 安装步骤

1. 克隆项目仓库

```bash
git clone https://github.com/10yihang/resume-ai-interview.git
cd resume-ai-interview
```

2. 安装依赖项

```bash
go mod download
```

3. 创建配置文件

```bash
cp .env.example .env
```

4. 编辑.env文件，填入你的API密钥

```
# 优先使用Grok 3 API
GROK3_API_KEY=your_grok3_api_key_here
# 或者使用OpenAI API（如果没有Grok 3 API密钥）
OPENAI_API_KEY=your_openai_api_key_here
```

5. 运行应用程序

```bash
go run cmd/server/main.go
```

6. 打开浏览器访问 http://localhost:8080

### OCR设置（可选）

如果需要处理PDF简历或职位描述，可以通过以下两种方式启用OCR功能：

#### 方式1：使用OCR.space API（推荐，无需本地安装）

1. 在[OCR.space](https://ocr.space/ocrapi)注册并获取免费API密钥
2. 在`.env`配置文件中设置：
```bash
OCR_SPACE_API_KEY=your_api_key_here
USE_OCR=true
```

#### 方式2：使用Tesseract OCR（本地处理）

1. 安装Tesseract OCR

在Windows上：
```powershell
winget install --id UB-Mannheim.TesseractOCR
```

在Ubuntu上：
```bash
sudo apt-get install tesseract-ocr
sudo apt-get install tesseract-ocr-chi-sim # 中文支持
```

在macOS上：
```bash
brew install tesseract
brew install tesseract-lang # 安装语言包
```

2. 在`.env`配置文件中设置：
```bash
TESSERACT_PATH=C:/Program Files/Tesseract-OCR/tesseract.exe # Windows路径示例
USE_OCR=true
```

## 使用方法

1. 上传你的简历（PDF或TXT格式）
2. 上传目标职位描述（PDF或TXT格式）
3. 点击"生成面试问题"按钮
4. 回答生成的面试问题
5. 查看评估结果和改进建议

## 项目结构

```
resume-ai-interview/
├── api/                # API处理程序
│   └── handlers/       # 请求处理函数
├── cmd/                # 应用程序入口
│   └── server/         # 服务器入口
├── config/             # 配置管理
├── internal/           # 内部包
│   ├── ai/             # AI问题生成
│   ├── interview/      # 面试评估
│   └── parser/         # 文件解析器
├── models/             # 数据模型
├── static/             # 静态资源
│   ├── css/            # 样式表
│   └── js/             # JavaScript文件
└── templates/          # HTML模板
```

## 许可证

[MIT](LICENSE)

## 致谢

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Grok 3 API](https://grok.x)
- [Bootstrap](https://getbootstrap.com/)
