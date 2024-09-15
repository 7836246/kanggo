package kanggo

import (
	"github.com/7836246/kanggo/constants"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// TextTemplateEngine 使用 Go 标准库 text/template 的模板引擎
type TextTemplateEngine struct {
	templates *template.Template
	lock      sync.RWMutex
	dir       string
	pattern   string
}

// NewTextTemplateEngine 创建一个新的 TextTemplateEngine 实例
func NewTextTemplateEngine(dir, pattern string) *TextTemplateEngine {
	return &TextTemplateEngine{dir: dir, pattern: pattern}
}

// Load 加载模板文件
func (e *TextTemplateEngine) Load() error {
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
func (e *TextTemplateEngine) Render(w http.ResponseWriter, name string, data interface{}) error {
	e.lock.RLock()
	defer e.lock.RUnlock()

	tmpl := e.templates.Lookup(name)
	if tmpl == nil {
		http.Error(w, "模板未找到", http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Type", constants.MIMETextPlain)
	return tmpl.Execute(w, data)
}
