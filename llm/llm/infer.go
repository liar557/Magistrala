package llm

import "fmt"

// AnalyzeMessages 是智慧农业 LLM 推理分析的主流程函数。
// 它负责：
// 1. 构建 prompt（调用 BuildPrompt）。
// 2. 调用大模型客户端推理（client.Infer）。
// 3. 剔除 <think> 标签内容（RemoveThinkSection）。
// 4. 去除首尾空白（TrimSpaceAll）。
// 5. 解析为结构体（JSONToStruct）。
// 返回结构化分析结果，供前端或业务代码直接使用。
func AnalyzeMessages(client LLMClient, messages []map[string]interface{}) (AnalysisResult, error) {
	// 1. 构建 Ollama 多模态消息体
	content, err := BuildOllamaMessages(messages)
	if err != nil {
		return AnalysisResult{}, err
	}

	// 2. 调用大模型客户端推理（假设 InferMultimodal 是多模态推理接口）
	result, err := client.InferMultimodal(content)
	if err != nil {
		return AnalysisResult{}, err
	}

	// fmt.Println("大模型原始输出：", result)

	// 3. 剔除 <think> 标签内容
	cleaned := RemoveThinkSection(result)

	// 4. 去除首尾空白
	cleaned = TrimSpaceAll(cleaned)

	// 5. 提取 JSON 子串
	jsonStr := ExtractJSON(cleaned)
	if jsonStr == "" {
		return AnalysisResult{}, fmt.Errorf("未检测到有效JSON，请检查大模型输出")
	}

	// 6. 解析为结构体
	return JSONToStruct(jsonStr)
}
