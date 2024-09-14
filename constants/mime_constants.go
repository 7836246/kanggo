package constants

// 定义常见的 MIME 类型的常量
const (
	MIMETextXML               = "text/xml"                          // XML 文本格式
	MIMETextHTML              = "text/html"                         // HTML 文本格式
	MIMETextPlain             = "text/plain"                        // 纯文本格式
	MIMEApplicationXML        = "application/xml"                   // XML 应用程序格式
	MIMEApplicationJSON       = "application/json"                  // JSON 应用程序格式
	MIMEApplicationJavaScript = "application/javascript"            // JavaScript 应用程序格式
	MIMEApplicationForm       = "application/x-www-form-urlencoded" // 表单 URL 编码格式
	MIMEOctetStream           = "application/octet-stream"          // 二进制流数据（任意文件类型）
	MIMEMultipartForm         = "multipart/form-data"               // 多部分表单数据格式（用于文件上传）

	// MIMETextXMLCharsetUTF8 定义常见的带有 UTF-8 字符集的 MIME 类型常量
	MIMETextXMLCharsetUTF8               = "text/xml; charset=utf-8"               // XML 文本格式，UTF-8 字符集
	MIMETextHTMLCharsetUTF8              = "text/html; charset=utf-8"              // HTML 文本格式，UTF-8 字符集
	MIMETextPlainCharsetUTF8             = "text/plain; charset=utf-8"             // 纯文本格式，UTF-8 字符集
	MIMEApplicationXMLCharsetUTF8        = "application/xml; charset=utf-8"        // XML 应用程序格式，UTF-8 字符集
	MIMEApplicationJSONCharsetUTF8       = "application/json; charset=utf-8"       // JSON 应用程序格式，UTF-8 字符集
	MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript; charset=utf-8" // JavaScript 应用程序格式，UTF-8 字符集
)
