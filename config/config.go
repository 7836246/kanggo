package config

import (
	"encoding/json"
	"fmt"
	"github.com/7836246/kanggo/version"
)

// Config 配置结构体，允许用户自定义 JSON 编解码器等
type Config struct {
	JSONEncoder func(v interface{}) ([]byte, error)
	JSONDecoder func(data []byte, v interface{}) error
	ShowBanner  bool // 新增字段，用于控制是否显示横幅
}

// DefaultConfig 返回默认的配置
func DefaultConfig() Config {
	return Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ShowBanner:  true, // 默认显示横幅
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
简洁高效的 Go Web 框架，专为快速开发与高性能设计。当前版本：%s
`, version.Version)
}
