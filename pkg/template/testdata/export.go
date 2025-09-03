package testdata

import (
	"embed"
)

//go:embed *.yaml
var files embed.FS

func TerraformConfig() string {
	return Content("dependabot-terraform.yaml")
}

func TerraformConfig2() string {
	return Content("dependabot-terraform2.yaml")
}

func GithubActionConfig() string {
	return Content("debendabot-githubactions.yaml")
}

func NodeJSConfig() string {
	return Content("dependabot-npm.yaml")
}

func PythonConfig() string {
	return Content("dependabot-python.yaml")
}

func NodeJSDailyConfig() string {
	return Content("dependabot-npm-daily.yaml")
}

func NodeJSMonthlyConfig() string {
	return Content("dependabot-npm-monthly.yaml")
}

func TerraformQuarterlyConfig() string {
	return Content("dependabot-terraform-quarterly.yaml")
}

func PythonYearlyConfig() string {
	return Content("dependabot-python-yearly.yaml")
}

func HeaderNodeJSConfig() string {
	return Content("dependabot-header-npm.yaml")
}

func HeaderTerraformConfig() string {
	return Content("dependabot-header-terraform.yaml")
}

func Content(name string) string {
	b, err := files.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(b)
}
