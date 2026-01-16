package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// orch, err := NewDefaultOrchestratorFromConfig()
	// if err != nil {
	// 	log.Fatalf("加载配置失败: %v", err)
	// }
	// mux.HandleFunc("/llm/plan-and-execute", func(w http.ResponseWriter, r *http.Request) {
	// 	PlanAndExecuteHandler(w, r, orch)
	// })
	ctrlAdapter, err := NewControlAdapterFromConfig()
	if err != nil {
		log.Fatalf("加载控制服务配置失败: %v", err)
	}
	mux.HandleFunc("/llm/plan-and-send", func(w http.ResponseWriter, r *http.Request) { // 新增路由
		PlanAndSendToControlHandler(w, r, ctrlAdapter)
	})
	log.Println("llm_api listening on :9000")
	if err := http.ListenAndServe(":9000", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
