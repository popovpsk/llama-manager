package config

import (
	"os"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	yamlContent := `
runs:
  - name: run1
    params:
      model_path: "Qwen3-32B-Q4_K_M.gguf"
      ngl: "65"
      context_size: "13824"
      flash_attn: true
      tensor_split: "45, 20"
      prio: "3"
      temp: "0.6"
      min_p: "0.0"
      top_p: "0.95"
      top_k: "20"
      host: "0.0.0.0"
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

	if len(cfg.Runs) != 1 {
		t.Errorf("Expected 1 run, got %d", len(cfg.Runs))
	}

	run := cfg.Runs[0]

	expectedCmd := "cd /home/aleksandr/repo/gguf/ && ../llama.cpp/build/bin/llama-server -m Qwen3-32B-Q4_K_M.gguf -ngl 65 -c 13824 --flash-attn --tensor-split \"45, 20\" --prio 3 --temp 0.6 --min-p 0.0 --top-p 0.95 --top-k 20 --host 0.0.0.0"
	if run.Params.BuildCommand() != expectedCmd {
		t.Errorf("Expected command: %s, got: %s", expectedCmd, run.Params.BuildCommand())
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	yamlContent := `
runs:
  - name: run1
    params:
      model_path: "Qwen3-32B-Q4_K_M.gguf"
      ngl "65"  # Invalid YAML: missing colon
      context_size: "13824"
      flash_attn: true
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
	if err == nil {
		t.Fatalf("Expected error for invalid YAML, got nil")
	}
	if cfg != nil {
		t.Errorf("Expected nil config on error, got non-nil")
	}
}

func TestLoad_MissingRequiredFields(t *testing.T) {
	yamlContent := `
runs:
  - name: run1
    params:
      ngl: "65"
      context_size: "13824"
      flash_attn: true
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
	if err == nil {
		t.Fatalf("Expected error for missing required fields, got nil")
	}
	if cfg != nil {
		t.Errorf("Expected nil config on error, got non-nil")
	}
}

func TestRunParams_BuildCommand_WithoutOptionalParams(t *testing.T) {
	yamlContent := `
runs:
  - name: run1
    params:
      model_path: "Qwen3-32B-Q4_K_M.gguf"
      ngl: "65"
      context_size: "13824"
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

	run := cfg.Runs[0]
	expectedCmd := "cd /home/aleksandr/repo/gguf/ && ../llama.cpp/build/bin/llama-server -m Qwen3-32B-Q4_K_M.gguf -ngl 65 -c 13824"
	if run.Params.BuildCommand() != expectedCmd {
		t.Errorf("Expected command: %s, got: %s", expectedCmd, run.Params.BuildCommand())
	}
}

func TestConfig_GetRun(t *testing.T) {
	yamlContent := `
runs:
  - name: run1
    params:
      model_path: "Qwen3-32B-Q4_K_M.gguf"
      ngl: "65"
      context_size: "13824"
  - name: run2
    params:
      model_path: "another-model.gguf"
      ngl: "100"
      context_size: "2048"
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

	run1 := cfg.GetRun("run1")
	if run1 == nil {
		t.Errorf("Expected to find run1")
	}
	if run1.Name != "run1" {
		t.Errorf("Expected run1, got %s", run1.Name)
	}

	run2 := cfg.GetRun("run2")
	if run2 == nil {
		t.Errorf("Expected to find run2")
	}
	if run2.Name != "run2" {
		t.Errorf("Expected run2, got %s", run2.Name)
	}

	missing := cfg.GetRun("missing")
	if missing != nil {
		t.Errorf("Expected nil for missing run")
	}
}
