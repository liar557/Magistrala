package main

import (
	"encoding/json"
	"fmt"
	"io"
	core "llm/core"
	llm "llm/llm"
	"net/http"
	"os"
	"path/filepath"
)

type result struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// 兼容旧调用，不覆盖 prompt。
func makeInfer(client llm.LLMClient) func([]map[string]interface{}) (string, error) {
	return makeInferWithPrompt(client, "")
}

// 对应旧的逻辑，直接调用设备去执行，已废弃
// 从 config.json 读取 Magistrala 与执行模块配置。
// 默认路径：agriDeviceExecutor/internal/config/config.json
// 结构示例：
//
//	{
//	  "magistrala": {"baseUrl":"http://localhost","userToken":"...","domainId":"...","channelId":"...","messagePort":9009},
//	  "executor": {"baseUrl":"http://127.0.0.1:8090"}
//	}
func NewDefaultOrchestratorFromConfig() (*core.Orchestrator, error) {
	// 读取 llm 模块的配置文件
	cfgPath := filepath.Clean("config/config.json")
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("打开配置失败: %w", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}
	var raw struct {
		Magistrala struct {
			MessagePort int `json:"messagePort"`
		} `json:"magistrala"`
		Executor struct {
			BaseURL string `json:"baseUrl"`
		} `json:"executor"`
		Mapping struct {
			Path string `json:"path"`
		} `json:"mapping"`
		LLM struct {
			Type     string `json:"type"`
			Model    string `json:"model"`
			Endpoint string `json:"endpoint"`
			APIKey   string `json:"apiKey"`
		} `json:"llm"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	if raw.Magistrala.MessagePort == 0 {
		return nil, fmt.Errorf("magistrala.messagePort 未配置")
	}
	if raw.Executor.BaseURL == "" {
		return nil, fmt.Errorf("executor.baseUrl 未配置")
	}
	if raw.Mapping.Path == "" {
		return nil, fmt.Errorf("mapping.path 未配置")
	}
	// 读取共享 baseUrl/userToken（必须存在）；这里不依赖 servicePort。
	type sharedMagCfg struct {
		BaseURL   string `json:"baseUrl"`
		UserToken string `json:"userToken"`
	}
	tryLoadSharedMag := func() *sharedMagCfg {
		candidates := []string{
			"data/magistrala.json",
			"../data/magistrala.json",
			"../../data/magistrala.json",
		}
		if exe, err := os.Executable(); err == nil {
			ed := filepath.Dir(exe)
			candidates = append(candidates,
				filepath.Join(ed, "data/magistrala.json"),
				filepath.Join(ed, "../data/magistrala.json"),
			)
		}
		for _, p := range candidates {
			b, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			var out sharedMagCfg
			if json.Unmarshal(b, &out) == nil {
				return &out
			}
		}
		return nil
	}
	sm := tryLoadSharedMag()
	if sm == nil || sm.BaseURL == "" || sm.UserToken == "" {
		return nil, fmt.Errorf("共享配置 data/magistrala.json 缺失 baseUrl/userToken")
	}
	// servicePort 仅取本地 config.json；共享配置不提供也不影响。
	orch := &core.Orchestrator{
		BaseURL:      sm.BaseURL,
		MessagePort:  raw.Magistrala.MessagePort,
		DomainID:     "", // 运行时传入
		ChannelID:    "", // 运行时传入
		Token:        sm.UserToken,
		ExecutorBase: raw.Executor.BaseURL,
		MappingPath:  raw.Mapping.Path,
		// Infer 在下方根据 LLMClient 设置
	}
	// 映射路径直接随 orchestrator 传递，不使用全局默认

	// 构建 LLM 客户端
	var client llm.LLMClient
	// 目前仅支持 Ollama 客户端
	if raw.LLM.Endpoint == "" {
		raw.LLM.Endpoint = "http://localhost:11434"
	}
	if raw.LLM.Model == "" {
		raw.LLM.Model = "qwen2.5:latest"
	}
	client = &llm.OllamaClient{Endpoint: raw.LLM.Endpoint, Model: raw.LLM.Model}
	orch.Infer = makeInfer(client)
	return orch, nil
}

// PlanAndExecuteHandler 一次性拉取→推理→执行。
// POST /llm/plan-and-execute {"limit":10}
func PlanAndExecuteHandler(w http.ResponseWriter, r *http.Request, orch *core.Orchestrator) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(result{Code: 405, Message: "method not allowed"})
		return
	}
	var body struct {
		Limit     int    `json:"limit"`
		DomainID  string `json:"domainId"`
		ChannelID string `json:"channelId"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Limit <= 0 {
		body.Limit = 10
	}

	if body.DomainID == "" || body.ChannelID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(result{Code: 400, Message: "domainId/channelId 不能为空"})
		return
	}
	// 覆盖域/通道用于本次请求
	ov := *orch
	ov.DomainID = body.DomainID
	ov.ChannelID = body.ChannelID

	cmds, err := ov.RunOnce(body.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(result{Code: 500, Message: err.Error()})
		return
	}
	// 下发执行
	for _, c := range cmds {
		_ = ov.SendToExecutor(c)
	}
	b, _ := json.Marshal(cmds)
	_ = json.NewEncoder(w).Encode(result{Code: 1000, Message: "ok", Data: b})
}

// 使用真实推理：基于 llm.AnalyzeRegionCommands 和配置创建的 LLMClient。
// 可传入 prompt 文本覆盖默认提示词。
func makeInferWithPrompt(client llm.LLMClient, prompt string) func([]map[string]interface{}) (string, error) {
	return func(messages []map[string]interface{}) (string, error) {
		return llm.AnalyzeRegionCommandsWithPrompt(client, messages, prompt)
	}
}

// NewControlAdapterFromConfig 从 config/config.json 读取配置并构造 ControlAdapter。
// 必填：magistrala.messagePort、mapping.path、controlService.baseUrl、共享 data/magistrala.json 的 baseUrl/userToken。
func NewControlAdapterFromConfig() (*core.ControlAdapter, error) {
	cfgPath := filepath.Clean("config/config.json")
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("打开配置失败: %w", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	var raw struct {
		Magistrala struct {
			MessagePort int `json:"messagePort"`
		} `json:"magistrala"`
		Mapping struct {
			Path string `json:"path"`
		} `json:"mapping"`
		Prompt struct {
			Path string `json:"path"`
		} `json:"prompt"`
		ControlService struct {
			BaseURL string `json:"baseUrl"`
		} `json:"controlService"`
		LLM struct {
			Type     string `json:"type"`
			Model    string `json:"model"`
			Endpoint string `json:"endpoint"`
			APIKey   string `json:"apiKey"`
		} `json:"llm"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	if raw.Magistrala.MessagePort == 0 {
		return nil, fmt.Errorf("magistrala.messagePort 未配置")
	}
	if raw.Mapping.Path == "" {
		return nil, fmt.Errorf("mapping.path 未配置")
	}
	if raw.ControlService.BaseURL == "" {
		return nil, fmt.Errorf("controlService.baseUrl 未配置")
	}

	// 读取共享 baseUrl/userToken（必须存在）；这里不依赖 messagePort。
	type sharedMagCfg struct {
		BaseURL   string `json:"baseUrl"`
		UserToken string `json:"userToken"`
	}
	tryLoadSharedMag := func() *sharedMagCfg {
		candidates := []string{
			"data/magistrala.json",
			"../data/magistrala.json",
			"../../data/magistrala.json",
		}
		if exe, err := os.Executable(); err == nil {
			ed := filepath.Dir(exe)
			candidates = append(candidates,
				filepath.Join(ed, "data/magistrala.json"),
				filepath.Join(ed, "../data/magistrala.json"),
			)
		}
		for _, p := range candidates {
			b, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			var out sharedMagCfg
			if json.Unmarshal(b, &out) == nil {
				return &out
			}
		}
		return nil
	}
	sm := tryLoadSharedMag()
	if sm == nil || sm.BaseURL == "" || sm.UserToken == "" {
		return nil, fmt.Errorf("共享配置 data/magistrala.json 缺失 baseUrl/userToken")
	}

	adapter := &core.ControlAdapter{
		BaseURL:     sm.BaseURL,
		MessagePort: raw.Magistrala.MessagePort,
		DomainID:    "", // 运行时传入
		ChannelID:   "", // 运行时传入
		Token:       sm.UserToken,
		MappingPath: raw.Mapping.Path,
		ControlBase: raw.ControlService.BaseURL,
		// Infer 在下方根据 LLMClient 设置
	}

	var promptText string
	if raw.Prompt.Path != "" {
		if pb, err := os.ReadFile(filepath.Clean(raw.Prompt.Path)); err == nil {
			promptText = string(pb)
		}
	}

	// 构建 LLM 客户端（默认 Ollama）
	if raw.LLM.Endpoint == "" {
		raw.LLM.Endpoint = "http://localhost:11434"
	}
	if raw.LLM.Model == "" {
		raw.LLM.Model = "qwen2.5:latest"
	}
	client := &llm.OllamaClient{Endpoint: raw.LLM.Endpoint, Model: raw.LLM.Model}
	adapter.Infer = makeInferWithPrompt(client, promptText)
	return adapter, nil
}

// PlanAndSendToControlHandler：拉取→推理→转任务→下发控制服务
// POST /llm/plan-and-send {"limit":10,"domainId":"...","channelId":"..."}
func PlanAndSendToControlHandler(w http.ResponseWriter, r *http.Request, baseAdapter *core.ControlAdapter) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(result{Code: 405, Message: "method not allowed"})
		return
	}
	var body struct {
		Limit     int    `json:"limit"`
		DomainID  string `json:"domainId"`
		ChannelID string `json:"channelId"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Limit <= 0 {
		body.Limit = 10
	}
	if body.DomainID == "" || body.ChannelID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(result{Code: 400, Message: "domainId/channelId 不能为空"})
		return
	}

	adapter := *baseAdapter
	adapter.DomainID = body.DomainID
	adapter.ChannelID = body.ChannelID

	tasks, err := adapter.RunTasks(body.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(result{Code: 500, Message: err.Error()})
		return
	}
	for _, t := range tasks {
		if err := adapter.PostTask(t); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(result{Code: 500, Message: err.Error()})
			return
		}
	}
	_ = json.NewEncoder(w).Encode(result{Code: 1000, Message: "ok", Data: mustJSON(tasks)})
}

// 辅助：序列化任务列表
func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
