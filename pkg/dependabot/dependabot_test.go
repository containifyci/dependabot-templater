package dependabot

import (
	"testing"

	. "github.com/containifyci/dependabot-templater/pkg/dependabot/testdata"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	for _, test := range []struct {
		name             string
		path             string
		kind             string
		expectedPath     string
		expectedTemplate string
	}{
		{
			name:             "folder in current directory",
			path:             "./test_path/",
			kind:             "terraform",
			expectedTemplate: TerraformConfig(),
		},
		{
			name:             "single folder in current directory",
			path:             "./test_path/projectb",
			kind:             "terraform",
			expectedTemplate: TerraformConfig(),
		},
		{
			name:             "folder relative to the current directory level 1",
			path:             "../dependabot/test_path/",
			kind:             "terraform",
			expectedTemplate: TerraformConfig(),
		},
		{
			name:             "folder relative to the current directory level 2",
			path:             "../../pkg/dependabot/test_path/",
			kind:             "terraform",
			expectedTemplate: TerraformConfig2(),
		},
		{
			name:             "Github Action folder relative to the current directory level 2",
			path:             "../../pkg/dependabot/test_path/",
			kind:             "gha",
			expectedTemplate: GithubActionConfig(),
		},
		{
			name:             "Npm folder relative to the current directory level 2",
			path:             "../../pkg/dependabot/test_path/",
			kind:             "npm",
			expectedTemplate: NodeJSConfig(),
		},
		{
			name:             "Python folder relative to the current directory level 2",
			path:             "../../pkg/dependabot/test_path/",
			kind:             "python",
			expectedTemplate: PythonConfig(),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			bot := New(WithKind(test.kind), WithRootPath("dependabot/test_path/"))
			_, tmpl := bot.GenarateConfigFile(test.path, "")

			assert.Equal(t, test.expectedTemplate, tmpl)
		})
	}
}

func TestReplacePrefix(t *testing.T) {
	for _, test := range []struct {
		name         string
		path         string
		expectedPath string
	}{
		{
			name:         "folder in current directory",
			path:         "./test_path/",
			expectedPath: "./test_path/",
		},
		{
			name:         "single folder in current directory",
			path:         "./test_path/projectb",
			expectedPath: "./test_path/projectb",
		},
		{
			name:         "folder relative to the current directory level 1",
			path:         "../dependabot/test_path/",
			expectedPath: "../dependabot/test_path/",
		},
		{
			name:         "folder of root path",
			path:         "dependabot/test_path/",
			expectedPath: ".",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := replacePrefix(test.path, "dependabot/test_path/", ".")

			assert.Equal(t, test.expectedPath, path)
		})
	}
}

func TestSearch(t *testing.T) {
	for _, test := range []struct {
		name             string
		folders          []string
		kind             string
		expectedPath     string
		expectedTemplate string
	}{
		{
			name:    "terraform",
			folders: []string{"test_path/projectb"},
			kind:    "terraform",
		}, {
			name:    "Github Action",
			folders: []string{"test_path/projectc"},
			kind:    "gha",
		}, {
			name:    "NodeJS",
			folders: []string{"test_path/projectd"},
			kind:    "npm",
		}, {
			name:    "Python",
			folders: []string{"test_path/projecte", "test_path/projectf"},
			kind:    "python",
		}, {
			name:    "GoLang",
			folders: []string{"test_path/projectgo"},
			kind:    "go",
		}, {
			name:    "Docker",
			folders: []string{"test_path/projecth"},
			kind:    "docker",
		}, {
			name:    "Maven",
			folders: []string{"test_path/projecti"},
			kind:    "maven",
		}, {
			name:    "Gradle",
			folders: []string{"test_path/projectj", "test_path/projectk"},
			kind:    "gradle",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			bot := New(WithKind(test.kind), WithRootPath("dependabot/test_path/"))
			result, err := bot.Search("./test_path/", test.kind)
			assert.NoError(t, err)
			assert.Equal(t, test.folders, result.Folders)
		})
	}
}

func TestNormalizeFolders(t *testing.T) {
	folders := []string{
		"test_path/projecta",
		"test_path/projectb",
		"test_path/projectc",
		"test_path/projectd",
		"test_path/projectc/node_modules/package/awesome",
		"test_path/projectc/.terraform/package/awesome",
		"test_path/projectc/vendor/package/awesome",
	}

	folders = normalizeFolders(folders)

	assert.Equal(t, []string{"test_path/projecta",
		"test_path/projectb",
		"test_path/projectc",
		"test_path/projectd"}, folders)
}
