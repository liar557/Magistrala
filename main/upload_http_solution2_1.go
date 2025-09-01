// 功能说明：
// 本程序用于通过 HTTP 上传 SenML 格式数据。
// 包含两种消息结构体：普通 SenML 数据（数值型），图片 SenML 数据（图片路径字符串）。
// 支持单次上传和定时上传。

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// 配置参数结构体
type Config struct {
	BaseURL      string
	DomainID     string
	ChannelID    string
	Subtopic     string
	ClientSecret string
	CACertPath   string
	BaseName     string
	Timeout      int
	UploadMode   string // "once" 或 "interval"
	Interval     int    // 定时上传间隔（秒）
}

// 普通 SenML 消息结构体（数值型）
type SenMLValueRecord struct {
	Bn string  `json:"bn,omitempty"` // Base Name
	N  string  `json:"n"`            // Name
	U  string  `json:"u,omitempty"`  // Unit
	V  float64 `json:"v,omitempty"`  // Value
}

// 图片 SenML 消息结构体（图片路径字符串）
type SenMLImageRecord struct {
	Bn string `json:"bn,omitempty"` // Base Name
	N  string `json:"n"`            // Name
	U  string `json:"u,omitempty"`  // Unit
	Vs string `json:"vs,omitempty"` // String Value（图片路径或URL）
}

// 生成普通 SenML 数据
func generateSenMLData(baseName string) []SenMLValueRecord {
	return []SenMLValueRecord{
		{Bn: baseName, N: "lumen", U: "CD", V: 53.1},
	}
}

// 生成图片 SenML 数据（此处 imagePath 应为图片URL）
func generateSenMLImageData(baseName, imageURL string) []SenMLImageRecord {
	return []SenMLImageRecord{
		{
			Bn: baseName,
			N:  "image",
			U:  "path",
			Vs: imageURL,
		},
	}
}

// 上传图片到 image_upload_server，返回图片访问 URL
func uploadImageToServer(imagePath string, uploadURL string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("无法打开图片: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(imagePath))
	if err != nil {
		return "", fmt.Errorf("创建表单失败: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("写入文件内容失败: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("上传失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
	}
	return string(respBody), nil // 返回图片访问 URL
}

// 上传数据到 HTTP 服务
func uploadData(cfg Config, data []byte, contentType string) bool {
	url := fmt.Sprintf("%s/http/m/%s/c/%s/%s", cfg.BaseURL, cfg.DomainID, cfg.ChannelID, cfg.Subtopic)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
		return false
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Client "+cfg.ClientSecret)

	// 加载CA证书
	caCertPool, err := loadCACert(cfg.CACertPath)
	if err != nil {
		log.Printf("加载CA证书失败: %v", err)
		return false
	}
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: true, // 跳过证书校验，适合本地测试
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("上传请求失败: %v", err)
		return false
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Printf("响应状态码: %d", resp.StatusCode)
	log.Printf("响应内容: %s", string(body))
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// 加载CA证书
func loadCACert(caCertPath string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}

func main() {
	cfg := Config{
		BaseURL:      "https://localhost",
		DomainID:     "562d704a-c442-499a-aff3-223f580bf6b3",
		ChannelID:    "b0ec13df-9ff0-48b9-9cb6-b3be072e7c99",
		Subtopic:     "temperature",
		ClientSecret: "9c3acbd2-ce7c-4ba3-b6ab-8792e310002c",
		CACertPath:   "CA/ca.crt",
		BaseName:     "ljp",
		Timeout:      10,
		UploadMode:   "once", // "once" 或 "interval"
		Interval:     30,
	}

	uploadType := "senml_image" // "senml" 普通数据，"senml_image" 图片路径数据
	imageDir := "."             // 当前代码文件所在目录
	imageName := "test.png"
	imagePath := fmt.Sprintf("%s/%s", imageDir, imageName)

	imageUploadURL := "http://localhost:18080/upload" // image_upload_server服务地址

	if cfg.UploadMode == "once" {
		log.Println("开始单次上传")
		var success bool
		switch uploadType {
		case "senml":
			senmlData := generateSenMLData(cfg.BaseName)
			data, _ := json.Marshal(senmlData)
			success = uploadData(cfg, data, "application/senml+json")
		case "senml_image":
			// 先上传图片，获取图片URL
			imgURL, err := uploadImageToServer(imagePath, imageUploadURL)
			if err != nil {
				log.Printf("图片上传失败: %v", err)
				os.Exit(1)
			}
			log.Printf("图片访问 URL: %s", imgURL)
			// 用图片URL生成SenML数据
			senmlImgData := generateSenMLImageData(cfg.BaseName, imgURL)
			data, _ := json.Marshal(senmlImgData)
			success = uploadData(cfg, data, "application/senml+json")
		default:
			log.Printf("未知上传类型: %s", uploadType)
			os.Exit(1)
		}
		log.Printf("上传结果: %v", success)
		if success {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	} else {
		log.Printf("开始定时上传，间隔: %d秒", cfg.Interval)
		for {
			var success bool
			switch uploadType {
			case "senml":
				senmlData := generateSenMLData(cfg.BaseName)
				data, _ := json.Marshal(senmlData)
				success = uploadData(cfg, data, "application/senml+json")
			case "senml_image":
				imgURL, err := uploadImageToServer(imagePath, imageUploadURL)
				if err != nil {
					log.Printf("图片上传失败: %v", err)
					continue
				}
				log.Printf("图片访问 URL: %s", imgURL)
				senmlImgData := generateSenMLImageData(cfg.BaseName, imgURL)
				data, _ := json.Marshal(senmlImgData)
				success = uploadData(cfg, data, "application/senml+json")
			default:
				log.Printf("未知上传类型: %s", uploadType)
				continue
			}
			log.Printf("上传结果: %v", success)
			time.Sleep(time.Duration(cfg.Interval) * time.Second)
		}
	}
}
