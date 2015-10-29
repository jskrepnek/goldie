# Goldie
#### A just right web framework for Go
Inspired by [Nancy](http://nancyfx.org/) and the [SDHP](https://github.com/NancyFx/Nancy/wiki/Introduction).  Uses [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux) behind the scenes.
## Install
```
go get github.com/jskrepnek/goldie
```
## Serve
```go
import (
  "github.com/jskrepnek/goldie"
)

func main() {
  
  goldie.Get["/hello/{name}"] = func (string name) string {
    return "Hello " + name
  }
  
  goldie.Serve()
}

```
## Supports
* Model binding
* Go Http templates

## Modules
Use modules to group related methods together and tie into the dependency injection system.

```go
package main

import (
	"github.com/jskrepnek/goldie"
	"net/http"
)

type TestModule struct {
}

func (this *TestModule) Construct(m *goldie.Module) {

	// the base route for all actions
	m.Path = "/test"

	// invoked before every module action
	m.Before.Add(func(req *http.Request) goldie.Response {
		log.Println("Before")
		return goldie.Response{}
	})

	// invoked after every module action
	m.After.Add(func(r *goldie.Response) {
		log.Println("After")
	})

	// registered with the path /test/spot
	m.Get["/spot"] = func() string {
		return "spot"
	}
}

func init() {
	goldie.AddModule(&TestModule{})
}
```