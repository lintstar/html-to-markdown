package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lintstar/html-to-markdown/v2/converter"
	"github.com/lintstar/html-to-markdown/v2/plugin/base"
	"github.com/lintstar/html-to-markdown/v2/plugin/commonmark"
	"github.com/lintstar/html-to-markdown/v2/plugin/yuque"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("用法: go run main.go <输入HTML文件> [输出MD文件]")
		fmt.Println("示例: go run main.go test.html output.md")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := "output.md"
	if len(os.Args) >= 3 {
		outputFile = os.Args[2]
	}

	// 读取输入HTML文件
	htmlContent, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	// 创建转换器，添加语雀插件
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			yuque.NewYuquePlugin(), // 添加语雀插件
		),
	)

	// 执行转换
	markdown, err := conv.ConvertString(string(htmlContent))
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}

	// 写入输出文件
	err = os.WriteFile(outputFile, []byte(markdown), 0644)
	if err != nil {
		log.Fatalf("写入输出文件失败: %v", err)
	}

	fmt.Printf("转换成功！\n")
	fmt.Printf("输入文件: %s\n", inputFile)
	fmt.Printf("输出文件: %s\n", outputFile)
	fmt.Printf("生成的 Markdown 内容:\n\n%s\n", markdown)
}
