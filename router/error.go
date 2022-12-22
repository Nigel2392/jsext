package router

import "fmt"

const (
	ErrCodeInvalid            = 400
	ErrCodeNoAuth             = 401
	ErrCodeForbidden          = 403
	ErrCodeNotFound           = 404
	ErrCodeNoMethod           = 405
	ErrCodeUnacceptable       = 406
	ErrCodeProxyAuthRequired  = 407
	ErrCodeRQTimeout          = 408
	ErrCodeTeapot             = 418
	ErrCodeInternal           = 500
	ErrCodeNYI                = 501
	ErrCodeBadGateway         = 502
	ErrCodeServiceUnavailable = 503
)

type RouterError struct {
	Message string
	Code    int
}

func (e RouterError) Error() string {
	return fmt.Sprintf("RouterError [%d]: %s", e.Code, e.Message)
}

func (e RouterError) String() string {
	return e.Error()
}

func (e RouterError) IsErr(errCode int) bool {
	return e.Code == errCode
}

func NewError(code int, msg ...string) RouterError {
	if len(msg) > 0 {
		return RouterError{Message: msg[0], Code: code}
	}
	var message string
	switch code {
	case ErrCodeInvalid:
		message = "Invalid request"
	case ErrCodeNoAuth:
		message = "Not authorized"
	case ErrCodeForbidden:
		message = "Forbidden"
	case ErrCodeNotFound:
		message = "Not found"
	case ErrCodeUnacceptable:
		message = "Unacceptable request format"
	case ErrCodeProxyAuthRequired:
		message = "Proxy authentication required"
	case ErrCodeNoMethod:
		message = "Method not allowed"
	case ErrCodeRQTimeout:
		message = "Request timeout"
	case ErrCodeTeapot:
		message = "I'm a teapot"
	case ErrCodeInternal:
		message = "Internal server error"
	case ErrCodeNYI:
		message = "Not yet implemented"
	case ErrCodeBadGateway:
		message = "Bad gateway"
	case ErrCodeServiceUnavailable:
		message = "Service unavailable"
	}
	return RouterError{Message: message, Code: code}
}
