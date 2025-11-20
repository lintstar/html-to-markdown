# Yuque Plugin for html-to-markdown

语雀编辑器 HTML 转 Markdown 插件

## 功能

- ✅ 转换语雀代码块 (`<card name="codeblock">`) 为 Markdown 代码块
- ✅ 转换语雀图片 (`<card name="image">`) 为 Markdown 图片
- ✅ 保留代码格式（换行、缩进、空格）
- ✅ 支持多种编程语言语法高亮
- ✅ 兼容模式：只处理语雀特有标签

## 快速开始

### 安装

```go
import (
    "github.com/JohannesKaufmann/html-to-markdown/v2/converter"
    "github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
    "github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
    "github.com/JohannesKaufmann/html-to-markdown/v2/plugin/yuque"
)
```

### 使用

```go
conv := converter.NewConverter(
    converter.WithPlugins(
        base.NewBasePlugin(),
        commonmark.NewCommonmarkPlugin(),
        yuque.NewYuquePlugin(), // 添加语雀插件
    ),
)

markdown, err := conv.ConvertString(htmlInput)
```

### 命令行工具

```bash
cd examples/yuque
go run main.go input.html output.md
```

## 示例

### 代码块

**输入：**
```html
<card type="inline" name="codeblock" 
      value="data:%7B%22mode%22%3A%22python%22%2C%22code%22%3A%22def%20hello()%3A%5Cn%20%20print(%5C%22Hello%5C%22)%22%7D">
</card>
```

**输出：**
````markdown
```python
def hello():
  print("Hello")
```
````

### 图片

**输入：**
```html
<card type="inline" name="image"
      value="data:%7B%22src%22%3A%22https%3A//example.com/img.png%22%2C%22title%22%3A%22Logo%22%7D">
</card>
```

**输出：**
```markdown
![Logo](https://example.com/img.png)
```

## 测试

```bash
go test -v
```

## 详细文档

查看项目根目录的 [YUQUE_PLUGIN.md](../../YUQUE_PLUGIN.md) 获取完整文档。

