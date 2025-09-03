package template

import (
	"fmt"
	"testing"

	. "github.com/containifyci/dependabot-templater/pkg/template/testdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndentYAML(t *testing.T) {
	tests := []struct {
		name     string
		spaces   int
		input    string
		expected string
	}{
		{
			name:   "No indentation",
			spaces: 0,
			input: `key: value
list:
  - item1
  - item2`,
			expected: `key: value
list:
  - item1
  - item2`,
		},
		{
			name:   "Basic indentation",
			spaces: 2,
			input: `key: value
list:
  - item1
  - item2`,
			expected: `  key: value
  list:
    - item1
    - item2`,
		},
		{
			name:     "Single line string",
			spaces:   4,
			input:    "simple: value",
			expected: `    simple: value`,
		},
		{
			name:   "Trailing newline handling",
			spaces: 2,
			input: `key: value
`,
			expected: `  key: value`,
		},
		{
			name:     "Empty input",
			spaces:   2,
			input:    "",
			expected: "  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := indentYAML(tt.spaces, tt.input)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestReadTemplate(t *testing.T) {
	cnt := readTemplate("dependabot-terraform.yml.tmpl")
	assert.NotNil(t, cnt)
}

func TestRenderHeader(t *testing.T) {
	for _, test := range []struct {
		kind             string
		expectedTemplate string
	}{
		{
			kind:             "terraform",
			expectedTemplate: HeaderTerraformConfig(),
		},
		{
			kind:             "npm",
			expectedTemplate: HeaderNodeJSConfig(),
		},
	} {
		t.Run(test.kind, func(t *testing.T) {
			tmpl, err := RenderHeader([]string{test.kind})
			require.NoError(t, err)
			assert.Equal(t, test.expectedTemplate, tmpl)
		})
	}
	kinds := []string{"npm", "terraform"}
	output, err := RenderHeader(kinds)
	assert.Nil(t, err)
	assert.NotNil(t, output)
}

func TestRenderDependaBot(t *testing.T) {
	for _, test := range []struct {
		name             string
		result           DependaBotResult
		expectedTemplate string
	}{
		{
			name: "npm-weekly-default",
			result: DependaBotResult{
				Folders:  []string{"projectd"},
				Template: "dependabot-npm.yml.tmpl",
				Registry: "npm-registry",
			},
			expectedTemplate: NodeJSConfig(),
		},
		{
			name: "terraform-weekly-default",
			result: DependaBotResult{
				Folders:  []string{"test_path/projectb"},
				Template: "dependabot-terraform.yml.tmpl",
			},
			expectedTemplate: TerraformConfig(),
		},
		{
			name: "npm-daily",
			result: DependaBotResult{
				Folders:  []string{"projectd"},
				Template: "dependabot-npm.yml.tmpl",
				Registry: "npm-registry",
				Interval: "daily",
			},
			expectedTemplate: NodeJSDailyConfig(),
		},
		{
			name: "npm-monthly",
			result: DependaBotResult{
				Folders:  []string{"projectd"},
				Template: "dependabot-npm.yml.tmpl",
				Registry: "npm-registry",
				Interval: "monthly",
			},
			expectedTemplate: NodeJSMonthlyConfig(),
		},
		{
			name: "terraform-quarterly",
			result: DependaBotResult{
				Folders:  []string{"test_path/projectb"},
				Template: "dependabot-terraform.yml.tmpl",
				Interval: "quarterly",
			},
			expectedTemplate: TerraformQuarterlyConfig(),
		},
		{
			name: "python-yearly",
			result: DependaBotResult{
				Folders:  []string{"projecte"},
				Template: "dependabot-python.yml.tmpl",
				Registry: "python-registry",
				Interval: "yearly",
			},
			expectedTemplate: PythonYearlyConfig(),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			tmpl, err := RenderDependaBot(test.result)
			require.NoError(t, err)
			assert.Equal(t, test.expectedTemplate, tmpl+"\n")
		})
	}
}

func TestRenderDependaBotIntervals(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		day      string
		wantDay  bool
	}{
		{"daily", "daily", "", false},
		{"weekly-default", "weekly", "", true}, // should get default sunday
		{"weekly-monday", "weekly", "monday", true},
		{"monthly", "monthly", "", false},
		{"quarterly", "quarterly", "", false},
		{"semiannually", "semiannually", "", false},
		{"yearly", "yearly", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DependaBotResult{
				Folders:  []string{"test"},
				Template: "dependabot-npm.yml.tmpl",
				Interval: tt.interval,
				Day:      tt.day,
			}
			
			tmpl, err := RenderDependaBot(result)
			require.NoError(t, err)
			
			// Verify interval is present
			assert.Contains(t, tmpl, fmt.Sprintf(`interval: "%s"`, tt.interval))
			
			// Verify day handling
			if tt.wantDay && tt.interval == "weekly" {
				expectedDay := tt.day
				if expectedDay == "" {
					expectedDay = "sunday" // default
				}
				assert.Contains(t, tmpl, fmt.Sprintf(`day: "%s"`, expectedDay))
			} else {
				// For non-weekly intervals, day should not appear
				assert.NotRegexp(t, `day:\s*"`, tmpl)
			}
		})
	}
}
