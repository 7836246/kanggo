package kanggo

import (
	"encoding/json"
	"fmt"
	"github.com/7836246/kanggo/version"
	"time"
)

// Config 配置结构体，包含多个配置选项，用户可以根据需要自定义这些选项
type Config struct {
	JSONEncoder          func(v interface{}) ([]byte, error)    // 自定义 JSON 编码器，默认使用标准库的 json.Marshal
	JSONDecoder          func(data []byte, v interface{}) error // 自定义 JSON 解码器，默认使用标准库的 json.Unmarshal
	ShowBanner           bool                                   // 是否在启动时显示欢迎横幅，默认显示
	PrintRoutes          bool                                   // 是否在启动时打印所有已注册的路由信息，默认打印
	ServerHeader         string                                 // 设置服务器响应头的 Server 字段，默认为 "KangGo"
	IdleTimeout          time.Duration                          // 服务器空闲连接的超时时间
	ReadTimeout          time.Duration                          // 服务器读取请求的超时时间
	WriteTimeout         time.Duration                          // 服务器写入响应的超时时间
	MaxRequestBodySize   int                                    // 最大请求体大小，默认为 4 MB
	CaseSensitiveRouting bool                                   // 路由是否区分大小写，默认区分
	StrictRouting        bool                                   // 是否启用严格路由模式，默认不启用
	UnescapePath         bool                                   // 是否对 URL 路径进行解码处理，默认不处理
}

// DefaultConfig 返回默认的配置
// 这是框架提供的默认配置，如果用户不提供自定义配置，则使用此配置
func DefaultConfig() Config {
	return Config{
		JSONEncoder:          json.Marshal,    // 使用标准库的 JSON 编码器
		JSONDecoder:          json.Unmarshal,  // 使用标准库的 JSON 解码器
		ShowBanner:           true,            // 启动时显示欢迎横幅
		PrintRoutes:          true,            // 启动时打印路由信息
		ServerHeader:         "KangGo",        // 设置默认的服务器响应头
		IdleTimeout:          0,               // 默认不设置空闲超时
		ReadTimeout:          0,               // 默认不设置读取超时
		WriteTimeout:         0,               // 默认不设置写入超时
		MaxRequestBodySize:   4 * 1024 * 1024, // 最大请求体大小为 4 MB
		CaseSensitiveRouting: true,            // 路由区分大小写
		StrictRouting:        false,           // 不启用严格路由模式
		UnescapePath:         false,           // 不对 URL 路径进行解码处理
	}
}

// PrintWelcomeBanner 显示欢迎横幅和版本信息
func PrintWelcomeBanner() {
	fmt.Printf(`
 __  __                       _______        
|  |/  |.---.-..-----..-----.|     __|.-----.
|     < |  _  ||     ||  _  ||    |  ||  _  |
|__|\__||___._||__|__||___  ||_______||_____|
                      |_____|
欢迎使用 KangGo - 一个简洁高效的 Go Web 框架，专为快速开发与高性能设计。当前版本：%s
`, version.Version)
}
