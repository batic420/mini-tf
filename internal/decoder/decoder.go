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
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

type Resource struct {
	Name        string            `yaml:"name"`
	DisplayName *string           `yaml:"displayName,omitempty"`
	Type        string            `yaml:"type"`
	DependsOn   *[]string         `yaml:"dependsOn,omitempty"`
	Properties  map[string]any    `yaml:"properties"`
	Outputs     map[string]string `yaml:"outputs"`
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
		if err := s.Check(); err != nil {
			return fmt.Errorf("resources.resource[%d]: %w", i, err)
		}
		if s.DependsOn != nil && !contains(names, s.DependsOn) {
			return fmt.Errorf("resources.resource[%d].dependsOn references unknown resource", i)
		}
	}

	return nil
}

func (r *Resource) Check() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.Type == "" {
		return fmt.Errorf("type is required")
	}
	switch r.Type {
	case "file":
		loc, ok := r.Properties["location"].(string)
		if !ok || loc == "" {
			return fmt.Errorf("resource %q: properties.location is required for type file", r.Name)
		}
		if ext, present := r.Properties["extension"]; present {
			if _, ok := ext.(string); !ok {
				return fmt.Errorf("resource %q: properties.extension must be a string", r.Name)
			}
		}
	case "folder":
		if _, ok := r.Properties["location"].(string); !ok {
			return fmt.Errorf("resource %q: properties.location is required for type folder", r.Name)
		}
	default:
		return fmt.Errorf("resource %q: unknown type %q", r.Name, r.Type)
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
