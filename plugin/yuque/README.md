# Yuque Plugin for html-to-markdown

语雀编辑器 HTML 转 Markdown 插件

## 功能

- ✅ 转换语雀代码块 (`<card name="codeblock">`) 为 Markdown 代码块
- ✅ 转换语雀图片 (`<card name="image">`) 为 Markdown 图片
- ✅ 保留代码格式（换行、缩进、空格）
- ✅ 支持多种编程语言语法高亮
- ✅ 兼容模式：只处理语雀特有标签

## 安装

```bash
go get github.com/lintstar/html-to-markdown/v2@main
```

## 快速开始

### 在代码中使用

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/lintstar/html-to-markdown/v2/converter"
    "github.com/lintstar/html-to-markdown/v2/plugin/base"
    "github.com/lintstar/html-to-markdown/v2/plugin/commonmark"
    "github.com/lintstar/html-to-markdown/v2/plugin/yuque"
)

func main() {
    html := `<card type="inline" name="codeblock" 
             value="data:%7B%22mode%22%3A%22python%22%2C%22code%22%3A%22print(%5C%22Hello%5C%22)%22%7D">
             </card>`
    
    conv := converter.NewConverter(
        converter.WithPlugins(
            base.NewBasePlugin(),
            commonmark.NewCommonmarkPlugin(),
            yuque.NewYuquePlugin(),
        ),
    )
    
    markdown, err := conv.ConvertString(html)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(markdown)
}
```

### 命令行工具

```bash
cd examples/yuque
go run main.go input.html output.md
```

## 转换示例

### 代码块

**语雀 HTML：**
```html
<card type="inline" name="codeblock" 
      value="data:%7B%22mode%22%3A%22python%22%2C%22code%22%3A%22def%20hello()%3A%5Cn%20%20print(%5C%22Hello%5C%22)%22%7D">
</card>
```

**Markdown 输出：**
````markdown
```python
def hello():
  print("Hello")
```
````

### 图片

**语雀 HTML：**
```html
<card type="inline" name="image"
      value="data:%7B%22src%22%3A%22https%3A//example.com/img.png%22%2C%22title%22%3A%22Logo%22%7D">
</card>
```

**Markdown 输出：**
```markdown
![Logo](https://example.com/img.png)
```

## 技术实现

### 工作原理

1. 检测 `<card>` 标签并识别 `name` 属性（`codeblock` 或 `image`）
2. URL 解码 `value` 属性（从 `data:%7B...` 解码为 JSON）
3. 解析 JSON 提取关键字段：
   - 代码块：`mode`（语言）和 `code`（代码内容）
   - 图片：`src`（URL）和 `title`（标题）
4. 生成标准 Markdown 格式

### 兼容性

- ✅ 只处理语雀特有的 `<card>` 标签
- ✅ 不影响其他 HTML 标签的转换
- ✅ 错误安全处理，解析失败会跳过

## 测试

```bash
# 运行单元测试
cd plugin/yuque
go test -v

# 测试转换
cd examples/yuque
go run main.go complex_test.html output.md
```

## 项目结构

```
plugin/yuque/
├── yuque.go          # 插件实现
└── yuque_test.go     # 单元测试

examples/yuque/
├── main.go           # 命令行工具
└── complex_test.html # 测试文件
```

## 常见问题

### Q: 如何在我的项目中使用？

在你的项目中：

```bash
go get github.com/lintstar/html-to-markdown/v2@main
```

在代码中导入：

```go
import "github.com/lintstar/html-to-markdown/v2/plugin/yuque"
```

### Q: 支持哪些语雀元素？

目前支持：
- ✅ 代码块（所有编程语言）
- ✅ 图片（带标题和不带标题）

### Q: 代码格式会丢失吗？

不会！插件完整保留：
- 换行符
- 缩进（空格和制表符）
- 空格

### Q: 如何扩展支持更多 card 类型？

编辑 `yuque.go`，在 `renderCard` 函数中添加新的 case：

```go
switch cardName {
case "codeblock":
    return y.renderCodeBlock(ctx, w, decodedData)
case "image":
    return y.renderImage(ctx, w, decodedData)
case "file":  // 新增
    return y.renderFile(ctx, w, decodedData)
}
```

## 相关链接

- [GitHub 仓库](https://github.com/lintstar/html-to-markdown)
- [原始项目](https://github.com/JohannesKaufmann/html-to-markdown)

