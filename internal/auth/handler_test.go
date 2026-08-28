package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginRejectsInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", NewHandler(nil, false).Login)

	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("{"))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestLoginSetsCookieAndReturnsToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := NewService(fakeUserRepository{user: newFakeUser(t)}, fakeTokenGenerator{})
	router := gin.New()
	router.POST("/login", NewHandler(service, true).Login)

	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"admin","password":"password"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	var body struct {
		Data UserResultWithToken `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Data.AccessToken != "access-token" {
		t.Fatalf("accessToken = %q, want %q", body.Data.AccessToken, "access-token")
	}
	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != AccessTokenCookie || cookies[0].Value != "access-token" {
		t.Fatalf("unexpected cookies: %+v", cookies)
	}
	if !cookies[0].HttpOnly || !cookies[0].Secure || cookies[0].SameSite != http.SameSiteLaxMode || cookies[0].Path != "/" {
		t.Fatalf("insecure cookie attributes: %+v", cookies[0])
	}
}

func TestLogoutClearsCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/logout", NewHandler(nil, false).Logout)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/logout", nil))

	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != AccessTokenCookie || cookies[0].MaxAge != -1 {
		t.Fatalf("cookie was not cleared: %+v", cookies)
	}
}
