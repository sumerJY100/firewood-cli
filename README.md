# Desc
    the client of firewood

# install

```
    go get github.com/xuyz/firewood-cli
    go install
```

# example

```
Usage of ./firewood-cli:
  -force
    	cover file when file is exist (default is false)
  -name string
    	the name of the struct
  -o string
    	the name of the file to write the output to (lower of name by default)
  -pkg string
    	the name of package (default "handler")
```

**create handler**

```
./firewood-cli -name Index -o index -pkg handler

package handler

import (
	"net/http"

	"github.com/mholt/binding"
	"github.com/xuyz/firewood"
)

type IndexReq struct {
	// dto.Index
}

// FieldMap implement binding FieldMap interface
func (r *IndexReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{}
}

//Do http handler
func (r *IndexReq) Do(res http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	resp := vo.Index{}

	return resp, 200, nil
}

//Index ...
func Index() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		firewood.New(res, req, &LoginReq{}, http.MethodPost)
	}
}
```


# TODO

 - [X] generate handler file
 - [ ] generate handler test file
 - [ ] generate firewood server by configure file
 - [ ] manager firewood server
