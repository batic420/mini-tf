package decoder

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Envelope struct {
	Parameters []Parameter `yaml:"parameters"`
	Variables  Variables   `yaml:"variables"`
	Resources  []Resource  `yaml:"resources"`
}

type Parameter struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Default string   `yaml:"default"`
	Values  []string `yaml:"values"`
}

type Variables struct {
	Variable string `yaml:"variables"`
}

type Resource struct {
	Name        string             `yaml:"name"`
	DisplayName string             `yaml:"displayName"`
	Type        string             `yaml:"type"`
	Properties  ResourceProperties `yaml:"properties"`
	Outputs     ResourceOutputs    `yaml:"outputs"`
}

type ResourceProperties struct {
	Location  string  `yaml:"location"`
	Extension *string `yaml:"displayName,omitempty"`
}

type ResourceOutputs struct {
	Path string `yaml:"path"`
}

func Load(path string) (*Envelope, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config: %w\n", err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f, yaml.Strict())

	var env Envelope
	if err := dec.Decode(&env); err != nil {
		return nil, fmt.Errorf("parsing config: %w\n", err)
	}

	return &env, nil
}

// TODO:
// - Add pointer receiver for validation of `env`
