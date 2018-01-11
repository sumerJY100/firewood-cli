package commands

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
	"github.com/xuyz/firewood-cli/utils"
)

func HandlerCmd() cli.Command {
	cmd := cli.Command{}
	cmd.Name = "handler"
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "the name of the struct",
		},
		cli.StringFlag{
			Name:  "pkg",
			Usage: "the name of package",
		},
		cli.StringFlag{
			Name:  "obj",
			Usage: "the name of the file to write the output to (lower of name by default)",
		},
		cli.BoolFlag{
			Name:  "force",
			Usage: "cover file when file is exist (default is false)",
		},
	}
	cmd.Action = Handler

	return cmd
}

type HandlerSt struct {
	Name string
	Pkg  string
}

var handler = `
package {{.Pkg}}

import (
	"net/http"

	"github.com/mholt/binding"
	"github.com/xuyz/firewood"
)

type {{.Name}}Req struct{
	// dto.{{.Name}}
}

// FieldMap implement binding FieldMap interface
func (r *{{.Name}}Req) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{}
}

//Do http handler
func (r *{{.Name}}Req) Do(res http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	resp := &{{.Name}}Req{}

	return resp, 200, nil
}

//{{.Name}} ...
func {{.Name}}() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		firewood.New(res, req, &{{.Name}}Req{}, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
	}
}
`

func Handler(ctx *cli.Context) error {
	t, err := template.New("handler").Parse(handler)
	if err != nil {
		fmt.Println(err)
		return err
	}

	name := ctx.String("name")
	pkg := ctx.String("pkg")
	force := ctx.Bool("force")
	obj := ctx.String("obj")

	if name == "" {
		fmt.Println("need name flag")
		return errors.New("need name flag")
	}

	if pkg == "" {
		pkg = "handler"
	}

	if obj == "" {
		obj = strings.ToLower(pkg + "/" + name + ".go")
	}

	m := HandlerSt{Name: name, Pkg: pkg}
	file := obj
	exist, err := utils.PathExists(file)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if exist && force == false {
		return nil
	}

	e, err := utils.PathExists(m.Pkg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !e {
		os.Mkdir(m.Pkg, os.ModePerm)
	}

	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	if err := t.Execute(f, &m); err != nil {
		return err
	}

	if err := exec.Command("goimports", "-w", file).Run(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
