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
					Name:    "environment",
					Type:    "string",
					Default: "dev",
					Values:  []string{"dev", "stage", "prod"},
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
						Properties: ResourceProperties{
							Location: "/etc/mini-tf",
							// --> No Extension property
						},
						Outputs: &ResourceOutputs{
							Path: "properties.location",
						},
					},
					{
						Name:        "config_file",
						DisplayName: sp("JSON example"),
						Type:        "file",
						Properties: ResourceProperties{
							Location:  "{{ .resources.config_folder.outputs.path }}",
							Extension: sp(".json"),
						},
						Outputs: &ResourceOutputs{
							Path: "properties.location",
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
					Name:    "environment",
					Type:    "string",
					Default: "dev",
					Values:  []string{"dev", "stage", "prod"},
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
						Properties: ResourceProperties{
							Location: "/etc/mini-tf",
							// --> No Extension property
						},
						Outputs: &ResourceOutputs{
							Path: "properties.location",
						},
					},
					{
						Name:        "config_file",
						DisplayName: sp("JSON example"),
						Type:        "file",
						Properties: ResourceProperties{
							Location:  "{{ .resources.config_folder.outputs.path }}",
							Extension: sp(".json"),
						},
						Outputs: &ResourceOutputs{
							Path: "properties.location",
						},
						DependsOn: &[]string{
							"config_folder",
						},
					},
				},
			},
			wantErr: "Required property on at least one resource is missing",
		},
		{
			name: "unmatching dependsOn property",
			env: Envelope{
				Parameters: &[]Parameter{{
					Name:    "environment",
					Type:    "string",
					Default: "dev",
					Values:  []string{"dev", "stage", "prod"},
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
						Properties: ResourceProperties{
							Location: "/etc/mini-tf",
							// --> No Extension property
						},
						Outputs: &ResourceOutputs{
							Path: "properties.location",
						},
					},
					{
						Name:        "config_file",
						DisplayName: sp("JSON example"),
						Type:        "file",
						Properties: ResourceProperties{
							Location:  "{{ .resources.config_folder.outputs.path }}",
							Extension: sp(".json"),
						},
						Outputs: &ResourceOutputs{
							Path: "properties.location",
						},
						DependsOn: &[]string{
							"unknown_resource",
						},
					},
				},
			},
			wantErr: "At least one dependsOn property points to an unknown resource",
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
