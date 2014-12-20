package main

import (
	"github.com/codegangsta/cli"
	"os"
)

const DOCKER_FILE string = "Dockerfile"

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "multi-docker-build"
	app.Usage = "execute multiple docker build"
	app.HideHelp = true
	app.Flags = Flags()
	app.Action = Command

	cli.AppHelpTemplate = HelpTemplate()
	app.Run(os.Args)
}
