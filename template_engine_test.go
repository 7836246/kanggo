package kanggo

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHTMLTemplateEngine(t *testing.T) {
	// 临时目录和文件名
	tempDir := "templates"
	tempFile := "templates/test.html"

	// 创建临时目录
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		t.Fatalf("创建目录失败: %v", err)
	}

	// 创建临时模板文件
	htmlContent := `<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    {{.Body}}
</body>
</html>`
	if err := os.WriteFile(tempFile, []byte(htmlContent), 0644); err != nil {
		t.Fatalf("创建 HTML 模板文件失败: %v", err)
	}

	// 测试完成后清理生成的文件和目录
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("清理临时文件失败: %v", err)
		}
	}()

	// 创建一个 HTML 模板引擎实例
	engine := NewHTMLTemplateEngine(tempDir, "*.html")

	// 加载模板
	if err := engine.Load(); err != nil {
		t.Fatalf("加载 HTML 模板失败: %v", err)
	}

	// 创建一个 httptest.ResponseRecorder 用于模拟 http.ResponseWriter
	recorder := httptest.NewRecorder()

	// 定义模板数据
	data := map[string]interface{}{
		"Title": "Hello",
		"Body":  "世界",
	}

	// 渲染模板
	if err := engine.Render(recorder, "test.html", data); err != nil {
		t.Fatalf("渲染 HTML 模板失败: %v", err)
	}

	// 检查渲染结果是否符合预期
	expected := `<html>
<head>
    <title>Hello</title>
</head>
<body>
    世界
</body>
</html>`

	// 比较实际结果和期望结果，去除首尾空格后比较
	if strings.TrimSpace(recorder.Body.String()) != strings.TrimSpace(expected) {
		t.Errorf("渲染结果不正确: 得到 %s，期望 %s", recorder.Body.String(), expected)
	}
}
