package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"go-boilerplate/internal/auth"
	"go-boilerplate/internal/http/response"
)

func Auth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if jwtManager == nil {
			response.Error(c, response.InternalServerError("jwt manager is not configured"))
			c.Abort()
			return
		}

		tokenString := bearerToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			response.Error(c, response.Unauthorized("missing bearer token"))
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseAccessToken(tokenString)
		if err != nil {
			response.Error(c, response.Unauthorized("invalid bearer token"))
			c.Abort()
			return
		}

		c.Set(auth.UserIDKey, claims.Subject)
		c.Next()
	}
}

func bearerToken(header string) string {
	if header == "" {
		return ""
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return parts[1]
}
