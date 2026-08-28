package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"go-boilerplate/internal/auth"
	"go-boilerplate/internal/http/response"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "success", gin.H{"ok": true})
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(requestIDHeader, "request-123")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if got := recorder.Header().Get(requestIDHeader); got != "request-123" {
		t.Fatalf("request ID header = %q, want %q", got, "request-123")
	}

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if _, exists := body["requestId"]; exists {
		t.Fatal("requestId must not be present in response body")
	}
}

func TestAuthAcceptsAccessTokenCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager, err := auth.NewJWTManager("test-secret", "test-issuer", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	token, err := jwtManager.GenerateAccessToken("user-123")
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.Use(Auth(jwtManager))
	router.GET("/", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "success", gin.H{"userId": c.GetString(auth.UserIDKey)})
	})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: auth.AccessTokenCookie, Value: token})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}
