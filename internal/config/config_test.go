package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://example.com")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "false")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("JWT_ISSUER", "test")
	t.Setenv("JWT_ACCESS_TOKEN_TTL", "30m")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "9090" || cfg.JWTAccessTokenTTL != 30*time.Minute || cfg.CORSAllowCredentials {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		key   string
		value string
		want  string
	}{
		{key: "PORT", value: "70000", want: "PORT"},
		{key: "CORS_ALLOW_CREDENTIALS", value: "maybe", want: "CORS_ALLOW_CREDENTIALS"},
		{key: "JWT_ACCESS_TOKEN_TTL", value: "0s", want: "JWT_ACCESS_TOKEN_TTL"},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			t.Setenv("PORT", "8080")
			t.Setenv("CORS_ALLOW_CREDENTIALS", "true")
			t.Setenv("JWT_ACCESS_TOKEN_TTL", "15m")
			t.Setenv(test.key, test.value)

			_, err := Load()
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want error containing %q", err, test.want)
			}
		})
	}
}
