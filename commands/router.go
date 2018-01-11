package commands

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"strings"

	"bufio"

	"github.com/urfave/cli"
	"github.com/xuyz/firewood-cli/utils"
)

func RouterCmd() cli.Command {
	cmd := cli.Command{}
	cmd.Name = "router"
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "the file of router conf",
		},
		cli.BoolFlag{
			Name:  "force",
			Usage: "cover file when file is exist (default is false)",
		},
	}
	cmd.Action = Router

	return cmd
}

type OneRouter struct {
	Pkg     string
	Path    string
	Handler string
}
type RouterSt struct {
	Pkg     string
	Routers []OneRouter
}

var router = `
package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	r := mux.NewRouter()

	{{range .Routers}}
	r.HandleFunc("{{.Path}}", {{.Pkg}}.{{.Handler}}())
    {{end}}

	return r
}
`

func Router(ctx *cli.Context) error {
	t, err := template.New("router").Parse(router)
	if err != nil {
		fmt.Println(err)
		return err
	}

	name := ctx.String("file")
	out := ctx.String("out")
	force := ctx.Bool("force")

	if name == "" {
		fmt.Println("need file")
		return errors.New("need file")
	}

	f, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	rs := make([]OneRouter, 0)
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')

		if err != nil || io.EOF == err {
			break
		}
		// path, pkg,handler
		ds := strings.Split(line, ",")
		if len(ds) != 3 {
			fmt.Printf("invalid content %s \n", line)
			return errors.New("invalid content")
		}

		rs = append(rs, OneRouter{
			Path:    ds[0],
			Pkg:     ds[1],
			Handler: ds[2][:len(ds[2])-1],
		})

	}
	if out == "" {
		out = strings.ToLower("app/router.go")
	}

	m := RouterSt{
		Pkg:     "app",
		Routers: rs,
	}
	exist, err := utils.PathExists(out)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if exist && force == false {
		return nil
	}

	fw, err := os.Create(out)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer fw.Close()

	if err := t.Execute(fw, &m); err != nil {
		return err
	}

	if err := exec.Command("goimports", "-w", out).Run(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
