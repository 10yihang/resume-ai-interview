package parser

import (
	"encoding/json"
	"regexp"
	"strings"
)

// sanitizeField 清理字符串字段，移除多余空格和特殊字符
func sanitizeField(s string) string {
	// 移除前后空白
	s = strings.TrimSpace(s)
	return s
}

// sanitizeStringArray 清理字符串数组
func sanitizeStringArray(arr []string) []string {
	if arr == nil {
		return []string{}
	}

	result := make([]string, 0, len(arr))
	for _, s := range arr {
		if s = sanitizeField(s); s != "" {
			result = append(result, s)
		}
	}
	return result
}

// tryFixJSONFormat 尝试修复常见的JSON格式错误
func tryFixJSONFormat(jsonStr string) string {
	// 移除可能导致错误的控制字符
	re := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	jsonStr = re.ReplaceAllString(jsonStr, "")

	// 确保JSON对象正确关闭
	bracketsCount := 0
	for _, c := range jsonStr {
		if c == '{' {
			bracketsCount++
		} else if c == '}' {
			bracketsCount--
		}
	}

	// 如果花括号不匹配，尝试修复
	if bracketsCount > 0 {
		// 缺少右括号，添加
		jsonStr += strings.Repeat("}", bracketsCount)
	} else if bracketsCount < 0 {
		// 缺少左括号，这种情况更复杂，只删除多余的右括号
		jsonStr = removeExtraBrackets(jsonStr, '}', '{', -bracketsCount)
	}

	// 尝试验证JSON是否有效
	var js interface{}
	if err := json.Unmarshal([]byte(jsonStr), &js); err != nil {
		// 如果仍然无效，尝试更激进的修复
		jsonStr = aggressiveJSONFix(jsonStr)
	}

	return jsonStr
}

// removeExtraBrackets 移除多余的括号
func removeExtraBrackets(jsonStr string, bracket rune, matchingBracket rune, count int) string {
	runes := []rune(jsonStr)
	removed := 0

	// 从末尾开始遍历，移除括号
	for i := len(runes) - 1; i >= 0 && removed < count; i-- {
		if runes[i] == bracket {
			// 检查是否是平衡的括号
			balanced := false
			nested := 0

			// 检查从这个位置向前是否能找到匹配的开括号
			for j := i - 1; j >= 0; j-- {
				if runes[j] == bracket {
					nested++
				} else if runes[j] == matchingBracket {
					if nested == 0 {
						balanced = true
						break
					}
					nested--
				}
			}

			// 如果没有匹配的开括号，这个括号是多余的
			if !balanced {
				runes = append(runes[:i], runes[i+1:]...)
				removed++
			}
		}
	}

	return string(runes)
}

// aggressiveJSONFix 采用更激进的手段修复JSON
func aggressiveJSONFix(jsonStr string) string {
	// 尝试提取有效的JSON部分
	validStart := strings.Index(jsonStr, "{")
	validEnd := strings.LastIndex(jsonStr, "}")

	if validStart >= 0 && validEnd > validStart {
		jsonStr = jsonStr[validStart : validEnd+1]
	}

	// 修复常见的引号不匹配问题
	return fixUnbalancedQuotes(jsonStr)
}

// fixUnbalancedQuotes 修复不平衡的引号
func fixUnbalancedQuotes(jsonStr string) string {
	runes := []rune(jsonStr)
	inQuote := false
	escaping := false

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' {
			escaping = !escaping
		} else {
			if runes[i] == '"' && !escaping {
				inQuote = !inQuote
			}
			escaping = false
		}
	}

	// 如果字符串结束时引号仍然是打开状态，添加一个关闭引号
	if inQuote {
		jsonStr += `"`
	}

	return jsonStr
}
