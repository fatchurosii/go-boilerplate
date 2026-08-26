package example

import (
	"context"
	"errors"
)

type Service struct{}

var ErrEmailExists = errors.New("email already exists")

type CreateInput struct {
	Name  string
	Email string
}

type Example struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PageMeta struct {
	Page       int
	Limit      int
	Total      int
	TotalPages int
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*Example, error) {
	if input.Email == "taken@example.com" {
		return nil, ErrEmailExists
	}

	return &Example{ID: 1, Name: input.Name, Email: input.Email}, nil
}

func (s *Service) List(ctx context.Context) ([]Example, error) {
	return []Example{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
	}, nil
}

func (s *Service) Paginate(ctx context.Context, page, limit int) ([]Example, PageMeta, error) {
	examples, err := s.List(ctx)
	if err != nil {
		return nil, PageMeta{}, err
	}
	total := len(examples)
	meta := PageMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: (total + limit - 1) / limit,
	}

	return examples, meta, nil
}
