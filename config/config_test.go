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
