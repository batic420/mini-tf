package decoder

import (
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/goccy/go-yaml"
)

type Envelope struct {
	Parameters *[]Parameter       `yaml:"parameters,omitempty"`
	Variables  *map[string]string `yaml:"variables,omitempty"`
	Resources  []Resource         `yaml:"resources"`
}

type Parameter struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Default string   `yaml:"default"`
	Values  []string `yaml:"values"`
}

type Resource struct {
	Name        string             `yaml:"name"`
	DisplayName *string            `yaml:"displayName,omitempty"`
	Type        string             `yaml:"type"`
	DependsOn   *[]string          `yaml:"dependsOn,omitempty"`
	Properties  ResourceProperties `yaml:"properties"`
	Outputs     *ResourceOutputs   `yaml:"outputs"`
}

type ResourceProperties struct {
	Location  string  `yaml:"location"`
	Extension *string `yaml:"extension,omitempty"`
}

type ResourceOutputs struct {
	Path string `yaml:"path"`
}

func Parse(r io.Reader) (*Envelope, error) {
	dec := yaml.NewDecoder(r, yaml.Strict())

	var env Envelope
	if err := dec.Decode(&env); err != nil {
		return nil, fmt.Errorf("parsing config: %w\n", err)
	}

	if err := env.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w\n", err)
	}

	return &env, nil
}

func Load(path string) (*Envelope, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config: %w", err)
	}
	defer f.Close()
	return Parse(f)
}

func (env *Envelope) Validate() error {
	if len(env.Resources) == 0 {
		return fmt.Errorf("`resources` must contain at least one resource")
	}

	names := make([]string, 0, len(env.Resources))
	for _, v := range env.Resources {
		names = append(names, v.Name)
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
		if s.DependsOn != nil && !contains(names, s.DependsOn) {
			return fmt.Errorf("resources.resource[%d].dependsOn references unknown resource", i)
		}
	}

	return nil
}

func contains(existing []string, targets *[]string) bool {
	for _, target := range *targets {
		if !slices.Contains(existing, target) {
			return false
		}
	}
	return true
}
