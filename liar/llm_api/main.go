package main

import (
	"liar/llm"
	"log"
	"net/http"
)

// CORS 包装器
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func main() {
	// 初始化 Ollama LLM 客户端（请根据实际模型名修改）
	client := llm.NewOllamaClient("http://localhost:11434", "gemma3:12b")

	// 使用 CORS 包装 AnalyzeHandler
	http.HandleFunc("/analyze", withCORS(AnalyzeHandler(client)))

	log.Println("LLM API服务已启动，监听端口: 8091")
	http.ListenAndServe(":8091", nil)
}
