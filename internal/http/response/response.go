package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "requestId"

// Base isinya field yang selalu ada di SETIAP response, sukses maupun error
type Base struct {
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}

type Response[T any] struct {
	Base
	Data T     `json:"data,omitempty"`
	Meta *Meta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Base
	Errors []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type AppError struct {
	Status  int
	Code    int
	Message string
	Errors  []FieldError
}

func newBase(c *gin.Context, message string) Base {
	return Base{
		Message:   message,
		RequestID: RequestID(c),
	}
}

func Success[T any](c *gin.Context, status int, message string, data T) {
	c.JSON(status, Response[T]{
		Base: newBase(c, message),
		Data: data,
	})
}

func SuccessPaginate[T any](c *gin.Context, status int, message string, data []T, meta Meta) {
	c.JSON(status, Response[[]T]{
		Base: newBase(c, message),
		Data: data,
		Meta: &meta,
	})
}

func Error(c *gin.Context, err *AppError) {
	if err == nil {
		err = InternalServerError("internal server error")
	}
	errs := err.Errors
	if errs == nil {
		errs = []FieldError{}
	}
	c.JSON(err.Status, ErrorResponse{
		Base:   newBase(c, err.Message),
		Errors: errs,
	})
}

func RequestID(c *gin.Context) string {
	value, exists := c.Get(RequestIDKey)
	if !exists {
		return ""
	}
	requestID, _ := value.(string)
	return requestID
}

func BadRequest(message string) *AppError {
	return newError(http.StatusBadRequest, message, nil)
}

func Unauthorized(message string) *AppError {
	return newError(http.StatusUnauthorized, message, nil)
}

func Forbidden(message string) *AppError {
	return newError(http.StatusForbidden, message, nil)
}

func NotFound(message string) *AppError {
	return newError(http.StatusNotFound, message, nil)
}

func InternalServerError(message string) *AppError {
	return newError(http.StatusInternalServerError, message, nil)
}

func ValidationError(errors []FieldError) *AppError {
	return newError(http.StatusBadRequest, "validation failed", errors)
}

func newError(status int, message string, errors []FieldError) *AppError {
	return &AppError{
		Status:  status,
		Code:    status,
		Message: message,
		Errors:  errors,
	}
}
