package main

import (
	"fmt"
	"llm/llm"
)

func testAnalyzeMessages() {
	client := llm.NewOllamaClient("http://localhost:11434", "qwen3:8b")
	messages := []map[string]interface{}{
		{"value": 25.1, "name": "温度", "unit": "°C"},
		{"value": 68.2, "name": "湿度", "unit": "%"},
		{"string_value": "作物叶片有轻微发黄", "name": "备注", "unit": ""},
		{"value": 1012, "name": "气压", "unit": "hPa"},
	}
	analysis, err := llm.AnalyzeMessages(client, messages)
	if err != nil {
		fmt.Println("推理失败:", err)
		return
	}

	fmt.Println("一句话总结:", analysis.Summary)
	fmt.Println("诊断:", analysis.Diagnosis)

	fmt.Print("风险: ")
	for i, risk := range analysis.Risks {
		if i > 0 {
			fmt.Print("，")
		}
		fmt.Print(risk)
	}
	fmt.Println()

	fmt.Print("建议: ")
	for i, sug := range analysis.Suggestions {
		if i > 0 {
			fmt.Print("，")
		}
		fmt.Print(sug)
	}
	fmt.Println()

	fmt.Println("详细分析:", analysis.RawAnalysis)
}

// func testImageToBase64() {
// 	url := "http://localhost:18080/uploaded/1754468898511455067.png"
// 	b64, mimeType, err := llm.UrlToBase64(url)
// 	if err != nil {
// 		fmt.Println("下载或编码失败:", err)
// 		return
// 	}
// 	fmt.Println("MIME类型:", mimeType)
// 	fmt.Println("Base64前100字符:", b64[:100])
// 	fmt.Printf("data:%s;base64,%s...\n", mimeType, b64[:100])
// }

// func testImageToBase64AndSaveHTML() {
// 	url := "http://localhost:18080/uploaded/1754468898511455067.png"
// 	b64, mimeType, err := llm.UrlToBase64(url)
// 	if err != nil {
// 		fmt.Println("下载或编码失败:", err)
// 		return
// 	}
// 	// 生成 HTML 内容
// 	html := fmt.Sprintf(`
// <!DOCTYPE html>
// <html>
// <body>
// <img src="data:%s;base64,%s" alt="测试图片"/>
// </body>
// </html>
// `, mimeType, b64)
// 	// 保存为文件
// 	err = os.WriteFile("test_image.html", []byte(html), 0644)
// 	if err != nil {
// 		fmt.Println("写入HTML文件失败:", err)
// 		return
// 	}
// 	fmt.Println("已生成 test_image.html，请用浏览器打开查看图片内容。")
// }

func main() {
	testAnalyzeMessages() // 如需测试 LLM 推理，取消注释
	// testImageToBase64() // 测试图片下载和编码
	// testImageToBase64AndSaveHTML()
}
