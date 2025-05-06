package config

import (
	"os"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	yamlContent := `
runs:
  - name: run1
    cmd: echo "Hello World"
  - name: run2
    cmd: echo "Goodbye World"
`
	tempFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	cfg, err := Load(tempFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(cfg.Runs) != 2 {
		t.Errorf("Expected 2 runs, got %d", len(cfg.Runs))
	}

	if cfg.Runs[0].Name != "run1" || cfg.Runs[0].Cmd != "echo \"Hello World\"" {
		t.Errorf("Expected run1 with cmd 'echo \"Hello World\"', got %v", cfg.Runs[0])
	}

	if cfg.Runs[1].Name != "run2" || cfg.Runs[1].Cmd != "echo \"Goodbye World\"" {
		t.Errorf("Expected run2 with cmd 'echo \"Goodbye World\"', got %v", cfg.Runs[1])
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("nonexistent.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	yamlContent := `invalid: -`
	tempFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	_, err = Load(tempFile.Name())
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestGetRun_Found(t *testing.T) {
	cfg := &Config{
		Runs: []Run{
			{Name: "run1", Cmd: "echo 'Hello'"},
			{Name: "run2", Cmd: "echo 'World'"},
		},
	}

	run := cfg.GetRun("run1")
	if run == nil {
		t.Error("Expected run1 to be found, got nil")
	}
	if run.Name != "run1" || run.Cmd != "echo 'Hello'" {
		t.Errorf("Expected run1 with cmd 'echo Hello', got %v", *run)
	}
}

func TestGetRun_NotFound(t *testing.T) {
	cfg := &Config{
		Runs: []Run{
			{Name: "run1", Cmd: "echo 'Hello'"},
		},
	}

	run := cfg.GetRun("run2")
	if run != nil {
		t.Errorf("Expected run2 to be not found, got %v", *run)
	}
}
