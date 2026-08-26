package example

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-boilerplate/internal/http/response"
	"go-boilerplate/internal/http/validation"
)

type Handler struct {
	service *Service
}

type CreateExampleRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.BadRequest("invalid request body"))
		return
	}

	if errs := validation.Validate(req); len(errs) > 0 {
		response.Error(c, response.ValidationError(errs))
		return
	}

	example, err := h.service.Create(c.Request.Context(), CreateInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			response.Error(c, response.BadRequest(err.Error()))
		} else {
			response.Error(c, response.InternalServerError("internal server error"))
		}
		return
	}

	response.Success(c, http.StatusCreated, "created", example)
}

func (h *Handler) List(c *gin.Context) {
	examples, err := h.service.List(c.Request.Context())
	if err != nil {
		response.Error(c, response.InternalServerError("internal server error"))
		return
	}

	response.Success(c, http.StatusOK, "success", examples)
}

func (h *Handler) Paginate(c *gin.Context) {
	page := queryInt(c, "page", 1)
	limit := queryInt(c, "limit", 10)
	examples, meta, err := h.service.Paginate(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, response.InternalServerError("internal server error"))
		return
	}

	response.SuccessPaginate(c, http.StatusOK, "success", examples, response.Meta{
		Page:       meta.Page,
		Limit:      meta.Limit,
		Total:      meta.Total,
		TotalPages: meta.TotalPages,
	})
}

func queryInt(c *gin.Context, key string, fallback int) int {
	value, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(fallback)))
	if err != nil || value < 1 {
		return fallback
	}

	return value
}
