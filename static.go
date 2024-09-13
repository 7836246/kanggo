package kanggo

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Static 配置结构体，定义静态文件服务的选项
type Static struct {
	Compress       bool                                     // 是否启用压缩，减少传输体积，默认值 false
	ByteRange      bool                                     // 是否支持字节范围请求，默认值 false
	Browse         bool                                     // 是否启用目录浏览，允许用户查看文件夹中的内容，默认值 false
	Download       bool                                     // 是否启用文件下载，启用后所有文件将以附件形式下载，默认值 false
	Index          string                                   // 用于提供目录的索引文件的名称，例如 "index.html"，默认值为 "index.html"
	CacheDuration  time.Duration                            // 非活动文件处理程序的缓存持续时间，使用负值禁用此选项，默认值 10 秒
	MaxAge         int                                      // 设置文件响应的 Cache-Control HTTP 头的值，MaxAge 以秒为单位，默认值 0
	ModifyResponse func(http.ResponseWriter, *http.Request) // 允许修改响应的自定义函数，用户可以在这里添加自定义的 HTTP 头等，默认值为 nil
	Next           func(Context) bool                       // 定义一个函数，当返回 true 时跳过此中间件，用于条件性地应用此中间件，默认值为 nil
}

// NewStatic 返回一个带有默认值的 Static 配置实例
func NewStatic() Static {
	return Static{
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
// prefix 是 URL 前缀，root 是文件系统中的根目录，config 是静态文件服务的配置
func (k *KangGo) Static(prefix, root string, config ...Static) *KangGo {
	// 使用默认的 Static 配置，如果用户提供了自定义配置，则使用用户的配置
	cfg := NewStatic()
	if len(config) > 0 {
		cfg = config[0]
	}

	// 去除前缀中的尾部斜杠，确保前缀统一
	prefix = strings.TrimSuffix(prefix, "/")

	// 创建文件服务器，使用 http.StripPrefix 处理 URL 前缀
	fileServer := http.StripPrefix(prefix, http.FileServer(http.Dir(root)))

	// 处理静态文件请求的处理函数
	handler := func(ctx Context) error {
		// 如果 Next 函数返回 true，跳过静态文件处理
		if cfg.Next != nil && cfg.Next(ctx) {
			return nil
		}

		// 获取请求的文件路径，去掉 URL 中的前缀部分
		file := strings.TrimPrefix(ctx.Request.URL.Path, prefix)
		file = filepath.Join(root, file)

		// 检查文件或目录的状态
		info, err := os.Stat(file)
		if os.IsNotExist(err) {
			return ctx.SendString("404 文件未找到") // 文件不存在，返回 404
		} else if err != nil {
			return ctx.SendString("500 服务器内部错误") // 服务器内部错误，返回 500
		}

		// 处理目录请求，如果请求的是一个目录而不是文件
		if info.IsDir() {
			indexFile := filepath.Join(file, cfg.Index) // 试图提供目录下的索引文件
			if _, err := os.Stat(indexFile); err == nil {
				// 如果索引文件存在，提供该文件
				http.ServeFile(ctx.Writer, ctx.Request, indexFile)
				return nil
			}
			if cfg.Browse {
				// 如果启用了目录浏览，提供目录内容列表（此处可以实现具体的目录浏览逻辑）
				return ctx.SendString("目录浏览功能尚未实现")
			}
			return ctx.SendString("403 禁止访问目录") // 禁止访问目录，返回 403
		}

		// 设置缓存控制头
		if cfg.MaxAge > 0 {
			ctx.Writer.Header().Set("Cache-Control", "max-age="+strconv.Itoa(cfg.MaxAge))
		}

		// 处理字节范围请求
		if cfg.ByteRange {
			ctx.Writer.Header().Set("Accept-Ranges", "bytes")
		}

		// 调用自定义响应修改函数（如果已定义）
		if cfg.ModifyResponse != nil {
			cfg.ModifyResponse(ctx.Writer, ctx.Request)
		}

		// 使用 fileServer 提供文件内容给客户端
		fileServer.ServeHTTP(ctx.Writer, ctx.Request)
		return nil
	}

	// 将 handler 函数转换为 HandlerFunc 类型的闭包
	wrappedHandler := func(ctx Context) error {
		return handler(ctx)
	}

	// 注册静态文件服务的路由，使用通配符来匹配所有子路径
	k.router.Handle("GET", prefix+"/*", wrappedHandler)

	// 将静态文件路由信息记录下来，用于打印
	k.router.RegisterStaticRoute(prefix + "/*")

	return k
}
