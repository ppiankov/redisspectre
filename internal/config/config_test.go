package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad_Found(t *testing.T) {
	dir := t.TempDir()
	content := `addr: redis.example.com:6379
db: 2
sample_size: 5000
idle_days: 60
big_key_size: 5242880
format: json
timeout: 10m
`
	if err := os.WriteFile(filepath.Join(dir, ".redisspectre.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Addr != "redis.example.com:6379" {
		t.Errorf("expected addr 'redis.example.com:6379', got %q", cfg.Addr)
	}
	if cfg.DB != 2 {
		t.Errorf("expected db 2, got %d", cfg.DB)
	}
	if cfg.SampleSize != 5000 {
		t.Errorf("expected sample_size 5000, got %d", cfg.SampleSize)
	}
	if cfg.IdleDays != 60 {
		t.Errorf("expected idle_days 60, got %d", cfg.IdleDays)
	}
	if cfg.BigKeySize != 5242880 {
		t.Errorf("expected big_key_size 5242880, got %d", cfg.BigKeySize)
	}
	if cfg.Format != "json" {
		t.Errorf("expected format 'json', got %q", cfg.Format)
	}
	if cfg.Timeout != "10m" {
		t.Errorf("expected timeout '10m', got %q", cfg.Timeout)
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Addr != "" {
		t.Errorf("expected empty addr, got %q", cfg.Addr)
	}
}

func TestLoad_YMLExtension(t *testing.T) {
	dir := t.TempDir()
	content := `format: sarif`
	if err := os.WriteFile(filepath.Join(dir, ".redisspectre.yml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "sarif" {
		t.Errorf("expected format 'sarif', got %q", cfg.Format)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	content := `{invalid yaml[[`
	if err := os.WriteFile(filepath.Join(dir, ".redisspectre.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestTimeoutDuration(t *testing.T) {
	cfg := Config{Timeout: "5m"}
	if cfg.TimeoutDuration() != 5*time.Minute {
		t.Errorf("expected 5m, got %v", cfg.TimeoutDuration())
	}

	cfg2 := Config{}
	if cfg2.TimeoutDuration() != 0 {
		t.Errorf("expected 0 for empty timeout, got %v", cfg2.TimeoutDuration())
	}
}
