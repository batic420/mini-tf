package decoder

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Envelope struct {
	Parameters *[]Parameter `yaml:"parameters,omitempty"`
	Variables  *Variables   `yaml:"variables,omitempty"`
	Resources  []Resource   `yaml:"resources"`
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
	DisplayName *string            `yaml:"displayName"`
	Type        string             `yaml:"type"`
	Properties  ResourceProperties `yaml:"properties"`
	Outputs     *ResourceOutputs   `yaml:"outputs"`
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

func (env *Envelope) Validate() error {
	if len(env.Resources) == 0 {
		return fmt.Errorf("`resources` must contain at least one resource")
	}

	for i, s := range env.Resources {
		if s.Name == "" {
			return fmt.Errorf("resources.resource[%d].name is required", i)
		}
		if s.Type == "" {
			return fmt.Errorf("resources.resource[%d].type is required", i)
		}
		if s.Type == "file" && s.Properties.Extension == nil {
			return fmt.Errorf("resources.resource[%d].properties.extension is required when using type `file`", i)
		}
	}

	return nil
}
