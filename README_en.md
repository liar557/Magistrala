## Usage Steps

### 1. Start the Magistrala Main Program

Please follow the official Magistrala documentation to start the main program, ensuring your environment and dependencies are properly configured.

### 2. Start the Image Upload Service

Navigate to the image upload service directory, build and run the Go service:

```bash
cd image_upload/image_upload_server_main
go run main.go
```
> Make sure the Go environment is installed and `go.mod` is properly configured.

### Additional Note: Image Upload Service Configuration

Currently, the image upload service can only be started via `liar/upload_main/main.go`.  
Before running, please modify the relevant information in `UploadConfig` according to your actual needs.  
Additionally, you need to create the corresponding channel and client configurations in advance to ensure the upload process works correctly.

Example startup process:

```bash
cd liar/upload_main
go run main.go
```

> Make sure the Go environment is installed and adjust the configuration parameters according to your scenario.

### 3. Start the LLM Service

Ensure Ollama is installed locally and pull the required model:

```bash
ollama pull gemma3:12b
ollama serve
```
> The LLM service is used for intelligent analysis modules. Please keep it running.

Navigate to the LLM directory, build and run the Go service:

```bash
cd llm/llm_api
go run main.go handle.go
```
> Also ensure the Go environment and `go.mod` are properly configured.

### 4. Start the Frontend UI (Optional)

Navigate to the frontend directory, install dependencies, and start the development server:

```bash
cd magistrala-ui
npm install
npm run dev
```
> After starting, access the frontend page in your browser using the address shown in the terminal (e.g., http://localhost:5173).

---

If you have any questions or encounter issues, please refer to the README files of each module or contact