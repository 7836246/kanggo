package kanggo

import (
	"github.com/7836246/kanggo/constants"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// HTMLTemplateEngine 使用 Go 标准库 html/template 的模板引擎
type HTMLTemplateEngine struct {
	templates *template.Template
	lock      sync.RWMutex
	dir       string
	pattern   string
}

// NewHTMLTemplateEngine 创建一个新的 HTMLTemplateEngine 实例
func NewHTMLTemplateEngine(dir, pattern string) *HTMLTemplateEngine {
	return &HTMLTemplateEngine{dir: dir, pattern: pattern}
}

// Load 加载模板文件
func (e *HTMLTemplateEngine) Load() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	tmpl, err := template.ParseGlob(filepath.Join(e.dir, e.pattern))
	if err != nil {
		return err
	}
	e.templates = tmpl
	return nil
}

// Render 渲染模板
func (e *HTMLTemplateEngine) Render(w http.ResponseWriter, name string, data interface{}) error {
	e.lock.RLock()
	defer e.lock.RUnlock()

	tmpl := e.templates.Lookup(name)
	if tmpl == nil {
		http.Error(w, "模板未找到", http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Type", constants.MIMETextHTML)
	return tmpl.Execute(w, data)
}
