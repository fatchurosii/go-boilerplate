package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const UserIDKey = "userId"

type JWTManager struct {
	secret         []byte
	issuer         string
	accessTokenTTL time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
}

func NewJWTManager(secret, issuer string, accessTokenTTL time.Duration) (*JWTManager, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}
	if accessTokenTTL <= 0 {
		return nil, errors.New("JWT_ACCESS_TOKEN_TTL must be positive")
	}

	return &JWTManager{
		secret:         []byte(secret),
		issuer:         issuer,
		accessTokenTTL: accessTokenTTL,
	}, nil
}

func (m *JWTManager) GenerateAccessToken(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenTTL)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

func (m *JWTManager) ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected token signing method")
		}

		return m.secret, nil
	}, jwt.WithIssuer(m.issuer))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
