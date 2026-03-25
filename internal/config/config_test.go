package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadReadsDotEnv(t *testing.T) {
	t.Setenv("X_API_KEY", "")
	t.Setenv("X_API_SECRET", "")
	t.Setenv("X_ACCESS_TOKEN", "")
	t.Setenv("X_ACCESS_TOKEN_SECRET", "")

	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := strings.Join([]string{
		"X_API_KEY=api-key",
		"X_API_SECRET=api-secret",
		"X_ACCESS_TOKEN=access-token",
		"X_ACCESS_TOKEN_SECRET=access-token-secret",
	}, "\n")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Credentials.APIKey != "api-key" {
		t.Fatalf("APIKey = %q, want %q", cfg.Credentials.APIKey, "api-key")
	}
	if cfg.Credentials.APISecret != "api-secret" {
		t.Fatalf("APISecret = %q, want %q", cfg.Credentials.APISecret, "api-secret")
	}
	if cfg.Credentials.AccessToken != "access-token" {
		t.Fatalf("AccessToken = %q, want %q", cfg.Credentials.AccessToken, "access-token")
	}
	if cfg.Credentials.AccessTokenSecret != "access-token-secret" {
		t.Fatalf("AccessTokenSecret = %q, want %q", cfg.Credentials.AccessTokenSecret, "access-token-secret")
	}
}

func TestLoadDefaultsToAsiaTokyoLocation(t *testing.T) {
	t.Setenv("X_API_KEY", "api-key")
	t.Setenv("X_API_SECRET", "api-secret")
	t.Setenv("X_ACCESS_TOKEN", "access-token")
	t.Setenv("X_ACCESS_TOKEN_SECRET", "access-token-secret")
	t.Setenv("TZ", "")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Location.String() != "Asia/Tokyo" {
		t.Fatalf("Location = %q, want %q", cfg.Location.String(), "Asia/Tokyo")
	}
}

func TestLoadUsesTZWhenProvided(t *testing.T) {
	t.Setenv("X_API_KEY", "api-key")
	t.Setenv("X_API_SECRET", "api-secret")
	t.Setenv("X_ACCESS_TOKEN", "access-token")
	t.Setenv("X_ACCESS_TOKEN_SECRET", "access-token-secret")
	t.Setenv("TZ", "UTC")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Location != time.UTC {
		t.Fatalf("Location = %q, want UTC", cfg.Location.String())
	}
}

func TestLoadPrefersExistingEnvironmentValues(t *testing.T) {
	t.Setenv("X_API_KEY", "from-env")
	t.Setenv("X_API_SECRET", "from-env-secret")
	t.Setenv("X_ACCESS_TOKEN", "from-env-token")
	t.Setenv("X_ACCESS_TOKEN_SECRET", "from-env-token-secret")

	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := strings.Join([]string{
		"X_API_KEY=from-file",
		"X_API_SECRET=from-file-secret",
		"X_ACCESS_TOKEN=from-file-token",
		"X_ACCESS_TOKEN_SECRET=from-file-token-secret",
	}, "\n")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Credentials.APIKey != "from-env" {
		t.Fatalf("APIKey = %q, want %q", cfg.Credentials.APIKey, "from-env")
	}
}
