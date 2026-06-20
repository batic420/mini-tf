package decoder

type Envelope struct {
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
