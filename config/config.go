package config

import (
	"io"

	"gopkg.in/yaml.v2"
)

type Run struct {
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
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
