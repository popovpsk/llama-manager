package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Run struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Params      RunParams `yaml:"params"`
}

type RunParams struct {
	ModelPath   string `yaml:"model_path"`
	NGLayers    string `yaml:"ngl"`
	ContextSize string `yaml:"context_size"`
	FlashAttn   bool   `yaml:"flash_attn"`
	TensorSplit string `yaml:"tensor_split,omitempty"`
	Priority    string `yaml:"prio,omitempty"`
	Temperature string `yaml:"temp,omitempty"`
	MinP        string `yaml:"min_p,omitempty"`
	TopP        string `yaml:"top_p,omitempty"`
	TopK        string `yaml:"top_k,omitempty"`
	Host        string `yaml:"host,omitempty"`
}

func (p *RunParams) BuildCommand() string {
	var parts []string
	parts = append(parts, "cd /home/aleksandr/repo/gguf/ && ../llama.cpp/build/bin/llama-server")
	parts = append(parts, "-m", p.ModelPath)
	parts = append(parts, "-ngl", p.NGLayers)
	parts = append(parts, "-c", p.ContextSize)
	if p.FlashAttn {
		parts = append(parts, "--flash-attn")
	}
	if p.TensorSplit != "" {
		parts = append(parts, "--tensor-split", fmt.Sprintf("\"%s\"", p.TensorSplit))
	}
	if p.Priority != "" {
		parts = append(parts, "--prio", p.Priority)
	}
	if p.Temperature != "" {
		parts = append(parts, "--temp", p.Temperature)
	}
	if p.MinP != "" {
		parts = append(parts, "--min-p", p.MinP)
	}
	if p.TopP != "" {
		parts = append(parts, "--top-p", p.TopP)
	}
	if p.TopK != "" {
		parts = append(parts, "--top-k", p.TopK)
	}
	if p.Host != "" {
		parts = append(parts, "--host", p.Host)
	}
	return join(parts)
}

func join(parts []string) string {
	return strings.Join(parts, " ")
}

type Config struct {
	Runs []Run `yaml:"runs"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetRun(name string) *Run {
	for i := range c.Runs {
		if c.Runs[i].Name == name {
			return &c.Runs[i]
		}
	}
	return nil
}
