// 功能说明：
// 本程序用于通过 HTTP 上传 SenML 格式数据和图片（图片以 base64 编码存储在 SenML 的 vd 字段）。
// 支持单次上传和定时上传，支持普通 SenML 数据和图片数据。
// 推荐用于小型图片或二进制内容的物联网场景。

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

// SenML消息结构体，支持图片（vd字段）
type SenMLRecord struct {
	Bn string  `json:"bn,omitempty"` // Base Name
	N  string  `json:"n"`            // Name
	U  string  `json:"u,omitempty"`  // Unit
	V  float64 `json:"v,omitempty"`  // Value
	Vd string  `json:"vd,omitempty"` // Data Value (base64图片)
}

// 生成普通 SenML 数据
func generateSenMLData(baseName string) []SenMLRecord {
	return []SenMLRecord{
		{Bn: baseName, N: "lumen", U: "CD", V: 53.1},
	}
}

// 生成 SenML 图片数据（vd字段存储图片base64）
func generateSenMLImageData(baseName, imagePath string) ([]SenMLRecord, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}
	imgBase64 := base64.StdEncoding.EncodeToString(data)
	return []SenMLRecord{
		{
			Bn: baseName,
			N:  "image",
			U:  "base64",
			Vd: imgBase64,
		},
	}, nil
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

	// log.Printf("请求URL: %s", url)
	// log.Printf("请求头: %v", req.Header)
	// log.Printf("请求体: %s", string(data))

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

	// 上传类型选择
	uploadType := "senml_image" // "senml" 普通数据，"senml_image" 图片数据
	imagePath := "test.png"     // 图片路径

	if cfg.UploadMode == "once" {
		log.Println("开始单次上传")
		var success bool
		switch uploadType {
		case "senml":
			senmlData := generateSenMLData(cfg.BaseName)
			data, _ := json.Marshal(senmlData)
			success = uploadData(cfg, data, "application/senml+json")
		case "senml_image":
			senmlImgData, err := generateSenMLImageData(cfg.BaseName, imagePath)
			if err != nil {
				log.Printf("图片读取失败: %v", err)
				os.Exit(1)
			}
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
				senmlImgData, err := generateSenMLImageData(cfg.BaseName, imagePath)
				if err != nil {
					log.Printf("图片读取失败: %v", err)
					continue
				}
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
