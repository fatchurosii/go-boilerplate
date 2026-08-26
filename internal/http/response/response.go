package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Base struct {
	Message  string `json:"message"`
	Success bool   `json:"success"`
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
	Code    int
	Message string
	Errors  []FieldError
}

func newBase(status bool, message string) Base {
	return Base{Message: message, Success: status}
}

func Success[T any](c *gin.Context, status int, message string, data T) {
	c.JSON(status, Response[T]{
		Base: newBase(true, message),
		Data: data,
	})
}

func SuccessPaginate[T any](c *gin.Context, status int, message string, data []T, meta Meta) {
	c.JSON(status, Response[[]T]{
		Base: newBase(true, message),
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

	c.JSON(err.Code, ErrorResponse{
		Base:   newBase(false, err.Message),
		Errors: errs,
	})
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

func newError(code int, message string, errors []FieldError) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Errors:  errors,
	}
}
