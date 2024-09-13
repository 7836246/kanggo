package kanggo

import (
	"fmt"
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
	Next           func(Context) bool                       // 定义一个函数，当返回 true 时跳过此中间件，默认值为 nil
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
	// 使用默认的 Static 配置，如果用户提供了自定义配置，则使用用户的配置
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

	// 处理静态文件请求的处理函数
	handler := func(ctx Context) error {
		// 如果 Next 函数返回 true，跳过静态文件处理
		if cfg.Next != nil && cfg.Next(ctx) {
			return nil
		}

		// 获取请求路径并去掉前缀部分
		relativePath := strings.TrimPrefix(ctx.Request.URL.Path, prefix)
		relativePath = strings.TrimPrefix(relativePath, "/") // 去掉可能的前导斜杠
		filePath := filepath.Join(root, relativePath)
		// 调试信息
		// fmt.Printf("请求的 URL 路径: %s\n", ctx.Request.URL.Path)
		// fmt.Printf("去掉前缀后的相对路径: %s\n", relativePath)
		// fmt.Printf("文件系统中的文件路径: %s\n", filePath)

		// 检查文件或目录是否存在
		info, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			// fmt.Println("文件未找到: ", filePath)  // 添加调试信息
			http.NotFound(ctx.Writer, ctx.Request) // 文件不存在，返回 404
			return nil
		} else if err != nil {
			// fmt.Println("无法获取文件信息: ", err)  // 添加调试信息
			http.Error(ctx.Writer, "500 服务器内部错误", http.StatusInternalServerError) // 服务器内部错误，返回 500
			return nil
		}

		// 如果请求的是一个目录，则处理索引文件
		if info.IsDir() {
			indexFile := filepath.Join(filePath, cfg.Index)
			// fmt.Printf("请求的是一个目录，尝试提供索引文件: %s\n", indexFile)
			if _, err := os.Stat(indexFile); err == nil {
				http.ServeFile(ctx.Writer, ctx.Request, indexFile)
				return nil
			}
			if cfg.Browse {
				// 提供目录浏览（此功能可扩展实现）
				return browseDirectory(ctx.Writer, filePath)
			}
			http.Error(ctx.Writer, "403 禁止访问目录", http.StatusForbidden) // 禁止访问目录，返回 403
			return nil
		}

		// 设置缓存控制头
		if cfg.MaxAge > 0 {
			ctx.Writer.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(cfg.MaxAge))
		}

		// 启用字节范围请求
		if cfg.ByteRange {
			ctx.Writer.Header().Set("Accept-Ranges", "bytes")
		}

		// 启用文件下载选项时设置 Content-Disposition 头
		if cfg.Download {
			ctx.Writer.Header().Set("Content-Disposition", "attachment")
		}

		// 调用自定义响应修改函数（如果已定义）
		if cfg.ModifyResponse != nil {
			cfg.ModifyResponse(ctx.Writer, ctx.Request)
		}

		// 使用 http.ServeFile 提供文件
		http.ServeFile(ctx.Writer, ctx.Request, filePath)
		return nil
	}

	// 将 handler 函数转换为 HandlerFunc 类型的闭包
	wrappedHandler := func(ctx Context) error {
		return handler(ctx)
	}

	// 注册静态文件服务的路由
	k.router.Handle("GET", prefix, wrappedHandler)

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
