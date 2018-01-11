package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/xuyz/firewood-cli/commands"
)

var (
	version     = "0.1"
	name        = "firewood-cli"
	description = "firewood 命令行工具"
)

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Version = version
	app.Description = description
	app.Usage = ""
	app.UsageText = "resthub-cli [global options] command [command options] [arguments...]"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "详细日志",
		},
	}

	app.Commands = commands.New()

	app.Run(os.Args)
}
