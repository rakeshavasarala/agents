package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RuleConfig struct {
	ID       string                 `yaml:"id"`
	Enabled  bool                   `yaml:"enabled"`
	Severity string                 `yaml:"severity,omitempty"`
	Params   map[string]interface{} `yaml:"params,omitempty"`
}

type Ruleset struct {
	Rules []RuleConfig `yaml:"rules"`
}

func LoadRulesConfig(path string) (*Ruleset, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var ruleset Ruleset
	err = yaml.Unmarshal(data, &ruleset)
	if err != nil {
		return nil, err
	}
	return &ruleset, nil
}