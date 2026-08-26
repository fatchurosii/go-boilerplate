package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go-boilerplate/internal/role"
	"go-boilerplate/internal/user"
)

type fakeUserRepository struct {
	user *user.User
	err  error
}

func (r fakeUserRepository) FindUserByUsername(context.Context, string) (*user.User, error) {
	return r.user, r.err
}

func (r fakeUserRepository) FindUserByID(context.Context, uuid.UUID) (*user.User, error) {
	return r.user, r.err
}

func (r fakeUserRepository) CreateUser(context.Context, *user.User) error {
	return r.err
}

type fakeTokenGenerator struct{}

func (fakeTokenGenerator) GenerateAccessToken(string) (string, error) {
	return "access-token", nil
}

func newFakeUser(t *testing.T) *user.User {
	t.Helper()
	password, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	return &user.User{
		ID:       uuid.New(),
		Username: "admin",
		Password: string(password),
		IsActive: true,
		Role: role.Role{
			ID:       uuid.New(),
			Name:     "Admin",
			Slug:     "admin",
			IsActive: true,
		},
	}
}

func TestServiceLogin(t *testing.T) {
	fakeUser := newFakeUser(t)
	service := NewService(fakeUserRepository{user: fakeUser}, fakeTokenGenerator{})

	result, err := service.Login(context.Background(), LoginInput{Username: "admin", Password: "password"})
	if err != nil {
		t.Fatal(err)
	}
	if result.AccessToken != "access-token" || result.User.ID != fakeUser.ID.String() {
		t.Fatalf("unexpected result: %+v", result)
	}

	_, err = service.Login(context.Background(), LoginInput{Username: "admin", Password: "wrong"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("error = %v, want ErrInvalidCredentials", err)
	}
}

func TestServiceLoginInvalidUser(t *testing.T) {
	service := NewService(fakeUserRepository{user: nil}, fakeTokenGenerator{})

	// Username gak ketemu di repository (fake return nil, nil).
	// Password isinya bebas / gak relevan karena Login harusnya reject
	// sebelum sempet nyampe tahap cek password.
	_, err := service.Login(context.Background(), LoginInput{Username: "admin1", Password: "irrelevant-password"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("error = %v, want ErrInvalidCredentials", err)
	}
}

func TestServiceMe(t *testing.T) {
	fakeUser := newFakeUser(t)
	service := NewService(fakeUserRepository{user: fakeUser}, fakeTokenGenerator{})

	result, err := service.Me(context.Background(), fakeUser.ID.String())
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != fakeUser.ID.String() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestServiceMeInvalidUser(t *testing.T) {
	service := NewService(fakeUserRepository{user: nil}, fakeTokenGenerator{})

	// User gak ketemu di repository; ID yang dikirim cukup valid UUID
	// (gak perlu nyambung ke fixture user manapun) karena yang dites
	// adalah kondisi "user not found", bukan data user itu sendiri.
	randomUserID := uuid.New()
	_, err := service.Me(context.Background(), randomUserID.String())
	if !errors.Is(err, ErrInvalidUser) {
		t.Fatalf("error = %v, want ErrInvalidUser", err)
	}
}
