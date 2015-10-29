package goldie

import (
	"github.com/gorilla/mux"
	"net/http"
)

func buildRouter() *mux.Router {
	r := mux.NewRouter()

	for route, action := range Get {
		r.HandleFunc(route.(string), newHandler(action))
	}
	return r
}

func Serve() {
	http.ListenAndServe(":8080", buildRouter())
}
