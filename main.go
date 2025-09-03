package main

import (
	"os"

	"github.com/containifyci/dependabot-templater/pkg/dependabot"
)

func main() {
	args := os.Args[1:]
	kind := args[0]
	path := args[1]
	prefix := arg(args, 2, "")
	bot := dependabot.New(dependabot.WithKind(kind))
	_, dependabot := bot.GenarateConfigFile(path, prefix)
	_, err := os.Stdout.WriteString(dependabot)
	if err != nil {
		panic(err)
	}
}

func arg(args []string, pos int, def string) string {
	if len(args) > pos {
		return args[pos]
	}
	return def
}
