package errs

var (
	// 请求成功
	OK = &Errno{HTTP: 200, Code: "", Message: ""}

	// （平台级）所有未知的服务器端错误
	InternalServerError = &Errno{HTTP: 500, Code: "InternalError", Message: "Internal server error."}

	// （平台级） 路由不匹配错误
	ErrPageNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.PageNotFound", Message: "Page not found."}
)
