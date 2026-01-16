## 使用步骤

### 1. 启动 Magistrala 主程序

请根据 Magistrala 官网的要求启动主程序，确保环境和依赖配置正确。

### 2. 启动图片上传服务

进入图片上传服务目录，编译并运行 Go 服务：

```bash
cd image_upload/image_upload_server_main
go run main.go
```
> 请确保已安装 Go 环境，并且 `go.mod` 配置正确。

### 补充说明：图片上传服务配置

目前图片上传服务只能通过 `liar/upload_main/main.go` 启动。  
运行前，请根据实际需求，修改 `UploadConfig` 中的相关信息。  
此外，需提前创建好对应的 channel、client 等配置，以保证上传流程正常。

示例启动流程：

```bash
cd liar/upload_main
go run main.go
```

> 请确保 Go 环境已安装，并根据自身场景调整配置。

### 3. 启动 LLM 服务

确保本地已安装 Ollama，并拉取所需模型：

```bash
ollama pull gemma3:12b
ollama serve
```
> LLM 服务用于智能分析模块，请保持其运行状态。

进入 LLM 目录，编译并运行 Go 服务：

```bash
cd llm/llm_api
go run main.go handle.go
```
> 同样需要确保 Go 环境和 `go.mod` 配置无误。

### 4. 启动前端 UI（可选）

进入前端目录，安装依赖并启动开发服务器：

```bash
cd magistrala-ui
npm install
npm run dev
```
> 启动后请根据终端提示的地址（如 http://localhost:5173）在浏览器中访问前端页面。

---

如有疑问或遇到问题，请查阅各模块的 README 文件或联系项目维护者
