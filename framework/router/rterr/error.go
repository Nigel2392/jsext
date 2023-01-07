package rterr

import "fmt"

const (
	ErrCodeInvalid            = 400 // Invalid request
	ErrCodeNoAuth             = 401 // No authentication
	ErrCodeForbidden          = 403 // Forbidden
	ErrCodeNotFound           = 404 // Not found
	ErrCodeNoMethod           = 405 // Method not allowed
	ErrCodeUnacceptable       = 406 // Unacceptable request format
	ErrCodeProxyAuthRequired  = 407 // Proxy authentication required
	ErrCodeRQTimeout          = 408 // Request timeout
	ErrCodeTeapot             = 418 // I'm a teapot
	ErrCodeInternal           = 500 // Internal server error
	ErrCodeNYI                = 501 // Not yet implemented
	ErrCodeBadGateway         = 502 // Bad gateway interface
	ErrCodeServiceUnavailable = 503 // Service unavailable
)

type ErrorThrower interface {
	Error(code int, message string) RouterError
	Throw(code int)
}

// RouterError is a custom error type for the router.
type RouterError struct {
	Message string
	Code    int
}

func (e RouterError) Error() string {
	return fmt.Sprintf("Error code %d: %s", e.Code, e.Message)
}

func (e RouterError) String() string {
	return e.Error()
}

// Check if the error code matches the specified error code.
func (e RouterError) IsCode(errCode ...int) bool {
	// If no error codes are specified, return true if the error code is 0.
	if len(errCode) == 0 {
		return e.Code == 0
	}
	var allTrue = true
	for _, code := range errCode {
		allTrue = allTrue && e.Code == code
	}
	return allTrue
}

// Check if the error is a RouterError and if the error code matches the specified error code.
func IsRouterError(err error, errCode ...int) bool {
	var routerErr, ok = err.(RouterError)
	if !ok {
		return false
	}
	return routerErr.IsCode(errCode...)
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
	default:
		message = "Unknown error has occurred."
	}
	return RouterError{Message: message, Code: code}
}
