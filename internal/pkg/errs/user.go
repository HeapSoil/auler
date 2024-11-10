package errs

var (
	ErrUserAlreadyExists = &Errno{
		HTTP:    400,
		Code:    "FailedOperation.UserAlreadyExist",
		Message: "User already exist.",
	}
)
