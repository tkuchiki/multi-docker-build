package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
)

const DEFAULT_TAG string = "latest"

var options string
var tag string
var config string
var image string
var quiet bool

func Print(msg string) {
	if !quiet {
		log.Fatal(msg)
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		cli.HelpFlag,
		cli.StringFlag{
			Name:  "options, o",
			Usage: "docker build options",
		},
		cli.StringFlag{
			Name:  "tag, t",
			Usage: "docker image tag",
			Value: DEFAULT_TAG,
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "config file",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "quiet mode",
		},
	}
}

func Command(c *cli.Context) {
	options = c.String("options")
	tag = c.String("tag")
	quiet = c.Bool("quiet")
	config = c.String("config")

	for _, path := range c.Args() {
		if FileNotExists(path) {
			Print(path + ": No such file or directory")
			continue
		}

		if IsDir(path) {
			ch := GoWalk(path)

			for p := range ch {
				Symlink(p)
				image := GetImage(p)
				tag := GetTag(tag, image)
				tagOption := fmt.Sprintf("%s:%s", image, tag)
				DockerBuild(p, tagOption, options, quiet)
				RemoveFile(DOCKER_FILE)
			}
		} else {
			Symlink(path)
			image := GetImage(path)
			tag := GetTag(tag, image)
			tagOption := fmt.Sprintf("%s:%s", image, tag)
			DockerBuild(path, tagOption, options, quiet)
			RemoveFile(DOCKER_FILE)
		}
	}
}

func GetTag(tag, image string) string {
	if tag == DEFAULT_TAG && FileExists(config) {
		configJson := DecodeJson(config).(map[string]interface{})

		if _, ok := configJson[image]; ok {
			return configJson[image].(string)
		} else {
			fmt.Println(ok)
		}
	}

	return tag
}

func HelpTemplate() string {
	return `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [options] [arguments...]

VERSION:
   {{.Version}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`
}
