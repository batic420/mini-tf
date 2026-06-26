package creator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/batic420/mini-tf/internal/decoder"
)

func topoSort(resources []decoder.Resource) ([]decoder.Resource, error) {
	index := make(map[string]*decoder.Resource, len(resources))
	inDegree := make(map[string]int, len(resources))

	for i := range resources {
		index[resources[i].Name] = &resources[i]
		inDegree[resources[i].Name] = 0
	}
	for _, r := range resources {
		if r.DependsOn != nil {
			inDegree[r.Name] += len(*r.DependsOn)
		}
	}

	queue := []string{}
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}

	sorted := make([]decoder.Resource, 0, len(resources))
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		sorted = append(sorted, *index[name])

		// find resources that depend on this one and unblock them
		for _, r := range resources {
			if r.DependsOn == nil {
				continue
			}
			for _, dep := range *r.DependsOn {
				if dep == name {
					inDegree[r.Name]--
					if inDegree[r.Name] == 0 {
						queue = append(queue, r.Name)
					}
				}
			}
		}
	}

	if len(sorted) != len(resources) {
		return nil, fmt.Errorf("circular dependency detected")
	}
	return sorted, nil
}

func buildTemplateData(resolved map[string]map[string]string) map[string]any {
	data := map[string]any{}
	for name, outputs := range resolved {
		data[name] = map[string]any{"outputs": outputs}
	}
	return data
}

func resolveString(val string, resolved map[string]map[string]string) (string, error) {
	if !strings.Contains(val, "{{") {
		return val, nil
	}
	tmpl, err := template.New("").Parse(val)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, map[string]any{"resources": buildTemplateData(resolved)}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func resolveProps(resource decoder.Resource, resolved map[string]map[string]string) (map[string]string, error) {
	result := map[string]string{}
	for key, val := range resource.Properties {
		strVal, ok := val.(string)
		if !ok {
			continue
		}
		resolvedVal, err := resolveString(strVal, resolved)
		if err != nil {
			return nil, fmt.Errorf("property %q: %w", key, err)
		}
		result[key] = resolvedVal
	}
	return result, nil
}

func CreateResource(env decoder.Envelope) error {
	sortedResources, err := topoSort(env.Resources)
	if err != nil {
		return fmt.Errorf("Got an error sorting the resources based on their dependencies: %w", err)
	}

	resolved := map[string]map[string]string{}

	for _, resource := range sortedResources {
		props, err := resolveProps(resource, resolved)
		if err != nil {
			return fmt.Errorf("resource %q: %w", resource.Name, err)
		}

		switch resource.Type {
		case "file":
			fileName := fmt.Sprintf("%s%s", resource.Name, props["extension"])
			filePath := props["location"]
			fullPath := filepath.Join(filePath, fileName)
			f, err := os.Create(fullPath)
			if err != nil {
				return fmt.Errorf("creating file %q: %w", fullPath, err)
			}
			if err := f.Close(); err != nil {
				return fmt.Errorf("closing file %q: %w", fullPath, err)
			}
		case "folder":
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("Got an error retrieving the current working directory: %w", err)
			}
			fullPath := filepath.Join(cwd, props["location"], resource.Name)
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return fmt.Errorf("Got an error creating the directory: %w", err)
			}
			resolved[resource.Name] = map[string]string{"path": fullPath}
		}
	}
	return nil
}
