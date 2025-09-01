package main

import "liar/upload"

func main() {
	cfg := upload.UploadConfig{
		BaseURL:      "https://localhost",
		DomainID:     "562d704a-c442-499a-aff3-223f580bf6b3",
		ChannelID:    "b0ec13df-9ff0-48b9-9cb6-b3be072e7c99",
		Subtopic:     "temperature",
		ClientSecret: "9c3acbd2-ce7c-4ba3-b6ab-8792e310002c",
		CACertPath:   "../CA/ca.crt",
		BaseName:     "ljp",
		Timeout:      10,
	}
	uploadType := "senml_image" // 或 "senml"
	imagePath := "test.jpg"
	imageUploadURL := "http://localhost:18080/upload"

	upload.Upload(cfg, uploadType, imagePath, imageUploadURL)
}
