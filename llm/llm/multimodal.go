package llm

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 判断是否为图片URL
func isImageURL(url string) bool {
	lower := strings.ToLower(url)
	return (strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")) &&
		(strings.HasSuffix(lower, ".png") || strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") || strings.HasSuffix(lower, ".webp"))
}

// 判断是否为视频URL
func isVideoURL(url string) bool {
	lower := strings.ToLower(url)
	return (strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")) &&
		(strings.HasSuffix(lower, ".mp4") || strings.HasSuffix(lower, ".mov") || strings.HasSuffix(lower, ".avi") || strings.HasSuffix(lower, ".mkv"))
}

// 下载文件并转base64
func urlToBase64(url string) (string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	mimeType := http.DetectContentType(data)
	b64 := base64.StdEncoding.EncodeToString(data)
	return b64, mimeType, nil
}

// 构建 Ollama 多模态消息体
func BuildOllamaMessages(messages []map[string]interface{}) ([]map[string]interface{}, error) {
	var content []map[string]interface{}
	// 分析要求
	content = append(content, map[string]interface{}{
		"type": "text",
		"text": `你是一个智慧农业专家，请分析以下传感器数据和图片，并严格以如下JSON格式输出（不要输出多余内容，不要markdown，不要解释）：
		{
			"summary": "一句话总结本次分析结论",
			"diagnosis": "对作物状态的诊断",
			"risks": ["风险1", "风险2"],
			"suggestions": ["建议1", "建议2"],
			"raw_analysis": "详细推理过程"
		}
		`,
	})
	for _, msg := range messages {
		if name, ok := msg["name"]; ok && name != "" {
			content = append(content, map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("名称：%v", name),
			})
		}
		if v, ok := msg["value"]; ok {
			content = append(content, map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("数值：%v", v),
			})
		}
		if unit, ok := msg["unit"]; ok && unit != "" {
			content = append(content, map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("单位：%v", unit),
			})
		}
		if sv, ok := msg["string_value"].(string); ok && sv != "" {
			switch {
			case isImageURL(sv):
				b64, mimeType, err := urlToBase64(sv)
				if err != nil {
					content = append(content, map[string]interface{}{
						"type": "text",
						"text": fmt.Sprintf("图片获取失败：%v", err),
					})
				} else {
					content = append(content, map[string]interface{}{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url": fmt.Sprintf("data:%s;base64,%s", mimeType, b64),
						},
					})
				}
			case isVideoURL(sv):
				b64, mimeType, err := urlToBase64(sv)
				if err != nil {
					content = append(content, map[string]interface{}{
						"type": "text",
						"text": fmt.Sprintf("视频获取失败：%v", err),
					})
				} else {
					content = append(content, map[string]interface{}{
						"type": "video_url",
						"video_url": map[string]interface{}{
							"url": fmt.Sprintf("data:%s;base64,%s", mimeType, b64),
						},
					})
				}
			default:
				content = append(content, map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("文本：%v", sv),
				})
			}
		}
	}
	return content, nil
}
