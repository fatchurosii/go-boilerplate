package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-boilerplate/internal/http/response"
	"go-boilerplate/internal/http/validation"
)

type Handler struct {
	service *Service
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.BadRequest("invalid request body"))
		return
	}
	if errs := validation.Validate(req); len(errs) > 0 {
		response.Error(c, response.ValidationError(errs))
		return
	}

	result, err := h.service.Login(c.Request.Context(), LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, serviceError(err))
		return
	}
	response.Success(c, http.StatusOK, "login success", result)
}

func (h *Handler) Logout(c *gin.Context) {
	response.Success[any](c, http.StatusOK, "logout success", nil)
}

func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		response.Error(c, response.Unauthorized("unauthorized"))
		return
	}
	userIDString, ok := userID.(string)
	if !ok || userIDString == "" {
		response.Error(c, response.Unauthorized("unauthorized"))
		return
	}

	result, err := h.service.Me(c.Request.Context(), userIDString)
	if err != nil {
		response.Error(c, serviceError(err))
		return
	}
	response.Success(c, http.StatusOK, "success", result)
}

func serviceError(err error) *response.AppError {
	if errors.Is(err, ErrInvalidCredentials) || errors.Is(err, ErrInvalidUser) {
		return response.Unauthorized(err.Error())
	}
	return response.InternalServerError("internal server error")
}
