package goldie

type Route interface{}
type Routes map[Route]Action

var (
	Get    Routes = Routes{}
	Post   Routes = Routes{}
	Put    Routes = Routes{}
	Delete Routes = Routes{}
)
