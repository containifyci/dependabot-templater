package dependabot

import (
	"bytes"
	"strings"

	"github.com/containifyci/dependabot-templater/pkg/search"
	"github.com/containifyci/dependabot-templater/pkg/template"
)

const nilStr = ""

func searchTerraform(path string) ([]string, string, error) {
	foundFolders, err := search.SearchForString(path, "backend")
	if err != nil {
		return nil, nilStr, err
	}
	return foundFolders, "dependabot-terraform.yml.tmpl", nil
}

func searchGithubActions(path string) ([]string, string, error) {
	foundFolders, err := search.SearchForFiles(path, "action.yml", "action.yaml")
	if err != nil {
		return nil, nilStr, err
	}

	foundFolders2, err := search.SearchForFolder(path, ".github/workflows")
	if err != nil {
		return nil, nilStr, err
	}

	foundFolders = append(foundFolders, foundFolders2...)
	return normalizeGithubActions(foundFolders), "dependabot-github-actions.yml.tmpl", nil
}

func normalizeGithubActions(folders []string) []string {
	for i, folder := range folders {
		if folder == ".github/workflows" {
			folders[i] = "/"
		}
	}
	return folders
}

func searchNPM(path string) ([]string, string, error) {
	foundFolders, err := search.SearchForFiles(path, "package.json")
	if err != nil {
		return nil, nilStr, err
	}
	return normalizeNPM(foundFolders), "dependabot-npm.yml.tmpl", nil
}

func normalizeNPM(folders []string) []string {
	validFolder := make([]string, len(folders))
	counter := 0
	for _, folder := range folders {
		if !strings.Contains(folder, "/node_modules/") {
			validFolder[counter] = folder
			counter++
		}
	}
	return validFolder[:counter]
}

func searchPython(path string) ([]string, string, error) {
	foundFolders, err := search.SearchForFiles(path, "requirements.txt", "pyproject.toml")
	if err != nil {
		return nil, nilStr, err
	}
	return foundFolders, "dependabot-python.yml.tmpl", nil
}

func searchGolang(path string) ([]string, string, error) {
	foundFolders, err := search.SearchForFiles(path, "go.mod")
	if err != nil {
		return nil, nilStr, err
	}
	return foundFolders, "dependabot-go.yml.tmpl", nil
}

func searchDocker(path string) ([]string, string, error) {
	foundFolders, err := search.SearchForFiles(path, "Dockerfile")
	if err != nil {
		return nil, nilStr, err
	}
	return foundFolders, "dependabot-docker.yml.tmpl", nil
}

func searchMaven(path string) ([]string, string, error) {
	foundFolders, err := search.SearchForFiles(path, "pom.xml")
	if err != nil {
		return nil, nilStr, err
	}
	return foundFolders, "dependabot-maven.yml.tmpl", nil
}

func searchGradle(path string) ([]string, string, error) {
	foundFolders, err := search.SearchForFiles(path, "build.gradle.kts", "build.gradle")
	if err != nil {
		return nil, nilStr, err
	}
	return foundFolders, "dependabot-gradle.yml.tmpl", nil
}

type DependaBot struct {
	kinds    []string
	rootPath string
	interval string
	day      string
}

type Option func(*DependaBot)

func WithRootPath(rootPath string) Option {
	return func(g *DependaBot) {
		g.rootPath = rootPath
	}
}

func WithKind(kind string) Option {
	return func(g *DependaBot) {
		var kinds []string
		if kind == "all" {
			kinds = []string{"gha", "docker", "terraform", "go", "gradle", "maven", "npm", "python"}
		} else {
			kinds = strings.Split(kind, ",")
		}
		g.kinds = kinds
	}
}

func WithInterval(interval string) Option {
	return func(g *DependaBot) {
		g.interval = interval
	}
}

func WithDay(day string) Option {
	return func(g *DependaBot) {
		g.day = day
	}
}

func New(opts ...Option) *DependaBot {

	bot := &DependaBot{}

	for _, opt := range opts {
		opt(bot)
	}

	return bot
}

func (d *DependaBot) Search(path, kind string) (template.DependaBotResult, error) {
	var folders []string
	var tmplfile string
	var err error
	switch kind {
	case "gha":
		folders, tmplfile, err = searchGithubActions(path)
	case "docker":
		folders, tmplfile, err = searchDocker(path)
	case "go":
		folders, tmplfile, err = searchGolang(path)
	case "gradle":
		folders, tmplfile, err = searchGradle(path)
	case "maven":
		folders, tmplfile, err = searchMaven(path)
	case "npm":
		folders, tmplfile, err = searchNPM(path)
	case "python":
		folders, tmplfile, err = searchPython(path)
	case "terraform":
		folders, tmplfile, err = searchTerraform(path)
	}
	return template.DependaBotResult{Folders: folders, Template: tmplfile, Registry: registry(kind), Interval: d.interval, Day: d.day}, err
}

func registry(kind string) string {
	switch kind {
	case "npm":
		return "npm-registry"
	case "python":
		return "python-registry"
	case "gha":
		fallthrough
	case "docker":
		fallthrough
	case "go":
		fallthrough
	case "gradle":
		fallthrough
	case "maven":
		fallthrough
	case "terraform":
		fallthrough
	default:
		return ""
	}
}

func (d *DependaBot) GenarateConfigFile(path string) ([]string, string) {
	var buffer bytes.Buffer
	packages := make([]string, 0)

	var foundKinds = make([]string, 0)
	for _, kind := range d.kinds {
		result, err := d.Search(path, kind)
		if err != nil {
			panic(err)
		}
		if len(result.Folders) <= 0 {
			continue
		}

		result.Folders = normalizeFolders(result.Folders)
		for i, folder := range result.Folders {
			result.Folders[i] = replacePrefix(folder, d.rootPath, ".")
		}

		packages = append(packages, kind)
		dependabot, err := template.RenderDependaBot(result)
		if err != nil {
			panic(err)
		}
		foundKinds = append(foundKinds, kind)
		buffer.WriteString(dependabot)
	}
	buffer.WriteString("\n")

	var buffer2 bytes.Buffer
	header, err := template.RenderHeader(foundKinds)
	if err != nil {
		panic(err)
	}
	buffer2.WriteString(strings.Trim(header, "\n"))
	// buffer2.WriteString("\n")
	buffer2.WriteString(buffer.String())
	return packages, buffer2.String()
}

func normalizeFolders(folders []string) []string {
	validFolder := make([]string, len(folders))
	counter := 0
	for _, folder := range folders {
		switch {
		case strings.Contains(folder, "/node_modules/"),
			strings.Contains(folder, ".terraform/"),
			strings.Contains(folder, "/vendor/"):
			continue
		}
		validFolder[counter] = folder
		counter++
	}
	return validFolder[:counter]
}

func replacePrefix(input, prefix, replacement string) string {
	str := strings.TrimPrefix(input, prefix)
	if len(str) <= 0 {
		return replacement
	}
	return str
}
