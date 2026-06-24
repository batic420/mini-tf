package creator

import (
	"fmt"
	"os"
	"path/filepath"

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

func CreateResource(env decoder.Envelope) error {
	sortedResources, err := topoSort(env.Resources)
	if err != nil {
		return fmt.Errorf("Got an error sorting the resources based on their dependencies: %w", err)
	}

	for _, resource := range sortedResources {
		switch resource.Type {
		case "file":
			fileName := fmt.Sprintf("%s%s", resource.Name, resource.Properties["extension"])
			fmt.Println(fileName)
		case "folder":
			loc, _ := resource.Properties["location"].(string)
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("Got an error retrieving the current working directory: %w", err)
			}
			fullPath := filepath.Join(cwd, loc, resource.Name)
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return fmt.Errorf("Got an error creating the directory: %w", err)
			}
		}
	}
	return nil
}
