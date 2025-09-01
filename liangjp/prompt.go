package llm

import (
	"fmt"
	"strings"
)

// 构建分析用的 prompt，强制要求只输出 JSON
func BuildPrompt(messages []map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString(`你是一个智慧农业专家，请分析以下传感器数据，并严格以如下JSON格式输出（不要输出多余内容，不要markdown，不要解释）：

	{
	"summary": "一句话总结本次分析结论",
	"diagnosis": "对作物状态的诊断",
	"risks": ["风险1", "风险2"],
	"suggestions": ["建议1", "建议2"],
	"raw_analysis": "详细推理过程"
	}

	数据如下：
	`)
	for i, msg := range messages {
		sb.WriteString(fmt.Sprintf("第%d条：", i+1))
		if v, ok := msg["value"]; ok {
			sb.WriteString(fmt.Sprintf("数值=%v；", v))
		}
		if sv, ok := msg["string_value"]; ok {
			sb.WriteString(fmt.Sprintf("文本/图片=%v；", sv))
		}
		if name, ok := msg["name"]; ok {
			sb.WriteString(fmt.Sprintf("名称=%v；", name))
		}
		if unit, ok := msg["unit"]; ok {
			sb.WriteString(fmt.Sprintf("单位=%v；", unit))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
