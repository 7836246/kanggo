package kanggo

import (
	"fmt"
	"github.com/7836246/kanggo/constants"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// StaticConfig 配置结构体，定义静态文件服务的选项
type StaticConfig struct {
	Compress       bool                                     // 是否启用压缩，减少传输体积，默认值 false
	ByteRange      bool                                     // 是否支持字节范围请求，默认值 false
	Browse         bool                                     // 是否启用目录浏览，允许用户查看文件夹中的内容，默认值 false
	Download       bool                                     // 是否启用文件下载，启用后所有文件将以附件形式下载，默认值 false
	Index          string                                   // 用于提供目录的索引文件的名称，例如 "index.html"，默认值为 "index.html"
	CacheDuration  time.Duration                            // 非活动文件处理程序的缓存持续时间，使用负值禁用此选项，默认值 10 秒
	MaxAge         int                                      // 设置文件响应的 Cache-Control HTTP 头的值，MaxAge 以秒为单位，默认值 0
	ModifyResponse func(http.ResponseWriter, *http.Request) // 自定义函数，允许修改响应，默认值为 nil
	Next           func(*Context) bool                      // 定义一个函数，当返回 true 时跳过此中间件，默认值为 nil
}

// NewStaticConfig 返回一个带有默认值的 StaticConfig 配置实例
func NewStaticConfig() StaticConfig {
	return StaticConfig{
		Compress:       false,
		ByteRange:      false,
		Browse:         false,
		Download:       false,
		Index:          "index.html",
		CacheDuration:  10 * time.Second,
		MaxAge:         0,
		ModifyResponse: nil,
		Next:           nil,
	}
}

// Static 注册一个静态文件服务路由
func (k *KangGo) Static(prefix, root string, config ...StaticConfig) *KangGo {
	cfg := NewStaticConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// 确保前缀总是以 '/' 开头
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	// 去除前缀中的尾部斜杠，确保前缀统一
	prefix = strings.TrimSuffix(prefix, "/")

	handler := func(ctx *Context) error {
		if cfg.Next != nil && cfg.Next(ctx) {
			return nil
		}

		// 获取相对路径
		relativePath := strings.TrimPrefix(ctx.Request.URL.Path, prefix)
		relativePath = strings.TrimPrefix(relativePath, "/")
		filePath := filepath.Join(root, relativePath)

		// 检查文件或目录是否存在
		info, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			http.NotFound(ctx.Writer, ctx.Request)
			return nil
		} else if err != nil {
			http.Error(ctx.Writer, "500 服务器内部错误", http.StatusInternalServerError)
			return nil
		}

		// 处理目录请求
		if info.IsDir() {
			indexFile := filepath.Join(filePath, cfg.Index)
			if _, err := os.Stat(indexFile); err == nil {
				http.ServeFile(ctx.Writer, ctx.Request, indexFile)
				return nil
			}
			if cfg.Browse {
				return browseDirectory(ctx.Writer, filePath)
			}
			http.Error(ctx.Writer, "403 禁止访问目录", http.StatusForbidden)
			return nil
		}

		// 处理文件请求
		if cfg.MaxAge > 0 {
			ctx.Writer.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(cfg.MaxAge))
		}
		if cfg.ByteRange {
			ctx.Writer.Header().Set("Accept-Ranges", "bytes")
		}
		if cfg.Download {
			ctx.Writer.Header().Set("Content-Disposition", "attachment")
		}
		if cfg.ModifyResponse != nil {
			cfg.ModifyResponse(ctx.Writer, ctx.Request)
		}

		http.ServeFile(ctx.Writer, ctx.Request, filePath)
		return nil
	}

	k.Router.Handle(constants.MethodGet, prefix+"/*", handler)
	return k
}

// browseDirectory 提供简单的目录浏览功能
func browseDirectory(w http.ResponseWriter, dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, "500 服务器内部错误", http.StatusInternalServerError)
		return err
	}

	// 构建简单的目录列表 HTML
	fmt.Fprintf(w, "<html><body><h1>目录浏览</h1><ul>")
	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			name += "/"
		}
		fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", name, name)
	}
	fmt.Fprint(w, "</ul></body></html>")
	return nil
}
