package errs

var (
	// 请求成功
	OK = &Errno{
		HTTP:    200,
		Code:    "",
		Message: "",
	}

	// （平台级）所有未知的服务器端错误
	InternalServerError = &Errno{
		HTTP:    500,
		Code:    "InternalError",
		Message: "Internal server error.",
	}

	// （平台级） 路由不匹配错误
	ErrPageNotFound = &Errno{
		HTTP:    404,
		Code:    "ResourceNotFound.PageNotFound",
		Message: "Page not found.",
	}

	// ErrBind 参数绑定错误
	ErrBind = &Errno{
		HTTP:    400,
		Code:    "InvalidParameter.BindError",
		Message: "Error occurred while binding the request body to the struct.",
	}

	// ErrInvalidParameter 表示所有验证失败的错误.
	ErrInvalidParameter = &Errno{
		HTTP:    400,
		Code:    "InvalidParameter",
		Message: "Parameter verification failed.",
	}

	// ErrSignToken 表示签发 JWT Token 时出错.
	ErrSignToken = &Errno{
		HTTP:    401,
		Code:    "AuthFailure.SignTokenError",
		Message: "Error occurred while signing the JSON web token.",
	}

	// ErrTokenInvalid 表示 JWT Token 格式错误.
	ErrTokenInvalid = &Errno{
		HTTP:    401,
		Code:    "AuthFailure.TokenInvalid",
		Message: "Token was invalid.",
	}

	// ErrUnauthorized 表示请求没有被授权.
	ErrUnauthorized = &Errno{
		HTTP:    401,
		Code:    "AuthFailure.Unauthorized",
		Message: "Unauthorized.",
	}
)
