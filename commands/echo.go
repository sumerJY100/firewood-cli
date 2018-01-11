package commands

import (
	"errors"

	"os"

	"github.com/urfave/cli"
)

func EchoCmd() cli.Command {
	cmd := cli.Command{}
	cmd.Name = "echo"
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input",
			Usage: "input data",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "save to file",
		},
	}
	cmd.Action = Echo

	return cmd
}

func Echo(ctx *cli.Context) error {
	input := ctx.String("input")
	file := ctx.String("file")

	if file == "" {
		return errors.New("need file")
	}

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(input + "\n")

	return nil
}
