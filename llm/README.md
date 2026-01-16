# LLM 模块说明

## 功能概述
- 拉取 Magistrala 指定域/通道的最新消息，补充分区信息，构造 LLM 输入。
- 调用本地/远程大模型输出“区域级（分区级）指令”。
- 将区域指令映射到具体执行设备（clientId），并通过执行模块下发控制。

## 目录结构
- llm/llm：核心逻辑与客户端，关键文件：
  - orchestrator：编排与映射逻辑，[llm/llm/orchestrator.go](llm/llm/orchestrator.go)
  - magistrala_reader：拉取通道消息，[llm/llm/magistrala_reader.go](llm/llm/magistrala_reader.go)
  - infer & ollama_client：LLM 调用与输出清洗，[llm/llm/infer.go](llm/llm/infer.go)、[llm/llm/ollama_client.go](llm/llm/ollama_client.go)
- llm/llm_api：对外 HTTP 服务（plan-and-execute），[llm/llm_api/main.go](llm/llm_api/main.go)

## 配置
  - magistrala.baseUrl / userToken / domainId / channelId / messagePort（必填，无默认值）
  - executor.baseUrl
  - mapping.path（供分区补全与设备映射使用，必填，无默认值）

## 设备注册表 schema
- 路径：data/device_registry.json
- 结构：domains[].channels[].partitions[]，分区下挂 sensors（publisher/clientId）与 executors（可执行 clientId）。
- 示例片段：
```json
{
  "domains": [
    {
      "domainId": "<domain>",
      "channels": [
        {
          "channelId": "<channel>",
          "partitions": [
            {
              "partitionId": "field-A",
              "partitionName": "A区",
              "sensors": ["<sensor-clientId>"],
              "executors": ["<executor-clientId>"]
            }
          ]
        }
      ]
    }
  ]
}