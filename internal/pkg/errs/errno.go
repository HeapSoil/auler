package errs

import "fmt"

// 错误码规范

type Errno struct {
	HTTP    int
	Code    string
	Message string
}

func (err *Errno) Error() string {
	return err.Message
}

func (err *Errno) SetMessage(format string, args ...interface{}) *Errno {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

// 解析业务的错误码和错误信息
func Decode(err error) (int, string, string) {
	if err == nil {
		return OK.HTTP, OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.HTTP, typed.Code, typed.Message
	default:
	}

	// 默认返回未知错误码和错误信息. 该错误代表服务端出错
	return InternalServerError.HTTP, InternalServerError.Code, err.Error()

}
