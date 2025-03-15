package ai_model

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// API请求结构体
type GLM4VRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type     string   `json:"type"` // "text" 或 "image_url"
	Text     string   `json:"text,omitempty"`
	ImageURL ImageURL `json:"image_url,omitempty"`
}

// 定义带有JSON标签的ImageURL结构体
type ImageURL struct {
	URL string `json:"url"`
}

// API响应结构体
type GLM4VResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// 封装图像处理函数
func ProcessImage(imagePath, apiKey, prompt string) (string, error) {
	// 读取并编码图片为base64
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("读取图片失败: %w", err)
	}
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// 构建请求体
	requestBody := GLM4VRequest{
		Model: "glm-4v-flash",
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: prompt,
					},
					{
						Type: "image_url",
						ImageURL: ImageURL{
							URL: fmt.Sprintf("data:image/png;base64,%s", base64Image), // 修改为image/png
						},
					},
				},
			},
		},
		MaxTokens: 1000,
	}

	// 序列化请求体
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 处理响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var response GLM4VResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API返回错误[%d]: %s", resp.StatusCode, response.Error.Message)
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("未收到有效回复")
}
