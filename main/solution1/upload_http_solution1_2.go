package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

// SenML消息结构体
type SenMLRecord struct {
	Bn string  `json:"bn,omitempty"`
	N  string  `json:"n"`
	U  string  `json:"u"`
	V  float64 `json:"v"`
}

// 图片消息结构体（JSON格式）
type ImagePayload struct {
	ImageBase64 string `json:"image_base64"`
	Filename    string `json:"filename"`
}
type ImageMessage struct {
	Channel   string       `json:"channel"`
	Created   int64        `json:"created"`
	Subtopic  string       `json:"subtopic"`
	Publisher string       `json:"publisher"`
	Protocol  string       `json:"protocol"`
	Payload   ImagePayload `json:"payload"`
}

func generateSenMLData(baseName string) []SenMLRecord {
	return []SenMLRecord{
		{Bn: baseName, N: "lumen", U: "CD", V: 53.1},
	}
}

func generateImageMessage(cfg Config, imagePath string) ([]ImageMessage, error) {
	data, err := os.ReadFile(imagePath) // 修改这里
	if err != nil {
		return nil, err
	}
	imgBase64 := base64.StdEncoding.EncodeToString(data)
	payload := ImagePayload{
		ImageBase64: imgBase64,
		Filename:    imagePath,
	}
	msg := ImageMessage{
		Channel:   cfg.ChannelID,
		Created:   time.Now().Unix(),
		Subtopic:  cfg.Subtopic,
		Publisher: "user1",
		Protocol:  "http",
		Payload:   payload,
	}
	return []ImageMessage{msg}, nil
}

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
				InsecureSkipVerify: true, // 跳过证书校验
			},
		},
	}

	// log.Printf("请求URL: %s", url)
	// log.Printf("请求头: %v", req.Header)
	// log.Printf("请求体: %s", string(data))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("上传请求失败: %v", err)
		return false
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body) // 修改这里
	log.Printf("响应状态码: %d", resp.StatusCode)
	log.Printf("响应内容: %s", string(body))
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func loadCACert(caCertPath string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caCertPath) // 修改这里
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

	// 上传类型选择
	uploadType := "image"   // "senml" 或 "image"
	imagePath := "test.png" // 图片路径

	if cfg.UploadMode == "once" {
		log.Println("开始单次上传")
		var success bool
		if uploadType == "senml" {
			senmlData := generateSenMLData(cfg.BaseName)
			data, _ := json.Marshal(senmlData)
			success = uploadData(cfg, data, "application/senml+json")
		} else {
			imgMsg, err := generateImageMessage(cfg, imagePath)
			if err != nil {
				log.Printf("图片读取失败: %v", err)
				os.Exit(1)
			}
			data, _ := json.Marshal(imgMsg)
			success = uploadData(cfg, data, "application/json")
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
			if uploadType == "senml" {
				senmlData := generateSenMLData(cfg.BaseName)
				data, _ := json.Marshal(senmlData)
				success = uploadData(cfg, data, "application/senml+json")
			} else {
				imgMsg, err := generateImageMessage(cfg, imagePath)
				if err != nil {
					log.Printf("图片读取失败: %v", err)
					continue
				}
				data, _ := json.Marshal(imgMsg)
				success = uploadData(cfg, data, "application/json")
			}
			log.Printf("上传结果: %v", success)
			time.Sleep(time.Duration(cfg.Interval) * time.Second)
		}
	}
}
