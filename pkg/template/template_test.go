package template

import (
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
		result           DependaBotResult
		expectedTemplate string
	}{
		{
			result: DependaBotResult{
				Folders:  []string{"projectd"},
				Template: "dependabot-npm.yml.tmpl",
				Registry: "npm-registry",
			},
			expectedTemplate: NodeJSConfig(),
		},
		{
			result: DependaBotResult{
				Folders:  []string{"test_path/projectb"},
				Template: "dependabot-terraform.yml.tmpl",
			},
			expectedTemplate: TerraformConfig(),
		},
	} {
		t.Run(test.result.Template, func(t *testing.T) {
			tmpl, err := RenderDependaBot(test.result)
			require.NoError(t, err)
			assert.Equal(t, test.expectedTemplate, tmpl+"\n")
		})
	}
	kinds := []string{"npm", "terraform"}
	output, err := RenderHeader(kinds)
	assert.Nil(t, err)
	assert.NotNil(t, output)
}
