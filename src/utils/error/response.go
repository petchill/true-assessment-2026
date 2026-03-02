package error

import (
	"errors"
	"net/http"
)

type ErrorCode string

const (
	ErrModelInference        ErrorCode = "01"
	ErrModelInferenceTimeout ErrorCode = "02"
	ErrGeneral               ErrorCode = "99"
)

func (e ErrorCode) Message() string {
	switch e {
	case ErrModelInference:
		return "model inference failed"
	case ErrModelInferenceTimeout:
		return "Recommendation generation exceeded timeout limit"
	default:
		return "error"
	}
}

type HttpResponseError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type AppError struct {
	StatusCode int
	ErrorCode  ErrorCode
	Message    string
	Cause      error
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return http.StatusText(e.StatusCode)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func New(statusCode int, errorCode ErrorCode, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func Wrap(statusCode int, errorCode ErrorCode, message string, cause error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
		Cause:      cause,
	}
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr != nil {
		return appErr
	}

	if err == nil {
		return &AppError{
			StatusCode: http.StatusInternalServerError,
			ErrorCode:  ErrGeneral,
			Message:    ErrGeneral.Message(),
		}
	}

	return &AppError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  ErrGeneral,
		Message:    err.Error(),
		Cause:      err,
	}
}
