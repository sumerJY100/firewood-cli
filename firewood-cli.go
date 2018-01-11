package main

import (
	"flag"
	"html/template"
	"os"
	"os/exec"
	"strings"
)

var handler = `
package {{.Pkg}}


type {{.Name}}Req struct{
	dto.{{.Name}}
}

// FieldMap implement binding FieldMap interface
func (r *{{.Name}}Req) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{}
}

//DoProcess 业务逻辑
func (r *{{.Name}}Req) Do(res http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	resp := vo.{{.Name}}{}

	return resp, 200, nil
}

//{{.Name}} ...
func {{.Name}}() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		firewood.New(res, req, &LoginReq{}, http.MethodPost)
	}
}
`

type MMM struct {
	Pkg  string
	Name string
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	t, err := template.New("handler").Parse(handler)
	if err != nil {
		panic(err)
	}

	name := flag.String("name", "", "the name of the struct")
	pkg := flag.String("pkg", "handler", "the name of package")
	force := flag.Bool("force", false, "cover file when file is exist (default is false)")
	obj := flag.String("o", "",
		"the name of the file to write the output to (lower of name by default)")
	flag.Parse()

	if *name == "" {
		flag.Usage()
		return
	}

	if *obj == "" {
		*obj = *name
	}

	m := MMM{Name: *name, Pkg: *pkg}
	file := strings.ToLower(m.Pkg + "/" + m.Name + ".go")
	exist, err := PathExists(file)
	if err != nil {
		panic(err)
	}

	if exist && *force == false {
		return
	}

	e, err := PathExists(m.Pkg)
	if err != nil {
		panic(err)
	}
	if !e {
		os.Mkdir(m.Pkg, os.ModePerm)
	}

	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := t.Execute(f, &m); err != nil {
		panic(err)
	}

	if err := exec.Command("goimports", "-w", "handler").Run(); err != nil {
		panic(err)
	}
}
