package yuque

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/JohannesKaufmann/dom"
	"github.com/lintstar/html-to-markdown/v2/converter"
	"github.com/lintstar/html-to-markdown/v2/marker"
	"golang.org/x/net/html"
)

type yuquePlugin struct{}

// NewYuquePlugin 创建一个新的语雀插件，用于处理语雀编辑器的特殊标签
func NewYuquePlugin() converter.Plugin {
	return &yuquePlugin{}
}

func (y *yuquePlugin) Name() string {
	return "yuque"
}

func (y *yuquePlugin) Init(conv *converter.Converter) error {
	// 注册 card 标签为内联标签
	conv.Register.RendererFor("card", converter.TagTypeInline, y.renderCard, converter.PriorityStandard)
	return nil
}

// CodeBlockData 代表语雀代码块的数据结构
type CodeBlockData struct {
	Mode string `json:"mode"` // 代码语言
	Code string `json:"code"` // 代码内容
}

// ImageData 代表语雀图片的数据结构
type ImageData struct {
	Src   string `json:"src"`   // 图片URL
	Title string `json:"title"` // 图片标题
}

func (y *yuquePlugin) renderCard(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
	// 获取 card 的 name 属性
	cardName := dom.GetAttributeOr(n, "name", "")
	if cardName == "" {
		return converter.RenderTryNext
	}

	// 获取 value 属性
	valueAttr := dom.GetAttributeOr(n, "value", "")
	if valueAttr == "" {
		return converter.RenderTryNext
	}

	// 检查是否是语雀的特殊格式（以 "data:" 开头）
	if !strings.HasPrefix(valueAttr, "data:") {
		return converter.RenderTryNext
	}

	// 移除 "data:" 前缀
	encodedData := strings.TrimPrefix(valueAttr, "data:")

	// URL 解码
	decodedData, err := url.QueryUnescape(encodedData)
	if err != nil {
		// 解码失败，跳过处理
		return converter.RenderTryNext
	}

	// 根据 card 的 name 属性处理不同类型
	switch cardName {
	case "codeblock":
		return y.renderCodeBlock(ctx, w, decodedData)
	case "image":
		return y.renderImage(ctx, w, decodedData)
	default:
		// 未知的 card 类型，跳过
		return converter.RenderTryNext
	}
}

func (y *yuquePlugin) renderCodeBlock(ctx converter.Context, w converter.Writer, jsonData string) converter.RenderStatus {
	var data CodeBlockData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return converter.RenderTryNext
	}

	// 获取代码内容和语言
	code := data.Code
	mode := data.Mode

	// 如果代码末尾有换行符，移除它
	code = strings.TrimSuffix(code, "\n")

	// 计算需要的代码围栏
	fenceChar := '`'
	maxCount := 0
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		count := 0
		for _, char := range line {
			if char == fenceChar {
				count++
			} else {
				if count > maxCount {
					maxCount = count
				}
				count = 0
			}
		}
		if count > maxCount {
			maxCount = count
		}
	}

	// 围栏至少要比代码中的连续反引号多一个
	fenceCount := maxCount + 1
	if fenceCount < 3 {
		fenceCount = 3
	}
	fence := strings.Repeat(string(fenceChar), fenceCount)

	// 将换行符替换为特殊标记，以保留代码块中的原始格式
	codeBytes := []byte(code)
	codeBytes = []byte(strings.ReplaceAll(string(codeBytes), "\n", string(marker.BytesMarkerCodeBlockNewline)))

	// 写入 markdown 代码块
	w.WriteString("\n\n")
	w.WriteString(fence)
	if mode != "" {
		w.WriteString(mode)
	}
	w.WriteRune('\n')
	w.Write(codeBytes)
	w.WriteRune('\n')
	w.WriteString(fence)
	w.WriteString("\n\n")

	return converter.RenderSuccess
}

func (y *yuquePlugin) renderImage(ctx converter.Context, w converter.Writer, jsonData string) converter.RenderStatus {
	var data ImageData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return converter.RenderTryNext
	}

	// 获取图片URL
	src := strings.TrimSpace(data.Src)
	if src == "" {
		return converter.RenderTryNext
	}

	// 处理绝对URL
	src = ctx.AssembleAbsoluteURL(ctx, "img", src)

	// 获取标题（alt）
	title := strings.ReplaceAll(data.Title, "\n", " ")
	alt := title // 使用 title 作为 alt 文本

	// 转义 alt 中的方括号
	alt = escapeAlt(alt)

	// 写入 markdown 图片语法
	w.WriteRune('!')
	w.WriteRune('[')
	w.WriteString(alt)
	w.WriteRune(']')
	w.WriteRune('(')
	w.WriteString(src)
	w.WriteRune(')')

	return converter.RenderSuccess
}

// escapeAlt 转义 alt 文本中的方括号
func escapeAlt(altString string) string {
	var result strings.Builder
	for i, char := range altString {
		if char == '[' || char == ']' {
			// 检查前一个字符是否已经是反斜杠
			if i > 0 && altString[i-1] != '\\' {
				result.WriteRune('\\')
			}
		}
		result.WriteRune(char)
	}
	return result.String()
}
