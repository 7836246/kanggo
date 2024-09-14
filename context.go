package kanggo

import (
	"encoding/json"
	"github.com/7836246/kanggo/constants"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
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
func NewContext(w http.ResponseWriter, req *http.Request, cfg Config) *Context {
	return &Context{
		Writer:      w,
		Request:     req,
		Params:      make(map[string]string),
		jsonEncoder: cfg.JSONEncoder,
		jsonDecoder: cfg.JSONDecoder,
	}
}

// Param 获取路径参数
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// Query 获取 URL 查询参数
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// DefaultQuery 获取 URL 查询参数并提供默认值
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value := c.Query(key); value != "" {
		return value
	}
	return defaultValue
}

// FormValue 获取 POST 表单数据
func (c *Context) FormValue(key string) string {
	return c.Request.FormValue(key)
}

// DefaultFormValue 获取 POST 表单数据并提供默认值
func (c *Context) DefaultFormValue(key, defaultValue string) string {
	if value := c.FormValue(key); value != "" {
		return value
	}
	return defaultValue
}

// BindJSON 解析 JSON 请求体到指定的对象
func (c *Context) BindJSON(obj interface{}) error {
	if c.jsonDecoder == nil {
		return json.NewDecoder(c.Request.Body).Decode(obj)
	}
	decoder := json.NewDecoder(c.Request.Body)
	return decoder.Decode(obj)
}

// BindForm 解析表单数据到指定的结构体
func (c *Context) BindForm(obj interface{}) error {
	if err := c.Request.ParseForm(); err != nil {
		return err
	}

	// 通过反射解析表单数据到结构体
	objValue := reflect.ValueOf(obj).Elem() // 获取结构体指针的反射值
	objType := objValue.Type()              // 获取结构体的类型

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)                  // 获取结构体的字段
		fieldValue := objValue.Field(i)            // 获取字段的值
		formKey := field.Tag.Get("form")           // 获取 `form` 标签
		if formKey == "" || !fieldValue.CanSet() { // 如果没有设置 `form` 标签，或字段不可设置，跳过
			continue
		}
		formValue := c.Request.FormValue(formKey) // 从请求中获取字段的值

		// 根据字段类型设置结构体字段的值
		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(formValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, err := strconv.ParseInt(formValue, 10, 64); err == nil {
				fieldValue.SetInt(v)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v, err := strconv.ParseUint(formValue, 10, 64); err == nil {
				fieldValue.SetUint(v)
			}
		case reflect.Float32, reflect.Float64:
			if v, err := strconv.ParseFloat(formValue, 64); err == nil {
				fieldValue.SetFloat(v)
			}
		case reflect.Bool:
			if v, err := strconv.ParseBool(formValue); err == nil {
				fieldValue.SetBool(v)
			}
		default:
			// 其他类型暂不处理
		}
	}

	return nil
}

// JSON 返回一个 JSON 响应
func (c *Context) JSON(code int, obj interface{}) error {
	data, err := c.jsonEncoder(obj)
	if err != nil {
		return err
	}
	c.Writer.Header().Set(constants.HeaderContentType, constants.MIMEApplicationJSON) // 使用常量
	c.Writer.WriteHeader(code)
	_, err = c.Writer.Write(data)
	return err
}

// JSONP 返回一个 JSONP 响应
func (c *Context) JSONP(callback string, obj interface{}) error {
	data, err := c.jsonEncoder(obj)
	if err != nil {
		return err
	}
	c.Writer.Header().Set(constants.HeaderContentType, constants.MIMEApplicationJavaScript) // 使用常量
	_, err = c.Writer.Write([]byte(callback + "(" + string(data) + ");"))
	return err
}

// SendString 返回一个纯文本响应
func (c *Context) SendString(msg string) error {
	c.Writer.Header().Set(constants.HeaderContentType, constants.MIMETextPlain) // 使用常量
	_, err := c.Writer.Write([]byte(msg))
	return err
}

// SendHTML 返回一个 HTML 响应
func (c *Context) SendHTML(html string) error {
	c.Writer.Header().Set(constants.HeaderContentType, constants.MIMETextHTML) // 使用常量
	_, err := c.Writer.Write([]byte(html))
	return err
}

// SendError 返回一个错误响应
func (c *Context) SendError(code int, message string) error {
	c.Writer.Header().Set(constants.HeaderContentType, constants.MIMETextPlain) // 使用常量
	c.Writer.WriteHeader(code)
	_, err := c.Writer.Write([]byte(message))
	return err
}

// SendFile 发送文件作为响应
func (c *Context) SendFile(filepath string, download bool) error {
	// 检查文件是否存在
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(c.Writer, "文件未找到", http.StatusNotFound)
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(c.Writer, "服务器错误", http.StatusInternalServerError)
		return err
	}

	// 设置响应头的 Content-Type
	mimeType := http.DetectContentType(make([]byte, 512))
	c.Writer.Header().Set(constants.HeaderContentType, mimeType)

	// 设置 Content-Disposition 头以启用下载
	if download {
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+fileInfo.Name()+"\"")
	}

	// 设置 Content-Length 头
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// 发送文件内容
	_, err = io.Copy(c.Writer, file)
	return err
}
