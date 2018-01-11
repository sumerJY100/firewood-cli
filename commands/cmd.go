package commands

import (
	"github.com/urfave/cli"
)

func New() []cli.Command {
	return []cli.Command{
		HandlerCmd(),
		RouterCmd(),
		EchoCmd(),
	}
}
