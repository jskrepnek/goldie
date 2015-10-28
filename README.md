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
