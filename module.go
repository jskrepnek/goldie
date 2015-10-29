package goldie

import (
	"log"
	"net/http"
)

type Response struct{}

type BeforeInterceptor func(req *http.Request) Response
type AfterInterceptor func(*Response)

type BeforeInterceptors []BeforeInterceptor
type AfterInterceptors []AfterInterceptor

func (c *BeforeInterceptors) Add(e BeforeInterceptor) {
	log.Println("Module: Adding before incterceptor")
	*c = append(*c, e)
}

func (c *AfterInterceptors) Add(e AfterInterceptor) {
	log.Println("Module: Adding after incterceptor")
	*c = append(*c, e)
}

type Module struct {
	Path   string
	Before BeforeInterceptors
	After  AfterInterceptors
	Get    Routes
	Post   Routes
	Put    Routes
	Delete Routes
}

func newModule() *Module {
	return &Module{
		Before: BeforeInterceptors{},
		After:  AfterInterceptors{},
		Get:    Routes{},
		Post:   Routes{},
		Put:    Routes{},
		Delete: Routes{},
	}
}

var modules []Module = []Module{}

type ModuleConstructor interface {
	Construct(*Module)
}

func AddModule(mc ModuleConstructor) {
	log.Printf("AddModule")
	module := newModule()
	mc.Construct(module)
	modules = append(modules, *module)
}
