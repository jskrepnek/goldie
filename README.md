# Goldie
#### A just right web framework for Go
Inspired by [Nancy](http://nancyfx.org/) and the [SDHP](https://github.com/NancyFx/Nancy/wiki/Introduction).  Uses [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux) behind the scenes.
## Warning
Extreme alpha stuff.  Do not use.
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
## Model Binding
Goldie will try to bind parameters from the route, query string, and request body.  Since Go does not retain function argument names during compliation, we can only
reliably bind to a single value type.
```
package main

import (
	"github.com/jskrepnek/goldie"
)

type Widget struct {
	Id int
	Type string
	Strength int
}

func init() {

	// bind to an integer
	goldie.Get["/widget/{id}"] = func(id int) Widget {
		return repo.Get(id)
	}

	// bind from the query string to a string 
	goldie.Get["/widget"] = func(type string) []Widget {
		return repo.GetByType(type)
	}

	// bind to a struct from the request body
	goldie.Post["/widget"] = func(widget Widget) Widget {
		return repo.Add(widget)
	}

	// bind from a route variable and the request to a struct
	goldie.Put["/widget/{id}"] = func(widget Widget) Widget {
		return repo.Update(widget)
	}

	// bind from a route variable to a value type and the request body to a struct
	goldie.Put["/widget/{id}"] = func(id int, widget Widget) Widget {
		return repo.Update(id, widget)
	}
}
```