package decoder

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	// Build helper function to return pointer for normal string value
	sp := func(s string) *string { return &s }

	cases := []struct {
		name    string
		env     Envelope
		wantErr string
	}{
		{
			name: "valid envelope",
			env: Envelope{
				Parameters: &[]Parameter{{
					Name:  "environment",
					Type:  "string",
					Value: "dev",
				}},
				Variables: &map[string]string{
					"defaultPermissions": "read",
				},
				Resources: []Resource{
					{
						Name: "config_folder",
						// --> No DisplayName property
						Type: "folder",
						// --> No DependsOn property
						Properties: map[string]any{
							"location": "/etc/mini-tf",
							// --> No "extension" property
						},
						Outputs: map[string]string{
							"path": "properties.location",
						},
					},
					{
						Name:        "config_file",
						DisplayName: sp("JSON example"),
						Type:        "file",
						Properties: map[string]any{
							"location":  "{{ .resources.config_folder.outputs.path }}",
							"extension": ".json",
						},
						Outputs: map[string]string{
							"path": "properties.location",
						},
						DependsOn: &[]string{
							"config_folder",
						},
					},
				},
			},
		},
		{
			name: "missing required values",
			env: Envelope{
				Parameters: &[]Parameter{{
					Name:  "environment",
					Type:  "string",
					Value: "dev",
				}},
				Variables: &map[string]string{
					"defaultPermissions": "read",
				},
				Resources: []Resource{
					{
						// --> No Name property - REQUIRED
						// --> No DisplayName property
						Type: "folder",
						// --> No DependsOn property
						Properties: map[string]any{
							"location": "/etc/mini-tf",
							// --> No Extension property
						},
						Outputs: map[string]string{
							"path": "properties.location",
						},
					},
					{
						Name:        "config_file",
						DisplayName: sp("JSON example"),
						Type:        "file",
						Properties: map[string]any{
							"location":  "{{ .resources.config_folder.outputs.path }}",
							"extension": ".json",
						},
						Outputs: map[string]string{
							"path": "properties.location",
						},
						DependsOn: &[]string{
							"config_folder",
						},
					},
				},
			},
			wantErr: "name is required",
		},
		{
			name: "unmatching dependsOn property",
			env: Envelope{
				Parameters: &[]Parameter{{
					Name:  "environment",
					Type:  "string",
					Value: "dev",
				}},
				Variables: &map[string]string{
					"defaultPermissions": "read",
				},
				Resources: []Resource{
					{
						Name: "config_folder",
						// --> No DisplayName property
						Type: "folder",
						// --> No DependsOn property
						Properties: map[string]any{
							"location": "/etc/mini-tf",
							// --> No Extension property
						},
						Outputs: map[string]string{
							"path": "properties.location",
						},
					},
					{
						Name:        "config_file",
						DisplayName: sp("JSON example"),
						Type:        "file",
						Properties: map[string]any{
							"location":  "{{ .resources.config_folder.outputs.path }}",
							"extension": ".json",
						},
						Outputs: map[string]string{
							"path": "properties.location",
						},
						DependsOn: &[]string{
							"unknown_resource",
						},
					},
				},
			},
			wantErr: "dependsOn references unknown resource",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.env.Validate()
			if tc.wantErr == "" && err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("expected error containing %q, got: %v", tc.wantErr, err)
			}
		})
	}
}
