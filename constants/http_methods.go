package constants

// 定义 HTTP 方法的常量
const (
	MethodGet     = "GET"     // RFC 7231, 4.3.1  GET 方法，用于请求指定资源的信息
	MethodHead    = "HEAD"    // RFC 7231, 4.3.2  HEAD 方法，与 GET 方法相同，但不返回主体部分
	MethodPost    = "POST"    // RFC 7231, 4.3.3  POST 方法，用于向指定资源提交数据
	MethodPut     = "PUT"     // RFC 7231, 4.3.4  PUT 方法，用于向指定资源上传数据
	MethodPatch   = "PATCH"   // RFC 5789         PATCH 方法，用于对指定资源进行部分修改
	MethodDelete  = "DELETE"  // RFC 7231, 4.3.5  DELETE 方法，用于删除指定资源
	MethodConnect = "CONNECT" // RFC 7231, 4.3.6  CONNECT 方法，用于建立到目标资源的隧道
	MethodOptions = "OPTIONS" // RFC 7231, 4.3.7  OPTIONS 方法，用于获取当前 URL 所支持的方法
	MethodTrace   = "TRACE"   // RFC 7231, 4.3.8  TRACE 方法，用于回显服务器收到的请求，主要用于测试或诊断
	MethodUse     = "USE"     // 自定义方法，用于中间件等操作的注册
)
