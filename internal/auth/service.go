package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go-boilerplate/internal/user"
)

type Service struct {
	repository userRepository
	tokens     tokenGenerator
}

type userRepository interface {
	FindUserByUsername(context.Context, string) (*user.User, error)
	FindUserByID(context.Context, uuid.UUID) (*user.User, error)
}

type tokenGenerator interface {
	GenerateAccessToken(string) (string, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidUser        = errors.New("invalid user")
)

type LoginInput struct {
	Username string
	Password string
}

type UserResultWithToken struct {
	AccessToken string     `json:"accessToken"`
	TokenType   string     `json:"tokenType"`
	User        UserResult `json:"user"`
}

type UserResult struct {
	ID       string     `json:"id"`
	Username string     `json:"username"`
	IsActive bool       `json:"isActive"`
	Role     RoleResult `json:"role"`
}

type RoleResult struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	IsActive bool   `json:"isActive"`
}

func NewService(repository userRepository, tokens tokenGenerator) *Service {
	return &Service{repository: repository, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*UserResultWithToken, error) {
	found, err := s.repository.FindUserByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if found == nil || !found.IsActive || !found.Role.IsActive {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokens.GenerateAccessToken(found.ID.String())
	if err != nil {
		return nil, err
	}

	return &UserResultWithToken{
		AccessToken: token,
		TokenType:   "Bearer",
		User:        toUserResult(found),
	}, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*UserResult, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidUser
	}

	found, err := s.repository.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if found == nil || !found.IsActive || !found.Role.IsActive {
		return nil, ErrInvalidUser
	}

	result := toUserResult(found)
	return &result, nil
}

func toUserResult(found *user.User) UserResult {
	return UserResult{
		ID:       found.ID.String(),
		Username: found.Username,
		IsActive: found.IsActive,
		Role: RoleResult{
			ID:       found.Role.ID.String(),
			Name:     found.Role.Name,
			Slug:     found.Role.Slug,
			IsActive: found.Role.IsActive,
		},
	}
}
