package template

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed *.tmpl
var templates embed.FS

// Custom function to indent YAML lines
func indentYAML(spaces int, yamlStr string) string {
	indentation := strings.Repeat(" ", spaces)
	lines := strings.Split(strings.TrimRight(yamlStr, "\n"), "\n")
	for i, line := range lines {
		lines[i] = indentation + line
	}
	return strings.Join(lines, "\n")
}

func readTemplate(name string) string {
	data, err := templates.ReadFile(name)
	if err != nil {
		fmt.Printf("Error reading template file %s\n", name)
		panic(err)
	}
	return string(data)
}

func RenderHeader(kinds []string) (string, error) {
	var tpl strings.Builder
	funcMap := template.FuncMap{
		"indent": indentYAML,
	}
	tmpl := template.Must(template.New("dependabot-header.yml.tmpl").Funcs(funcMap).Parse(readTemplate("dependabot-header.yml.tmpl")))

	var regs strings.Builder
	for _, kind := range kinds {
		if reg, ok := registries[kind]; ok {
			regs.WriteString(reg)
		}
	}

	entry := DependaBotEntry{
		Registries: regs.String(),
	}
	err := tmpl.Execute(&tpl, entry)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

type DependaBotResult struct {
	Folders  []string
	Template string
	Registry string
	Interval string
	Day      string
}

type DependaBotEntry struct {
	Directory  string
	Registries string
	Interval   string
	Day        string
}

func RenderDependaBot(result DependaBotResult) (string, error) {
	var tpl strings.Builder
	var entries = make([]DependaBotEntry, 0)
	
	// Set default values if not provided for backward compatibility
	interval := result.Interval
	if interval == "" {
		interval = "weekly"
	}
	
	day := result.Day
	if day == "" && interval == "weekly" {
		day = "sunday"
	}
	
	for _, folder := range result.Folders {
		entries = append(entries, DependaBotEntry{
			Directory:  folder, 
			Registries: result.Registry,
			Interval:   interval,
			Day:        day,
		})
	}
	tmpl := template.Must(template.New(result.Template).Parse(readTemplate(result.Template)))
	err := tmpl.Execute(&tpl, entries)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

var registries map[string]string = map[string]string{
	"npm": `npm-registry:
  type: npm-registry
  url: https://europe-west3-npm.pkg.dev/xxxxxxx/npm-registry
  username: "_json_key_base64" # <- Note the username
  password: ${{ secrets.ARTIFACTORY_REGISTRY_SERVICE_ACCOUNT_KEY_BASE64 }} # base64 encoded service account key stored as Github secret. SA must have reader permissions in npm repository.
`,
	"python": `python-registry:
  type: python-index
  url: https://europe-west3-python.pkg.dev/xxxxxxx/python-registry
  username: "_json_key_base64" # <- Note the username
  password: ${{ secrets.ARTIFACTORY_REGISTRY_SERVICE_ACCOUNT_KEY_BASE64 }} # base64 encoded service account key stored as Github secret. SA must have reader permissions in npm repository.
`,
}
