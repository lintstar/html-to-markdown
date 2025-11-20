package yuque

import (
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
)

func TestYuqueCodeBlock(t *testing.T) {
	input := `<card type="inline" name="codeblock" value="data:%7B%22mode%22%3A%22python%22%2C%22code%22%3A%22def%20hello()%3A%5Cn%20%20%20%20print(%5C%22Hello%5C%22)%22%7D"></card>`

	expected := "```python\ndef hello():\n    print(\"Hello\")\n```"

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			NewYuquePlugin(),
		),
	)

	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	if markdown != expected {
		t.Errorf("期望:\n%s\n\n实际:\n%s", expected, markdown)
	}
}

func TestYuqueImage(t *testing.T) {
	input := `<card type="inline" name="image" value="data:%7B%22src%22%3A%22https%3A//example.com/image.png%22%2C%22title%22%3A%22Test%20Image%22%7D"></card>`

	expected := "![Test Image](https://example.com/image.png)"

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			NewYuquePlugin(),
		),
	)

	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	if markdown != expected {
		t.Errorf("期望:\n%s\n\n实际:\n%s", expected, markdown)
	}
}

func TestYuqueCodeBlockJSON(t *testing.T) {
	input := `<card type="inline" name="codeblock" value="data:%7B%22mode%22%3A%22json%22%2C%22code%22%3A%22%7B%5Cn%20%20%5C%22name%5C%22%3A%20%5C%22test%5C%22%5Cn%7D%22%7D"></card>`

	expected := "```json\n{\n  \"name\": \"test\"\n}\n```"

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			NewYuquePlugin(),
		),
	)

	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	if markdown != expected {
		t.Errorf("期望:\n%s\n\n实际:\n%s", expected, markdown)
	}
}

func TestNonYuqueCard(t *testing.T) {
	// 测试非语雀的 card 标签应该被跳过
	input := `<card type="inline" name="unknown" value="test"></card>`

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			NewYuquePlugin(),
		),
	)

	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	// 非语雀格式的 card 应该被忽略，生成空白内容
	if markdown != "" {
		t.Errorf("期望空字符串，实际: %s", markdown)
	}
}

func TestMixedContent(t *testing.T) {
	input := `
	<h1>标题</h1>
	<p>这是一段文本</p>
	<card type="inline" name="codeblock" value="data:%7B%22mode%22%3A%22go%22%2C%22code%22%3A%22package%20main%5Cn%5Cnfunc%20main()%20%7B%7D%22%7D"></card>
	<p>更多文本</p>
	<card type="inline" name="image" value="data:%7B%22src%22%3A%22https%3A//example.com/img.png%22%2C%22title%22%3A%22%22%7D"></card>
	`

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			NewYuquePlugin(),
		),
	)

	markdown, err := conv.ConvertString(input)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	// 验证输出包含所有预期元素
	expected := []string{
		"# 标题",
		"这是一段文本",
		"```go",
		"package main",
		"func main() {}",
		"```",
		"更多文本",
		"![](https://example.com/img.png)",
	}

	for _, exp := range expected {
		if !contains(markdown, exp) {
			t.Errorf("输出不包含期望的内容: %s\n完整输出:\n%s", exp, markdown)
		}
	}
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(substr) == 0 || indexOfSubstring(str, substr) >= 0)
}

func indexOfSubstring(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
