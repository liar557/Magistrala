package main

import (
	"encoding/json"
	"io"
	"liar/llm"
	"net/http"
)

// AnalyzeHandler 返回一个用于智能分析的 HTTP 处理函数
func AnalyzeHandler(client *llm.OllamaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Messages []map[string]interface{} `json:"messages"`
		}
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "请求格式错误", http.StatusBadRequest)
			return
		}
		result, err := llm.AnalyzeMessages(client, req.Messages)
		if err != nil {
			http.Error(w, "模型推理失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 直接返回结构化结果
		json.NewEncoder(w).Encode(result)
	}
}
