package kanggo

import "net/http"

// TemplateEngine 定义模板引擎接口
type TemplateEngine interface {
	Load() error                                                       // 加载模板文件
	Render(w http.ResponseWriter, name string, data interface{}) error // 渲染模板
}
