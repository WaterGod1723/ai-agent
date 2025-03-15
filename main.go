package main

import (
	ai_model "ai-agent/ai-model"
	"fmt"
	"os"
)

func main() {
	d, err := os.ReadFile("model_api.key")
	if err != nil {
		fmt.Println(err)
		return
	}
	// API配置
	apiKey := string(d)      // 替换为你的API密钥
	imagePath := "image.png" // 图片路径

	// 调用封装的函数
	reply, err := ai_model.ProcessImage(imagePath, apiKey, "返回人物头像的坐标，返回json格式，不要解释，长度单位使用百分比，例如 [{name:\"光头男士\",x: 0.3, y: 0.4, width: 0.1, height: 0.1}]")
	if err != nil {
		fmt.Println("处理图像失败:", err)
		return
	}

	fmt.Println("模型回复:", reply)
}
