package KangGo

import (
	"kanggo/kanggo/config"
	"net/http"
)

// Context 代表 HTTP 请求的上下文
type Context struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	Params      map[string]string
	jsonEncoder func(v interface{}) ([]byte, error)
	jsonDecoder func(data []byte, v interface{}) error
}

// NewContext 创建一个新的 Context 实例
func NewContext(w http.ResponseWriter, req *http.Request, cfg config.Config) Context {
	return Context{
		Writer:      w,
		Request:     req,
		Params:      make(map[string]string),
		jsonEncoder: cfg.JSONEncoder,
		jsonDecoder: cfg.JSONDecoder,
	}
}

// Param 获取路径参数
func (c Context) Param(key string) string {
	return c.Params[key]
}

// JSON 返回一个 JSON 响应
func (c Context) JSON(code int, obj interface{}) error {
	data, err := c.jsonEncoder(obj)
	if err != nil {
		return err
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	_, err = c.Writer.Write(data)
	return err
}

// SendString 返回一个纯文本响应
func (c Context) SendString(msg string) error {
	_, err := c.Writer.Write([]byte(msg))
	return err
}
