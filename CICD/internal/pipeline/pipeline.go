package pipeline

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Step is a single execution unit.
type Step struct {
	Name    string   `yaml:"name"`
	Cmd     string   `yaml:"cmd"`
	Timeout string   `yaml:"timeout,omitempty"` // "30s", "5m"
	Env     []string `yaml:"env,omitempty"`     // "KEY=value"
	WorkDir string   `yaml:"work_dir,omitempty"`
}

// Stage groups steps that run in parallel.
type Stage struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

// NotifyConfig holds notification settings.
type NotifyConfig struct {
	DingTalkURL string `yaml:"dingtalk_url,omitempty"`
	EmailSMTP   string `yaml:"email_smtp,omitempty"`
	EmailTo     string `yaml:"email_to,omitempty"`
}

// Pipeline is the root definition of a CI/CD pipeline.
type Pipeline struct {
	Name   string       `yaml:"name"`
	Stages []Stage      `yaml:"stages"`
	Notify NotifyConfig `yaml:"notify,omitempty"`
}

// Load reads a YAML file and returns a Pipeline.
func Load(path string) (*Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read pipeline file: %w", err)
	}
	var p Pipeline
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse pipeline: %w", err)
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Pipeline) validate() error {
	if p.Name == "" {
		return fmt.Errorf("pipeline name is required")
	}
	if len(p.Stages) == 0 {
		return fmt.Errorf("at least one stage is required")
	}
	for i, stage := range p.Stages {
		if stage.Name == "" {
			return fmt.Errorf("stage[%d]: name is required", i)
		}
		if len(stage.Steps) == 0 {
			return fmt.Errorf("stage[%d] %q: at least one step is required", i, stage.Name)
		}
		for j, step := range stage.Steps {
			if step.Name == "" {
				return fmt.Errorf("stage[%d] %q, step[%d]: name is required", i, stage.Name, j)
			}
			if step.Cmd == "" {
				return fmt.Errorf("stage[%d] %q, step[%d] %q: cmd is required", i, stage.Name, j, step.Name)
			}
		}
	}
	return nil
}
