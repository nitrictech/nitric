package client

import (
	"testing"
)

func TestNewResourceNameNormalizer(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		description string
	}{
		{
			name:        "valid name",
			input:       "my-bucket",
			wantErr:     false,
			description: "should accept valid name",
		},
		{
			name:        "empty name",
			input:       "",
			wantErr:     true,
			description: "should reject empty name",
		},
		{
			name:        "starts with number",
			input:       "1bucket",
			wantErr:     true,
			description: "should reject name starting with number",
		},
		{
			name:        "starts with special char",
			input:       "@bucket",
			wantErr:     true,
			description: "should reject name starting with special character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewResourceNameNormalizer(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewResourceNameNormalizer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourceNameNormalizer_Methods(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantParts   []string
		wantUnmod   string
		wantPascal  string
		wantCamel   string
		wantSnake   string
		wantKebab   string
		description string
	}{
		{
			name:        "simple name",
			input:       "my-bucket",
			wantParts:   []string{"my", "bucket"},
			wantUnmod:   "my-bucket",
			wantPascal:  "MyBucket",
			wantCamel:   "myBucket",
			wantSnake:   "my_bucket",
			wantKebab:   "my-bucket",
			description: "should handle simple hyphenated name",
		},
		{
			name:        "complex name",
			input:       "my-super-cool-bucket",
			wantParts:   []string{"my", "super", "cool", "bucket"},
			wantUnmod:   "my-super-cool-bucket",
			wantPascal:  "MySuperCoolBucket",
			wantCamel:   "mySuperCoolBucket",
			wantSnake:   "my_super_cool_bucket",
			wantKebab:   "my-super-cool-bucket",
			description: "should handle complex hyphenated name",
		},
		{
			name:        "mixed separators",
			input:       "my_super-cool.bucket",
			wantParts:   []string{"my", "super", "cool", "bucket"},
			wantUnmod:   "my_super-cool.bucket",
			wantPascal:  "MySuperCoolBucket",
			wantCamel:   "mySuperCoolBucket",
			wantSnake:   "my_super_cool_bucket",
			wantKebab:   "my-super-cool-bucket",
			description: "should handle mixed separators",
		},
		{
			name:        "multiple separators",
			input:       "my__super--cool..bucket",
			wantParts:   []string{"my", "super", "cool", "bucket"},
			wantUnmod:   "my__super--cool..bucket",
			wantPascal:  "MySuperCoolBucket",
			wantCamel:   "mySuperCoolBucket",
			wantSnake:   "my_super_cool_bucket",
			wantKebab:   "my-super-cool-bucket",
			description: "should handle multiple consecutive separators",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalizer, err := NewResourceNameNormalizer(tt.input)
			if err != nil {
				t.Fatalf("NewResourceNameNormalizer() error = %v", err)
			}

			// Test Parts
			gotParts := normalizer.Parts()
			if !slicesEqual(gotParts, tt.wantParts) {
				t.Errorf("Parts() = %v, want %v", gotParts, tt.wantParts)
			}

			// Test Unmodified
			if got := normalizer.Unmodified(); got != tt.wantUnmod {
				t.Errorf("Unmodified() = %v, want %v", got, tt.wantUnmod)
			}

			// Test PascalCase
			if got := normalizer.PascalCase(); got != tt.wantPascal {
				t.Errorf("PascalCase() = %v, want %v", got, tt.wantPascal)
			}

			// Test CamelCase
			if got := normalizer.CamelCase(); got != tt.wantCamel {
				t.Errorf("CamelCase() = %v, want %v", got, tt.wantCamel)
			}

			// Test SnakeCase
			if got := normalizer.SnakeCase(); got != tt.wantSnake {
				t.Errorf("SnakeCase() = %v, want %v", got, tt.wantSnake)
			}

			// Test KebabCase
			if got := normalizer.KebabCase(); got != tt.wantKebab {
				t.Errorf("KebabCase() = %v, want %v", got, tt.wantKebab)
			}
		})
	}
}

// Helper function to compare string slices
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
