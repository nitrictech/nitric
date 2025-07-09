package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ValidationTest struct {
	name        string
	jsonData    string
	expectError bool
	err         string
}

func TestApplicationSchemaValidationErrors(t *testing.T) {
	tests := []ValidationTest{
		{
			name:        "missing required name field",
			jsonData:    `{"description": "test application", "targets": ["aws"]}`,
			expectError: true,
			err:         "Invalid application configuration: The `name` property is required",
		},
		{
			name:        "missing required targets field",
			jsonData:    `{"name": "test-app", "description": "test application"}`,
			expectError: true,
			err:         "Invalid application configuration: The `targets` property is required",
		},
		{
			name:        "missing required description field",
			jsonData:    `{"name": "test-app", "targets": ["aws"]}`,
			expectError: true,
			err:         "Invalid application configuration: The `description` property is required",
		},
		{
			name:        "valid application",
			jsonData:    `{"name": "test-app", "description": "test application", "targets": ["aws"]}`,
			expectError: false,
			err:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, result, err := ApplicationFromJson(tt.jsonData)
			if err != nil {
				assert.Fail(t, "Failed to parse JSON: %v", err)
			}

			if tt.expectError {
				assert.NotNil(t, result.Errors())
				assert.False(t, result.Valid())

				errorStr := FormatErrors(result)
				assert.Contains(t, errorStr, tt.err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
				if result != nil {
					assert.Nil(t, result.Errors())
					assert.True(t, result.Valid())
				}
			}
		})
	}
}

func TestServiceSchemaValidationErrors(t *testing.T) {
	tests := []ValidationTest{
		{
			name: "missing container configuration",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"services": {
					"my-service": {}
				}
			}`,
			expectError: true,
			err:         "service my-service has an invalid config: The `container` property is required",
		},
		{
			name: "missing both docker and image in container",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {}
					}
				}
			}`,
			expectError: true,
			err:         "service my-service has an invalid container property: Must provide either a valid docker or image configuration. But not both",
		},
		{
			name: "both docker and image specified",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {
							"docker": {"dockerfile": "Dockerfile"},
							"image": {"id": "nginx:latest"}
						}
					}
				}
			}`,
			expectError: true,
			err:         "service my-service has an invalid container property: Must provide either a valid docker or image configuration. But not both",
		},
		{
			name: "docker image missing required id",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {
							"image": {}
						}
					}
				}
			}`,
			expectError: true,
			err:         " - service my-service has an invalid container property (image): The `id` property is required",
		},
		{
			name: "valid service with docker",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {
							"docker": {
								"dockerfile": "Dockerfile",
								"context": "."
							}
						}
					}
				}
			}`,
			expectError: false,
		},
		{
			name: "valid service with image",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {
							"image": {
								"id": "nginx:latest"
							}
						}
					}
				}
			}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, result, err := ApplicationFromJson(tt.jsonData)
			if err != nil {
				assert.Fail(t, "Failed to parse JSON: %v", err)
			}

			if tt.expectError {
				assert.NotNil(t, result.Errors())
				assert.False(t, result.Valid())

				errorStr := FormatErrors(result)
				assert.Contains(t, errorStr, tt.err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
				if result != nil {
					assert.True(t, result.Valid())
					assert.Nil(t, result.Errors())
				}
			}
		})
	}
}

func TestServiceTriggerValidationErrors(t *testing.T) {
	tests := []ValidationTest{
		{
			name: "missing required path in trigger",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {
							"docker": {"dockerfile": "Dockerfile"}
						},
						"triggers": {
							"my-trigger": {
								"schedule": {
									"cron_expression": "0 0 * * *"
								}
							}
						}
					}
				}
			}`,
			expectError: true,
			err:         "service my-service has an invalid triggers property (my-trigger): The `path` property is required",
		},
		{
			name: "missing schedule in trigger",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {
							"docker": {"dockerfile": "Dockerfile"}
						},
						"triggers": {
							"my-trigger": {
								"path": "/trigger"
							}
						}
					}
				}
			}`,
			expectError: true,
			err:         "service my-service has an invalid triggers property (my-trigger): The `schedule` property is required",
		},
		{
			name: "missing cron_expression in schedule",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {
							"docker": {"dockerfile": "Dockerfile"}
						},
						"triggers": {
							"my-trigger": {
								"path": "/trigger",
								"schedule": {}
							}
						}
					}
				}
			}`,
			expectError: true,
			err:         "service my-service has an invalid my-trigger property (schedule): The `cron_expression` property is required",
		},
		{
			name: "valid trigger with schedule",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {
							"docker": {
								"dockerfile": "Dockerfile",
								"context": "."
							}
						},
						"triggers": {
							"my-trigger": {
								"path": "/trigger",
								"schedule": {
									"cron_expression": "0 0 * * *"
								}
							}
						}
					}
				}
			}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, result, err := ApplicationFromJson(tt.jsonData)
			if err != nil {
				assert.Fail(t, "Failed to parse JSON: %v", err)
			}

			if tt.expectError {
				assert.NotNil(t, result.Errors())
				assert.False(t, result.Valid())

				errorStr := FormatErrors(result)
				assert.Contains(t, errorStr, tt.err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
				if result != nil {
					assert.True(t, result.Valid())
					assert.Nil(t, result.Errors())
				}
			}
		})
	}
}

func TestEntrypointValidationErrors(t *testing.T) {
	tests := []ValidationTest{
		{
			name: "missing routes in entrypoint",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"entrypoints": {
					"my-entrypoint": {}
				}
			}`,
			expectError: true,
			err:         "entrypoint my-entrypoint has an invalid config: The `routes` property is required",
		},
		{
			name: "invalid route name pattern (missing trailing slash)",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"entrypoints": {
					"my-entrypoint": {
						"routes": {
							"api": {
								"name": "my-service"
							}
						}
					}
				}
			}`,
			expectError: true,
			err:         "entrypoint my-entrypoint has an invalid routes property: Missing trailing slash for route api",
		},
		{
			name: "missing required name in route",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"entrypoints": {
					"my-entrypoint": {
						"routes": {
							"api/": {}
						}
					}
				}
			}`,
			expectError: true,
			err:         "entrypoint my-entrypoint has an invalid routes property (api/): The `name` property is required",
		},
		{
			name: "valid entrypoint with routes",
			jsonData: `{
				"name": "test-app",
				"description": "test application",
				"targets": ["aws"],
				"entrypoints": {
					"my-entrypoint": {
						"routes": {
							"api/": {
								"name": "my-service"
							}
						}
					}
				}
			}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, result, err := ApplicationFromJson(tt.jsonData)
			if err != nil {
				assert.Fail(t, "Failed to parse JSON: %v", err)
			}

			if tt.expectError {
				assert.NotNil(t, result.Errors())
				assert.False(t, result.Valid())

				errorStr := FormatErrors(result)
				assert.Contains(t, errorStr, tt.err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
				if result != nil {
					assert.True(t, result.Valid())
					assert.Nil(t, result.Errors())
				}
			}
		})
	}
}
func TestApplicationIsValidErrors(t *testing.T) {
	tests := []ValidationTest{
		{
			name: "duplicate resource names across different types",
			jsonData: `{
				"name": "test-app",
				"description": "test application duplicate resource names different types",
				"targets": ["aws"],
				"services": {
					"my-resource": {
						"container": {
							"docker": {"dockerfile": "Dockerfile"}
						}
					}
				},
				"buckets": {
					"my-resource": {}
				}
			}`,
			expectError: true,
			err:         "bucket name my-resource is already in use by a service",
		},
		{
			name: "duplicate service names",
			jsonData: `{
				"name": "test-app",
				"description": "test application duplicate service names",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {
							"docker": {"dockerfile": "Dockerfile"}
						}
					},
					"my-service": {
						"container": {
							"image": {"id": "nginx:latest"}
						}
					}
				}
			}`,
			expectError: false, // TODO: we should not allow duplicate service names, but currently it just takes the last one in the parsing
			err:         "service name my-service is already in use",
		},
		{
			name: "valid application with unique resource names",
			jsonData: `{
				"name": "test-app",
				"description": "test application valid application",
				"targets": ["aws"],
				"services": {
					"my-service": {
						"container": {
							"docker": {"dockerfile": "Dockerfile"}
						}
					}
				},
				"buckets": {
					"my-bucket": {}
				},
				"entrypoints": {
					"my-entrypoint": {
						"routes": {
							"api/": {
								"name": "my-service"
							}
						}
					}
				}
			}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _, err := ApplicationFromJson(tt.jsonData)
			if err != nil {
				assert.Fail(t, "Failed to parse JSON: %v", err)
			}

			err = app.IsValid()

			if tt.expectError {
				assert.NotNil(t, err)

				assert.Contains(t, err.Error(), tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
